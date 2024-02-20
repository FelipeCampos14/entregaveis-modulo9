import paho.mqtt.client as mqtt

topicArray = ["RED","OX","NH3"]

# Callback quando uma mensagem é recebida do servidor.
def on_message(client, userdata, message):
    print(f"\nRecebido: {message.payload.decode()} no tópico {message.topic}")

# Callback para quando o cliente recebe uma resposta CONNACK do servidor.
def on_connect(client, userdata, flags, rc):
    print("Conectado com código de resultado "+str(rc))
    client.subscribe(f"MiCS/#")

# Configuração do cliente2
client = mqtt.Client("Sensor de poluição no ar(subscriber)")
for topic in topicArray:
    client.on_connect = on_connect
    client.on_message = on_message

# Conecte ao broker
client.connect("localhost", 1892, 60)

# Loop para manter o cliente executando e escutando por mensagens
client.loop_forever()