# **Projeto de Circuit Breaker com Retry em Go**

Este projeto demonstra o uso de circuit breaker com retry em microserviços utilizando Golang.

## **1. Instalando o GoLang no Linux**

Para instalar o GoLang em sistemas baseados em Linux, execute o seguinte comando:

```bash
sudo apt install golang-go -y
```

## **2. Iniciando o microserviço 'service-b'**
1. Navegue até o diretório do microserviço service-b:
```bash
cd service-b
```

2. Inicie o microserviço executando o comando abaixo:
```bash
go run main.go
```

3. O microserviço estará rodando em http://localhost:9091

## **3. Iniciando o microserviço 'service-a'**
1. Navegue até o diretório do microserviço service-a:
```bash
cd service-a
```

2. Inicie o microserviço executando o comando abaixo:
```bash
go run main.go
```

3. O microserviço estará rodando em http://localhost:9090

## **4. Testando o Comportamento de Retry**
1. Acesse a página http://localhost:9090
2. No campo de entrada, digite o valor 123 e envie a requisição.
3. Pare o microserviço service-b. Para isso, você pode usar Ctrl+C no terminal onde o service-b está rodando.
4. Atualize a página no navegador do microserviço service-a. Você verá que o contador de tentativas de retry começará a aumentar.
5. Após algumas tentativas, inicie novamente o microserviço service-b, e ele começará a retornar respostas antes de exibir um erro, simulando o comportamento de retry.


## **5. Concluindo**
Agora que você testou o comportamento de retry, pode ajustar a configuração do seu circuit breaker e experimentar diferentes cenários.
Divirta-se explorando o testes de circuit breaker e retry no seu microserviço!