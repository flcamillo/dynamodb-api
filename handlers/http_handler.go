package handlers

import (
	"api/interfaces"
	"api/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Configuração do HttpHandler.
type HttpHandlerConfig struct {
	// log de aplicação
	Log interfaces.Log
	// repositório de dados
	Repository interfaces.Repository
}

// Estrutura do HttpHandler.
type HttpHandler struct {
	// configuração do handler
	config *HttpHandlerConfig
	// configura o tracer
	tracer trace.Tracer
	// metricas de requisições
	postCounter   metric.Int64Counter
	getCounter    metric.Int64Counter
	putCounter    metric.Int64Counter
	deleteCounter metric.Int64Counter
	findCounter   metric.Int64Counter
}

// Cria uma nova instância do HttpHandler.
func NewHttpHandler(config *HttpHandlerConfig) *HttpHandler {
	h := &HttpHandler{
		config: config,
		tracer: otel.Tracer("http.handler"),
	}
	// configura as metricas
	meter := otel.Meter("http.handler")
	if counter, err := meter.Int64Counter("post.requests",
		metric.WithDescription("The number POST executed"),
		metric.WithUnit("{requests}")); err == nil {
		h.postCounter = counter
	} else {
		panic(err)
	}
	if counter, err := meter.Int64Counter("get.requests",
		metric.WithDescription("The number GET executed"),
		metric.WithUnit("{requests}")); err == nil {
		h.getCounter = counter
	} else {
		panic(err)
	}
	if counter, err := meter.Int64Counter("put.requests",
		metric.WithDescription("The number PUT executed"),
		metric.WithUnit("{requests}")); err == nil {
		h.putCounter = counter
	} else {
		panic(err)
	}
	if counter, err := meter.Int64Counter("delete.requests",
		metric.WithDescription("The number DELETE executed"),
		metric.WithUnit("{requests}")); err == nil {
		h.deleteCounter = counter
	} else {
		panic(err)
	}
	if counter, err := meter.Int64Counter("find.requests",
		metric.WithDescription("The number FIND executed"),
		metric.WithUnit("{requests}")); err == nil {
		h.deleteCounter = counter
	} else {
		panic(err)
	}
	return h
}

// Middleware para logar as requisições HTTP.
func (p *HttpHandler) logHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.config.Log.Info("received {%s} request from {%s} on {%s} agent {%s}", r.Method, r.RemoteAddr, r.URL.Path, r.UserAgent())
		h(w, r)
	}
}

// Registra os handlers HTTP no roteador fornecido.
func (p *HttpHandler) HandleRequest(router *http.ServeMux) {
	router.HandleFunc("GET /health", p.handleHealth)
	router.HandleFunc("GET /eventos", p.logHandlerFunc(p.handleFind))
	router.HandleFunc("GET /eventos/{id}", p.logHandlerFunc(p.handleGet))
	router.HandleFunc("POST /eventos", p.logHandlerFunc(p.handlePost))
	router.HandleFunc("PUT /eventos/{id}", p.logHandlerFunc(p.handlePut))
	router.HandleFunc("DELETE /eventos/{id}", p.logHandlerFunc(p.handleDelete))
}

// Processa requisições para checagem de saúde da aplicação.
func (p *HttpHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Processa requisições GET.
func (p *HttpHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	ctx, span := p.tracer.Start(r.Context(), "get")
	defer span.End()
	p.getCounter.Add(ctx, 1)
	id := r.PathValue("id")
	if id == "" {
		span.AddEvent("id not provided")
		p.toJson(ctx, w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event, err := p.config.Repository.Get(r.Context(), id)
	if err != nil {
		span.RecordError(err)
		p.toJson(ctx, w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	if event == nil {
		p.toJson(ctx, w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Not Found",
			Status:   http.StatusNotFound,
			Detail:   "Event not found",
			Instance: r.URL.String(),
		}, http.StatusNotFound)
		return
	}
	p.toJson(ctx, w, event, http.StatusOK)
}

// Processa requisições POST.
func (p *HttpHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	ctx, span := p.tracer.Start(r.Context(), "post")
	defer span.End()
	p.postCounter.Add(ctx, 1)
	event := &models.Event{}
	if err := p.fromJson(r.Context(), w, r, event); err != nil {
		return
	}
	if err := event.Validate(); err != nil {
		span.RecordError(err)
		p.toJson(r.Context(), w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Invalid Body",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event.Id = uuid.New().String()
	err := p.config.Repository.Save(r.Context(), event)
	if err != nil {
		span.RecordError(err)
		p.toJson(r.Context(), w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	p.toJson(r.Context(), w, event, http.StatusCreated)
}

// Processa requisições PUT.
func (p *HttpHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	ctx, span := p.tracer.Start(r.Context(), "put")
	defer span.End()
	p.putCounter.Add(ctx, 1)
	id := r.PathValue("id")
	if id == "" {
		span.AddEvent("id not provided")
		p.toJson(r.Context(), w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event := &models.Event{}
	if err := p.fromJson(r.Context(), w, r, event); err != nil {
		return
	}
	if err := event.Validate(); err != nil {
		span.RecordError(err)
		p.toJson(r.Context(), w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Invalid Body",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event.Id = id
	err := p.config.Repository.Save(r.Context(), event)
	if err != nil {
		span.RecordError(err)
		p.toJson(r.Context(), w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	p.toJson(r.Context(), w, event, http.StatusCreated)
}

// Processa requisições DELETE.
func (p *HttpHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx, span := p.tracer.Start(r.Context(), "delete")
	defer span.End()
	p.deleteCounter.Add(ctx, 1)
	id := r.PathValue("id")
	if id == "" {
		span.AddEvent("id not provided")
		p.toJson(r.Context(), w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event, err := p.config.Repository.Delete(r.Context(), id)
	if err != nil {
		span.RecordError(err)
		p.toJson(r.Context(), w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	if event == nil {
		p.toJson(r.Context(), w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Not Found",
			Status:   http.StatusNotFound,
			Detail:   "Event not found",
			Instance: r.URL.String(),
		}, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Processa requisições GET com filtro.
func (p *HttpHandler) handleFind(w http.ResponseWriter, r *http.Request) {
	ctx, span := p.tracer.Start(r.Context(), "find")
	defer span.End()
	p.putCounter.Add(ctx, 1)
	// deve fazer o parser para validar se não há erros no formulario
	err := r.ParseForm()
	if err != nil {
		span.RecordError(err)
		p.toJson(r.Context(), w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Invalid Body",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	// configura valores default caso sejam informados
	from := time.Now().Add(-1 * time.Hour)
	to := time.Now()
	statusCode := 0
	// trata os valores informados atualizando o default
	// se necessário
	if v := strings.TrimSpace(r.Form.Get("from")); v != "" {
		from, err = time.Parse(time.RFC3339, v)
		if err != nil {
			span.RecordError(err)
			p.toJson(r.Context(), w, models.ErrorResponse{
				Type:     "about:blank",
				Title:    "Invalid Body",
				Status:   http.StatusBadRequest,
				Detail:   fmt.Sprintf("parameter {from} invalid, %s", err.Error()),
				Instance: r.URL.String(),
			}, http.StatusBadRequest)
			return
		}
	}
	if v := strings.TrimSpace(r.Form.Get("to")); v != "" {
		to, err = time.Parse(time.RFC3339, v)
		if err != nil {
			span.RecordError(err)
			p.toJson(r.Context(), w, models.ErrorResponse{
				Type:     "about:blank",
				Title:    "Invalid Body",
				Status:   http.StatusBadRequest,
				Detail:   fmt.Sprintf("parameter {to} invalid, %s", err.Error()),
				Instance: r.URL.String(),
			}, http.StatusBadRequest)
			return
		}
	}
	if v := strings.TrimSpace(r.Form.Get("statusCode")); v != "" {
		statusCode, err = strconv.Atoi(v)
		if err != nil {
			span.RecordError(err)
			p.toJson(r.Context(), w, models.ErrorResponse{
				Type:     "about:blank",
				Title:    "Invalid Body",
				Status:   http.StatusBadRequest,
				Detail:   fmt.Sprintf("parameter {statusCode} invalid, %s", err.Error()),
				Instance: r.URL.String(),
			}, http.StatusBadRequest)
			return
		}
	}
	events, err := p.config.Repository.FindByDateAndReturnCode(ctx, from, to, statusCode)
	if err != nil {
		span.RecordError(err)
		p.toJson(r.Context(), w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	p.toJson(r.Context(), w, events, http.StatusOK)
}

// Converte o corpo da requisição de JSON para o objeto fornecido.
func (p *HttpHandler) fromJson(ctx context.Context, w http.ResponseWriter, r *http.Request, object interface{}) error {
	ctx, span := p.tracer.Start(ctx, "fromJson")
	defer span.End()
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(object)
	if err != nil {
		span.RecordError(err)
		return p.toJson(ctx, w, models.ErrorResponse{
			Type:     "https://www.rfc-editor.org/rfc/rfc8259",
			Title:    "Invalid JSON",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
	}
	return nil
}

// Converte o objeto para JSON e escreve na resposta HTTP.
func (p *HttpHandler) toJson(ctx context.Context, w http.ResponseWriter, object interface{}, statusCode int) error {
	ctx, span := p.tracer.Start(ctx, "toJson")
	defer span.End()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(object)
}
