# Load Testing Guide - k6

Guide completo para executar testes de carga na DynamoDB API usando k6.

## O que é k6?

k6 é uma ferramenta moderna de teste de carga construída para equipes de DevOps, desenvolvida em Go.

**Características:**
- Sintaxe simples em JavaScript
- Resultados detalhados em tempo real
- Suporte a distributed testing
- Integração com Datadog, InfluxDB, Prometheus
- Cloud testing na plataforma Grafana Cloud

## Instalação

### macOS
```bash
brew install k6
```

### Linux (Ubuntu/Debian)
```bash
sudo apt-get install software-properties-common
sudo add-apt-repository ppa:k6-dev/k6-release
sudo apt-get update
sudo apt-get install k6
```

### Windows
```bash
choco install k6  # Com Chocolatey
# Ou baixe de: https://github.com/grafana/k6/releases
```

### Docker
```bash
docker pull grafana/k6:latest
```

## Executar Testes

### Teste Básico
```bash
k6 run load-test.js
```

### Com URL customizada
```bash
k6 run load-test.js --env BASE_URL=http://seu-host:7000
```

### Com variáveis de ambiente
```bash
k6 run load-test.js --env BASE_URL=http://localhost:7000 --env TEST_RUN=my-test
```

### Com saída em arquivo
```bash
k6 run load-test.js --out json=results.json
k6 run load-test.js --out csv=results.csv
```

### Com Datadog
```bash
k6 run load-test.js \
  --out datadog \
  --env DD_API_KEY=your_api_key \
  --env DD_SITE=datadoghq.com
```

### Com Prometheus
```bash
k6 run load-test.js \
  --out prometheus-remote \
  --env PROMETHEUS_URL=http://localhost:9090
```

### Cloud Testing (Grafana Cloud)
```bash
# Criar conta em: https://app.grafana.com
k6 cloud run load-test.js
```

## Interpretando Resultados

### Métricas Principais

```
http_req_duration - Latência da requisição
  avg: Média
  min: Mínima
  max: Máxima
  p(50): Mediana
  p(95): Percentil 95
  p(99): Percentil 99

http_req_failed - Taxa de falha (erros / total)
http_req_blocked - Tempo esperando conexão
http_req_connecting - Tempo de conexão
http_req_tls_handshaking - Tempo TLS
http_req_sending - Tempo enviando request
http_req_waiting - Tempo aguardando resposta
http_req_receiving - Tempo recebendo response

vus - Virtual Users (usuários simultâneos)
vus_max - Máximo de VUs
iterations - Total de iterações
```

### Exemplo de Saída

```
     http_requests......................: 50000 ok=50000 failed=0
     http_req_duration..................: avg=150ms   min=10ms    med=120ms   max=1500ms p(90)=250ms p(95)=350ms
     http_req_failed....................: 0%
     iterations..........................: 10000
     vus................................: 50

     check_status_is_200.................: 99.8%
     check_has_event_id..................: 99.8%
```

### Análise de Performance

- **p95 < 500ms**: ✅ Excelente
- **p95 < 1s**: ✅ Bom
- **p95 < 2s**: ⚠️ Aceitável
- **p95 > 2s**: ❌ Ruim
- **Failed > 1%**: ❌ Problema

## Customizar Testes

### Modificar Estágios

```javascript
export const options = {
  stages: [
    { duration: '30s', target: 20 },   // Ramp-up
    { duration: '1m30s', target: 20 }, // Stay at 20
    { duration: '30s', target: 0 },    // Ramp-down
  ],
};
```

### Modificar Thresholds

```javascript
export const options = {
  thresholds: {
    'http_req_duration': ['p(95)<200'],  // 95% < 200ms
    'http_req_failed': ['rate<0.05'],    // Taxa falha < 5%
  },
};
```

### Adicionar Setup/Teardown

```javascript
export function setup() {
  console.log('Executando setup...');
  return { /* dados */ };
}

export function teardown(data) {
  console.log('Executando teardown...');
}
```

## Casos de Teste Avançados

### Test 1: Teste de Pico (Spike Test)
```javascript
export const options = {
  stages: [
    { duration: '2m', target: 10 },     // Baseline
    { duration: '1m', target: 100 },    // Spike
    { duration: '3m', target: 100 },    // Sustain spike
    { duration: '2m', target: 10 },     // Recover
  ],
  thresholds: {
    'http_req_duration': ['p(99)<1500'],
  },
};
```

### Test 2: Teste de Soak (Longa Duração)
```javascript
export const options = {
  stages: [
    { duration: '5m', target: 50 },      // Ramp-up
    { duration: '30m', target: 50 },     // Soak
    { duration: '5m', target: 0 },       // Ramp-down
  ],
};
```

### Test 3: Teste de Stress (Limite)
```javascript
export const options = {
  stages: [
    { duration: '2m', target: 100 },
    { duration: '5m', target: 200 },
    { duration: '5m', target: 300 },
    { duration: '5m', target: 400 },
    { duration: '5m', target: 500 },
    { duration: '10m', target: 0 },
  ],
  thresholds: {
    'http_req_duration': ['p(99)<3000'],
  },
};
```

### Test 4: Teste de Endurance
```javascript
export const options = {
  stages: [
    { duration: '1m', target: 30 },      // Ramp-up
    { duration: '8h', target: 30 },      // 8 horas constante
    { duration: '1m', target: 0 },       // Ramp-down
  ],
};
```

## Integração com CI/CD

### GitHub Actions
```yaml
name: Load Tests
on: [push, pull_request]

jobs:
  load-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: grafana/setup-k6-action@v1
      - run: k6 run load-test.js
        env:
          BASE_URL: http://localhost:7000
```

### GitLab CI
```yaml
load_test:
  image: grafana/k6:latest
  script:
    - k6 run load-test.js
  variables:
    BASE_URL: "http://localhost:7000"
```

## Comparar Resultados

### Gerar relatório
```bash
# JSON detalhado
k6 run load-test.js --out json=baseline.json

# Depois comparar
k6 run load-test.js --out json=new-test.json
```

### Analisar com Datadog
```bash
# Todos os testes com Datadog
k6 run load-test.js --out datadog --env DD_API_KEY=xxx
# Ver resultados em: https://app.datadoghq.com
```

## Troubleshooting

### "Connection Refused"
```bash
# Verificar se API está rodando
curl http://localhost:7000/health

# Testar com URL correcta
k6 run load-test.js --env BASE_URL=http://seu-host:7000
```

### "Too Many Requests"
```javascript
// Aumentar sleep
sleep(2);  // ao invés de sleep(1)

// Ou reduzir VUs
export const options = {
  stages: [
    { duration: '1m', target: 10 },  // Menor valor
  ],
};
```

### "High Duration / Slow Response"
- Verificar backend (ver logs da API)
- Verificar DynamoDB (latência)
- Verificar rede
- Reduzir número de VUs

## Best Practices

1. **Sempre fazer warmup** antes de testes reais
2. **Testar ambiente isolado** (não produção)
3. **Usar thresholds realisticamente** (baseado em SLA)
4. **Rodar múltiplas vezes** para consistência
5. **Analisar logs** durante teste (Datadog/Logs)
6. **Documentar resultados** para comparação futura
7. **Graduar o aumento de carga** (ramp-up)

## Referências

- [k6 Documentação](https://k6.io/docs/)
- [k6 GitHub](https://github.com/grafana/k6)
- [Grafana Cloud](https://grafana.com/auth/sign-up/create-account)
- [Test Planning Guide](https://k6.io/docs/test-types/introduction/)

---

**Dica:** Execute `k6 run load-test.js` regularmente para acompanhar performance!

