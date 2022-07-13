# NAST Adapter Streaming tools

## Stack of project

* golang v 1.17
* nast.io streaming driver
* docker-ce
* docker-composer

## Nats Service

### Options

```go
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

}

```

### New Instance

Start a new instrace of nast service.

```go
ns, err := nats.NewService(opts...)
if err != nil {
    panic(err)
}
defer func() {
    _ = ns.Close()
}()
```

### Publish

Create publish client statement

```go
for i := 0; i < 5; i++ {
    payload := fmt.Sprintf("%s-%s", "hello", strconv.Itoa(i))
    fmt.Println("Producer> Send message:", payload)
    msg := nats.Message{Payload: []byte(payload), Subject: "hello"}
    if err := ns.Publish(msg); err != nil {
    fmt.Println("Publish Error:", err)
    }
    time.Sleep(1 * time.Second)
}
```

### Subscriber

Create a subscriber statement

```go
fmt.Println("Subscriber> Connected to NATS at:", ns.ConnectedUrl())
sub, err := ns.Subscribe("hello", ns.MessageHandler())
if err != nil {
    fmt.Println("Subscribe Error:", err)
}
defer func() { _ = sub.Close() }()
```

## Streaming

### New Instace

Create a new instance

```go
// Create a new Instance wrapper streaming
	nt, err := nats.NewService(opts...)
	if err != nil {
		panic(err)
	}

	st := streaming.NewService(nt)
	// st.Register("anchor", nt)
```

### Register a channels

Subscriber the channels

```go
// Create a list of subscribers for register and listening
subscribers := []interface{}{
    stm.Subscribe(Candle1.String()),
    stm.Subscribe(Candle2.String()),
    stm.Subscribe(Candle3.String()),
    stm.Subscribe(Candle4.String()),
}

// Register all subscribers
stm.SubscribeAll(subscribers)

```

### Publish message

Push a message to broker

```go
for i := 1; i <= 10; i++ {
    fmt.Printf("producer: %s> channel: %s\n", serverName, channel)
    obj := fmt.Sprintf("%s-%d", "hello", i)
    payload := nats.Message{Payload: []byte(obj), Subject: channel}
    if err := stream.Publish(serverName, payload); err != nil {
        fmt.Println(err)
    }
}
```

### Read payload from streaming

```go
for {
    select {
        case cl1 := <-stm.GetMessage(Candle1.String()):
        fmt.Println("Candle1 =", string(cl1))

        case cl2 := <-stm.GetMessage(Candle2.String()):
        fmt.Println("Candle2 =", string(cl2))

        case cl3 := <-stm.GetMessage(Candle3.String()):
        fmt.Println("Candle3 =", string(cl3))

        case cl4 := <-stm.GetMessage(Candle4.String()):
        fmt.Println("Candle4 =", string(cl4))

        case <-time.After(time.Duration(5) * time.Second):
        fmt.Println("Exit listening...")
        break loop
    }
}

```

### Unsubscribe the channel

```go
if err := stream.Unsubscribe(Dev.String(), channel); err != nil {
    fmt.Println(err)
}
```

## Examples code

The folder `examples`has some examples of code ready to running

## TODO

* [X] PoC nats runneble
* [X] PoC Streaming runneble
* [ ] Test Unit

## Versioning and license

We use SemVer for versioning. You can see the versions available by checking the tags on this repository.

For more details about our license model, please take a look at the LICENSE file

---

2020, thiagozs
