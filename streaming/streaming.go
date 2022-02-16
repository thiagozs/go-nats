package streaming

import (
	"fmt"
	"nats-go/nats"

	"github.com/nats-io/stan.go"
)

const (
	SERVER = "default"
)

type repoStream struct {
	Nats map[string]nats.NatsServiceRepo
	Chan map[string]map[string]chan []byte
	Subs map[string]map[string]stan.Subscription
}

// NewSevice start a streming wrapper
func NewService(natsService nats.NatsServiceRepo) StreamingServiceRepo {
	nats := make(map[string]nats.NatsServiceRepo)
	chans := make(map[string]map[string]chan []byte)
	subs := make(map[string]map[string]stan.Subscription)

	nats[SERVER] = natsService

	return &repoStream{
		Nats: nats,
		Chan: chans,
		Subs: subs,
	}
}

// Register method for subscription a new nast instance
func (s *repoStream) Register(natsService nats.NatsServiceRepo) *repoStream {
	s.Nats[SERVER] = natsService
	return s
}

// Subscribe register a cliente in a channel broadcasting
func (s *repoStream) Subscribe(channel string) error {
	_, ok := s.Chan[SERVER]
	if !ok {
		s.Chan[SERVER] = make(map[string]chan []byte)
		_, ok := s.Chan[SERVER][channel]
		if !ok {
			s.Chan[SERVER][channel] = make(chan []byte, 100)
		}
	} else {
		_, ok := s.Chan[SERVER][channel]
		if !ok {
			s.Chan[SERVER][channel] = make(chan []byte, 100)
		}
	}
	sub, err := s.Nats[SERVER].Subscribe(channel, func(msg *stan.Msg) {
		s.Chan[SERVER][channel] <- msg.Data
	})
	if err != nil {
		return err
	}

	_, ok = s.Subs[SERVER]
	if !ok {
		s.Subs[SERVER] = make(map[string]stan.Subscription)
		_, ok := s.Subs[SERVER][channel]
		if !ok {
			s.Subs[SERVER][channel] = sub
		}
	} else {
		_, ok := s.Subs[SERVER][channel]
		if !ok {
			s.Subs[SERVER][channel] = sub
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
	return s.Subs[SERVER][channel].Unsubscribe()
}

// Publishsend a payload to channel
func (s *repoStream) Publish(payload nats.Message) error {
	return s.Nats[SERVER].Publish(payload)
}

// GetMessage return a message from channel broadcast
func (s *repoStream) GetMessage(channel string) chan []byte {
	return s.Chan[SERVER][channel]
}
