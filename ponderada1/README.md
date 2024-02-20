# Ponderada 1

## Vídeo

https://github.com/FelipeCampos14/entregaveis-modulo9/assets/99193547/eb28d047-bc5c-466f-9626-95ea54d02e57

## Como usar
Devemos fazer algumas instalações para rodar esse sistema MQTT

- Lib **Mosquitto** para o broker:
1. Para usuários de Linux, basta realizar o segundo comando
```
sudo apt-get install mosquitto mosquitto-clients
```
2. Para usuários de macOS: 
```
brew install mosquitto
```
3. Para usuários de Windows é necessário baixar o instalador nesse [link](https://mosquitto.org/download/)

- Lib **Eclipse-paho** para o cliente:
```
pip install -r requirements.txt
```

### Rodando
Em um terminal dê o comando:
```
mosquitto -c mosquitto.conf
```
Em outro terminal escreva: 
```
python3 ./src/publisher.py
```
E e m um terceiro terminal dê o último comando:
```
python3 ./src/publisher.py
```
## Explicação

Para essa ponderada foi usado eclipse-paho(python) e mosquitto para desenvolver, respectivamente, os clientes e broker mqtt. Foi selecionado o sensor de gases MiCs e, como não informa a taxa de registro, ô fiz publicar de 5 em 5 segundos. Além disso, escolhei acompanhar 3 gases, assim criando 3 tópicos diferentes para cada sensor.
