package main

import (
	"encoding/binary"
	"fmt"
	"math"
	publisher "ponderada2/publisher"
	subscriber "ponderada2/subscriber"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	fmt.Println("Connected")
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

var MessagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Recebido: %f do tópico: %s\n", math.Float64frombits(binary.LittleEndian.Uint64(msg.Payload())), msg.Topic())
}

func main() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetDefaultPublishHandler(MessagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	subscriber.Subscribe("sensor/+", client, MessagePubHandler)
	for i:=0; i<4;i++{
		publisher.Publish(client, 1)
	}

	client.Disconnect(250)
}
