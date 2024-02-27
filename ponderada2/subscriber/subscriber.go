package subscriber

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Subscribe(topic string, client MQTT.Client, handler MQTT.MessageHandler) {
	token := client.Subscribe(topic, 1, handler)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("Failed to subscribe to topic: %v", token.Error())
		panic(token.Error())
	}
	fmt.Printf("Subscribed to topic: %s \n", topic)
}
