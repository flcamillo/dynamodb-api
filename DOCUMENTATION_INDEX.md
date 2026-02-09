# 📚 Índice de Documentação - DynamoDB API

Bem-vindo! Escolha o documento que melhor se adequa ao seu objetivo:

## 👤 Por Tipo de Usuário

### 🚀 Desenvolvedor iniciando agora
1. **[QUICKSTART.md](QUICKSTART.md)** ⚡ - Setup em 5 minutos
2. **[.env.example](.env.example)** 📝 - Template de variáveis
3. Comece a codificar!

### 👨‍💻 Desenvolvedor continuando
1. **[README.md](README.md)** 📖 - Overview rápido
2. **[README_DETALHADO.md](README_DETALHADO.md)** 📊 - Arquitetura e API completa
3. **[ENV_VARIABLES.md](ENV_VARIABLES.md)** 🔧 - Configuração avançada

### 🏗️ Arquiteto de Sistemas
1. **[README_DETALHADO.md](README_DETALHADO.md#arquitetura-detalhada)** 📐 - Diagramas Mermaid
2. **[README_DETALHADO.md](README_DETALHADO.md#fluxo-de-requisição-detalhado)** 🔄 - Fluxos de dados
3. **[README_DETALHADO.md](README_DETALHADO.md#datadog-integration)** 🐶 - Observabilidade

### 🔬 DevOps / SRE
1. **[ENV_VARIABLES.md](ENV_VARIABLES.md)** 🔧 - Todas as variáveis de ambiente
2. **[README_DETALHADO.md](README_DETALHADO.md#datadog-integration)** 📊 - Integração Datadog
3. **[QUICKSTART.md](QUICKSTART.md)** - Docker Compose setup
4. **[extra/docker-compose.yml](extra/docker-compose.yml)** 🐳 - Orquestração

### 📖 Operador / Suporte
1. **[README_DETALHADO.md](README_DETALHADO.md#troubleshooting)** 🐛 - Troubleshooting
2. **[README_DETALHADO.md](README_DETALHADO.md#datadog-integration)** 🐶 - Observabilidade
3. **[ENV_VARIABLES.md](ENV_VARIABLES.md#troubleshooting-rápido)** ⚠️ - Problemas comuns

---

## 📂 Por Tópico

### Setup e Instalação
- **[QUICKSTART.md](QUICKSTART.md)** - 5 minutos de setup
- **[ENV_VARIABLES.md](ENV_VARIABLES.md)** - Configuração de variáveis
- **[.env.example](.env.example)** - Template
- **[README_DETALHADO.md](README_DETALHADO.md#instalação-e-configuração)** - Setup detalhado

### Arquitetura e Design
- **[README_DETALHADO.md](README_DETALHADO.md#arquitetura-detalhada)** - Diagramas Mermaid
- **[README_DETALHADO.md](README_DETALHADO.md#fluxo-de-requisição-detalhado)** - Sequência de requisições
- **[README_DETALHADO.md](README_DETALHADO.md#ciclo-de-vida-de-um-evento)** - Lifecycle de eventos
- **[README_DETALHADO.md](README_DETALHADO.md#arquitetura-de-logs-e-traces)** - Telemetria

### Uso da API
- **[README_DETALHADO.md](README_DETALHADO.md#endpoints-da-api)** - Documentação de endpoints
- **[README_DETALHADO.md](README_DETALHADO.md#exemplos-com-curl)** - Exemplos de cURL
- **[README.md](README.md#endpoints-da-api)** - Quick reference

### Observabilidade
- **[README_DETALHADO.md](README_DETALHADO.md#telemetria-e-observabilidade)** - OTEL overview
- **[README_DETALHADO.md](README_DETALHADO.md#datadog-integration)** - Datadog setup
- **[ENV_VARIABLES.md](ENV_VARIABLES.md#opentelemetry-configuration)** - OTEL vars
- **[extra/datadog.txt](extra/datadog.txt)** - Instruções Datadog

### Troubleshooting
- **[README_DETALHADO.md](README_DETALHADO.md#troubleshooting)** - Guia completo
- **[ENV_VARIABLES.md](ENV_VARIABLES.md#troubleshooting-rápido)** - Problemas comuns

### Deployment
- **[README_DETALHADO.md](README_DETALHADO.md#executando-a-aplicação)** - Opções de execução
- **[README_DETALHADO.md](README_DETALHADO.md#checklist-de-deploy)** - Checklist
- **[ENV_VARIABLES.md](ENV_VARIABLES.md#presets-de-configuração)** - Presets (dev/staging/prod)

---

## 🔍 Busca Rápida

### Quero fazer X...

#### ...começar a desenvolver
→ [QUICKSTART.md](QUICKSTART.md)

#### ...entender a arquitetura
→ [README_DETALHADO.md#arquitetura-detalhada](README_DETALHADO.md#arquitetura-detalhada)

#### ...integrar com Datadog
→ [README_DETALHADO.md#datadog-integration](README_DETALHADO.md#datadog-integration)

#### ...configurar variáveis de ambiente
→ [ENV_VARIABLES.md](ENV_VARIABLES.md)

#### ...testar a API
→ [README_DETALHADO.md#exemplos-com-curl](README_DETALHADO.md#exemplos-com-curl)

#### ...fazer deploy em produção
→ [README_DETALHADO.md#checklist-de-deploy](README_DETALHADO.md#checklist-de-deploy) + [ENV_VARIABLES.md#preset-4-produção](ENV_VARIABLES.md#preset-4-produção)

#### ...resolver um problema
→ [README_DETALHADO.md#troubleshooting](README_DETALHADO.md#troubleshooting)

#### ...entender como funcionam os logs
→ [README_DETALHADO.md#arquitetura-de-logs-e-traces](README_DETALHADO.md#arquitetura-de-logs-e-traces)

#### ...executar testes
→ [README_DETALHADO.md#testes](README_DETALHADO.md#testes)

#### ...testar performance e carga
→ [LOAD_TESTING.md](LOAD_TESTING.md)

---

## 📊 Mapa de Documentação

```
DynamoDB API
├── README.md (Este arquivo)
│   └── Overview rápido e quick start
│
├── QUICKSTART.md ⭐ Comece aqui
│   └── 5 minutos de setup
│
├── README_DETALHADO.md 📚 Documentação Completa
│   ├── Arquitetura (diagramas Mermaid)
│   ├── Instalação e Configuração
│   ├── Endpoints da API
│   ├── Exemplos com cURL
│   ├── Telemetria e Observabilidade
│   ├── Datadog Integration (com screenshots)
│   ├── Troubleshooting
│   └── Testes
│
├── ENV_VARIABLES.md 🔧 Configuração
│   ├── Todas as variáveis de ambiente
│   ├── Presets (dev/staging/prod)
│   ├── Scripts de setup
│   ├── Secrets seguros
│   └── Verificação de configuração
│
├── .env.example 📝 Template
│   └── Copie e customize
│
└── extra/
    ├── docker-compose.yml 🐳 Infraestrutura
    ├── otel-collector.yaml 📡 OTEL Config
    ├── prometheus.yaml 📊 Prometheus Config
    ├── datadog.txt 🐶 Datadog Instructions
    └── datadog_*.png 📸 Screenshots Datadog
```

---

## 🎯 Learning Path Recomendado

### Para Iniciantes (0-2h)
1. ✅ Ler [QUICKSTART.md](QUICKSTART.md) (10 min)
2. ✅ Executar o setup (5 min)
3. ✅ Testar com cURL (5 min)
4. ✅ Ler [README.md](README.md) overview (10 min)
5. ✅ Explorar code base (20 min)

### Para Intermediários (2-6h)
1. ✅ Ler [README_DETALHADO.md](README_DETALHADO.md#arquitetura-detalhada) (30 min)
2. ✅ Estudar diagramas Mermaid (20 min)
3. ✅ Entender fluxo de requisições (20 min)
4. ✅ Praticar exemplos de cURL (30 min)
5. ✅ Configurar OTEL/Datadog (1h)
6. ✅ Executar testes (30 min)

### Para Avançados (6h+)
1. ✅ Code review completo
2. ✅ Implementar melhorias
3. ✅ Setup production (ENV_VARIABLES.md)
4. ✅ Integração com CI/CD
5. ✅ Performance tuning (LOAD_TESTING.md)
6. ✅ Security hardening

---

## 🤝 Contribuindo

Ao adicionar novos recursos:
1. Atualize [README_DETALHADO.md](README_DETALHADO.md)
2. Adicione exemplos em [README_DETALHADO.md#exemplos-com-curl](README_DETALHADO.md#exemplos-com-curl)
3. Documente variáveis em [ENV_VARIABLES.md](ENV_VARIABLES.md)
4. Atualize [QUICKSTART.md](QUICKSTART.md) se relevante

---

## 📞 Suporte

- **Problema técnico?** → [Troubleshooting](README_DETALHADO.md#troubleshooting)
- **Configuração?** → [ENV_VARIABLES.md](ENV_VARIABLES.md)
- **Dúvida sobre API?** → [Endpoints](README_DETALHADO.md#endpoints-da-api)
- **Help com Datadog?** → [Datadog Integration](README_DETALHADO.md#datadog-integration)

---

**Última atualização:** Fevereiro 2026
**Versão:** 2.0.0

