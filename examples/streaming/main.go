package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/thiagozs/go-nats/nats"
	"github.com/thiagozs/go-nats/streaming"
)

type ServerID int
type CandleID int
type BrokerID int

const (
	Dev ServerID = iota
)

const (
	Candle1 CandleID = iota
	Candle2
	Candle3
	Candle4
)

const (
	Broker1 BrokerID = iota
)

func (s ServerID) String() string {
	return [...]string{"Dev"}[s]
}
func (c CandleID) String() string {
	return [...]string{"Candle1", "Candle2", "Candle3", "Candle4"}[c]
}
func (b BrokerID) String() string {
	return [...]string{"Broker1"}[b]
}

var counters = map[string]int{}

func counter(channel string, mutex *sync.RWMutex) {
	mutex.RLock()
	defer mutex.RUnlock()
	if _, ok := counters[channel]; !ok {
		counters[channel] = 0
	}
	counters[channel]++

}

func main() {
	// load a nast server, if you can use more than 1 server
	//just load other config and call this in one variable
	server, err := loadBroker(Broker1.String())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a new Instance wrapper streaming
	stm := streaming.NewService(server)
	// Register the servers here

	// Create a list of subscribers for register and listening
	subscribers := []interface{}{
		stm.Subscribe(Candle1.String()),
		stm.Subscribe(Candle2.String()),
		stm.Subscribe(Candle3.String()),
		stm.Subscribe(Candle4.String()),
	}

	// Register all subscribers
	stm.SubscribeAll(subscribers)

	// Simulate all producer content for candles
	go producer(Candle1.String(), stm)
	go producer(Candle2.String(), stm)
	go producer(Candle3.String(), stm)
	go producer(Candle4.String(), stm)

	// just for count data streaming
	mutex := sync.RWMutex{}

	// Get messages from broker in distict channels
loop:
	for {
		select {
		case cl1 := <-stm.GetMessage(Candle1.String()):
			fmt.Println("Candle1 =", string(cl1))
			counter(Candle1.String(), &mutex)

		case cl2 := <-stm.GetMessage(Candle2.String()):
			fmt.Println("Candle2 =", string(cl2))
			counter(Candle2.String(), &mutex)

		case cl3 := <-stm.GetMessage(Candle3.String()):
			fmt.Println("Candle3 =", string(cl3))
			counter(Candle3.String(), &mutex)

		case cl4 := <-stm.GetMessage(Candle4.String()):
			fmt.Println("Candle4 =", string(cl4))
			counter(Candle4.String(), &mutex)

		case <-time.After(time.Duration(5) * time.Second):
			fmt.Println("Exit listening...")
			break loop
		}
		fmt.Printf("stats %v\n", counters)
	}

}

func producer(channel string,
	stream streaming.StreamingServiceRepo) {
	for i := 1; i <= 100; i++ {
		fmt.Printf("producer> channel: %s\n", channel)
		obj := fmt.Sprintf("%s-%d", "hello", i)
		payload := nats.Message{Payload: []byte(obj), Subject: channel}
		if err := stream.Publish(payload); err != nil {
			fmt.Println(err)
		}
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(1000)
		time.Sleep(time.Duration(r) * time.Millisecond)
	}

	if err := stream.Unsubscribe(channel); err != nil {
		fmt.Println(err)
	}
}

func startBroker1() (nats.NatsServiceRepo, error) {
	opts := []nats.Options{
		nats.ClusterName("test-cluster"),
		nats.ClusterID("clusterID"),
		nats.NastPort("4223"),
		nats.NastURL("localhost"),
		nats.PubAckWait(30 * time.Second),
		nats.MaxInFlight(20),
		nats.PingParams(nats.PingParameter{
			TTLA: 100,
			TTLB: 250,
		}),
	}
	ns, err := nats.NewService(opts...)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func loadBroker(serverName string) (nats.NatsServiceRepo, error) {
	var server nats.NatsServiceRepo
	var err error
	switch serverName {
	case Broker1.String():
		server, err = startBroker1()
	}

	return server, err
}
