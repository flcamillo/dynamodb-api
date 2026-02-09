package main

import (
	"api/apis"
	"api/repositories"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

var (
	// configuração da aplicação
	applicationConfig *Config
	// encerramento do OTel SDK
	otelShutdown func(ctx context.Context) error
)

// inicializa recursos essenciais da aplicação
func init() {
	// inicializa o log padrão
	slog.SetDefault(otelslog.NewLogger(os.Getenv("DD_SERVICE")))
	// inicializa as configurações da aplicação
	currentDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		os.Exit(1)
	}
	applicationConfig, err = LoadConfig(filepath.Join(currentDirectory, "config.json"))
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		os.Exit(1)
	}
	// inicializa a telemetria
	os.Environ()
	otelShutdown, err = setupOTelSDK(context.Background())
	if err != nil {
		slog.Error(fmt.Sprintf("failed to setup OTel SDK: %s", err))
		os.Exit(1)
	}
	// inicializa o client do DynamoDB
	sdkConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		os.Exit(1)
	}
	applicationConfig.DynamoDBClient = dynamodb.NewFromConfig(sdkConfig)
	// inicializa o repositório
	applicationConfig.Repository = repositories.NewMemoryDB(&repositories.MemoryDBConfig{
		TTL: time.Duration(applicationConfig.RecordTTLMinutes) * time.Minute,
	})
	// applicationConfig.Repository = repositories.NewDynamoDBRepository(&repositories.DynamoDBConfig{
	// 	Log:    applicationConfig.Log,
	// 	Client: applicationConfig.DynamoDBClient,
	// 	Table:  "eventos",
	// 	TTL:    time.Duration(applicationConfig.RecordTTLMinutes) * time.Minute,
	// })
	if err := applicationConfig.Repository.Create(context.Background()); err != nil {
		slog.Error(fmt.Sprintf("failed to create repository: %s", err))
		os.Exit(1)
	}
}

// inicia a aplicação
func main() {
	// inicia a API
	api := apis.NewHttpApi(&apis.HttpApiConfig{
		Address:    applicationConfig.Address,
		Port:       applicationConfig.Port,
		Repository: applicationConfig.Repository,
	})
	api.Run()
	// encerra a telemetria
	err := otelShutdown(context.Background())
	if err != nil {
		slog.Error(fmt.Sprintf("failed to shutdown OTel SDK: %s", err))
	}
}
