package main

import (
	// "encoding/binary"
	"fmt"
	// "math"
	publisher "ponderada2/publisher"
	subscriber "ponderada2/subscriber"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestMain(t *testing.T) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetDefaultPublishHandler(MessagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	t.Run("TestReceiveData", func(t *testing.T) {
		clientTest1 := mqtt.NewClient(opts)
		if token := clientTest1.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		subscriber.Subscribe("sensor/+", clientTest1, func(client mqtt.Client, msg mqtt.Message) {
			resultado := msg.Payload()
			topic := msg.Topic()

			if resultado == nil {
				t.Errorf("Message not received.")
			} else {
				fmt.Printf("\nMessage received from topic %s\n ", topic)
			}
		})
		publisher.Publish(clientTest1, 1)
		clientTest1.Disconnect(1000)

	})

	// t.Run("TestMatchData", func(t *testing.T) {

	// 	publisher.Values["RED"] = publisher.MapValues[0]
	// 	publisher.Values["OX"] = publisher.MapValues[1]
	// 	publisher.Values["NH3"] = publisher.MapValues[2]

	// 	var Topics = map[string]float64{
	// 		"RED":publisher.MapValues[0],
	// 		"OX":publisher.Values["OX"],
	// 		"NH3":publisher.Values["NH3"],
	// 	}

	// 	clientTest2 := mqtt.NewClient(opts)
	// 	if token := clientTest2.Connect(); token.Wait() && token.Error() != nil {
	// 		panic(token.Error())
	// 	}
	// 	subscriber.Subscribe("sensor/+", clientTest2, func(client mqtt.Client, msg mqtt.Message) {
	// 		resultado := math.Float64frombits(binary.LittleEndian.Uint64(msg.Payload()))
	// 		expected := Topics[msg.Topic()]

	// 		if resultado == expected {
	// 			fmt.Print("Message matches the expected\n")
	// 		} else {
	// 			t.Errorf("Message %f from topic %s, different from expected %f\n ", resultado, msg.Topic(), expected)
	// 		}
	// 	})
	// 	publisher.Publish(clientTest2, 1)
	// 	clientTest2.Disconnect(1000)

	// })

	t.Run("TestDataFrequency", func(t *testing.T) {
		clientTest3 := mqtt.NewClient(opts)
		if token := clientTest3.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		var sampleSize = 3
		var rate = 1.0
		var tolerance = 0.1
		var expectedLess, expectedPlus = rate - (rate * tolerance), rate + (rate * tolerance)

		var startTime = time.Now()

		for i := 1; i <= sampleSize; i++ {
			publisher.Publish(clientTest3, 1)
		}
		var totalTime = time.Since(startTime).Seconds()

		AverageTime := totalTime / float64(sampleSize)

		switch {
		case AverageTime > expectedPlus:
			t.Errorf("Messages are taking longer than expected to be published. Took:%f, expected:%f", AverageTime, expectedPlus)
		case AverageTime < expectedLess:
			t.Errorf("Messages are taking less than expected to be published. Took:%f, expected:%f", AverageTime, expectedLess)
		}
	})
}
