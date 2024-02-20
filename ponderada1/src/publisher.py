import paho.mqtt.client as mqtt
import time
import random

topicArray = ["RED","OX","NH3"]

# Configuração do cliente
client = mqtt.Client("Sensor de poluição no ar(Publisher)")

# Conecte ao broker 
client.connect("localhost", 1892, 60)

def verify_status(status_code, topic):
    if status_code == 0:
        print(f"Publicado: {message}, no tópico: {topic}  \n")
    else:
        print("Falhou em publicar mensagem")

# Loop para publicar mensagens continuamente
try:
    while True:
        sensorValueArray = [random.uniform(1,1000), random.uniform(0.05,10), random.uniform(1,300)]
        for topic, sensorValue in zip(topicArray,sensorValueArray):
            message = sensorValue  
            result = client.publish(f"MiCS/{topic}", message)
            verify_status(result[0], topic)
        time.sleep(5)
except KeyboardInterrupt:
    print("Publicação encerrada")

client.disconnect()