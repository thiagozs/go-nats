package streaming

import (
	"fmt"

	"github.com/thiagozs/go-nats/nats"

	"github.com/nats-io/stan.go"
)

type ServerType int

const (
	DEFAULT_NAME ServerType = iota
)

func (s ServerType) String() string {
	return [...]string{"DEFAULT"}[s]
}

type StreamOpts func(s *StreamOptsParams) error

type StreamOptsParams struct {
	serverName string
}

func NewStreamOpts(opts ...StreamOpts) StreamOptsParams {
	opt := StreamOptsParams{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func SetServerName(name string) StreamOpts {
	return func(o *StreamOptsParams) error {
		o.serverName = name
		return nil
	}
}

type repoStream struct {
	nats   map[string]nats.NatsServiceRepo
	chann  map[string]map[string]chan interface{}
	subs   map[string]map[string]stan.Subscription
	params StreamOptsParams
}

// NewSevice start a streming wrapper
func NewService(natsService nats.NatsServiceRepo,
	opts ...StreamOpts) *repoStream {
	params := NewStreamOpts(opts...)

	if params.serverName == "" {
		params.serverName = DEFAULT_NAME.String()
	}

	nats := make(map[string]nats.NatsServiceRepo)
	nats[params.serverName] = natsService

	return &repoStream{
		nats:   nats,
		chann:  make(map[string]map[string]chan interface{}),
		subs:   make(map[string]map[string]stan.Subscription),
		params: params,
	}
}

// Register method for subscription a new nast instance
func (s *repoStream) Register(natsService nats.NatsServiceRepo) *repoStream {
	s.nats[s.params.serverName] = natsService
	return s
}

// Subscribe register a cliente in a channel broadcasting
func (s *repoStream) Subscribe(channel string) error {
	_, ok := s.chann[s.params.serverName]
	if !ok {
		s.chann[s.params.serverName] = make(map[string]chan interface{})
		_, ok := s.chann[s.params.serverName][channel]
		if !ok {
			s.chann[s.params.serverName][channel] = make(chan interface{})
		}
	} else {
		_, ok := s.chann[s.params.serverName][channel]
		if !ok {
			s.chann[s.params.serverName][channel] = make(chan interface{})
		}
	}

	cfg := s.nats[s.params.serverName].GetOpts()
	subOpts := s.getSubsOpts(cfg)

	var sub stan.Subscription
	var err error

	if cfg.SubsOpts.SubWrapper != nil {
		// With Wrapper NewRelic
		sub, err = s.nats[s.params.serverName].Subscribe(channel,
			cfg.SubsOpts.SubWrapper(s.chann[s.params.serverName][channel]),
			subOpts...)
	} else {
		// Without NewRelic
		sub, err = s.nats[s.params.serverName].Subscribe(channel, func(msg *stan.Msg) {
			s.chann[s.params.serverName][channel] <- nats.MessageCh{Msg: msg, Data: msg.Data}

			if cfg.SubsOpts.ManualAckMode == nil {
				msg.Ack()
			}
		}, subOpts...)
	}

	if err != nil {
		return err
	}

	_, ok = s.subs[s.params.serverName]
	if !ok {
		s.subs[s.params.serverName] = make(map[string]stan.Subscription)
		_, ok := s.subs[s.params.serverName][channel]
		if !ok {
			s.subs[s.params.serverName][channel] = sub
		}
	} else {
		_, ok := s.subs[s.params.serverName][channel]
		if !ok {
			s.subs[s.params.serverName][channel] = sub
		}
	}

	return nil
}

// SubscribeAll register a array of subscribers
func (s *repoStream) SubscribeAll(arr []interface{}) {
	// TODO: improvement this method
	for _, f := range arr {
		if v, ok := f.(SubscribeType); ok {
			if err := v; err != nil {
				fmt.Println(err)
			}
		}
	}
}

// Unsubscribe unregister a client of channel and server
func (s *repoStream) Unsubscribe(channel string) error {
	return s.subs[s.params.serverName][channel].Unsubscribe()
}

// Publishsend a payload to channel
func (s *repoStream) Publish(payload nats.Message) error {
	return s.nats[s.params.serverName].Publish(payload)
}

// GetMessage return a message from channel broadcast
func (s *repoStream) GetMessage(channel string) chan nats.MessageCh {
	rr := make(chan nats.MessageCh)
	go func() {
		for {
			select {
			case msg := <-s.chann[s.params.serverName][channel]:
				msgch, ok := msg.(nats.MessageCh)
				if !ok {
					continue
				}
				rr <- msgch
			default:
				continue
			}
		}
	}()

	return rr
}

func (s *repoStream) getSubsOpts(opts *nats.Config) []stan.SubscriptionOption {

	subOpts := []stan.SubscriptionOption{}

	if opts.SubsOpts.ManualAckMode != nil {
		subOpts = append(subOpts, opts.SubsOpts.ManualAckMode)
	}

	if opts.SubsOpts.StartWithLastReceived != nil {
		subOpts = append(subOpts, opts.SubsOpts.StartWithLastReceived)
	}

	if opts.SubsOpts.DurableName != nil {
		subOpts = append(subOpts, opts.SubsOpts.DurableName)
	}

	if opts.SubsOpts.MaxInFlight != nil {
		subOpts = append(subOpts, opts.SubsOpts.MaxInFlight)
	}

	if opts.SubsOpts.StartAtSequence != nil {
		subOpts = append(subOpts, opts.SubsOpts.StartAtSequence)
	}

	return subOpts
}

func (s *repoStream) Close() error {
	for ch := range s.subs[s.params.serverName] {
		s.subs[s.params.serverName][ch].Unsubscribe()
	}
	return s.nats[s.params.serverName].Close()
}
