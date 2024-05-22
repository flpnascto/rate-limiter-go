# Rate Limiter with Go

## :notebook_with_decorative_cover: Sobre o Projeto

Aplicação desenvolvida no curso **Pós Graduação Go Expert - Full Cycle** na linguagem Go.
A aplicação contém um _rate limiter_ que opera como um middleware em um webserver sendo capaz de limitar o número de requisições com base em dois critérios:
- Endereço IP: restringir requisições recebidas de um único endereço IP após atingir o limite de requisições por segundo.
- Token de Acesso: restringir requisições baseadas em um token de acesso único após atingir o limite de requisições por segundo
- As configurações de limite do token de acesso se sobrepões as do IP

## :computer: Tecnologias Aplicadas
* Go
* Net/Http
* Redis
* Docker

## :arrow_forward: Executando a aplicação

Para executar a aplicação siga as instruções abaixo.

### Pré-requisitos

Primeiramente é necessário que possua instalado as seguintes ferramentas: Go, Git, Docker.
Além disto é bom ter um editor para trabalhar com o código como VSCode.

### Instalação

1. Faça uma cópia do repositório (HTTPS ou SSH)
   ```sh
   git clone https://github.com/flpnascto/rate-limiter-go
   ```
   ```sh
   git clone git@github.com:flpnascto/rate-limiter-go.git
   ```
2. Acessar a pasta do repositório local e instanciar o Docker
   - Execute o comando `docker-compose up -d`
   - O webserver já estará em execução na porta 8080 (padrão)

3. Realizar requisições via terminal
  - Requisição padrão: `curl localhost:8080/api/ping`
  - Requisição com token: `curl -H "API_KEY: ab123" http://localhost:8080/api/ping`

4. Realizar requisições com REST CLient
  - Na raiz do repositório existe um arquivo api.http com dois exemplos de requisição, considerando o caso de haver ou não um token.
  - Para utilizar este arquivo é necessário a extensão **REST Client** instalada no VSCode

5. Realizar o teste do teste de carga
  - No Docker Desktop abra o container `goserver` e acesse o terminal do contêiner
  - Execute no terminal o comando `go test audit/test_test.go`

6. Alterando as configurações do rate limiter
  - Acesse o arquivo ./cmd/config.json para alterar
    - A quantidade máxima de requisições por segundo por IP ou TOKEN
    - O tempo, em segundos, de bloqueio para requisições que estão acima do limite
  - Caso realize alterações neste arquivo é necessário recriar o docker compose
    - Utilize o comando `docker-compose up -d --build`