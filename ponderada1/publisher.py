import paho.mqtt.client as mqtt
import time
import random

# Configuração do cliente
client = mqtt.Client("Sensor de radiação(Publisher)")

# Conecte ao broker
client.connect("localhost", 1892, 60)

def verify_status(status_code):
    if status_code == 0:
        print(f"Publicado: {message}")
    else:
        print("Falhou em publicar mensagem")



# Loop para publicar mensagens continuamente
try:
    while True:
        message = f"{random.randint(1,1280)}"   
        result = client.publish("test/radiation", message)
        verify_status(result[0])
        time.sleep(2)
except KeyboardInterrupt:
    print("Publicação encerrada")

client.disconnect()