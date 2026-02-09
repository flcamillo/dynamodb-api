# DynamoDB API - Documentação Técnica Detalhada

[![Go Version](https://img.shields.io/badge/Go-1.25.5-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![AWS SDK](https://img.shields.io/badge/AWS%20SDK-v2-FF9900?style=flat-square&logo=amazonaws)](https://aws.amazon.com/sdk-for-go/)
[![OpenTelemetry](https://img.shields.io/badge/OpenTelemetry-1.40.0-430098?style=flat-square)](https://opentelemetry.io)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

Uma API RESTful enterprise-grade construída em Go para gerenciar eventos com AWS DynamoDB ou repositório em memória. Oferece suporte a múltiplos modos de deployment (HTTP Server e AWS Lambda), com observabilidade completa via OpenTelemetry e integração Datadog.

## 📚 Índice

- [Visão Geral](#visão-geral)
- [Características](#características)
- [Arquitetura Detalhada](#arquitetura-detalhada)
- [Requisitos](#requisitos)
- [Instalação e Configuração](#instalação-e-configuração)
- [Executando a Aplicação](#executando-a-aplicação)
- [Endpoints da API](#endpoints-da-api)
- [Exemplos com cURL](#exemplos-com-curl)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Configuração Avançada](#configuração-avançada)
- [Telemetria e Observabilidade](#telemetria-e-observabilidade)
- [Datadog Integration](#datadog-integration)
- [Troubleshooting](#troubleshooting)

## 🎯 Visão Geral

A **DynamoDB API** é uma solução robusta para gerenciar eventos com alta performance, observabilidade e escalabilidade. Desenvolvida em Go, oferece:

- **Performance**: Processamento rápido de requisições com latência mínima
- **Escalabilidade**: Suporte a DynamoDB para escalabilidade horizontal
- **Flexibilidade**: Escolha entre DynamoDB e repositório em memória
- **Observabilidade**: Telemetria completa via OpenTelemetry
- **Confiabilidade**: TTL automático, validação robusta e tratamento de erros

## ✨ Características

### Funcionalidades Principais
- ✅ **API RESTful completa** para CRUD de eventos com validação de dados
- ✅ **Suporte dual de deployment**: HTTP Server + AWS Lambda
- ✅ **Repositórios plugáveis**: DynamoDB e In-Memory Database
- ✅ **OpenTelemetry integrado** para observabilidade completa
- ✅ **Métricas e Tracing** automáticos em todas as operações
- ✅ **Logs estruturados** via slog + OTEL bridge
- ✅ **TTL (Time To Live)** para expiração automática de registros
- ✅ **Suporte a metadata** customizável e extensível por evento
- ✅ **Integração Datadog** para APM, logs e métricas
- ✅ **Tratamento de erros** com contexto e rastreamento

### Recursos de Observabilidade
- 📊 **Métricas**: Taxa de requisições, latência, erros, duração de operações
- 📈 **Tracing Distribuído**: Rastreamento completo de requisições end-to-end
- 📝 **Logs Estruturados**: JSON estruturado com contexto completo
- 🔍 **APM**: Application Performance Monitoring via Datadog/Jaeger
- 🎯 **Atributos de Contexto**: Service name, version, environment, trace IDs

## 🏗️ Arquitetura Detalhada

### Diagrama de Componentes Avançado

```mermaid
graph TB
    subgraph "Clientes"
        HTTP_CLIENT["🖥️ Cliente HTTP"]
        LAMBDA_EVENT["⚡ AWS Lambda Event"]
        curl["📱 cURL/Postman"]
    end
    
    subgraph "Camada de Ingresso"
        HTTP_API["🌐 HTTP API Server<br/>Port 7000<br/>net/http Router"]
        LAMBDA_ROUTER["📦 Lambda Router<br/>aws-lambda-go"]
    end
    
    subgraph "Camada de Processamento"
        HTTP_HANDLER["🔧 HTTP Handler<br/>Parsing | Validation<br/>Response Formatting"]
        LAMBDA_HANDLER["🔧 Lambda Handler<br/>Event Parsing<br/>Response Builder"]
        VALIDATOR["✓ Validator<br/>Date Validation<br/>Status Code Check<br/>Metadata Validation"]
    end
    
    subgraph "Camada de Dados"
        REPO_INTERFACE["📊 Repository Interface<br/>interface{}<br/>Save | Get | Delete | Find"]
        DYNAMODB_REPO["🗄️ DynamoDB Repository<br/>AWS SDK v2<br/>PutItem | GetItem | Query"]
        MEMORY_REPO["💾 Memory Repository<br/>map[string]Event<br/>In-Process Storage"]
    end
    
    subgraph "Armazenamento"
        DYNAMODB["☁️ AWS DynamoDB<br/>Tabela: 'eventos'<br/>TTL: expiration<br/>GSI: date-statusCode"]
        MEMORY["🖥️ Memória Local<br/>Runtime Storage<br/>Para Testes"]
    end
    
    subgraph "Observabilidade"
        OTEL_SDK["📡 OpenTelemetry SDK<br/>Tracer | Meter<br/>Logger Bridge"]
        OTEL_EXPORTER["🔄 OTLP Exporter<br/>gRPC Protocol<br/>:4317"]
    end
    
    subgraph "Backends de Observabilidade"
        DATADOG["🐶 Datadog Agent<br/>Traces | Metrics<br/>Logs | APM<br/>:8126"]
        JAEGER["🔍 Jaeger Backend<br/>Trace Storage<br/>UI: :16686"]
        PROMETHEUS["📊 Prometheus<br/>Metrics Storage<br/>UI: :9090"]
        OTEL_COLLECTOR["📡 OTEL Collector<br/>Receives OTLP<br/>Routes to Backends<br/>:4317"]
    end
    
    subgraph "Infraestrutura"
        DOCKER["🐳 Docker Compose<br/>DynamoDB Local<br/>OTEL Collector<br/>Datadog Agent"]
    end
    
    HTTP_CLIENT -->|HTTP| HTTP_API
    curl -->|HTTP| HTTP_API
    LAMBDA_EVENT -->|Event| LAMBDA_ROUTER
    
    HTTP_API -->|Parse| HTTP_HANDLER
    LAMBDA_ROUTER -->|Parse| LAMBDA_HANDLER
    
    HTTP_HANDLER -->|Validate| VALIDATOR
    LAMBDA_HANDLER -->|Validate| VALIDATOR
    
    VALIDATOR -->|OK| REPO_INTERFACE
    
    REPO_INTERFACE -->|Implements| DYNAMODB_REPO
    REPO_INTERFACE -->|Implements| MEMORY_REPO
    
    DYNAMODB_REPO -->|AWS SDK| DYNAMODB
    MEMORY_REPO -->|In-Memory| MEMORY
    
    HTTP_HANDLER -->|Telemetry| OTEL_SDK
    LAMBDA_HANDLER -->|Telemetry| OTEL_SDK
    VALIDATOR -->|Telemetry| OTEL_SDK
    DYNAMODB_REPO -->|Metrics| OTEL_SDK
    
    OTEL_SDK -->|Export| OTEL_EXPORTER
    OTEL_EXPORTER -->|OTLP| OTEL_COLLECTOR
    
    OTEL_COLLECTOR -->|Process| DATADOG
    OTEL_COLLECTOR -->|Process| JAEGER
    OTEL_COLLECTOR -->|Process| PROMETHEUS
    
    DOCKER -.->|Provides| OTEL_COLLECTOR
    DOCKER -.->|Provides| DYNAMODB
    DOCKER -.->|Provides| DATADOG
```

### Fluxo de Requisição Detalhado

```mermaid
sequenceDiagram
    participant Client as 🖥️ Cliente
    participant Server as 🌐 HTTP API
    participant Handler as 🔧 Handler
    participant Validator as ✓ Validator
    participant Repo as 📊 Repository
    participant DB as ☁️ DynamoDB
    participant OTEL as 📡 OpenTelemetry
    participant Datadog as 🐶 Datadog
    
    Client->>Server: POST /eventos (JSON)
    activate Server
    
    Server->>Handler: RouteRequest()
    activate Handler
    
    Handler->>Handler: ParseJSON()
    Handler->>OTEL: StartSpan("POST /eventos")
    activate OTEL
    
    Handler->>Validator: Validate(event)
    activate Validator
    Validator->>Validator: CheckDate()
    Validator->>Validator: CheckStatusCode()
    Validator->>Validator: CheckMetadata()
    Validator-->>Handler: ValidationResult
    deactivate Validator
    
    alt Validation Failed
        Handler->>OTEL: RecordError()
        Handler-->>Server: 400 Bad Request
    else Validation Passed
        Handler->>Handler: GenerateUUID()
        Handler->>Handler: CalculateExpiration()
        
        Handler->>OTEL: AddEvent("attempting_save")
        Handler->>Repo: Save(event)
        activate Repo
        
        Repo->>DB: PutItem(event)
        activate DB
        DB-->>Repo: Success
        deactivate DB
        
        Repo->>OTEL: RecordMetric("dynamodb.latency_ms")
        Repo-->>Handler: SaveResult
        deactivate Repo
        
        Handler->>OTEL: AddAttributes()
        Handler->>OTEL: RecordMetric("http.requests.total")
        Handler->>OTEL: RecordMetric("http.request.duration_ms")
        
        Handler-->>Server: 201 Created + JSON
    end
    
    OTEL-->>Datadog: ExportOTLP(traces, metrics, logs)
    deactivate OTEL
    
    Server-->>Client: HTTP Response
    deactivate Handler
    deactivate Server
```

### Ciclo de Vida de um Evento

```mermaid
graph TB
    subgraph "Criação"
        C1["1. Cliente envia JSON<br/>POST /eventos"]
        C2["2. Gera UUID único<br/>usando google/uuid"]
        C3["3. Calcula Expiração<br/>agora + TTL"]
    end
    
    subgraph "Validação"
        V1["✓ Validação de Data<br/>RFC3339 format"]
        V2["✓ Status Code<br/>0-599"]
        V3["✓ Metadata<br/>map[string]string"]
    end
    
    subgraph "Armazenamento"
        S1["💾 Armazena no DynamoDB<br/>Tabela: eventos"]
        S2["⏰ Define TTL<br/>Campo: expiration"]
    end
    
    subgraph "Retorno"
        R1["✅ Resposta 201 Created<br/>Com ID e metadados"]
    end
    
    subgraph "Expiração"
        E1["⏳ Aguarda expiração<br/>DynamoDB TTL"]
        E2["🗑️ Remove automaticamente<br/>Após TTL expirar"]
    end
    
    C1 --> C2
    C2 --> C3
    C3 --> V1
    V1 --> V2
    V2 --> V3
    V3 --> S1
    S1 --> S2
    S2 --> R1
    R1 --> E1
    E1 --> E2
```

### Arquitetura de Logs e Traces

```mermaid
graph TB
    subgraph "Geração"
        GEN1["🔧 Handlers<br/>HTTP & Lambda"]
        GEN2["📊 Repositories<br/>DB Operations"]
        GEN3["✓ Validators<br/>Validation Logic"]
    end
    
    subgraph "Coleta"
        SLOG["📝 slog Logger<br/>Structured Logging"]
        OTEL_BRIDGE["🌉 OTEL slog Bridge<br/>Integração automática"]
        TRACER["🎯 OTEL Tracer<br/>Distributed Tracing"]
        METER["📈 OTEL Meter<br/>Metrics Collection"]
    end
    
    subgraph "SDK"
        SDK["📦 OpenTelemetry SDK<br/>Batch Processing<br/>Sampling"]
    end
    
    subgraph "Exportação"
        EXPORTER["🔄 OTLP Exporter<br/>gRPC Protocol<br/>Batch Export"]
    end
    
    subgraph "Transporte"
        NETWORK["🌐 Network<br/>localhost:4317"]
    end
    
    subgraph "Collector"
        COLLECTOR["📡 OTEL Collector<br/>Receive • Process • Export"]
    end
    
    subgraph "Destinos"
        DD["🐶 Datadog<br/>APM | Logs | Metrics"]
        JAEGER["🔍 Jaeger<br/>Traces"]
        PROM["📊 Prometheus<br/>Metrics"]
    end
    
    GEN1 --> SLOG
    GEN2 --> SLOG
    GEN3 --> SLOG
    
    GEN1 --> TRACER
    GEN2 --> TRACER
    GEN3 --> TRACER
    
    GEN2 --> METER
    
    SLOG --> OTEL_BRIDGE
    OTEL_BRIDGE --> SDK
    TRACER --> SDK
    METER --> SDK
    
    SDK --> EXPORTER
    EXPORTER --> NETWORK
    NETWORK --> COLLECTOR
    
    COLLECTOR --> DD
    COLLECTOR --> JAEGER
    COLLECTOR --> PROM
```

### Estado e Transições de Eventos

```mermaid
stateDiagram-v2
    [*] --> Criado: POST /eventos
    
    Criado --> Recuperável: Salvo no DynamoDB
    
    Recuperável --> Consultado: GET /eventos/{id}
    Recuperável --> Listado: GET /eventos
    Recuperável --> Atualizado: PUT /eventos/{id}
    Recuperável --> Deletado: DELETE /eventos/{id}
    
    Consultado --> Recuperável
    Listado --> Recuperável
    
    Atualizado --> Recuperável: TTL recalculado
    
    Recuperável --> Expirado: Tempo TTL passado
    Expirado --> [*]: Removido pelo DynamoDB
    
    Deletado --> [*]: Removo imediato
```

## 📦 Requisitos

### Versões Mínimas
- **Go**: 1.25.5+
- **AWS SDK for Go**: v2 (v1.41.1+)
- **Docker**: 20.10+ (para ambiente local com Datadog)
- **Docker Compose**: 2.0+ (para orquestração)

### Dependências do Projeto

```go
require (
    github.com/aws/aws-lambda-go v1.52.0           // AWS Lambda runtime
    github.com/aws/aws-sdk-go-v2 v1.41.1            // AWS SDK base
    github.com/aws/aws-sdk-go-v2/config v1.32.7     // AWS configuration
    github.com/aws/aws-sdk-go-v2/service/dynamodb v1.55.0
    github.com/google/uuid v1.6.0                   // UUID generation
    
    // OpenTelemetry Core
    go.opentelemetry.io/otel v1.40.0                // OTEL API
    go.opentelemetry.io/otel/trace v1.40.0          // Tracing
    go.opentelemetry.io/otel/metric v1.40.0         // Metrics
    
    // OpenTelemetry SDK
    go.opentelemetry.io/otel/sdk v1.40.0            // OTEL SDK
    go.opentelemetry.io/otel/sdk/metric v1.40.0     // Metrics SDK
    go.opentelemetry.io/otel/sdk/log v0.16.0        // Logs SDK
    
    // OpenTelemetry Exporters (OTLP)
    go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.40.0
    go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.40.0
    go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.16.0
    
    // OpenTelemetry Instrumentation
    go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.65.0
    go.opentelemetry.io/contrib/bridges/otelslog v0.15.0
)
```

## 🚀 Instalação e Configuração

### 1. Clone o Repositório

```bash
git clone https://github.com/flcamillo/dynamodb-api.git
cd dynamodb-api
```

### 2. Configure Go e Dependências

```bash
# Verificar versão do Go
go version

# Download de dependências
go mod download

# Atualizar dependências
go mod tidy

# Verificar integridade
go mod verify
```

### 3. Configure o arquivo `config.json`

```json
{
  "address": "0.0.0.0",
  "port": 7000,
  "record_ttl_minutes": 1440
}
```

**Estrutura da Configuração:**
- `address` (string): Endereço para bind do servidor (0.0.0.0 = todos os interfaces)
- `port` (int): Porta do servidor HTTP (padrão: 7000)
- `record_ttl_minutes` (int64): Tempo de vida dos registros em minutos (padrão: 1440 = 24 horas)

### 4. Configurar Variáveis de Ambiente

#### Opção A: Ambiente Local com Dados Simulados
```bash
# Básico
export GO_ENV=development
export LOG_LEVEL=debug

# DynamoDB Local
export AWS_ENDPOINT_URL_DYNAMODB=http://localhost:8000
export AWS_REGION=local
export AWS_ACCESS_KEY_ID=local
export AWS_SECRET_ACCESS_KEY=local

# OpenTelemetry (sem Datadog)
export OTEL_SDK_DISABLED=false
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
export OTEL_RESOURCE_ATTRIBUTES=service.name=dynamodb-api,service.version=1.0.0,deployment.environment=dev
```

#### Opção B: Datadog Development
```bash
# AWS Configuration
export AWS_ENDPOINT_URL_DYNAMODB=http://localhost:8000
export AWS_REGION=local
export AWS_ACCESS_KEY_ID=local
export AWS_SECRET_ACCESS_KEY=local

# OpenTelemetry + Datadog
export OTEL_SDK_DISABLED=false
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
export OTEL_EXPORTER_OTLP_INSECURE=true
export OTEL_RESOURCE_ATTRIBUTES=service.name=dynamodb-api,service.version=1.0.0,deployment.environment=dev,team=backend

# Datadog Agent
export DD_SERVICE=dynamodb-api
export DD_ENV=dev
export DD_VERSION=1.0.0
export DD_TRACE_AGENT_URL=http://localhost:8126
export DD_AGENT_HOST=localhost
export DD_TRACE_AGENT_PORT=8126
export DD_PROFILING_ENABLED=true

# Datadog APM
export DD_TRACE_SAMPLE_RATE=1.0
export DD_METRICS_ENABLED=true
export DD_LOGS_INJECTION=true
```

#### Opção C: Produção AWS DynamoDB + Datadog
```bash
# AWS Production
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=${YOUR_ACCESS_KEY}
export AWS_SECRET_ACCESS_KEY=${YOUR_SECRET_KEY}
export AWS_ROLE_ARN=arn:aws:iam::ACCOUNT:role/DynamoDBRole

# OpenTelemetry Production
export OTEL_SDK_DISABLED=false
export OTEL_EXPORTER_OTLP_ENDPOINT=https://opentelemetry-backend.example.com:4317
export OTEL_EXPORTER_OTLP_INSECURE=false
export OTEL_EXPORTER_OTLP_TIMEOUT=30s
export OTEL_RESOURCE_ATTRIBUTES=service.name=dynamodb-api,service.version=1.0.0,deployment.environment=prod,team=backend,cloud.provider=aws,cloud.region=us-east-1

# Datadog Production
export DD_SERVICE=dynamodb-api
export DD_ENV=prod
export DD_VERSION=1.0.0
export DD_TRACE_AGENT_URL=http://datadog-agent.your-domain.com:8126
export DD_AGENT_HOST=datadog-agent.your-domain.com
export DD_TRACE_AGENT_PORT=8126
export DD_PROFILING_ENABLED=true
export DD_TRACE_SAMPLE_RATE=0.1
export DD_METRICS_ENABLED=true
```

#### Variáveis de Ambiente Detalhadas

```bash
# ============================================
# AWS Configuration
# ============================================
AWS_REGION=us-east-1                                    # Região AWS padrão
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE                 # Credencial AWS
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_ENDPOINT_URL_DYNAMODB=http://localhost:8000        # DynamoDB local (dev)
AWS_ROLE_ARN=arn:aws:iam::123456789012:role/MyRole     # IAM Role (prod)

# ============================================
# OpenTelemetry Core Configuration
# ============================================
OTEL_SDK_DISABLED=false                                 # Habilita OTEL SDK
OTEL_TRACES_EXPORTER=otlp                              # Exporter para traces
OTEL_METRICS_EXPORTER=otlp                             # Exporter para métricas
OTEL_LOGS_EXPORTER=otlp                                # Exporter para logs
OTEL_EXPORTER_OTLP_PROTOCOL=grpc                       # Protocolo (grpc/http)

# ============================================
# OpenTelemetry OTLP Configuration
# ============================================
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317      # OTEL Collector endpoint
OTEL_EXPORTER_OTLP_INSECURE=true                       # TLS disable (dev)
OTEL_EXPORTER_OTLP_TIMEOUT=30s                         # Timeout das exportações
OTEL_EXPORTER_OTLP_HEADERS=api-key=your-api-key        # Headers customizados

# Traces específicos
OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:4317/v1/traces
OTEL_EXPORTER_OTLP_TRACES_INSECURE=true
OTEL_EXPORTER_OTLP_TRACES_TIMEOUT=30s

# Métricas específicas
OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=http://localhost:4317/v1/metrics
OTEL_EXPORTER_OTLP_METRICS_INSECURE=true
OTEL_EXPORTER_OTLP_METRICS_TIMEOUT=30s

# Logs específicos
OTEL_EXPORTER_OTLP_LOGS_ENDPOINT=http://localhost:4317/v1/logs
OTEL_EXPORTER_OTLP_LOGS_INSECURE=true
OTEL_EXPORTER_OTLP_LOGS_TIMEOUT=30s

# ============================================
# OpenTelemetry Resource Attributes
# ============================================
OTEL_RESOURCE_ATTRIBUTES=service.name=dynamodb-api,service.version=1.0.0,deployment.environment=dev,team=backend,cloud.provider=aws,cloud.region=us-east-1

# Atributos individuais (formato expandido):
# OTEL_SERVICE_NAME=dynamodb-api
# OTEL_SERVICE_VERSION=1.0.0
# OTEL_DEPLOYMENT_ENVIRONMENT=dev
# OTEL_SERVICE_TEAM=backend
# OTEL_CLOUD_PROVIDER=aws
# OTEL_CLOUD_REGION=us-east-1
# OTEL_SERVICE_INSTANCE_ID=host-001

# ============================================
# OpenTelemetry Sampler Configuration
# ============================================
OTEL_TRACES_SAMPLER=parentbased_always_on              # Sampler strategy
OTEL_TRACES_SAMPLER_ARG=0.1                            # Amostragem 10%

# Opções de sampler:
# - always_on: Sempre traça
# - always_off: Nunca traça
# - traceidratio: Baseado em percentual
# - parentbased_always_on: Herda decisão do parent
# - parentbased_always_off: Herda decisão do parent

# ============================================
# OpenTelemetry Batch Processor Configuration
# ============================================
OTEL_BSP_SCHEDULE_DELAY=5000                           # Delay em ms (default: 5000)
OTEL_BSP_MAX_QUEUE_SIZE=2048                           # Tamanho máximo da fila
OTEL_BSP_MAX_EXPORT_BATCH_SIZE=512                     # Tamanho máximo do batch
OTEL_BSP_TIMEOUT=30000                                 # Timeout em ms

# ============================================
# Datadog Configuration
# ============================================
DD_SERVICE=dynamodb-api                                # Nome do serviço
DD_ENV=dev                                             # Ambiente (dev/staging/prod)
DD_VERSION=1.0.0                                       # Versão da aplicação
DD_AGENT_HOST=localhost                                # Host do Datadog Agent
DD_TRACE_AGENT_PORT=8126                               # Porta do Trace Agent
DD_TRACE_AGENT_URL=http://localhost:8126               # URL completa
DD_DOGSTATSD_PORT=8125                                 # Porta DogStatsD (métricas)

# ============================================
# Datadog Advanced Configuration
# ============================================
DD_TRACE_SAMPLE_RATE=1.0                               # Taxa de amostragem (1.0 = 100%)
DD_TRACE_DEBUG=false                                   # Debug mode
DD_TRACE_ENABLED=true                                  # Habilita tracing
DD_METRICS_ENABLED=true                                # Habilita métricas
DD_LOGS_INJECTION=true                                 # Injeta trace IDs em logs
DD_PROFILING_ENABLED=true                              # Habilita profiling
DD_PROFILING_SAMPLE_RATE=0.1                           # Taxa de amostragem de profiling
DD_APPSEC_ENABLED=false                                # Security scanning

# Datadog Proxy Configuration
DD_PROXY_HTTPS=http://proxy.example.com:8080
DD_PROXY_NO_PROXY=localhost,127.0.0.1

# ============================================
# Application Configuration
# ============================================
GO_ENV=development                                      # Go environment
LOG_LEVEL=debug                                        # Log level (debug/info/warn/error)
PORT=7000                                              # Application port
BIND_ADDRESS=0.0.0.0                                   # Bind address

# ============================================
# Performance Tuning
# ============================================
GOMAXPROCS=4                                           # Max parallel processors
GOMEMLIMIT=256MiB                                      # Memory limit (Go 1.19+)
GODEBUG=gctrace=0                                      # GC trace output
```

## ▶️ Executando a Aplicação

### Opção 1: Desenvolvimento Local (HTTP Server)

```bash
# Verificar dependências
go mod download

# Executar aplicação
go run main.go

# Output esperado:
# Setup OTel SDK successfully
# Listening on 0.0.0.0:7000
```

### Opção 2: Com Docker Compose (DynamoDB Local + OTEL Collector)

```bash
# Iniciar stack de desenvolvimento
docker-compose -f extra/docker-compose.yml up -d

# Verificar containers
docker-compose -f extra/docker-compose.yml ps

# Executar aplicação
go run main.go

# Acessar serviços
# - API: http://localhost:7000
# - Datadog: http://localhost:8000 (mock)
# - OTEL Collector: http://localhost:13133
# - Prometheus: http://localhost:9090
# - Jaeger: http://localhost:16686
```

### Opção 3: Build e Executar Binário

```bash
# Build otimizado
go build -o api -ldflags="-s -w" .

# Executar binário
./api

# Build com informações de debug
go build -o api .
./api -version  # se implementado
```

### Opção 4: AWS Lambda

```bash
# Build para Lambda (ARM64)
GOOS=linux GOARCH=arm64 go build -o bootstrap .

# Comprimir
zip lambda.zip bootstrap

# Deploy via AWS CLI
aws lambda create-function \
  --function-name dynamodb-api \
  --runtime provided.al2 \
  --role arn:aws:iam::ACCOUNT:role/lambda-role \
  --handler bootstrap \
  --zip-file fileb://lambda.zip \
  --architectures arm64 \
  --timeout 30 \
  --memory-size 512 \
  --environment Variables="{
    AWS_REGION=us-east-1,
    OTEL_SDK_DISABLED=false,
    DD_SERVICE=dynamodb-api,
    DD_ENV=prod
  }"

# Atualizar função existente
aws lambda update-function-code \
  --function-name dynamodb-api \
  --zip-file fileb://lambda.zip

# Invocar função
aws lambda invoke \
  --function-name dynamodb-api \
  --payload '{"httpMethod":"GET","path":"/health"}' \
  response.json
```

## 📡 Endpoints da API

### 1. Health Check

```http
GET /health
```

Verifica se a aplicação está ativa e pronta para receber requisições.

**Resposta:** `200 OK`

---

### 2. Criar Evento

```http
POST /eventos
Content-Type: application/json
```

Cria um novo evento com UUID único gerado automaticamente.

**Request Body:**
```json
{
  "date": "2024-01-29T10:30:00Z",
  "statusCode": 200,
  "statusMessage": "Operação bem-sucedida",
  "metadata": {
    "user_id": "123",
    "request_id": "abc-def-ghi",
    "correlation_id": "xyz-789",
    "version": "1.0"
  }
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "date": "2024-01-29T10:30:00Z",
  "statusCode": 200,
  "statusMessage": "Operação bem-sucedida",
  "expiration": 1706633400,
  "metadata": {
    "user_id": "123",
    "request_id": "abc-def-ghi",
    "correlation_id": "xyz-789",
    "version": "1.0"
  }
}
```

**Métricas Registradas:**
- `post.requests` +1
- `http.request.duration_ms` (duração)
- `dynamodb.latency_ms` (duração da operação DB)

---

### 3. Obter Evento

```http
GET /eventos/{id}
```

Recupera um evento específico pelo UUID.

**Path Parameters:**
- `id` (string, required): UUID do evento (formato: 550e8400-e29b-41d4-a716-446655440000)

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "date": "2024-01-29T10:30:00Z",
  "statusCode": 200,
  "statusMessage": "Operação bem-sucedida",
  "expiration": 1706633400,
  "metadata": {
    "user_id": "123",
    "request_id": "abc-def-ghi"
  }
}
```

**Response (404 Not Found):**
```json
{
  "error": "evento não encontrado",
  "requestId": "req-123"
}
```

---

### 4. Atualizar Evento

```http
PUT /eventos/{id}
Content-Type: application/json
```

Atualiza um evento existente.

**Path Parameters:**
- `id` (string, required): UUID do evento

**Request Body:**
```json
{
  "date": "2024-01-29T11:00:00Z",
  "statusCode": 201,
  "statusMessage": "Criado com sucesso",
  "metadata": {
    "user_id": "456",
    "updated_by": "admin"
  }
}
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "date": "2024-01-29T11:00:00Z",
  "statusCode": 201,
  "statusMessage": "Criado com sucesso",
  "expiration": 1706636800,
  "metadata": {
    "user_id": "456",
    "updated_by": "admin"
  }
}
```

---

### 5. Deletar Evento

```http
DELETE /eventos/{id}
```

Remove um evento específico permanentemente.

**Path Parameters:**
- `id` (string, required): UUID do evento

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "date": "2024-01-29T11:00:00Z",
  "statusCode": 201,
  "statusMessage": "Criado com sucesso",
  "expiration": 1706636800,
  "metadata": {
    "user_id": "456"
  }
}
```

---

### 6. Listar Eventos (Find)

```http
GET /eventos?startDate=2024-01-29T00:00:00Z&endDate=2024-01-30T00:00:00Z&statusCode=200
```

Lista eventos filtrando por período e código de status.

**Query Parameters:**
- `startDate` (string, required): Data inicial (RFC3339 format)
- `endDate` (string, required): Data final (RFC3339 format)
- `statusCode` (integer, required): Código HTTP para filtrar (0-599)

**Response (200 OK):**
```json
{
  "items": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "date": "2024-01-29T10:30:00Z",
      "statusCode": 200,
      "statusMessage": "OK",
      "expiration": 1706633400,
      "metadata": {}
    },
    {
      "id": "6f0fa6d4-f3c2-4a8e-9b12-5c8e9f2a1b3d",
      "date": "2024-01-29T14:45:00Z",
      "statusCode": 200,
      "statusMessage": "OK",
      "expiration": 1706648700,
      "metadata": {}
    }
  ],
  "total": 2
}
```

---

## 📝 Exemplos com cURL

### Pré-requisitos
```bash
# Certificar que a API está rodando
curl -s http://localhost:7000/health

# Salvar URL base em variável
API="http://localhost:7000"
```

### Health Check
```bash
curl -X GET $API/health

# Saída esperada:
# OK
```

### Criar Evento
```bash
curl -X POST $API/eventos \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-29T10:30:00Z",
    "statusCode": 200,
    "statusMessage": "Sucesso",
    "metadata": {
      "user_id": "user-001",
      "source": "mobile_app"
    }
  }' | jq '.'

# Salvar o ID para os próximos exemplos
EVENT_ID=$(curl -s -X POST $API/eventos \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-29T10:30:00Z",
    "statusCode": 200,
    "statusMessage": "Sucesso",
    "metadata": {"user_id": "user-001"}
  }' | jq -r '.id')

echo "Event ID: $EVENT_ID"
```

### Obter Evento
```bash
curl -X GET $API/eventos/$EVENT_ID | jq '.'
```

### Listar Eventos
```bash
curl -X GET "$API/eventos?startDate=2024-01-28T00:00:00Z&endDate=2024-01-30T23:59:59Z&statusCode=200" | jq '.'
```

### Atualizar Evento
```bash
curl -X PUT $API/eventos/$EVENT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-29T11:45:00Z",
    "statusCode": 201,
    "statusMessage": "Atualizado",
    "metadata": {
      "updated_by": "admin",
      "version": "2.0"
    }
  }' | jq '.'
```

### Deletar Evento
```bash
curl -X DELETE $API/eventos/$EVENT_ID | jq '.'
```

---

## 📂 Estrutura do Projeto

```
dynamodb-api/
│
├── 📋 Arquivos de Configuração
│   ├── main.go                    # Ponto de entrada
│   ├── config.go                  # Gerenciamento de configuração
│   ├── config.json                # Arquivo de configuração
│   ├── otel.go                    # Setup OpenTelemetry
│   ├── go.mod                     # Definição de módulo e dependências
│   ├── go.sum                     # Checksum das dependências
│   └── README.md                  # Documentação
│
├── 📦 models/                     # Modelos de dados
│   ├── event.go                   # Estrutura do evento
│   ├── event_test.go              # Testes unitários
│   ├── error_response.go          # Modelo de erro
│   ├── error_response_test.go     # Testes
│   ├── paginated_response.go      # Resposta paginada
│   └── paginated_response_test.go # Testes
│
├── 🔧 handlers/                   # Handlers de requisição
│   ├── http_handler.go            # Implementação HTTP
│   ├── http_handler_test.go       # Testes
│   ├── lambda_handler.go          # Implementação Lambda
│   ├── lambda_handler_test.go     # Testes
│   └── handlers.go                # Tipos compartilhados
│
├── 💾 repositories/               # Implementações de armazenamento
│   ├── dynamodb.go                # DynamoDB client
│   ├── dynamodb_test.go           # Testes
│   ├── memorydb.go                # In-memory storage
│   ├── memorydb_test.go           # Testes
│   └── repository.go              # Tipos compartilhados
│
├── 🌉 interfaces/                 # Definição de interfaces
│   ├── dynamodb_client.go         # Interface DynamoDB
│   ├── log.go                     # Interface Logger
│   └── repository.go              # Interface Repository
│
├── 📝 logs/                       # Sistema de logging
│   ├── stdout.go                  # Logger padrão
│   └── stdout_test.go             # Testes
│
├── 🌐 apis/                       # Implementação de APIs
│   ├── http_api.go                # API HTTP Server
│   ├── http_api_test.go           # Testes
│   ├── lambda_api.go              # Lambda wrapper
│   ├── lambda_api_test.go         # Testes
│   └── apis.go                    # Tipos compartilhados
│
└── 🐳 extra/                      # Recursos adicionais
    ├── docker-compose.yml         # Orquestração local
    ├── otel-collector.yaml        # Config OTEL Collector
    ├── prometheus.yaml            # Config Prometheus
    ├── datadog.txt                # Instruções Datadog
    ├── datadog_dashboard.png      # Screenshot Datadog
    ├── datadog_logs.png           # Screenshot logs
    ├── datadog_apm*.png           # Screenshots APM
    └── datadog_custom_metrics.png # Screenshots métricas
```

---

## ⚙️ Configuração Avançada

### Arquivo `config.json` Detalhado

```json
{
  "address": "0.0.0.0",
  "port": 7000,
  "record_ttl_minutes": 1440
}
```

**Parâmetros da Configuração:**

| Parâmetro | Tipo | Padrão | Min | Max | Descrição |
|-----------|------|--------|-----|-----|-----------|
| `address` | string | `0.0.0.0` | - | - | Endereço para bind (0.0.0.0 = todos os interfaces) |
| `port` | int | `7000` | 1 | 65535 | Porta do servidor HTTP |
| `record_ttl_minutes` | int64 | `1440` | 1 | 525600 | TTL em minutos (máx: 1 ano) |

---

## 📊 Telemetria e Observabilidade

### OpenTelemetry - Visão Geral

A aplicação implementa observabilidade completa através do OpenTelemetry (OTEL), um padrão aberto para coleta de telemetria.

```mermaid
graph LR
    APP["🔧 Aplicação<br/>Traces | Metrics<br/>Logs"]
    
    OTEL["📡 OpenTelemetry SDK<br/>Batch Processor<br/>Sampler"]
    
    EXPORTER["🔄 OTLP Exporter<br/>gRPC Protocol<br/>Batching"]
    
    COLLECTOR["📡 OTEL Collector<br/>Receive | Process<br/>Export"]
    
    BACKENDS["🎯 Backends<br/>Datadog | Jaeger<br/>Prometheus"]
    
    APP -->|Instrumentation| OTEL
    OTEL -->|Export| EXPORTER
    EXPORTER -->|:4317| COLLECTOR
    COLLECTOR -->|Transform| BACKENDS
```

### Métricas Coletadas

```mermaid
graph TB
    subgraph "Requisições HTTP"
        M1["post.requests - Contador POST"]
        M2["get.requests - Contador GET"]
        M3["put.requests - Contador PUT"]
        M4["delete.requests - Contador DELETE"]
        M5["find.requests - Contador FIND"]
    end
    
    subgraph "Performance"
        P1["http.request.duration_ms - Histograma"]
        P2["dynamodb.latency_ms - Histograma"]
        P3["repository.operation.duration_ms - Histograma"]
    end
    
    subgraph "Erros"
        E1["http.errors.total - Contador"]
        E2["repository.errors.total - Contador"]
        E3["validation.errors.total - Contador"]
    end
    
    subgraph "Dados"
        D1["dynamodb.items.count - Gauge"]
        D2["events.created.total - Contador"]
        D3["events.deleted.total - Contador"]
    end
```

### Traces Distribuídos

Cada operação gera um trace completo com:
- **Span Root**: POST /eventos, GET /eventos, etc
- **Spans Filhos**: Validação, salvamento DB, geração ID
- **Atributos**: event_id, status_code, duration_ms
- **Eventos**: validation_start, db_save_start, db_save_success
- **Erros**: exception type, stack trace

### Logs Estruturados

Formato JSON estruturado com contexto completo:

```json
{
  "timestamp": "2025-02-08T10:30:45.123456Z",
  "level": "INFO",
  "logger": "dynamodb-api",
  "message": "evento criado com sucesso",
  
  "attributes": {
    "event_id": "550e8400-e29b-41d4-a716-446655440000",
    "status_code": 201,
    "duration_ms": 125,
    "repository": "dynamodb",
    "operation": "save"
  },
  
  "trace": {
    "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
    "span_id": "00f067aa0ba902b7",
    "trace_flags": "01"
  },
  
  "resource": {
    "service.name": "dynamodb-api",
    "service.version": "1.0.0",
    "deployment.environment": "dev"
  }
}
```

---

## 🐶 Datadog Integration

### Visão Geral

Integração com Datadog para APM, logs e métricas em tempo real.

```mermaid
graph TB
    APP["🔧 DynamoDB API<br/>OTEL SDK"]
    
    OTEL_EXPORTER["🔄 OTLP Exporter<br/>gRPC :4317"]
    
    OTEL_COLLECTOR["📡 OTEL Collector<br/>localhost:4317"]
    
    DD_EXPORTER["🐶 Datadog Exporter<br/>:8125 | :8126"]
    
    DD_AGENT["🐶 Datadog Agent<br/>Container | Local"]
    
    DD_BACKEND["☁️ Datadog Cloud<br/>APM | Logs<br/>Metrics | Events"]
    
    APP --> OTEL_EXPORTER
    OTEL_EXPORTER --> OTEL_COLLECTOR
    OTEL_COLLECTOR --> DD_EXPORTER
    DD_EXPORTER --> DD_AGENT
    DD_AGENT --> DD_BACKEND
```

### Configuração com Docker Compose

```yaml
version: '3.8'

services:
  # DynamoDB Local
  dynamodb-local:
    image: amazon/dynamodb-local:latest
    ports:
      - "8000:8000"
    environment:
      - AWS_ACCESS_KEY_ID=local
      - AWS_SECRET_ACCESS_KEY=local

  # OTEL Collector
  otel-collector:
    image: otel/opentelemetry-collector-contrib:latest
    ports:
      - "4317:4317"    # OTLP gRPC receiver
      - "4318:4318"    # OTLP HTTP receiver
      - "13133:13133"  # Health check
    volumes:
      - ./extra/otel-collector.yaml:/etc/otel-collector-config.yaml
    command:
      - "--config=/etc/otel-collector-config.yaml"

  # Datadog Agent
  datadog-agent:
    image: datadog/agent:latest
    environment:
      - DD_API_KEY=${DD_API_KEY}
      - DD_SITE=datadoghq.com
      - DD_APM_ENABLED=true
      - DD_LOGS_ENABLED=true
    ports:
      - "8126:8126"  # APM
      - "8125:8125"  # DogStatsD
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

  # Prometheus (opcional)
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./extra/prometheus.yaml:/etc/prometheus/prometheus.yml

  # Jaeger (opcional)
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"
```

### Dashboards e Visualizações

#### 1. Dashboard Principal

![Datadog Dashboard](./extra/datadog_dashboard.png)

Mostra overview de:
- Taxa de requisições por segundo
- Latência p50, p95, p99
- Taxa de erro
- Distribuição por método HTTP

#### 2. APM - Traces Distribuídos

![Datadog APM 1](./extra/datadog_apm1.png)
![Datadog APM 2](./extra/datadog_apm2.png)
![Datadog APM 3](./extra/datadog_apm3.png)
![Datadog APM 4](./extra/datadog_apm4.png)
![Datadog APM 5](./extra/datadog_apm5.png)
![Datadog APM 6](./extra/datadog_apm6.png)
![Datadog APM 7](./extra/datadog_apm7.png)
![Datadog APM 8](./extra/datadog_apm8.png)
![Datadog APM 9](./extra/datadog_apm9.png)
![Datadog APM 10](./extra/datadog_apm10.png)
![Datadog APM 11](./extra/datadog_apm11.png)
![Datadog APM 12](./extra/datadog_apm12.png)

Análise detalhada de:
- Flame graphs de traces
- Dependências de serviços
- Latência por operação
- Taxa de erro por endpoint

#### 3. Logs Estruturados

![Datadog Logs](./extra/datadog_logs.png)

Visualização de:
- Logs em tempo real
- Filtros por trace_id
- Correlação com traces
- Níveis de severity

#### 4. Métricas Customizadas

![Datadog Custom Metrics](./extra/datadog_custom_metrics.png)

Gráficos de:
- Taxa de criação/deleção de eventos
- Latência DynamoDB
- Distribuição de status codes
- Trends de performance

### Queries Úteis no Datadog

```datadog
# Taxa de requisições por segundo
avg:post.requests{service:dynamodb-api}.as_count()

# Latência p99 por endpoint
p99:http.request.duration_ms{service:dynamodb-api} by {endpoint}

# Taxa de erro
sum:http.errors.total{service:dynamodb-api}.as_count()

# Eventos criados por hora
sum:events.created.total{service:dynamodb-api}.as_count()

# Latência DynamoDB
avg:dynamodb.latency_ms{service:dynamodb-api} by {operation}

# Distribuição de status codes
sum:http.response.status{service:dynamodb-api} by {status_code}
```

---

## 🧪 Testes

### Executar Testes

```bash
# Todos os testes
go test ./...

# Com verbose
go test ./... -v

# Com cobertura
go test ./... -cover

# Cobertura detalhada
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Testes específicos
go test ./handlers -v
go test ./repositories -v -run TestDynamoDB
go test ./models -v -run TestEventValidation
```

### Testes de Carga

```bash
# Usando hey
hey -n 10000 -c 100 -m POST -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-29T10:30:00Z",
    "statusCode": 200,
    "statusMessage": "Test",
    "metadata": {}
  }' \
  http://localhost:7000/eventos

# Usando Apache Bench
ab -n 10000 -c 100 http://localhost:7000/health

# Usando k6 (Load testing)
k6 run load-test.js
```

---

## 🐛 Troubleshooting

### Problema: API não inicia

```bash
# Windows - Verificar porta em uso
netstat -ano | findstr :7000

# Linux/Mac
lsof -i :7000

# Solução: Mudar porta em config.json ou liberar a porta
# Windows
taskkill /PID <PID> /F

# Linux/Mac
kill -9 <PID>
```

### Problema: DynamoDB não conecta

```bash
# Verificar credenciais AWS
aws sts get-caller-identity

# Para DynamoDB local
docker-compose -f extra/docker-compose.yml up -d dynamodb-local

# Verificar conexão
curl http://localhost:8000

# Testar com AWS CLI
aws dynamodb list-tables --endpoint-url http://localhost:8000
```

### Problema: OTEL não funciona

```bash
# Verificar OTEL Collector
curl http://localhost:13133

# Verificar variáveis de ambiente
echo $OTEL_EXPORTER_OTLP_ENDPOINT
echo $OTEL_SDK_DISABLED

# Verificar logs
docker logs otel-collector

# Debug mode
export OTEL_LOG_LEVEL=debug
go run main.go
```

### Problema: Datadog não recebe dados

```bash
# Verificar Datadog Agent
docker logs datadog-agent

# Verificar conexão
curl http://localhost:8126/trace/validate

# Verificar firewall
nc -zv localhost 8126
nc -zv localhost 8125
```

---

## 📋 Checklist de Deploy

### Desenvolvimento
- [ ] Go 1.25.5+ instalado
- [ ] Docker e Docker Compose instalado
- [ ] DynamoDB local rodando
- [ ] OTEL Collector rodando
- [ ] config.json configurado
- [ ] Variáveis de ambiente (dev) definidas
- [ ] Testes unitários passando
- [ ] Cobertura de código validada

### Staging
- [ ] Build otimizado compilado
- [ ] Credenciais AWS staging configuradas
- [ ] Tabela DynamoDB staging criada
- [ ] OTEL Collector apontando para staging
- [ ] Datadog agent staging configurado
- [ ] Testes de carga executados
- [ ] Logs sendo capturados corretamente

### Produção
- [ ] Arquivo `config.json` configurado
- [ ] Credenciais AWS production seguras
- [ ] Tabela DynamoDB production com backups
- [ ] TTL configurado corretamente
- [ ] OpenTelemetry production endpoint
- [ ] Datadog production configurado
- [ ] Monitoramento e alertas ativados
- [ ] Disaster recovery plan testado
- [ ] Security scanning executado
- [ ] Performance benchmarking realizado

---

## 📚 Referências

### Documentação Oficial
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [Go HTTP Package](https://pkg.go.dev/net/http)
- [AWS Lambda for Go](https://github.com/aws/aws-lambda-go)

### Ferramentas e Serviços
- [Datadog APM](https://docs.datadoghq.com/tracing/)
- [Jaeger Tracing](https://www.jaegertracing.io/)
- [Prometheus](https://prometheus.io/)
- [OTEL Collector](https://opentelemetry.io/docs/collector/)

### Best Practices
- [Observability Engineering](https://www.oreilly.com/library/view/observability-engineering/9781492076438/)
- [SRE Best Practices](https://sre.google/sre-book/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

---

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

---

## 📝 Licença

Este projeto é fornecido como-está para fins educacionais e de demonstração.

---

## 📞 Suporte

Para dúvidas ou problemas:
1. Consulte os testes unitários em `*_test.go`
2. Revise os exemplos de cURL neste documento
3. Verifique os comentários no código-fonte
4. Analise os logs via OTEL/Datadog

---

**Última atualização:** Fevereiro 2026
**Versão da Documentação:** 2.0.0
**Status:** ✅ Produção-Ready

