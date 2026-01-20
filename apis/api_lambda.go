package apis

import (
	"api/handlers"
	"api/interfaces"

	"github.com/aws/aws-lambda-go/lambda"
)

// Estrutura da API para AWS Lambda.
type LambdaApi struct {
	repository interfaces.Repository
}

// Cria uma nova instância da API para AWS Lambda.
func NewLambdaApi(repository interfaces.Repository) *LambdaApi {
	return &LambdaApi{
		repository: repository,
	}
}

// Inicia a API para AWS Lambda.
func (p *LambdaApi) Run() {
	handler := handlers.NewLambdaHandler(p.repository)
	lambda.Start(handler.HandleRequest)
}
