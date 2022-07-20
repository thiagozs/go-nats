package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
		nats.ClusterID("clusterID-consumer"),
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	for {
		select {
		case msg := <-stream.GetMessage("hello"):
			fmt.Printf("Listening> message: %s sequence: %d \n", string(msg.Data), msg.Msg.Sequence)
			msg.Msg.Ack()

		case <-quit:
			fmt.Println("Listening> exit")
			stream.Close()
			return
		}
	}
}
