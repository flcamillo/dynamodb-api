package handlers

import (
	"api/interfaces"
	"api/models"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

// Estrutura do LambdaHandler.
type LambdaHandler struct {
	repository interfaces.Repository
}

// Cria uma nova instância do LambdaHandler.
func NewLambdaHandler(repository interfaces.Repository) *LambdaHandler {
	return &LambdaHandler{
		repository: repository,
	}
}

// Identifica o método HTTP da requisição e direciona para o handler apropriado.
func (p *LambdaHandler) HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	switch request.RequestContext.HTTP.Method {
	case "GET":
		return p.handleGet(ctx, request)
	case "POST":
		return p.handlePost(ctx, request)
	case "PUT":
		return p.handlePut(ctx, request)
	case "DELETE":
		return p.handleDelete(ctx, request)
	default:
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 405,
			Body:       "Method Not Allowed",
		}, nil
	}
}

// Processa requisições GET.
func (p *LambdaHandler) handleGet(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	id := request.PathParameters["id"]
	if id == "" {
		return p.toJson(models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	event, err := p.repository.Get(ctx, id)
	if err != nil {
		return p.toJson(models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusInternalServerError)
	}
	if event == nil {
		return p.toJson(models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Not Found",
			Status:   http.StatusNotFound,
			Detail:   "Event not found",
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusNotFound)
	}
	return p.toJson(event, http.StatusOK)
}

// Processa requisições POST.
func (p *LambdaHandler) handlePost(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	event := &models.Event{}
	if err := json.NewDecoder(strings.NewReader(request.Body)).Decode(event); err != nil {
		return p.toJson(models.ErrorResponse{
			Type:     "https://www.rfc-editor.org/rfc/rfc8259",
			Title:    "Invalid JSON",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	if err = event.Validate(); err != nil {
		return p.toJson(models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Invalid Body",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	event.Id = uuid.New().String()
	err = p.repository.Save(ctx, event)
	if err != nil {
		return p.toJson(models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusInternalServerError)
	}
	return p.toJson(event, http.StatusCreated)
}

// Processa requisições PUT.
func (p *LambdaHandler) handlePut(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	id := request.PathParameters["id"]
	if id == "" {
		return p.toJson(models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	event := &models.Event{}
	if err := json.NewDecoder(strings.NewReader(request.Body)).Decode(event); err != nil {
		return p.toJson(models.ErrorResponse{
			Type:     "https://www.rfc-editor.org/rfc/rfc8259",
			Title:    "Invalid JSON",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	if err = event.Validate(); err != nil {
		return p.toJson(models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Invalid Body",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	event.Id = id
	err = p.repository.Save(ctx, event)
	if err != nil {
		return p.toJson(models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusInternalServerError)
	}
	return p.toJson(event, http.StatusCreated)
}

// Processa requisições DELETE.
func (p *LambdaHandler) handleDelete(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	id := request.PathParameters["id"]
	if id == "" {
		return p.toJson(models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	event, err := p.repository.Delete(ctx, id)
	if err != nil {
		return p.toJson(models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusInternalServerError)
	}
	if event == nil {
		return p.toJson(models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Not Found",
			Status:   http.StatusNotFound,
			Detail:   "Event not found",
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusNotFound)
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}

// Converte o objeto para JSON e escreve na resposta HTTP.
func (p *LambdaHandler) toJson(object interface{}, statusCode int) (response events.APIGatewayV2HTTPResponse, err error) {
	data, err := json.Marshal(object)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Body:       string(data),
	}, nil
}
