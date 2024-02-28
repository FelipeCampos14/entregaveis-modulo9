package publisher

import (
	"encoding/binary"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"math"
	"math/rand/v2"
	"time"
)

var Topics = [3]string{"RED", "OX", "NH3"}

func randFloats(min, max float64) float64 {
	res := min + rand.Float64()*(max-min)
	return res
}

var Values = [3]float64{randFloats(1.0, 1000.0), randFloats(0.05, 10.0), randFloats(1.0, 300.0)}

func Publish(client MQTT.Client) {
	for i := 0; i <= len(Topics)-1; i++ {
		topicStringf := fmt.Sprintf("sensor/%s", Topics[i])
		ValuesBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(ValuesBytes, math.Float64bits(Values[i]))
		token := client.Publish(topicStringf, 0, false, ValuesBytes)
		token.Wait()
		if token.Error() != nil {
			fmt.Printf("Failed to publish to topic: %s", Topics[i])
			panic(token.Error())
		}

	}
	time.Sleep(5 * time.Second)
}
