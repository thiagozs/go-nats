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

type RepoStream struct {
	Nats    map[string]nats.NatsServiceRepo
	Chann   map[string]map[string]chan interface{}
	Subs    map[string]map[string]stan.Subscription
	Params  StreamOptsParams
	GetChan chan nats.MessageCh
}

// NewSevice start a streming wrapper
func NewService(natsService nats.NatsServiceRepo,
	opts ...StreamOpts) *RepoStream {
	params := NewStreamOpts(opts...)

	if params.serverName == "" {
		params.serverName = DEFAULT_NAME.String()
	}

	natsr := make(map[string]nats.NatsServiceRepo)
	natsr[params.serverName] = natsService

	return &RepoStream{
		Nats:    natsr,
		Chann:   make(map[string]map[string]chan interface{}),
		Subs:    make(map[string]map[string]stan.Subscription),
		Params:  params,
		GetChan: make(chan nats.MessageCh),
	}
}

// Register method for subscription a new nast instance
func (s *RepoStream) Register(natsService nats.NatsServiceRepo) *RepoStream {
	s.Nats[s.Params.serverName] = natsService
	return s
}

// Subscribe register a cliente in a channel broadcasting
func (s *RepoStream) Subscribe(channel string) error {
	_, ok := s.Chann[s.Params.serverName]
	if !ok {
		s.Chann[s.Params.serverName] = make(map[string]chan interface{})
		_, ok := s.Chann[s.Params.serverName][channel]
		if !ok {
			s.Chann[s.Params.serverName][channel] = make(chan interface{})
		}
	} else {
		_, ok := s.Chann[s.Params.serverName][channel]
		if !ok {
			s.Chann[s.Params.serverName][channel] = make(chan interface{})
		}
	}

	cfg := s.Nats[s.Params.serverName].GetOpts()
	subOpts := s.getSubsOpts(cfg)

	var sub stan.Subscription
	var err error

	if cfg.SubsOpts.SubWrapper != nil {
		// With Wrapper NewRelic
		sub, err = s.Nats[s.Params.serverName].Subscribe(channel,
			cfg.SubsOpts.SubWrapper(s.Chann[s.Params.serverName][channel]),
			subOpts...)
	} else {
		// Without NewRelic
		sub, err = s.Nats[s.Params.serverName].Subscribe(channel, func(msg *stan.Msg) {
			s.Chann[s.Params.serverName][channel] <- nats.MessageCh{Msg: msg, Data: msg.Data}

			if cfg.SubsOpts.ManualAckMode == nil {
				msg.Ack()
			}
		}, subOpts...)
	}

	if err != nil {
		return err
	}

	_, ok = s.Subs[s.Params.serverName]
	if !ok {
		s.Subs[s.Params.serverName] = make(map[string]stan.Subscription)
		_, ok := s.Subs[s.Params.serverName][channel]
		if !ok {
			s.Subs[s.Params.serverName][channel] = sub
		}
	} else {
		_, ok := s.Subs[s.Params.serverName][channel]
		if !ok {
			s.Subs[s.Params.serverName][channel] = sub
		}
	}

	return nil
}

// QueueSubscribe register a client in a channel broadcasting
func (s *RepoStream) QueueSubscribe(channel, qgroup string) error {
	_, ok := s.Chann[s.Params.serverName]
	if !ok {
		s.Chann[s.Params.serverName] = make(map[string]chan interface{})
		_, ok := s.Chann[s.Params.serverName][channel]
		if !ok {
			s.Chann[s.Params.serverName][channel] = make(chan interface{})
		}
	} else {
		_, ok := s.Chann[s.Params.serverName][channel]
		if !ok {
			s.Chann[s.Params.serverName][channel] = make(chan interface{})
		}
	}

	cfg := s.Nats[s.Params.serverName].GetOpts()
	subOpts := s.getSubsOpts(cfg)

	var sub stan.Subscription
	var err error

	if cfg.SubsOpts.SubWrapper != nil {
		// With Wrapper NewRelic
		sub, err = s.Nats[s.Params.serverName].QueueSubscribe(channel, qgroup,
			cfg.SubsOpts.SubWrapper(s.Chann[s.Params.serverName][channel]),
			subOpts...)
	} else {
		// Without NewRelic
		sub, err = s.Nats[s.Params.serverName].QueueSubscribe(channel, qgroup, func(msg *stan.Msg) {
			s.Chann[s.Params.serverName][channel] <- nats.MessageCh{Msg: msg, Data: msg.Data}

			if cfg.SubsOpts.ManualAckMode == nil {
				msg.Ack()
			}
		}, subOpts...)
	}

	if err != nil {
		return err
	}

	_, ok = s.Subs[s.Params.serverName]
	if !ok {
		s.Subs[s.Params.serverName] = make(map[string]stan.Subscription)
		_, ok := s.Subs[s.Params.serverName][channel]
		if !ok {
			s.Subs[s.Params.serverName][channel] = sub
		}
	} else {
		_, ok := s.Subs[s.Params.serverName][channel]
		if !ok {
			s.Subs[s.Params.serverName][channel] = sub
		}
	}

	return nil
}

// SubscribeAll register a array of subscribers
func (s *RepoStream) SubscribeAll(arr []interface{}) {
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
func (s *RepoStream) Unsubscribe(channel string) error {
	return s.Subs[s.Params.serverName][channel].Unsubscribe()
}

// Publishsend a payload to channel
func (s *RepoStream) Publish(payload nats.Message) error {
	return s.Nats[s.Params.serverName].Publish(payload)
}

// GetMessage return a message from channel broadcast
func (s *RepoStream) GetMessage(channel string) <-chan nats.MessageCh {
	go s.monitoringCh(channel)
	return s.GetChan
}

func (s *RepoStream) monitoringCh(channel string) {
	for {
		select {
		case msg := <-s.Chann[s.Params.serverName][channel]:
			msgch, ok := msg.(nats.MessageCh)
			if !ok {
				continue
			}
			go func(in nats.MessageCh) {
				s.GetChan <- in
			}(msgch)
		}
	}
}

func (s *RepoStream) getSubsOpts(opts *nats.Config) []stan.SubscriptionOption {
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

func (s *RepoStream) Close() error {
	for ch := range s.Subs[s.Params.serverName] {
		s.Subs[s.Params.serverName][ch].Unsubscribe()
	}
	return s.Nats[s.Params.serverName].Close()
}
