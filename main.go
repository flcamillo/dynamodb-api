package main

import (
	"api/apis"
	"api/logs"
	"api/repositories"
	"context"
	"fmt"
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
	log := logs.NewStdoutLog()
	// inicializa as configurações da aplicação
	currentDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error("%s", err)
		os.Exit(1)
	}
	applicationConfig, err = LoadConfig(filepath.Join(currentDirectory, "config.json"))
	if err != nil {
		log.Error("%s", err)
		os.Exit(1)
	}
	// inicializa a telemetria
	otelShutdown, err = setupOTelSDK(context.Background(), "dynamodb-api")
	if err != nil {
		log.Error(fmt.Sprintf("failed to setup OTel SDK: %s", err))
		os.Exit(1)
	}
	applicationConfig.Log = otelslog.NewLogger("api")
	// inicializa o client do DynamoDB
	sdkConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Error("%s", err)
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
		log.Error(fmt.Sprintf("failed to create repository: %s", err))
		os.Exit(1)
	}
}

// inicia a aplicação
func main() {
	// inicia a API
	api := apis.NewHttpApi(&apis.HttpApiConfig{
		Log:        applicationConfig.Log,
		Address:    applicationConfig.Address,
		Port:       applicationConfig.Port,
		Repository: applicationConfig.Repository,
	})
	api.Run()
	// encerra a telemetria
	err := otelShutdown(context.Background())
	if err != nil {
		applicationConfig.Log.Error(fmt.Sprintf("failed to shutdown OTel SDK: %s", err))
	}
}
