package streaming

import (
	"nats-go/nats"

	"github.com/nats-io/stan.go"
)

type SubscribeType func(serverName string, channel string) error

type Streaming struct {
	Nats map[string]nats.NatsServiceRepo
	Chan map[string]map[string]chan []byte
	Subs map[string]map[string]stan.Subscription
}
