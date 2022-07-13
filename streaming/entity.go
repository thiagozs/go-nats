package streaming

import (
	"github.com/nats-io/stan.go"
	"github.com/thiagozs/go-nats/nats"
)

type SubscribeType func(serverName string, channel string) error

type Streaming struct {
	Nats map[string]nats.NatsServiceRepo
	Chan map[string]map[string]chan nats.MessageCh
	Subs map[string]map[string]stan.Subscription
}
