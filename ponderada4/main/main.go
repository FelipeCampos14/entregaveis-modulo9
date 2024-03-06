package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	publisher "ponderada2/publisher"
	subscriber "ponderada2/subscriber"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	godotenv "github.com/joho/godotenv"
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

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	if username == "" || password == "" {
		// GitHub Secrets not found, try loading from .env file
		err := godotenv.Load("../.env")
		if err != nil {
			fmt.Println("Error loading .env file")
			return
		}

		username = os.Getenv("HIVE_USER")
		password = os.Getenv("HIVE_PSWD")
	}
	var broker = os.Getenv("BROKER_ADDR")
	var port = 8883
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", broker, port))
	opts.SetClientID("Ponderada4")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(MessagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	subscriber.Subscribe("sensor/+", client, MessagePubHandler)
	for i := 0; i < 4; i++ {
		publisher.Publish(client, 1)
	}

	client.Disconnect(250)
}
