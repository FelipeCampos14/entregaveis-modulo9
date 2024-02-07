import paho.mqtt.client as mqtt

topicArray = ["RED","OX","NH3"]

# Callback quando uma mensagem é recebida do servidor.
def on_message(client, userdata, message):
    print(f"\nRecebido: {message.payload.decode()} no tópico {message.topic}")
    # print(message)
    # check_radiation_level(message.payload.decode())

# def check_radiation_level(radiation):
#     if float(radiation) > 1200:
#         print("Radiação acima do normal, tome cuidado.")
#     else:
#         print("Radiação em nível normal.")

# Callback para quando o cliente recebe uma resposta CONNACK do servidor.
def on_connect(client, userdata, flags, rc):
    print("Conectado com código de resultado "+str(rc))
    # Inscreva no tópico aqui, ou se perder a conexão e se reconectar, então as
    # subscrições serão renovadas.
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