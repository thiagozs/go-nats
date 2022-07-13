package nats

import (
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/nats-io/nuid"
	"github.com/nats-io/stan.go"
)

type repoNats struct {
	mutex sync.RWMutex
	opts  Config
}

func NewService(opts ...Options) (NatsServiceRepo, error) {
	s := repoNats{}
	s.opts = NewOptions(opts...)

	// NOTE: change a validation for more options

	s.opts.ConnLostHandler = func(con stan.Conn, reason error) {
		log.Fatalf("Connection lost, reason: %v", reason)
	}

	// if clusterID are empty fill up the id
	if s.opts.ClusterID == "" {
		s.opts.ClusterID = nuid.Next()
	}

	if err := s.conn(); err != nil {
		return &s, err
	}

	return &s, nil
}

// ReConnect turn back the connection with other config
func (r *repoNats) ReConnect(opt ...Options) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// start a new config
	s := repoNats{}

	// save a new options
	opts := NewOptions(opt...)
	s.opts = opts

	// Check if connection are able to reconnect
	if !r.checkConnIsValid() {
		return errors.New("connection lost or invalid")
	}

	// fill up the clusterID if not exist a string
	// if clusterID are empty fill up the id
	if s.opts.ClusterID == "" {
		s.opts.ClusterID = nuid.Next()
	}

	if err := s.Close(); err != nil {
		return err
	}

	if err := s.conn(); err != nil {
		return err
	}

	return nil
}

func (r *repoNats) checkConnIsValid() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if r.opts.ConnNats.NatsConn() != nil &&
		r.opts.ConnNats.NatsConn().IsConnected() {
		return true
	}
	return false
}
func (r *repoNats) conn() error {
	opts := []stan.Option{
		stan.Pings(r.opts.PingParams.TTLA, r.opts.PingParams.TTLB),
		stan.MaxPubAcksInflight(r.opts.MaxInFlight),
		stan.PubAckWait(r.opts.PubAckWait),
		stan.NatsURL(strings.Join([]string{r.opts.NatsURL, r.opts.NatsPort}, ":")),
		stan.SetConnectionLostHandler(r.opts.ConnLostHandler),
	}

	sc, err := stan.Connect(
		r.opts.ClusterName,
		r.opts.ClusterID,
		opts...,
	)

	if err != nil {
		return err
	}
	r.opts.ConnNats = sc
	return nil
}

// Close turn off the connection
func (r *repoNats) Close() error {
	return r.opts.ConnNats.Close()
}

// QueueSubscribe return the process for listening the queue
func (r *repoNats) QueueSubscribe(subject, qgroup string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	return r.opts.ConnNats.QueueSubscribe(subject, qgroup, cb, opts...)
}

// Subscribe listening the subject channel
func (r *repoNats) Subscribe(subject string, cb stan.MsgHandler, opts ...stan.SubscriptionOption) (stan.Subscription, error) {
	return r.opts.ConnNats.Subscribe(subject, cb, opts...)
}

// Publish a message payload on subject
func (r *repoNats) Publish(obj Message) error {
	if !r.checkConnIsValid() {
		return errors.New("publish connection lost or invalid")
	}
	return r.opts.ConnNats.Publish(obj.Subject, obj.Payload)
}

// MessageHandler return a handler for request
func (r *repoNats) MessageHandler() stan.MsgHandler {
	return r.opts.MessageHandler
}

// SetMessageHandler setup a request message handler
func (r *repoNats) SetMessageHandler(fn func(msg *stan.Msg)) {
	r.opts.MessageHandler = fn
}

// ConnectedUrl return a url connected
func (r *repoNats) ConnectedUrl() string {
	return r.opts.ConnNats.NatsConn().ConnectedUrl()
}

// ConnectedAddr return a address connection
func (r *repoNats) ConnectedAddr() string {
	return r.opts.ConnNats.NatsConn().ConnectedAddr()
}

// ConnectedServerId return a id server
func (r *repoNats) ConnectedServerId() string {
	return r.opts.ConnNats.NatsConn().ConnectedServerId()
}

// DeliverAllAvailable delivery all message are avaliable options
func (r *repoNats) DeliverAllAvailable() stan.SubscriptionOption {
	return stan.DeliverAllAvailable()
}

// Get Opts configs
func (r *repoNats) GetOpts() *Config {
	return &r.opts
}
