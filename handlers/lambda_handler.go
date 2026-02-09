package handlers

import (
	"api/interfaces"
	"api/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Configuração do LambdaHandler.
type LambdaHandlerConfig struct {
	// log de aplicação
	Log interfaces.Log
	// repositório de dados
	Repository interfaces.Repository
}

// Estrutura do LambdaHandler.
type LambdaHandler struct {
	// configuração do handler
	config *LambdaHandlerConfig
	// configura o tracer
	tracer trace.Tracer
}

// Cria uma nova instância do LambdaHandler.
func NewLambdaHandler(config *LambdaHandlerConfig) *LambdaHandler {
	return &LambdaHandler{
		config: config,
		tracer: otel.Tracer("lambda.handler"),
	}
}

// Identifica o método HTTP da requisição e direciona para o handler apropriado.
func (p *LambdaHandler) HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	start := time.Now()
	switch request.RequestContext.HTTP.Method {
	case "GET":
		response, err = p.handleGet(ctx, request)
	case "POST":
		response, err = p.handlePost(ctx, request)
	case "PUT":
		response, err = p.handlePut(ctx, request)
	case "DELETE":
		response, err = p.handleDelete(ctx, request)
	default:
		response, err = events.APIGatewayV2HTTPResponse{
			StatusCode: 405,
			Body:       "Method Not Allowed",
		}, nil
	}
	duration := time.Since(start)
	p.config.Log.Info(
		fmt.Sprintf("request duration {%dms} status code {%d} method {%s} path {%s} remote address {%s} agent {%s}",
			duration.Milliseconds(),
			response.StatusCode,
			request.RequestContext.HTTP.Method,
			request.RequestContext.HTTP.Path,
			request.RequestContext.HTTP.SourceIP,
			request.RequestContext.HTTP.UserAgent,
		),
	)
	return response, err
}

// Processa requisições GET.
func (p *LambdaHandler) handleGet(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	ctx, span := p.tracer.Start(
		ctx,
		"handleGet",
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()
	id := request.PathParameters["id"]
	if id == "" {
		span.AddEvent("{id} not provided")
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	event, err := p.config.Repository.Get(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to get record from repository")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to get record from repository, %s", err))
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusInternalServerError)
	}
	if event == nil {
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Not Found",
			Status:   http.StatusNotFound,
			Detail:   "Event not found",
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusNotFound)
	}
	return p.toJson(ctx, event, http.StatusOK)
}

// Processa requisições POST.
func (p *LambdaHandler) handlePost(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	ctx, span := p.tracer.Start(
		ctx,
		"handlePost",
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()
	event := &models.Event{}
	if err := json.NewDecoder(strings.NewReader(request.Body)).Decode(event); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to decode json")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to decode json, %s", err))
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "https://www.rfc-editor.org/rfc/rfc8259",
			Title:    "Invalid JSON",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	if err = event.Validate(); err != nil {
		span.AddEvent(
			"record validation failed",
			trace.WithAttributes(attribute.String("error", err.Error())),
		)
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Invalid Body",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	event.Id = uuid.New().String()
	err = p.config.Repository.Save(ctx, event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to save record in repository")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to save record in repository, %s", err))
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusInternalServerError)
	}
	return p.toJson(ctx, event, http.StatusCreated)
}

// Processa requisições PUT.
func (p *LambdaHandler) handlePut(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	ctx, span := p.tracer.Start(
		ctx,
		"handlePut",
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()
	id := request.PathParameters["id"]
	if id == "" {
		span.AddEvent("{id} not provided")
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	event := &models.Event{}
	if err := json.NewDecoder(strings.NewReader(request.Body)).Decode(event); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to decode json")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to decode json, %s", err))
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "https://www.rfc-editor.org/rfc/rfc8259",
			Title:    "Invalid JSON",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	if err = event.Validate(); err != nil {
		span.AddEvent(
			"record validation failed",
			trace.WithAttributes(attribute.String("error", err.Error())),
		)
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Invalid Body",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	event.Id = id
	err = p.config.Repository.Save(ctx, event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to save record in repository")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to save record in repository, %s", err))
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusInternalServerError)
	}
	return p.toJson(ctx, event, http.StatusCreated)
}

// Processa requisições DELETE.
func (p *LambdaHandler) handleDelete(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	ctx, span := p.tracer.Start(
		ctx,
		"handleDelete",
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()
	id := request.PathParameters["id"]
	if id == "" {
		span.AddEvent("{id} not provided")
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusBadRequest)
	}
	event, err := p.config.Repository.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to delete record from repository")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to delete record from repository, %s", err))
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusInternalServerError)
	}
	if event == nil {
		span.AddEvent("record not found")
		return p.toJson(ctx, models.ErrorResponse{
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

// Processa requisições GET com filtro.
func (p *LambdaHandler) handleFind(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	ctx, span := p.tracer.Start(
		ctx,
		"handleFind",
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()
	// configura valores default caso sejam informados
	from := time.Now().Add(-1 * time.Hour)
	to := time.Now()
	statusCode := 0
	// trata os valores informados atualizando o default
	// se necessário
	if v := strings.TrimSpace(request.QueryStringParameters["from"]); v != "" {
		from, err = time.Parse(time.RFC3339, v)
		if err != nil {
			span.AddEvent(
				"unable to parse value of {from} to RFC3339",
				trace.WithAttributes(attribute.String("error", err.Error())),
			)
			return p.toJson(ctx, models.ErrorResponse{
				Type:     "about:blank",
				Title:    "Invalid Body",
				Status:   http.StatusBadRequest,
				Detail:   fmt.Sprintf("parameter {from} invalid, %s", err.Error()),
				Instance: request.RequestContext.HTTP.Path,
			}, http.StatusBadRequest)
		}
	}
	if v := strings.TrimSpace(request.QueryStringParameters["to"]); v != "" {
		to, err = time.Parse(time.RFC3339, v)
		if err != nil {
			span.AddEvent(
				"unable to parse value of {to} to RFC3339",
				trace.WithAttributes(attribute.String("error", err.Error())),
			)
			return p.toJson(ctx, models.ErrorResponse{
				Type:     "about:blank",
				Title:    "Invalid Body",
				Status:   http.StatusBadRequest,
				Detail:   fmt.Sprintf("parameter {to} invalid, %s", err.Error()),
				Instance: request.RequestContext.HTTP.Path,
			}, http.StatusBadRequest)
		}
	}
	if v := strings.TrimSpace(request.QueryStringParameters["statusCode"]); v != "" {
		statusCode, err = strconv.Atoi(v)
		if err != nil {
			span.AddEvent(
				"unable to parse value of {statusCode} to interger",
				trace.WithAttributes(attribute.String("error", err.Error())),
			)
			return p.toJson(ctx, models.ErrorResponse{
				Type:     "about:blank",
				Title:    "Invalid Body",
				Status:   http.StatusBadRequest,
				Detail:   fmt.Sprintf("parameter {statusCode} invalid, %s", err.Error()),
				Instance: request.RequestContext.HTTP.Path,
			}, http.StatusBadRequest)
		}
	}
	events, err := p.config.Repository.FindByDateAndReturnCode(ctx, from, to, statusCode)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to find record in repository")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to find record in repository, %s", err))
		return p.toJson(ctx, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: request.RequestContext.HTTP.Path,
		}, http.StatusInternalServerError)
	}
	return p.toJson(ctx, events, http.StatusOK)
}

// Converte o objeto para JSON e escreve na resposta HTTP.
func (p *LambdaHandler) toJson(ctx context.Context, object interface{}, statusCode int) (response events.APIGatewayV2HTTPResponse, err error) {
	ctx, span := p.tracer.Start(ctx, "toJson")
	defer span.End()
	data, err := json.Marshal(object)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to marshal object to json")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to marshal object to json, %s", err))
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
