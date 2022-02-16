package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/thiagozs/go-nats/nats"
)

func main() {

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
		nats.MessageHandler(func(msg *stan.Msg) {
			fmt.Println("Subscriber> Got task request on:", msg.Subject)
			fmt.Println("Subscriber> Message:", string(msg.Data))
		}),

		// If you need to customize this services,
		// just set a group default handlers wrapper.
		//
		// nats.OptsHandler(
		// 	nats.SetAckWait(60),
		// 	nats.SetDurableName("durable-name"),
		// 	nats.SetMaxInFlight(20),
		// 	nats.SetManualAckMode(true),
		// ),
	}

	ns, err := nats.NewService(opts...)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = ns.Close()
	}()

	done := make(chan struct{})

	go func() {
		fmt.Println("Subscriber> Connected to NATS at:", ns.ConnectedUrl())
		sub, err := ns.Subscribe("hello", ns.MessageHandler())
		if err != nil {
			fmt.Println("Subscribe Error:", err)
		}
		defer func() { _ = sub.Close() }()
		<-done
	}()

	fmt.Println("Producer> Connected to NATS at:", ns.ConnectedUrl())
	for i := 0; i < 5; i++ {
		payload := fmt.Sprintf("%s-%s", "hello", strconv.Itoa(i))
		fmt.Println("Producer> Send message:", payload)
		msg := nats.Message{Payload: []byte(payload), Subject: "hello"}
		if err := ns.Publish(msg); err != nil {
			fmt.Println("Publish Error:", err)
		}
		time.Sleep(1 * time.Second)
	}
	done <- struct{}{}

}
