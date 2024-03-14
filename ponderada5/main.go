package main

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	godotenv "github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	fmt.Println("Connected")
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

var MessageHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Recebido: %f do tópico: %s\n", math.Float64frombits(binary.LittleEndian.Uint64(msg.Payload())), msg.Topic())
	addToDb(msg, msg.Topic())
}

var Topics = [3]string{"RED", "OX", "NH3"}

var MapValues = [3]float64{randFloats(1.0, 1000.0), randFloats(0.05, 10.0), randFloats(1.0, 300.0)}

var Values = map[string]float64{
	"RED": MapValues[0],
	"OX":  MapValues[1],
	"NH3": MapValues[2],
}

func main() {

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	if username == "" || password == "" {
		// GitHub Secrets not found, try loading from .env file
		err := godotenv.Load(".env")
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
	opts.SetClientID("Ponderada5")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(MessageHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	Subscribe("sensor/+", client, MessageHandler)
	Publish(client, 1)

	client.Disconnect(250)
}

func randFloats(min, max float64) float64 {
	res := min + rand.Float64()*(max-min)
	return res
}

func Publish(client MQTT.Client, repTime time.Duration) {

	for i := 0; i <= len(Topics)-1; i++ {
		topicStringf := fmt.Sprintf("sensor/%s", Topics[i])
		ValuesBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(ValuesBytes, math.Float64bits(Values[Topics[i]]))
		token := client.Publish(topicStringf, 0, false, ValuesBytes)
		token.Wait()
		if token.Error() != nil {
			fmt.Printf("Failed to publish to topic: %s", Topics[i])
			panic(token.Error())
		}

	}
	time.Sleep(repTime * time.Second)
}

func Subscribe(topic string, client MQTT.Client, handler MQTT.MessageHandler) {
	token := client.Subscribe(topic, 1, handler)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("Failed to subscribe to topic: %v", token.Error())
		panic(token.Error())
	}
	fmt.Printf("\nSubscribed to topic: %s\n", topic)
}

func addToDb(msg MQTT.Message, table string) {
	db, _ := sql.Open("sqlite3", "./database/ponderada5.db")
	defer db.Close() // Defer Closing the database

	tableParts := strings.Split(table, "/")

	// Criando a tabla
	sqlStmt := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s
	(id INTEGER PRIMARY KEY, sensorValue FLOAT, time DATETIME)
	`, tableParts[1])
	// Preparando o sql statement de forma segura
	command, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Executando o comando sql
	command.Exec()

	// Criando uma função para inserir usuários
	insertData := func(db *sql.DB, data float64) {
		stmt := fmt.Sprintf(`INSERT INTO %s(sensorValue, time) VALUES (?, ?)`, tableParts[1])
		statement, err := db.Prepare(stmt)
		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = statement.Exec(data, time.Now())
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	insertData(db, math.Float64frombits(binary.LittleEndian.Uint64(msg.Payload())))
}
