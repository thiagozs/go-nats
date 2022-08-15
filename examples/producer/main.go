package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/thiagozs/go-nats/streaming"

	"github.com/thiagozs/go-nats/nats"
)

func main() {
	subOpts := []nats.SubsOptHandler{
		nats.SetManualAckMode(true),
	}

	if os.Getenv("SEQUENCE") != "" {
		u, _ := strconv.ParseUint(os.Getenv("SEQUENCE"), 0, 64)
		subOpts = append(subOpts, nats.SetStartAtSequence(u))
	}

	opts := []nats.Options{
		nats.ClusterName("test-cluster"),
		nats.ClusterID("clusterID-producer"),
		nats.NatsPort("4223"),
		nats.NatsURL("localhost"),
		nats.PubAckWait(30 * time.Second),
		nats.MaxInFlight(1),
		nats.PingParams(nats.PingParameter{
			TTLA: 100,
			TTLB: 250,
		}),
		nats.NewSubsOpts(subOpts...),
	}

	nts, err := nats.NewService(opts...)
	if err != nil {
		panic(err)
	}
	stream := streaming.NewService(nts)

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			msg := []byte(fmt.Sprintf("Hello %d", i))
			fmt.Println("Publishing:", string(msg))
			payload := nats.Message{Payload: msg, Subject: "hello"}
			if err := stream.Publish(payload); err != nil {
				panic(err)
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		msg := []byte(fmt.Sprintf("Hello %d", i))
		fmt.Println("Publishing:", string(msg))
		payload := nats.Message{Payload: msg, Subject: "hello"}
		if err := stream.Publish(payload); err != nil {
			panic(err)
		}
	}
}
