package main

import (
	"fmt"
	"os"
	"strconv"
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
		nats.ClusterID("clusterID"),
		nats.NastPort("4223"),
		nats.NastURL("localhost"),
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

	if err := stream.Subscribe("hello"); err != nil {
		panic(err)
	}

	wait := make(chan bool, 1)

	if os.Getenv("NATS_STREAMING_PUB") == "true" {
		go func() {
			defer func() { wait <- true }()
			for i := 0; i < 50; i++ {
				msg := []byte(fmt.Sprintf("Hello %d", i))
				fmt.Println("Publishing:", string(msg))
				payload := nats.Message{Payload: msg, Subject: "hello"}
				if err := stream.Publish(payload); err != nil {
					panic(err)
				}
			}
		}()
	}

	count := 0

Loop:
	for {
		select {
		case msg := <-stream.GetMessage("hello"):
			if count > 5 && os.Getenv("NATS_STREAMING_PUB") == "true" {
				continue
			}
			count++
			fmt.Println("Listening> Got message:", string(msg.Data), msg.Msg.Sequence)
			fmt.Println("AckManualMode")
			msg.Msg.Ack()

		case <-time.After(time.Duration(3) * time.Second):
			fmt.Println("Listening> exit Timeout")
			wait <- true
		case <-wait:
			fmt.Println("Listening> exit")
			break Loop
		}
	}

	fmt.Println("exitFor")
	stream.Close()
	fmt.Println("Finish wait, done")
}
