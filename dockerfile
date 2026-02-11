# ----------------------------------------------------------------------------
# Efetua o build da aplicação Go usando uma imagem base do Golang
# ----------------------------------------------------------------------------
FROM golang:1.25.7-alpine3.23 AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o app

# ----------------------------------------------------------------------------
# Cria uma imagem final leve usando Alpine e copia o binário da aplicação
# ----------------------------------------------------------------------------
FROM alpine:3.23.3

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot
USER nonroot:nonroot

WORKDIR /app

COPY --from=builder /app/app /app/app

ENTRYPOINT ["/app/app"]

EXPOSE 7000
# ----------------------------------------------------------------------------
# Informações adicionais
# ----------------------------------------------------------------------------
# Para construir a imagem:
# docker build -t dynamodb-api .
# Para rodar o container:
# docker run -p 7000:7000 dynamodb-api
# docker run -it --entrypoint sh dynamodb-api