# Desafio Rage Limiter

Este repositório foi criado exclusivamente para hospedar o código do desenvolvimento do Desfio de implementação do Rate Limiter da Pós Go Expert, ministrado pela Full Cycle.

## Descrição do Desafio

A seguir estão os dados fornecidos na descrição do desafio.


### Objetivo

Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

### Descrição

O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

- **Endereço IP:** O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.

- **Token de Acesso:** O rate limiter também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
    - **API_KEY:** \<TOKEN\>

- As configurações de limite do **token de acesso devem se sobrepor as do IP**. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

### Requisitos

- O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web

- O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.

- O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.

- As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo ".env" na pasta raiz.

- Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.

- O sistema deve responder adequadamente quando o limite é excedido:
    - Código HTTP: **429**
    - Mensagem: **you have reached the maximum number of requests or actions allowed within a certain time frame**

- Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.

- Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.

- A lógica do limiter deve estar separada do middleware.


### Exemplos


1. **Limitação por IP:** Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.

2. **Limitação por Token:** Se um token **abc123** tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.

3. Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.

### Dicas

- Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.

## Entrega

- O código-fonte completo da implementação.

- Documentação explicando como o rate limiter funciona e como ele pode ser configurado.

- Testes automatizados demonstrando a eficácia e a robustez do rate limiter.

- Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.

- O servidor web deve responder na porta 8080.


## Execução do Desafio


### Descrição


Para a execução do desafio foi criado uma api no path /test que exibe apenas uma mensagem simples: `"message": "request successful"`. Essa api está configurada em um servidor web gin e utiliza o viper para carregar as variáveis de ambiente. 

As configurações do rate limiter e da porta do servidor são realizadas nas variáveis de ambiente, que estão configuradas no arquivo `docker-compose.yaml`, presente na raiz do projeto. A seguir estão as descrições das variáveis:

- **MAX_REQUESTS:** Quantidade máxima de requests por IP no período informado.
- **BLOCK_DURATION:** Bloco de tempo (em segundos) para o cálculo da quantidade de requests por IP.
- **MAX_REQUESTS_TOKEN:** Quantidade máxima de requests por TOKEN no período informado.
- **BLOCK_DURATION_TOKEN:** Bloco de tempo (em segundos) para o cálculo da quantidade de requests por TOKEN.
- **REDIS_ADDRESS:** Endereço para acesso ao REDIS
- **REDIS_PASSWORD:** Senha utilizada para conexão com o REDIS
- **PORT:** Porta TCP utilizada pelo WEBSERVER.

### Execução da Aplicação e Dependências

Para executar o sistema, basta utilizar o comando `docker-compose up --build -d` na raiz do projeto que serão geradas as imagens docker e, em seguida, os containers serão iniciados. Abaixo segue um exemplo dos containers em execução:

```bash

wander@bsnote283:~/desafio-rate-limiter$ docker-compose up --build -d
[+] Building 0.8s (10/10) FINISHED                                                                                                                                                            docker:default
 => [rate-limiter internal] load build definition from Dockerfile                                                                                          0.0s
 => => transferring dockerfile: 243B                                                                                                                       0.0s
 => WARN: FromAsCasing: 'as' and 'FROM' keywords' casing do not match (line 1)                                                                             0.0s
 => [rate-limiter internal] load metadata for docker.io/library/golang:latest                                                                              0.6s
 => [rate-limiter internal] load .dockerignore                                                                                                             0.0s
 => => transferring context: 2B                                                                                                                            0.0s
 => [rate-limiter builder 1/4] FROM docker.io/library/golang:latest@sha256:ad5c126b5cf501a8caef751a243bb717ec204ab1aa56dc41dc11be089fafcb4f                0.0s
 => [rate-limiter internal] load build context                                                                                                             0.0s
 => => transferring context: 14.25kB                                                                                                                       0.0s
 => CACHED [rate-limiter builder 2/4] WORKDIR /app                                                                                                         0.0s
 => CACHED [rate-limiter builder 3/4] COPY . .                                                                                                             0.0s
 => CACHED [rate-limiter builder 4/4] RUN GOOS=linux CGO_ENABLED=0 go build -C "cmd/server" -ldflags="-w -s" -o  server .                                  0.0s
 => CACHED [rate-limiter stage-1 1/1] COPY --from=builder /app/cmd/server .                                                                                0.0s
 => [rate-limiter] exporting to image                                                                                                                      0.0s
 => => exporting layers                                                                                                                                    0.0s
 => => writing image sha256:eafdcac4e6cfb55113e7d77a69b59f1d80da176d26d40609d0db395d4e8f80e9                                                               0.0s
 => => naming to docker.io/library/desafio-rate-limiter-rate-limiter                                                                                       0.0s
 Container redis  Creating
 Container redis  Created
 Container rate-limiter  Creating
 Container rate-limiter  Created
 Container redis  Starting
 Container redis  Started
 Container rate-limiter  Starting
 Container rate-limiter  Started
wander@bsnote283:~/desafio-rate-limiter$ 
wander@bsnote283:~/desafio-rate-limiter$ docker ps
CONTAINER ID   IMAGE                               COMMAND                  CREATED         STATUS         PORTS                                       NAMES
ee58def613b3   desafio-rate-limiter-rate-limiter   "./server"               8 seconds ago   Up 7 seconds   0.0.0.0:8080->8080/tcp, :::8080->8080/tcp   rate-limiter
fad9eaa73e7e   redis                               "docker-entrypoint.s…"   8 seconds ago   Up 7 seconds   0.0.0.0:6379->6379/tcp, :::6379->6379/tcp   redis
wander@bsnote283:~/desafio-rate-limiter$ 


```

Com essa execução, são inicializados dois containers: um do redis e um do rate-limiter:


```bash

wander@bsnote283:~/desafio-rate-limiter$ docker ps
CONTAINER ID   IMAGE                               COMMAND                  CREATED         STATUS         PORTS                                       NAMES
ee58def613b3   desafio-rate-limiter-rate-limiter   "./server"               8 seconds ago   Up 7 seconds   0.0.0.0:8080->8080/tcp, :::8080->8080/tcp   rate-limiter
fad9eaa73e7e   redis                               "docker-entrypoint.s…"   8 seconds ago   Up 7 seconds   0.0.0.0:6379->6379/tcp, :::6379->6379/tcp   redis
wander@bsnote283:~/desafio-rate-limiter$ 


```

A partir deste ponto o sistema já está disponível e os testes já podem ser realizados.


### Testes de funcionalidade da Aplicação

Para a realização dos testes funcionais foi criado um arquivo http no caminho `test/teste.http` para a realização dentro do próprio VScode. Neste arquivo foram inseridas duas chamadas para o endpoint `/test` da API: uma sem token e outra com o header de token `API_KEY`. 

Dessa forma, para realizar os testes, basta iniciar os containers (como descrito nos passos anteriores) e executar as chamadas. A seguir está um exemplo de execução das execuções dessas chamadas presentes no arquivo:

![teste01.png](/.img/teste01.png)



### Testes de Carga


Para a realização do teste de carga, utlizamos o **hey**, que é pequena ferramenta de benchmark para HTTP. Para instalar, basta executar o seguinte comando:

```bash

go install github.com/rakyll/hey@latest

```

Feito isso, já é possível executar o teste de carga. Primeiramente vamos executar o teste de carga utilizando o token. Estes serão os parâmetros utilizados para executar o teste:

- **-n 1000:** Número total de requisições a serem enviadas.
- **-c 100:** Número de requisições concorrentes.
- **-H "API_KEY: abc123":** Cabeçalho HTTP para incluir o token de acesso.

A seguir está o exemplo da execução do teste com os parâmetros descritos acima:

```bash

wander@bsnote283:~/desafio-rate-limiter$ hey -n 1000 -c 100 -H "API_KEY: abc123" http://localhost:8080/test

Summary:
  Total:	0.0847 secs
  Slowest:	0.0354 secs
  Fastest:	0.0002 secs
  Average:	0.0077 secs
  Requests/sec:	11803.3913
  
  Total data:	100400 bytes
  Size/request:	100 bytes

Response time histogram:
  0.000 [1]	|
  0.004 [86]	|■■■■■
  0.007 [668]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.011 [146]	|■■■■■■■■■
  0.014 [10]	|■
  0.018 [2]	|
  0.021 [24]	|■
  0.025 [10]	|■
  0.028 [21]	|■
  0.032 [6]	|
  0.035 [26]	|■■


Latency distribution:
  10% in 0.0040 secs
  25% in 0.0057 secs
  50% in 0.0061 secs
  75% in 0.0071 secs
  90% in 0.0107 secs
  95% in 0.0251 secs
  99% in 0.0341 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0003 secs, 0.0002 secs, 0.0354 secs
  DNS-lookup:	0.0002 secs, 0.0000 secs, 0.0039 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0007 secs
  resp wait:	0.0074 secs, 0.0001 secs, 0.0305 secs
  resp read:	0.0000 secs, 0.0000 secs, 0.0014 secs

Status code distribution:
  [200]	100 responses
  [429]	900 responses



wander@bsnote283:~/desafio-rate-limiter$ 


```

Neste exemplo de execução, a variável `MAX_REQUESTS_TOKEN` e a variável `BLOCK_DURATION_TOKEN` estavam configuradas com os valores de 100 e 60 (respectivamente) e por isso, houveram 100 respostas de código 200 e as outras 900 foram 429.

Com este resultado conseguimos comprovar que o rate limiter está funcionando perfeitamente, mesmo sobre uma alta demanda.

Para executar o teste utilizando os valores por IP, basta remover o parâmentro `-H "API_KEY: abc123` do comando e executar novamente. Abaixo segue o exemplo da requisição considerando o IP:

```bash

wander@bsnote283:~/desafio-rate-limiter$ hey -n 1000 -c 100 http://localhost:8080/test

Summary:
  Total:	0.0702 secs
  Slowest:	0.0145 secs
  Fastest:	0.0002 secs
  Average:	0.0065 secs
  Requests/sec:	14245.5297
  
  Total data:	106480 bytes
  Size/request:	106 bytes

Response time histogram:
  0.000 [1]	|
  0.002 [10]	|■
  0.003 [2]	|
  0.004 [1]	|
  0.006 [229]	|■■■■■■■■■■■■■■
  0.007 [656]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.009 [27]	|■■
  0.010 [14]	|■
  0.012 [20]	|■
  0.013 [30]	|■■
  0.015 [10]	|■


Latency distribution:
  10% in 0.0057 secs
  25% in 0.0059 secs
  50% in 0.0062 secs
  75% in 0.0066 secs
  90% in 0.0074 secs
  95% in 0.0112 secs
  99% in 0.0134 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0003 secs, 0.0002 secs, 0.0145 secs
  DNS-lookup:	0.0002 secs, 0.0000 secs, 0.0040 secs
  req write:	0.0001 secs, 0.0000 secs, 0.0017 secs
  resp wait:	0.0061 secs, 0.0001 secs, 0.0127 secs
  resp read:	0.0000 secs, 0.0000 secs, 0.0010 secs

Status code distribution:
  [200]	20 responses
  [429]	980 responses



wander@bsnote283:~/desafio-rate-limiter$ 


```

Observe que, ao remover o header `API_KEY`, os limites utilizados para bloquear as requsições foram alterados e passaram a respeitar as variávies `MAX_REQUESTS` e `BLOCK_DURATION`, que no momento da execução dos testes estavam configuradas com os valores 20 e 60 respectivamente.

Com a execução destes testes podemos confirmar que o sistema está apto a ser executado mesmo sobre alta demanda sem compromoter a funcionalidade de rate limite.
