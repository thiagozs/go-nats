package streaming

import "nats-go/nats"

type StreamingServiceRepo interface {
	//Register(name string, service nats.NatsServiceRepo) *repoStream
	Subscribe(channel string) error
	SubscribeAll(arr []interface{})
	Unsubscribe(channel string) error
	Publish(payload nats.Message) error
	GetMessage(channel string) chan []byte
}
