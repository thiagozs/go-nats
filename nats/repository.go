package nats

import "github.com/nats-io/stan.go"

type NatsServiceRepo interface {
	ReConnect(opt ...Options) error
	Close() error
	QueueSubscribe(subject, qgroup string, cb stan.MsgHandler,
		opts ...stan.SubscriptionOption) (stan.Subscription, error)
	Subscribe(subject string, cb stan.MsgHandler,
		opts ...stan.SubscriptionOption) (stan.Subscription, error)
	Publish(obj Message) error
	SetMessageHandler(fn func(msg *stan.Msg))
	MessageHandler() stan.MsgHandler
	ConnectedUrl() string
	ConnectedAddr() string
	ConnectedServerId() string
	DeliverAllAvailable() stan.SubscriptionOption
	GetOpts() *Config
}
