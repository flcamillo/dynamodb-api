package apis

import (
	"api/handlers"
	"api/interfaces"

	"github.com/aws/aws-lambda-go/lambda"
)

// Configuração da API para AWS Lambda.
type LambdaApiConfig struct {
	// repositório de dados
	Repository interfaces.Repository
}

// Estrutura da API para AWS Lambda.
type LambdaApi struct {
	// repositório de dados
	config *LambdaApiConfig
}

// Cria uma nova instância da API para AWS Lambda.
func NewLambdaApi(config *LambdaApiConfig) *LambdaApi {
	return &LambdaApi{
		config: config,
	}
}

// Inicia a API para AWS Lambda.
func (p *LambdaApi) Run() {
	handler := handlers.NewLambdaHandler(&handlers.LambdaHandlerConfig{
		Repository: p.config.Repository,
	})
	lambda.Start(handler.HandleRequest)
}
