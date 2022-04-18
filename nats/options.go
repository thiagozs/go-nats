package nats

import (
	"time"

	"github.com/nats-io/stan.go"
)

type Options func(*Config)
type OptionsHandler func(*ConfigOptionsHandler)

type ConfigOptionsHandler struct {
	DurableName, MaxInFlight,
	AckWait, ManualAckMode stan.SubscriptionOption
}

type Config struct {
	ClusterName     string
	ClusterID       string
	PingParams      PingParameter
	PubAckWait      time.Duration
	NatsURL         string
	NatsPort        string
	MaxInFlight     int
	ConnLostHandler func(conn stan.Conn, reason error)
	MessageHandler  func(msg *stan.Msg)
	ConnNats        stan.Conn
	CfgHandler      ConfigOptionsHandler
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

func NastURL(url string) Options {
	return func(o *Config) {
		o.NatsURL = url
	}
}

func NastPort(port string) Options {
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
 * ConfigOptionsHandler ***************************
 * ------------------------------------------------
 * Configuration for a custom handler. First
 * implementation request and settings.
 */

func OptsHandler(opts ...OptionsHandler) Options {
	return func(o *Config) {
		opt := ConfigOptionsHandler{}
		for _, oo := range opts {
			oo(&opt)
		}
	}
}

func SetAckWait(duration time.Duration) OptionsHandler {
	return func(o *ConfigOptionsHandler) {
		o.AckWait = stan.AckWait(duration)
	}
}

func SetDurableName(name string) OptionsHandler {
	return func(o *ConfigOptionsHandler) {
		o.DurableName = stan.DurableName(name)
	}
}

func SetMaxInFlight(numMessage int) OptionsHandler {
	return func(o *ConfigOptionsHandler) {
		o.MaxInFlight = stan.MaxInflight(numMessage)
	}
}

func SetManualAckMode(isManual bool) OptionsHandler {
	return func(o *ConfigOptionsHandler) {
		if isManual {
			o.ManualAckMode = stan.SetManualAckMode()
		}
	}
}
