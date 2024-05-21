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
   - Execute o comando `docker-compose up`

3. Verificar os resultados do teste de carga
  - No Docker Desktop abra o container `goserver` e acesse o terminal

4. Alterando as configurações do rate limiter
  - Acesse o arquivo ./cmd/config.json para alterar
    - A quantidade máxima de requisições por segundo por IP ou TOKEN
    - O tempo, em segundos, de bloqueio para requisições que estão acima do limite


