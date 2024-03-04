# Ponderada 3

## Enunciado

Atividade prática em grupo (recomendação: grupos de 2 ou 3 alunos), de subir um broker remoto e um broker local do MQTT para conduzir cenários de análise de vulnerabilidade (dentro do CIA Triad), identificando situações onde pode ser comprometido cada um dos três pilares: Confiabilidade, Integridade e Disponibilidade.

## Análise

Ao subir um broker à alguns pontos de vulnerabilidade que devem ser observados. Os primeiros dizem respeito à **confiabilidade** dos dados:

- Em um broker local, apesar de estar limitado à rede do computador, é possível invadí-lo utilizando uma malware ou explorando vulnerabilidades de alguma rede que o computador esteja conectado.
- Quanto à um broker remoto, as informações ficam um pouco mais expostas, apesar de poder providenciar cadastros únicos aos usuários. O motivo desta vunrabilidade é exatamente estas credenciais que podem ser vazadas, além de, por ser uma empresa terceira que gerencia este broker, ele ainda fica sucetível á essa empresa ter seus dados vazados ou servidores invadidos, facilitando o acesso de pessoas mal-intencionadas ao broker.

Segundo ponto de atenção é relativo à **integridade** destes dados:
- Em um broker local, é possível que mais de uma aplicação esteja apontando para o mesmo endpoint, enviando dados não correspondentes aos desejados.
- Em um broker remoto, é possível que mais de uma pessoa esteja visando o mesmo cluster, isso pode levar a diferentes padrões de dados, o que gera um enviesamento no banco de dados. 

Outro cenário, que é comum aos dois, é um excesso de tráfego, como um cliente que deveria publicar uma vez a cada 15 minuto, pode começar a publicar numa frequência de 15 milissegundos, o que pode levar a aplicação à cair.

Finalmente, em relação à disponibilidade dos dados:
- Em um broker local, os dados são praticamente restritos ao usuário daquele computador, dificultando que qualquer outra pessoa consiga acessar esse dados.

Em ambos, local e remoto, é possível que a aplicação fique lenta caso haja um execesso de tráfico, até impossibilitando um desenvolvedor ou usuário de acessar os dados. 

## Simulação

Para simular um cenário que testase ou evidenciasse estas vulnerabilidades, criaria-se dois brokers, um local e um remoto, e testaria o seguinte, de acordo com cada valor da tríade:
1. Confiabilidade:
- Testaria fazer conexão à aplicação sem autorização necessária, baseada em credenciais e informações do do computador de um usuário cofiável;

2. Integridade:
- Testaria enviar informações que não são de acordo com as esperadas para aquela aplicação ou aquele tópico;
- Testaria o envio em demasia de uma informação para ver como a aplicaçã lida;

3. Disponibilidade:
- Criaria um caso onde há um excesso de tráfego na aplicação e testaria fazer uma requisição para testar a capacidade e tempo de resposta do servidor;