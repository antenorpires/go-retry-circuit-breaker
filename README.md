# Retry + Circuit Breaker

Este repositório tem como objetivo **demonstrar, de forma didática**, o uso dos padrões de **Retry** e **Circuit Breaker** em uma arquitetura de microserviços utilizando **Golang**.

O projeto foi criado para **estudo, experimentação e contribuição da comunidade**, simulando falhas reais como timeout, erro HTTP 5xx e recuperação automática de serviços.

---

## 📐 Arquitetura

O projeto é composto por dois microserviços:

* **service-a** → consumidor

  * Implementa **Retry** e **Circuit Breaker**
  * Possui interface web simples para testes
* **service-b** → provedor

  * Simula latência, falhas e respostas de sucesso
  * Não conhece conceitos de resiliência

Fluxo da requisição:

```mermaid
Service A
   |
   |  HTTP POST
   v
Service B
   |
   |-- 200 OK        → conta sucesso
   |-- 500 Error     → conta falha
   |-- Timeout (2s)  → conta falha
```

---

## 🧰 Tecnologias utilizadas

* Go 1.17+
* net/http
* github.com/sony/gobreaker
* github.com/hashicorp/go-retryablehttp

---

## 1️⃣ Instalando o Go no Linux

Em sistemas baseados em Debian/Ubuntu:

```bash
sudo apt install golang-go -y
```

Verifique a instalação:

```bash
go version
```

---

## 2️⃣ Iniciando o microserviço **service-a**

1. Navegue até o diretório:

```bash
cd service-a
```

2. Execute o serviço:

```bash
go run main.go
```

3. O serviço estará disponível em:

```
http://localhost:9090
```

Ele exibe uma página HTML simples com um campo de entrada para testes.

---

## 3️⃣ Iniciando o microserviço **service-b**

1. Navegue até o diretório:

```bash
cd service-b
```

2. Execute o serviço:

```bash
go run main.go
```

3. O serviço estará disponível em:

```
http://localhost:9091
```

⚠️ Aceita **apenas requisições HTTP POST**.

---

## 4️⃣ Testando o comportamento de Retry e Circuit Breaker

### Casos de teste simulados

| ID enviado    | Comportamento           |
| ------------- | ----------------------- |
| `123`         | ✅ Sucesso (200 + JSON)  |
| `fail`        | ❌ Erro HTTP 500         |
| outro valor   | ⚠️ 200 + `"failed"`     |
| muitos `fail` | 🔥 Circuit Breaker abre |

---

### Teste de Retry

1. Acesse:

```
http://localhost:9090
```

2. Envie o valor `123`
3. Pare o **service-b** (`Ctrl + C`)
4. Envie novamente a requisição pelo **service-a**
5. Observe as tentativas de retry acontecendo
6. Inicie novamente o **service-b**

O serviço volta a responder antes de um erro definitivo.

---

### Teste de Circuit Breaker

1. Envie várias requisições com:

```
id = fail
```

2. Após aproximadamente 10 requisições com **50% ou mais de falha**:

```
Circuit Breaker entra em estado OPEN
```

3. As requisições passam a falhar imediatamente
4. Após o tempo de timeout:

```
OPEN → HALF-OPEN → CLOSED (se sucesso)
```

---

## 🧠 Observações importantes

* O **Circuit Breaker pertence ao consumidor** (service-a)
* O **provedor (service-b) não conhece resiliência**
* Retry resolve falhas transitórias
* Circuit Breaker protege contra falhas persistentes

Essa separação reflete o que é usado em ambientes reais de produção.

---

## 📄 Licença

Este projeto é licenciado sob a **MIT License**.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
