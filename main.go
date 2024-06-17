package main

import (
	"fmt"
	"log"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/listeners/auth"
)

func main() {
	fmt.Println("Hello Elius")
	options := mqtt.Options{
		BufferSize:      512 * 1024,
		BufferBlockSize: 16 * 1024,
	}
	
	server := mqtt.NewServer(&options)
	server.Events.OnConnect = func(c events.Client, p events.Packet) {
		fmt.Printf("On connect Client %s connected  ", c.ID)
	}
	server.Events.OnDisconnect = func(c events.Client, err error) {
		fmt.Printf("On disconnect %s Client ", c.ID)
	}

	server.Events.OnMessage = func(cl events.Client, pk events.Packet) (pkx events.Packet, err error) {
		log.Printf("Received message on topic %s: %s", pk.TopicName, string(pk.Payload))
		go processMessage(cl, pk)
		return pk, nil
	}
	server.Events.OnSubscribe = func(filter string, cl events.Client, qos byte) {
		fmt.Printf("Client %s subscribed to topic: %s", cl.ID, filter)
		// You can also log the QoS level if needed
		fmt.Printf(" with QoS: %d\n", qos)
	}
	tcp := listeners.NewTCP("t1", ":1884")
	err := server.AddListener(tcp, &listeners.Config{
		Auth: new(auth.Allow),
	})
	if err != nil {
		log.Fatal("Error")
	}
	err = server.Serve()
	if err != nil {
		log.Fatal("Error 2")
	}

	log.Println("Asynchronous MQTT server is running on port 1883")

	// Keep the server running
	stop := make(chan struct{})
	<-stop

}

func processMessage(cl events.Client, pk events.Packet) {
	// Asynchronous message processing logic
	log.Printf("Processing message on topic %s: %s : %s", pk.TopicName, string(pk.Payload), cl.ID)
}
