# DynamoDB API

Uma API REST robusta desenvolvida em Go para gerenciar eventos com suporte a AWS DynamoDB e AWS Lambda. O projeto oferece múltiplas formas de deployment e é totalmente testado com cobertura acima de 90%.

## 📋 Tabela de Conteúdos

- [Visão Geral](#visão-geral)
- [Requisitos](#requisitos)
- [Instalação](#instalação)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Bibliotecas Utilizadas](#bibliotecas-utilizadas)
- [Configuração](#configuração)
- [Uso](#uso)
- [API REST](#api-rest)
- [Testes](#testes)
- [Cobertura de Código](#cobertura-de-código)
- [Deploy](#deploy)
- [Contribuindo](#contribuindo)

## 🎯 Visão Geral

Este projeto é uma API REST completa para gerenciamento de eventos com as seguintes características:

- **Dual Deployment**: Funciona como servidor HTTP standalone ou como AWS Lambda function
- **Armazenamento Flexível**: Suporta armazenamento em memória (desenvolvimento) ou DynamoDB (produção)
- **RFC 9457 Compliance**: Respostas de erro segue o padrão RFC 9457 (Problem Details for HTTP APIs)
- **Testes Abrangentes**: Cobertura de código > 95% com testes unitários completos
- **Validação de Dados**: Validação automática de eventos com data e status code

## 📦 Requisitos

- **Go**: 1.25.5 ou superior
- **AWS CLI**: (opcional, para configurar credenciais da AWS)
- **Docker**: (opcional, para containerizar a aplicação)

### Dependências de Produção

- `github.com/aws/aws-sdk-go-v2`: AWS SDK v2 para Go
- `github.com/aws/aws-sdk-go-v2/service/dynamodb`: Cliente DynamoDB
- `github.com/aws/aws-lambda-go`: Framework para funções Lambda
- `github.com/google/uuid`: Geração de UUIDs

## 🚀 Instalação

### 1. Clone o repositório

```bash
git clone <repository-url>
cd dynamodb-api
```

### 2. Instale as dependências

```bash
go mod download
```

### 3. Configure as variáveis de ambiente (opcional)

```bash
export API_PORT=8080  # Porta padrão é 8080
```

### 4. Compile o projeto

```bash
go build -o api
```

## 📁 Estrutura do Projeto

```
dynamodb-api/
├── apis/                      # Camada de entrada da API
│   ├── api_http.go           # HTTP server
│   ├── api_http_test.go      # Testes do HTTP server
│   ├── api_lambda.go         # AWS Lambda handler
│   └── api_lambda_test.go    # Testes do Lambda handler
│
├── handlers/                  # Camada de lógica de requisições
│   ├── handler_http.go       # Handlers HTTP
│   ├── handler_http_test.go  # Testes dos handlers HTTP
│   ├── handler_lambda.go     # Handlers Lambda
│   └── handler_lambda_test.go # Testes dos handlers Lambda
│
├── repositories/             # Camada de persistência
│   ├── dynamodb.go           # Implementação DynamoDB
│   ├── dynamodb_test.go      # Testes DynamoDB (95.1% cobertura)
│   ├── memorydb.go           # Implementação em memória
│   └── memorydb_test.go      # Testes MemoryDB
│
├── models/                   # Modelos de dados
│   ├── event.go              # Modelo de Evento
│   ├── event_test.go         # Testes do modelo Event
│   ├── error_response.go     # Modelo de resposta de erro (RFC 9457)
│   └── error_response_test.go # Testes do ErrorResponse
│
├── interfaces/               # Contatos/Interfaces
│   ├── repository.go         # Interface Repository
│   └── dynamodb_client.go    # Interface DynamoDBClient
│
├── main.go                   # Ponto de entrada da aplicação
├── go.mod                    # Definição do módulo
├── go.sum                    # Checksums das dependências
└── README.md                 # Este arquivo
```

## 📚 Bibliotecas Utilizadas

### Dependências Diretas

| Biblioteca | Versão | Propósito |
|-----------|--------|----------|
| `github.com/aws/aws-lambda-go` | v1.52.0 | Framework para AWS Lambda |
| `github.com/aws/aws-sdk-go-v2` | v1.41.1 | AWS SDK para Go |
| `github.com/aws/aws-sdk-go-v2/config` | v1.32.7 | Configuração AWS SDK |
| `github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue` | v1.20.30 | Conversão de atributos DynamoDB |
| `github.com/aws/aws-sdk-go-v2/service/dynamodb` | v1.53.6 | Cliente DynamoDB |
| `github.com/google/uuid` | v1.6.0 | Geração de UUIDs |

### Dependências Indiretas

As dependências indiretas são gerenciadas automaticamente pelo `go mod` e incluem suporte a credenciais AWS, serviços de configuração e autenticação.

## ⚙️ Configuração

### Variáveis de Ambiente

```bash
# Porta da API (padrão: 8080)
export API_PORT=8080

# Região AWS (padrão: conforme configuração AWS)
export AWS_REGION=us-east-1

# Tabela DynamoDB (configurável no código)
# DYNAMODB_TABLE=eventos

# Profile AWS
export AWS_PROFILE=default
```

### Configuração de Credenciais AWS

#### Usando arquivo ~/.aws/credentials

```ini
[default]
aws_access_key_id = YOUR_ACCESS_KEY
aws_secret_access_key = YOUR_SECRET_KEY
```

#### Usando variáveis de ambiente

```bash
export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY
export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
```

#### Usando IAM Role (para Lambda)

Configure as permissões de execução da função Lambda para ter acesso ao DynamoDB.

## 💻 Uso

### Iniciar o Servidor HTTP

```bash
./api
```

O servidor iniciará na porta 8080 (ou conforme `API_PORT`).

```
starting server on :8080
```

### Endpoints Disponíveis

#### Health Check

```bash
GET /health
```

**Resposta:**
```
200 OK
OK
```

#### Criar Evento

```bash
POST /eventos
Content-Type: application/json

{
  "date": "2024-01-15T10:30:00Z",
  "statusCode": 200,
  "statusMessage": "Success",
  "metadata": {
    "userId": "123",
    "action": "create"
  }
}
```

**Respostas:**

- `201 Created`: Evento criado com sucesso
- `400 Bad Request`: Dados inválidos ou data/statusCode ausentes
- `500 Internal Server Error`: Erro ao salvar o evento

#### Obter Evento

```bash
GET /eventos/{id}
```

**Exemplos:**

```bash
curl http://localhost:8080/eventos/550e8400-e29b-41d4-a716-446655440000
```

**Respostas:**

- `200 OK`: Evento encontrado
- `404 Not Found`: Evento não existe
- `400 Bad Request`: ID ausente ou inválido

#### Atualizar Evento

```bash
PUT /eventos/{id}
Content-Type: application/json

{
  "date": "2024-01-15T10:30:00Z",
  "statusCode": 201,
  "statusMessage": "Updated",
  "metadata": {}
}
```

**Respostas:**

- `201 Created`: Evento atualizado com sucesso
- `400 Bad Request`: Dados inválidos ou ID ausente
- `500 Internal Server Error`: Erro ao salvar o evento

#### Deletar Evento

```bash
DELETE /eventos/{id}
```

**Exemplos:**

```bash
curl -X DELETE http://localhost:8080/eventos/550e8400-e29b-41d4-a716-446655440000
```

**Respostas:**

- `204 No Content`: Evento deletado com sucesso
- `404 Not Found`: Evento não existe
- `400 Bad Request`: ID ausente ou inválido

### Exemplo de Resposta de Erro

```json
{
  "type": "about:blank",
  "title": "Bad Request",
  "status": 400,
  "detail": "Missing event ID in URL",
  "instance": "/eventos/",
  "code": "INVALID_REQUEST"
}
```

## 🧪 Testes

### Executar Todos os Testes

```bash
go test ./...
```

### Executar Testes de um Pacote Específico

```bash
# Testes dos handlers
go test -v ./handlers

# Testes dos repositórios
go test -v ./repositories

# Testes dos modelos
go test -v ./models

# Testes das APIs
go test -v ./apis
```

### Executar com Verbosidade

```bash
go test -v ./...
```

### Executar Teste Específico

```bash
go test -run TestDynamoDBCreateSuccess ./repositories
```

### Testes com Timeout

```bash
go test -timeout 30s ./...
```

## 📊 Cobertura de Código

### Gerar Relatório de Cobertura

```bash
# Gerar arquivo de cobertura
go test -coverprofile=coverage.out ./...

# Exibir cobertura em cada função
go tool cover -func=coverage.out

# Gerar relatório HTML
go tool cover -html=coverage.out -o coverage.html
```

### Cobertura Atual por Pacote

| Pacote | Cobertura | Status |
|--------|-----------|--------|
| `api/models` | 100.0% | ✅ Completo |
| `api/repositories` | 95.1% | ✅ Excelente |
| `api/handlers` | 88.3% | ✅ Muito Bom |
| `api/apis` | 21.1% | ⚠️ Necessário melhorar |
| `api/interfaces` | N/A | - |

**Nota**: A cobertura do pacote `apis` é limitada porque a função `Run()` inicia um servidor HTTP que não pode ser testado facilmente em testes unitários.

### Testes por Pacote

#### Models (100% - 2 arquivos)
- Event: Validação de data e status code
- ErrorResponse: Estrutura RFC 9457

#### Repositories (95.1% - 21+ testes)
- **DynamoDB** (92.3% - 18 testes):
  - Create, Save, Get, Delete
  - FindByDateAndReturnCode
  - Casos de erro e edge cases
  
- **MemoryDB** (100% - 14 testes):
  - Operações CRUD completas
  - Validação de expiração TTL
  - Casos de erro

#### Handlers (88.3% - 42+ testes)
- **HTTP** (88% - 24 testes):
  - Todos os métodos HTTP (GET, POST, PUT, DELETE)
  - Health check
  - Validação de entrada
  - Tratamento de erros
  
- **Lambda** (88% - 18+ testes):
  - Todos os métodos HTTP
  - Routing
  - Serialização JSON
  - Tratamento de erros

#### APIs (21.1%)
- Configuração do servidor HTTP
- Injeção de dependências

## 🚢 Deploy

### Deploy Local

```bash
# Compilar
go build -o api

# Executar
./api

# Com porta customizada
API_PORT=9000 ./api
```

### Deploy em Docker

```dockerfile
FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
```

**Build e run:**

```bash
docker build -t dynamodb-api .
docker run -p 8080:8080 dynamodb-api
```

### Deploy em AWS Lambda

1. Compile o binário para Linux:

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap ./main.go
zip lambda-function.zip bootstrap
```

2. Crie uma função Lambda com o binário compilado
3. Configure a variável de ambiente `AWS_REGION`
4. Configure IAM Role com permissões de DynamoDB

### Deploy em AWS ECS

1. Build a imagem Docker
2. Faça push para ECR
3. Crie uma task definition
4. Crie um serviço ECS

## 🔍 Modelos de Dados

### Event

```go
type Event struct {
    Id            string            `json:"id"`
    Date          time.Time         `json:"date"`
    StatusCode    int               `json:"statusCode"`
    StatusMessage string            `json:"statusMessage"`
    Metadata      map[string]string `json:"metadata,omitempty"`
    Expiration    int64             `json:"-"`
}
```

**Validação:**
- `Date`: Obrigatório, não pode ser zero
- `StatusCode`: Obrigatório, deve ser >= 0

### ErrorResponse (RFC 9457)

```go
type ErrorResponse struct {
    Type     string `json:"type"`
    Status   int    `json:"status"`
    Title    string `json:"title"`
    Detail   string `json:"detail"`
    Instance string `json:"instance"`
    Code     string `json:"code,omitempty"`
}
```

## 🔐 Segurança

### Boas Práticas Implementadas

- ✅ Validação de entrada em todos os endpoints
- ✅ Headers de segurança padrão
- ✅ Timeouts de requisição (30s read/write, 60s idle)
- ✅ Limite de tamanho de header (1MB)
- ✅ Autenticação via AWS IAM (Lambda)
- ✅ Geração de IDs com UUID v4

### Recomendações

1. **Autenticação**: Adicione API Gateway com autenticação
2. **CORS**: Configure CORS se necessário
3. **Rate Limiting**: Implemente rate limiting
4. **HTTPS**: Use HTTPS em produção
5. **WAF**: Considere usar AWS WAF

## 🐛 Troubleshooting

### Erro: "connection refused"

**Causa**: Servidor não está rodando na porta configurada

**Solução**:
```bash
# Verificar se a porta está em uso
lsof -i :8080

# Usar outra porta
API_PORT=9000 ./api
```

### Erro: "NoCredentialsError"

**Causa**: Credenciais AWS não configuradas

**Solução**:
```bash
# Configure credenciais
aws configure

# Ou use variáveis de ambiente
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
```

### Erro: "ResourceNotFoundException"

**Causa**: Tabela DynamoDB não existe

**Solução**:
```bash
# A tabela será criada automaticamente na primeira execução
# Se não funcionar, crie manualmente via AWS Console
```

### Testes falhando

**Causa**: Dependências não instaladas

**Solução**:
```bash
go mod tidy
go mod download
go test ./...
```

## 📈 Performance

### Benchmarks

Para rodar benchmarks (a adicionar):

```bash
go test -bench=. ./...
```

### Otimizações

- Usar MemoryDB para desenvolvimento (em memória)
- Usar DynamoDB para produção (escalável)
- Connection pooling automático do AWS SDK
- Timeouts configuráveis

## 📝 Logging

O projeto usa o package `log` padrão do Go. Logs são enviados para stdout:

```
starting server on :8080
```

Para melhorar o logging, considere usar:
- `github.com/sirupsen/logrus`
- `go.uber.org/zap`
- AWS CloudWatch Logs

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Checklist para Contribuições

- [ ] Testes unitários adicionados
- [ ] Cobertura de código mantida > 90%
- [ ] `go fmt` executado
- [ ] `go vet` sem erros
- [ ] README atualizado se necessário

## 📄 Licença

Este projeto está licenciado sob a MIT License - veja o arquivo LICENSE para detalhes.

## 📞 Suporte

Para reportar problemas ou sugerir melhorias, abra uma issue no repositório.

## 🎓 Aprendizados e Boas Práticas

Este projeto demonstra:

1. **Arquitetura Limpa**: Separação clara entre camadas (handlers, repositories, models)
2. **Interface Segregation**: Uso de interfaces para desacoplamento
3. **Dependency Injection**: Injeção de dependências para testabilidade
4. **Testes Abrangentes**: Unit tests com mocks e table-driven tests
5. **Error Handling**: Tratamento robusto de erros
6. **RFC Compliance**: Seguindo padrões web (RFC 9457)
7. **Multi-deployment**: Flexibilidade entre HTTP e Lambda
8. **Configuration Management**: Configuração via variáveis de ambiente

## 🔗 Recursos Úteis

- [Go Documentation](https://golang.org/doc/)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)
- [AWS Lambda Go](https://github.com/aws/aws-lambda-go)
- [RFC 9457 - Problem Details](https://www.rfc-editor.org/rfc/rfc9457)
- [DynamoDB Documentation](https://docs.aws.amazon.com/dynamodb/)

---

**Versão**: 1.0.0  
**Última atualização**: Janeiro 2026  
**Linguagem**: Go 1.25.5
