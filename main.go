package main

import (
	"api/apis"
	"api/interfaces"
	"api/repositories"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// define as variaveis globais
var (
	dynamodbClient interfaces.DynamoDBClient
	repository     interfaces.Repository
	port           int
)

// inicializa recursos essenciais da aplicação
func init() {
	// inicializa o client do DynamoDB
	sdkConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	dynamodbClient = dynamodb.NewFromConfig(sdkConfig)
	// inicializa o repositório
	repository = repositories.NewMemoryDB(1 * time.Hour)
	if err := repository.Create(context.Background()); err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}
	// identifica a porta da API assumindo como default a 8080
	port = 8080
	portEnv := os.Getenv("API_PORT")
	if portEnv != "" {
		n, err := strconv.Atoi(portEnv)
		if err == nil && n > 0 && n < 65536 {
			port = n
		} else {
			log.Printf("invalid API_PORT value: %s, using default %d\n", portEnv, port)
		}
	}
}

// inicia a aplicação
func main() {
	api := apis.NewHttpApi(port, repository)
	api.Run()
}
