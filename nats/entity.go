package nats

type PingParameter struct {
	TTLA int
	TTLB int
}
type Message struct {
	Payload []byte
	Subject string
}
