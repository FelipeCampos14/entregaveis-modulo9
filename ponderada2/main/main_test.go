package main

import (
	"encoding/binary"
	"fmt"
	"math"
	publisher "ponderada2/publisher"
	subscriber "ponderada2/subscriber"
	"sync"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func TestMain(t *testing.T) {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetDefaultPublishHandler(MessagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	var Topics = [3]string{"RED", "OX", "NH3"}

	var messageChannels = make(map[string]chan MQTT.Message)

	for i := 0; i <= len(Topics)-1; i++ {
		topicStringf := fmt.Sprintf("sensor/%s", Topics[i])
		messageChannels[topicStringf] = make(chan MQTT.Message)
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	t.Run("TestReceiveData", func(t *testing.T) {
		var wg sync.WaitGroup

		for i := 0; i <= len(Topics)-1; i++ {
			wg.Add(1)
			topicStringf := fmt.Sprintf("sensor/%s", Topics[i])
			go func(topic string) {
				defer wg.Done()

				subscriber.Subscribe(topic, client, func(client MQTT.Client, msg MQTT.Message) {
					messageChannels[topic] <- msg
				})

				go func() {
					time.Sleep(1 * time.Second)
					publisher.Publish(client)
				}()
			}(topicStringf)
		}

		wg.Wait()

	})

	t.Run("TestMatchData", func(t *testing.T) {
		var wg sync.WaitGroup

		for i := 0; i <= len(Topics)-1; i++ {
			wg.Add(1)
			topicStringf := fmt.Sprintf("sensor/%s", Topics[i])
			go func(topic string) {
				defer wg.Done()

				subscriber.Subscribe(topic, client, func(client MQTT.Client, msg MQTT.Message) {
					resultado := math.Float64frombits(binary.LittleEndian.Uint64(msg.Payload()))
					expected := publisher.Values[i]

					if resultado != expected {
						t.Errorf("Message received in topic %s does not match data expected, received: %f; expected: %f.", topic, resultado, expected)
					}
				})

				go func() {
					time.Sleep(1 * time.Second)
					publisher.Publish(client)
				}()
			}(topicStringf)
		}

		wg.Wait()

	})

	t.Run("TestDataFrequency", func(t *testing.T) {
		var sampleSize = 3
		var rate = 5.0
		var tolerance = 0.1
		var expectedLess, expectedPlus = rate - (rate * tolerance), rate + (rate * tolerance)

		var startTime = time.Now()

		for i := 1; i <= sampleSize; i++ {
			publisher.Publish(client)
		}
		var totalTime = time.Since(startTime).Seconds()

		AverageTime := totalTime / float64(sampleSize)

		switch {
		case AverageTime > expectedPlus:
			t.Errorf("Messages are taking longer than expected to be published. Took:%f, expected:%f", AverageTime, expectedPlus)
		case AverageTime < expectedLess:
			t.Errorf("Messages are taking longer than expected to be published. Took:%f, expected:%f", AverageTime, expectedLess)
		}
	})
}
