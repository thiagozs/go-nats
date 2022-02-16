package main

import (
	"fmt"
	"nats-go/nats"
	"nats-go/streaming"
	"time"
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
	}

	ss, err := nats.NewService(opts...)
	if err != nil {
		panic(err)
	}

	st := streaming.NewService(ss)
	//st.Register("candle", ss)

	if err := st.Subscribe("candle"); err != nil {
		panic(err)
	}

	payload := nats.Message{Payload: []byte("teste"), Subject: "hello"}
	if err := st.Publish(payload); err != nil {
		panic(err)
	}

	fmt.Println(string(<-st.GetMessage("candle")))

	if err := st.Unsubscribe("candle"); err != nil {
		panic(err)
	}

}
