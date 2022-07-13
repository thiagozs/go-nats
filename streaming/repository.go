package streaming

import "github.com/thiagozs/go-nats/nats"

type StreamingServiceRepo interface {
	Register(service nats.NatsServiceRepo) *repoStream
	Subscribe(channel string) error
	SubscribeAll(arr []interface{})
	Unsubscribe(channel string) error
	Publish(payload nats.Message) error
	GetMessage(channel string) chan nats.MessageCh
	Close() error
}
