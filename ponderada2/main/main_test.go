package main

import (
	"encoding/binary"
	"fmt"
	"math"
	publisher "ponderada2/publisher"
	subscriber "ponderada2/subscriber"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func TestMain(t *testing.T) {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetDefaultPublishHandler(MessagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	var MessageChannel = make(chan MQTT.Message)

	var Topics = [3]string{"RED", "OX", "NH3"}

	t.Run("TestReceiveData", func(t *testing.T) {
		// FUNÇÃO ASSÍNCRONA
		go func() {
			time.Sleep(1 * time.Second)
			publisher.Publish(client)
		}()

		for i := 0; i <= len(Topics)-1; i++ {
			topicStringf := fmt.Sprintf("sensor/%s", Topics[i])
			subscriber.Subscribe(topicStringf, client, func(client MQTT.Client, msg MQTT.Message) {
				MessageChannel <- msg
			})

			select {
			case msg := <-MessageChannel:
				if msg != nil {
					fmt.Printf("Message received  at topic: %s\n", Topics[i])
				}
			case <-time.After(5 * time.Second):
				t.Errorf("Timeout waiting for message in topic: %s", Topics[i])
			}
		}
	})

	t.Run("TestMatchData", func(t *testing.T) {
		// FUNÇÃO ASSÍNCRONA
		go func() {
			time.Sleep(1 * time.Second)
			publisher.Publish(client)
		}()

		for i := 0; i <= len(Topics)-1; i++ {
			topicStringf := fmt.Sprintf("sensor/%s", Topics[i])
			subscriber.Subscribe(topicStringf, client, func(client MQTT.Client, msg MQTT.Message) {
				MessageChannel <- msg
			})

			select {
			case msg := <-MessageChannel:
				resultado := math.Float64frombits(binary.LittleEndian.Uint64(msg.Payload()))
				expected := publisher.Values[i]

				if resultado == expected {
					fmt.Printf("Received message: %f \n", resultado)
				} else {
					t.Errorf("Message received in topic %s does not match data expected, received: %f; expected: %f.", Topics[i], resultado, expected)
				}
			case <-time.After(5 * time.Second):
				t.Errorf("Timeout waiting for message in topic: %s", Topics[i])
			}
		}
	})

	t.Run("TestDataFrequency", func(t *testing.T) {

		// FUNÇÃO ASSÍNCRONA
		go func() {
			time.Sleep(1 * time.Second)
			publisher.Publish(client)
		}()

		for i := 0; i <= len(Topics)-1; i++ {
			topicStringf := fmt.Sprintf("sensor/%s", Topics[i])
			subscriber.Subscribe(topicStringf, client, func(client MQTT.Client, msg MQTT.Message) {
				MessageChannel <- msg
			})

			// select {
			// case msg := <-MessageChannel:
			// 	resultado := math.Float64frombits(binary.LittleEndian.Uint64(msg.Payload()))
			// 	expected := publisher.Values[i]

			// 	if resultado == expected {
			// 		fmt.Printf("Received message: %f \n", resultado)
			// 	} else {
			// 		t.Errorf("Message receive in topic %s does not match data expected, received: %f; expected: %f.", Topics[i], resultado, expected)
			// 	}
			// case <-time.After(5 * time.Second):
			// 	t.Errorf("Timeout waiting for message in topic: %s", Topics[i])
			// }
			for MessageChannel {

			}
		}
	})
}
