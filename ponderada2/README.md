# Ponderada2

## Vídeo 

[Link](https://youtu.be/joE_RidlMHs])

## Explicação

Para esta ponderada foi demandado que os alunos criassem uma série de testes que comprovassem a funcionalidade de três aspectos do código:
- Recebimento de mensagens
- Correspodência entre mensagem enviada e recebida
- Publicação fosse feita dentro da frequência determinada

Para isso foi desenvolvida a seguinte estrutra de pastas:
.
├── go.mod
├── go.sum
├── main
│   ├── main.go
│   └── main_test.go
├── mosquitto.conf
├── publisher
│   └── publisher.go
├── README.md
└── subscriber
    └── subscriber.go

Onde a execução do MQTT foi dividido em publisher, subscriber e main, sendo o último o responsável por coordenar o código, então por este ser o principal, os testes foram feitos visando os fluxos deste código.


# Execução

Para executar o código: 
```
go run main/main.go
```

Para testar:
```
go test ponderada2/main -run TestMain -v
```