package nats

import (
	"time"

	"github.com/nats-io/stan.go"
	"github.com/newrelic/go-agent/v3/integrations/nrstan"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Options func(*Config)

type SubsOptHandler func(*SubsOptionsHandler)

type SubsOptionsHandler struct {
	DurableName, MaxInFlight,
	AckWait, ManualAckMode,
	StartWithLastReceived,
	StartAtSequence stan.SubscriptionOption
	SubWrapper func(ch chan<- interface{}) func(msg *stan.Msg)
}

type Config struct {
	ClusterName     string
	ClusterID       string
	PingParams      PingParameter
	PubAckWait      time.Duration
	NatsServer      string
	NatsURL         string
	NatsPort        string
	MaxInFlight     int
	ConnLostHandler func(conn stan.Conn, reason error)
	MessageHandler  func(msg *stan.Msg)
	ConnNats        stan.Conn
	SubsOpts        SubsOptionsHandler
}

func NewOptions(opts ...Options) Config {
	opt := Config{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func ClusterName(name string) Options {
	return func(o *Config) {
		o.ClusterName = name
	}
}

func ClusterID(id string) Options {
	return func(o *Config) {
		o.ClusterID = id
	}
}

func PingParams(params PingParameter) Options {
	return func(o *Config) {
		o.PingParams = params
	}
}

func PubAckWait(time time.Duration) Options {
	return func(o *Config) {
		o.PubAckWait = time
	}
}

func NatsURL(url string) Options {
	return func(o *Config) {
		o.NatsURL = url
	}
}

func NatsServer(server string) Options {
	return func(o *Config) {
		o.NatsServer = server
	}
}

func NatsPort(port string) Options {
	return func(o *Config) {
		o.NatsPort = port
	}
}

func ConnLostHandler(fn func(con stan.Conn, reason error)) Options {
	return func(o *Config) {
		o.ConnLostHandler = fn
	}
}

func ConnNats(conn stan.Conn) Options {
	return func(o *Config) {
		o.ConnNats = conn
	}
}

func MessageHandler(fn func(msg *stan.Msg)) Options {
	return func(o *Config) {
		o.MessageHandler = fn
	}
}

func MaxInFlight(num int) Options {
	return func(o *Config) {
		o.MaxInFlight = num
	}
}

/*
 * NewSubsOpts ***************************
 * ------------------------------------------------
 * Configuration for a custom SubsOptionsHandler. First
 * implementation request and settings.
 */

func NewSubsOpts(opts ...SubsOptHandler) Options {
	return func(o *Config) {
		opt := SubsOptionsHandler{}
		for _, oo := range opts {
			oo(&opt)
		}
		o.SubsOpts = opt
	}
}

func SetAckWait(duration time.Duration) SubsOptHandler {
	return func(o *SubsOptionsHandler) {
		o.AckWait = stan.AckWait(duration)
	}
}

func SetDurableName(name string) SubsOptHandler {
	return func(o *SubsOptionsHandler) {
		o.DurableName = stan.DurableName(name)
	}
}

func SetMaxInFlight(numMessage int) SubsOptHandler {
	return func(o *SubsOptionsHandler) {
		o.MaxInFlight = stan.MaxInflight(numMessage)
	}
}

func SetManualAckMode(isManual bool) SubsOptHandler {
	return func(o *SubsOptionsHandler) {
		if isManual {
			o.ManualAckMode = stan.SetManualAckMode()
			return
		}
		o.ManualAckMode = nil
	}
}

func SetSubWrapper(nr *newrelic.Application) SubsOptHandler {
	return func(o *SubsOptionsHandler) {
		o.SubWrapper = func(ch chan<- interface{}) func(msg *stan.Msg) {
			return nrstan.StreamingSubWrapper(nr, func(msg *stan.Msg) {
				if o.ManualAckMode == nil {
					msg.Ack()
				}
				ch <- MessageCh{msg, msg.Data}
			})
		}
	}
}

func SetStartWithLastReceived(isManual bool) SubsOptHandler {
	return func(o *SubsOptionsHandler) {
		if isManual {
			o.StartWithLastReceived = stan.StartWithLastReceived()
			return
		}
		o.StartWithLastReceived = nil
	}
}

func SetStartAtSequence(sequence uint64) SubsOptHandler {
	return func(o *SubsOptionsHandler) {
		if sequence > 0 {
			o.StartAtSequence = stan.StartAtSequence(sequence)
			return
		}
		o.StartAtSequence = nil
	}
}
