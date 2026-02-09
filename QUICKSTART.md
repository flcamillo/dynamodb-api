# Quick Start Guide - DynamoDB API

Um guia rápido para começar a desenvolver com a DynamoDB API.

## 5 Minutos para Começar

### 1. Clone e Prepare

```bash
git clone https://github.com/flcamillo/dynamodb-api.git
cd dynamodb-api
go mod download
```

### 2. Inicie o Docker Compose

```bash
docker-compose -f extra/docker-compose.yml up -d
```

Isso inicia:
- ✅ DynamoDB Local (porta 8000)
- ✅ OTEL Collector (porta 4317)
- ✅ Prometheus (porta 9090)
- ✅ Jaeger (porta 16686)

### 3. Configure Variáveis de Ambiente

```bash
# DynamoDB Local
export AWS_ENDPOINT_URL_DYNAMODB=http://localhost:8000
export AWS_REGION=local
export AWS_ACCESS_KEY_ID=local
export AWS_SECRET_ACCESS_KEY=local

# OpenTelemetry
export OTEL_SDK_DISABLED=false
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
export OTEL_RESOURCE_ATTRIBUTES=service.name=dynamodb-api,service.version=1.0.0,deployment.environment=dev
```

### 4. Execute a Aplicação

```bash
go run main.go

# Saída esperada:
# Setup OTel SDK successfully
# Listening on 0.0.0.0:7000
```

### 5. Teste a API

```bash
# Health Check
curl http://localhost:7000/health

# Criar Evento
curl -X POST http://localhost:7000/eventos \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-29T10:30:00Z",
    "statusCode": 200,
    "statusMessage": "Sucesso",
    "metadata": {"user_id": "123"}
  }' | jq '.'
```

## Acessar Dashboards

- **Prometheus**: http://localhost:9090
- **Jaeger**: http://localhost:16686
- **API**: http://localhost:7000

## Integração com Datadog (Opcional)

### 1. Configure Datadog Agent

```bash
# Instale o Datadog Agent
# https://docs.datadoghq.com/agent/basic_agent_usage/

# Configure variáveis
export DD_API_KEY=your_api_key_here
export DD_SERVICE=dynamodb-api
export DD_ENV=dev
export DD_AGENT_HOST=localhost
export DD_TRACE_AGENT_PORT=8126
```

### 2. Execute com Datadog

```bash
docker-compose -f extra/docker-compose.yml up -d datadog-agent
go run main.go
```

### 3. Veja os dados no Datadog

- APM Traces: https://app.datadoghq.com/apm/traces
- Logs: https://app.datadoghq.com/logs
- Metrics: https://app.datadoghq.com/metric

## Comandos Úteis

### Executar Testes
```bash
go test ./...           # Todos os testes
go test ./... -v        # Com detalhes
go test ./... -cover    # Com cobertura
```

### Build
```bash
go build -o api .       # Build completo
go build -ldflags="-s -w" -o api .  # Build otimizado
```

### Limpar
```bash
# Remover containers
docker-compose -f extra/docker-compose.yml down

# Remover volumes (limpar dados)
docker-compose -f extra/docker-compose.yml down -v

# Remover cache Go
go clean -cache
go clean -modcache
```

## Troubleshooting Rápido

| Problema | Solução |
|----------|---------|
| Porta 7000 em uso | `lsof -i :7000` depois `kill -9 <PID>` |
| DynamoDB não conecta | `docker-compose up -d dynamodb-local` |
| OTEL não funciona | Verificar: `echo $OTEL_SDK_DISABLED` |
| Sem dados no Datadog | Verificar Datadog Agent: `docker logs datadog-agent` |

## Próximos Passos

1. Leia a [Documentação Detalhada](README_DETALHADO.md)
2. Explore os [Exemplos de cURL](README_DETALHADO.md#exemplos-com-curl)
3. Verifique a [Arquitetura](README_DETALHADO.md#arquitetura-detalhada)
4. Configure [Datadog Integration](README_DETALHADO.md#datadog-integration)

---

**Pronto para começar?** Execute os 5 passos acima em ~2 minutos! 🚀

