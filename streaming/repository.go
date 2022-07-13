package streaming

import "github.com/thiagozs/go-nats/nats"

type StreamingServiceRepo interface {
	Register(service nats.NatsServiceRepo) *RepoStream
	Subscribe(channel string) error
	QueueSubscribe(channel, qgroup string) error
	SubscribeAll(arr []interface{})
	Unsubscribe(channel string) error
	Publish(payload nats.Message) error
	GetMessage(channel string) chan nats.MessageCh
	Close() error
}
