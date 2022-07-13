package nats

import "github.com/nats-io/stan.go"

type PingParameter struct {
	TTLA int
	TTLB int
}
type Message struct {
	Payload []byte
	Subject string
}

type MessageCh struct {
	Msg  *stan.Msg
	Data []byte
}
