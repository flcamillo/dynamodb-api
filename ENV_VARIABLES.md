# Environment Variables Configuration

Este arquivo documenta todas as variáveis de ambiente suportadas pela aplicação.

## Arquivo .env para Desenvolvimento

Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo:

```bash
# ============================================
# AWS Configuration
# ============================================
AWS_REGION=local
AWS_ENDPOINT_URL_DYNAMODB=http://localhost:8000
AWS_ACCESS_KEY_ID=local
AWS_SECRET_ACCESS_KEY=local

# ============================================
# OpenTelemetry Configuration
# ============================================
OTEL_SDK_DISABLED=false
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
OTEL_EXPORTER_OTLP_INSECURE=true
OTEL_TRACES_EXPORTER=otlp
OTEL_METRICS_EXPORTER=otlp
OTEL_LOGS_EXPORTER=otlp

# ============================================
# OpenTelemetry Resource Attributes
# ============================================
OTEL_RESOURCE_ATTRIBUTES=service.name=dynamodb-api,service.version=1.0.0,deployment.environment=dev,team=backend

# ============================================
# Datadog Configuration (Opcional)
# ============================================
DD_SERVICE=dynamodb-api
DD_ENV=dev
DD_VERSION=1.0.0
DD_AGENT_HOST=localhost
DD_TRACE_AGENT_PORT=8126
DD_TRACE_SAMPLE_RATE=1.0
DD_METRICS_ENABLED=true
DD_LOGS_INJECTION=true

# ============================================
# Application Configuration
# ============================================
GO_ENV=development
LOG_LEVEL=debug
PORT=7000
BIND_ADDRESS=0.0.0.0
```

## Carregar .env com direnv

Se usar [direnv](https://direnv.net/):

```bash
# Instalar direnv
brew install direnv  # macOS
apt-get install direnv  # Linux

# Criar .envrc
echo 'export $(cat .env | xargs)' > .envrc

# Autorizar
direnv allow

# Automático em cada cd para o diretório!
```

## Carregar .env Manualmente

### Linux/macOS
```bash
set -a
source .env
set +a

# Verificar
env | grep AWS_
env | grep OTEL_
env | grep DD_
```

### Windows PowerShell
```powershell
# Ler arquivo .env e exportar
Get-Content .env | ForEach-Object {
    if ($_ -match '^\s*[^#]') {
        $name, $value = $_.Split('=')
        [Environment]::SetEnvironmentVariable($name, $value)
    }
}

# Verificar
$env:AWS_REGION
$env:OTEL_SDK_DISABLED
```

### Windows CMD
```batch
setx AWS_REGION local
setx AWS_ENDPOINT_URL_DYNAMODB http://localhost:8000
setx OTEL_SDK_DISABLED false
setx OTEL_EXPORTER_OTLP_ENDPOINT http://localhost:4317
```

## Presets de Configuração

### Preset 1: Desenvolvimento Local
```bash
export AWS_REGION=local
export AWS_ENDPOINT_URL_DYNAMODB=http://localhost:8000
export AWS_ACCESS_KEY_ID=local
export AWS_SECRET_ACCESS_KEY=local
export OTEL_SDK_DISABLED=false
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
export OTEL_RESOURCE_ATTRIBUTES=service.name=dynamodb-api,service.version=1.0.0,deployment.environment=dev
```

### Preset 2: Development com Datadog
```bash
export AWS_REGION=local
export AWS_ENDPOINT_URL_DYNAMODB=http://localhost:8000
export AWS_ACCESS_KEY_ID=local
export AWS_SECRET_ACCESS_KEY=local

export OTEL_SDK_DISABLED=false
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
export OTEL_RESOURCE_ATTRIBUTES=service.name=dynamodb-api,service.version=1.0.0,deployment.environment=dev

export DD_SERVICE=dynamodb-api
export DD_ENV=dev
export DD_VERSION=1.0.0
export DD_AGENT_HOST=localhost
export DD_TRACE_AGENT_PORT=8126
```

### Preset 3: Staging
```bash
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=${STAGING_AWS_KEY}
export AWS_SECRET_ACCESS_KEY=${STAGING_AWS_SECRET}

export OTEL_SDK_DISABLED=false
export OTEL_EXPORTER_OTLP_ENDPOINT=https://otel-staging.example.com:4317
export OTEL_EXPORTER_OTLP_INSECURE=false
export OTEL_RESOURCE_ATTRIBUTES=service.name=dynamodb-api,service.version=1.0.0,deployment.environment=staging

export DD_SERVICE=dynamodb-api
export DD_ENV=staging
export DD_VERSION=1.0.0
export DD_AGENT_HOST=datadog-agent-staging.example.com
export DD_TRACE_SAMPLE_RATE=0.5
```

### Preset 4: Produção
```bash
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=${PROD_AWS_KEY}
export AWS_SECRET_ACCESS_KEY=${PROD_AWS_SECRET}

export OTEL_SDK_DISABLED=false
export OTEL_EXPORTER_OTLP_ENDPOINT=https://otel-prod.example.com:4317
export OTEL_EXPORTER_OTLP_INSECURE=false
export OTEL_RESOURCE_ATTRIBUTES=service.name=dynamodb-api,service.version=1.0.0,deployment.environment=prod,cloud.provider=aws,cloud.region=us-east-1

export DD_SERVICE=dynamodb-api
export DD_ENV=prod
export DD_VERSION=1.0.0
export DD_AGENT_HOST=datadog-agent-prod.example.com
export DD_TRACE_SAMPLE_RATE=0.1
export DD_PROFILING_ENABLED=true
```

## Scripts de Setup

### setup-dev.sh
```bash
#!/bin/bash
set -e

echo "🚀 Configurando ambiente de desenvolvimento..."

# Carregar variáveis
set -a
source .env
set +a

# Verificar dependências
echo "✓ Verificando Go..."
go version

echo "✓ Verificando Docker..."
docker --version

echo "✓ Verificando Docker Compose..."
docker-compose --version

# Download de dependências
echo "⬇️  Baixando dependências Go..."
go mod download
go mod tidy

# Iniciar Docker Compose
echo "🐳 Iniciando Docker Compose..."
docker-compose -f extra/docker-compose.yml up -d

# Aguardar serviços
echo "⏳ Aguardando serviços..."
sleep 5

echo "✅ Ambiente de desenvolvimento pronto!"
echo ""
echo "Próximos passos:"
echo "  1. Execute: go run main.go"
echo "  2. Teste: curl http://localhost:7000/health"
echo "  3. Veja: http://localhost:16686 (Jaeger)"
```

### setup-datadog.sh
```bash
#!/bin/bash
set -e

echo "🐶 Configurando Datadog..."

# Pedir API Key
read -p "Digite sua Datadog API Key: " DD_API_KEY

# Exportar
export DD_API_KEY=$DD_API_KEY
export DD_SERVICE=dynamodb-api
export DD_ENV=dev

# Iniciar Datadog Agent
docker-compose -f extra/docker-compose.yml up -d datadog-agent

echo "✅ Datadog Agent iniciado!"
echo ""
echo "Acesse: https://app.datadoghq.com"
```

### setup-prod.sh
```bash
#!/bin/bash
set -e

echo "⚠️  Configurando ambiente de PRODUÇÃO..."

# Validações
if [ -z "$PROD_AWS_KEY" ]; then
    echo "❌ Erro: PROD_AWS_KEY não definida"
    exit 1
fi

if [ -z "$DD_API_KEY" ]; then
    echo "❌ Erro: DD_API_KEY não definida"
    exit 1
fi

echo "✓ Verificações de segurança passadas"

# Build otimizado
echo "🔨 Building aplicação..."
go build -ldflags="-s -w" -o api .

# Testes
echo "🧪 Executando testes..."
go test ./...

echo "✅ Pronto para deploy em produção!"
```

## Verificação de Configuração

Execute este script para verificar se tudo está configurado corretamente:

```bash
#!/bin/bash

echo "🔍 Verificando configuração..."
echo ""

# AWS
echo "AWS Configuration:"
echo "  AWS_REGION: ${AWS_REGION:-❌ NOT SET}"
echo "  AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID:+✓ SET}"
echo "  AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY:+✓ SET}"
echo ""

# OpenTelemetry
echo "OpenTelemetry Configuration:"
echo "  OTEL_SDK_DISABLED: ${OTEL_SDK_DISABLED:-❌ NOT SET}"
echo "  OTEL_EXPORTER_OTLP_ENDPOINT: ${OTEL_EXPORTER_OTLP_ENDPOINT:-❌ NOT SET}"
echo "  OTEL_RESOURCE_ATTRIBUTES: ${OTEL_RESOURCE_ATTRIBUTES:+✓ SET}"
echo ""

# Datadog
echo "Datadog Configuration:"
echo "  DD_SERVICE: ${DD_SERVICE:-❌ NOT SET}"
echo "  DD_ENV: ${DD_ENV:-❌ NOT SET}"
echo "  DD_AGENT_HOST: ${DD_AGENT_HOST:-❌ NOT SET}"
echo ""

# Testes de Conectividade
echo "Connectivity Tests:"

# DynamoDB
if [ "$AWS_ENDPOINT_URL_DYNAMODB" ]; then
    if curl -s "${AWS_ENDPOINT_URL_DYNAMODB}" > /dev/null 2>&1; then
        echo "  DynamoDB: ✓ OK"
    else
        echo "  DynamoDB: ❌ UNREACHABLE"
    fi
fi

# OTEL Collector
if [ "$OTEL_EXPORTER_OTLP_ENDPOINT" ]; then
    if curl -s "${OTEL_EXPORTER_OTLP_ENDPOINT}" > /dev/null 2>&1; then
        echo "  OTEL Collector: ✓ OK"
    else
        echo "  OTEL Collector: ❌ UNREACHABLE"
    fi
fi

# Datadog Agent
if [ "$DD_AGENT_HOST" ]; then
    if curl -s "http://${DD_AGENT_HOST}:${DD_TRACE_AGENT_PORT}/trace/validate" > /dev/null 2>&1; then
        echo "  Datadog Agent: ✓ OK"
    else
        echo "  Datadog Agent: ❌ UNREACHABLE"
    fi
fi

echo ""
echo "✅ Verificação concluída!"
```

## Secrets Seguros (Production)

Para produção, use um gerenciador de secrets:

### AWS Secrets Manager
```bash
# Store
aws secretsmanager create-secret \
  --name dynamodb-api/prod \
  --secret-string '{
    "aws_access_key": "...",
    "aws_secret_key": "...",
    "dd_api_key": "..."
  }'

# Retrieve no Go
import "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
```

### HashiCorp Vault
```bash
# Store
vault kv put secret/dynamodb-api \
  aws_key=xxx \
  aws_secret=yyy \
  dd_api_key=zzz

# Retrieve no Go
import "github.com/hashicorp/vault-client-go"
```

### Datadog Secrets Management
```bash
# Use Datadog Agent integrado para gerenciar secrets
# Veja: https://docs.datadoghq.com/agent/guide/secrets-management/
```

---

**Dica:** Nunca commite `.env` com dados sensíveis! Use `.env.example` para documentar a estrutura.

