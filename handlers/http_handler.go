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

	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Estrutura do ResponseWriter para capturar o status code das requisições.
type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

// Sobrescreve o método WriteHeader para capturar o status code.
func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.statusCode = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

// Sobrescreve o método Write para garantir que o status code seja capturado mesmo quando WriteHeader não é chamado explicitamente.
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// Configuração do HttpHandler.
type HttpHandlerConfig struct {
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
	requestCounter   metric.Int64Counter
	requestHistogram metric.Float64Histogram
}

// Cria uma nova instância do HttpHandler.
func NewHttpHandler(config *HttpHandlerConfig) *HttpHandler {
	h := &HttpHandler{
		config: config,
		tracer: otel.Tracer("http.handler"),
	}
	// configura as metricas
	meter := otel.Meter("http.server.metrics")
	if counter, err := meter.Int64Counter("custom.http.requests.total",
		metric.WithDescription("The number of HTTP requests executed"),
		metric.WithUnit("{requests}")); err == nil {
		h.requestCounter = counter
	} else {
		panic(err)
	}
	if histogram, err := meter.Float64Histogram("custom.http.requests.duration",
		metric.WithDescription("The duration of HTTP requests"),
		metric.WithUnit("ms")); err == nil {
		h.requestHistogram = histogram
	} else {
		panic(err)
	}
	return h
}

// helper para adicionar metricas na rota.
func (p *HttpHandler) routeHandler(route string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		h.ServeHTTP(rw, r)
		duration := time.Since(start)
		attrs := []attribute.KeyValue{
			attribute.String("http.method", strings.ToUpper(r.Method)),
			attribute.String("http.route", route),
			attribute.String("http.status_code", fmt.Sprintf("%d", rw.statusCode)),
		}
		p.requestCounter.Add(r.Context(), 1, metric.WithAttributes(attrs...))
		p.requestHistogram.Record(r.Context(), float64(duration.Milliseconds()), metric.WithAttributes(attrs...))
	})
}

// Registra os handlers HTTP no roteador fornecido.
func (p *HttpHandler) HandleRequest(router *http.ServeMux) {
	router.Handle("GET /health", otelhttp.NewHandler(p.routeHandler("/health", http.HandlerFunc(p.handleHealth)), ""))
	router.Handle("GET /eventos", otelhttp.NewHandler(p.routeHandler("/eventos", http.HandlerFunc(p.handleFind)), ""))
	router.Handle("GET /eventos/{id}", otelhttp.NewHandler(p.routeHandler("/eventos/{id}", http.HandlerFunc(p.handleGet)), ""))
	router.Handle("POST /eventos", otelhttp.NewHandler(p.routeHandler("/eventos", http.HandlerFunc(p.handlePost)), ""))
	router.Handle("PUT /eventos/{id}", otelhttp.NewHandler(p.routeHandler("/eventos/{id}", http.HandlerFunc(p.handlePut)), ""))
	router.Handle("DELETE /eventos/{id}", otelhttp.NewHandler(p.routeHandler("/eventos/{id}", http.HandlerFunc(p.handleDelete)), ""))
}

// Processa requisições para checagem de saúde da aplicação.
func (p *HttpHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	_, span := p.tracer.Start(r.Context(), "handleHealth")
	defer span.End()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Processa requisições GET.
func (p *HttpHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	ctx, span := p.tracer.Start(r.Context(), "handleGet")
	defer span.End()
	id := r.PathValue("id")
	if id == "" {
		span.AddEvent("{id} not provided")
		p.toJson(ctx, w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event, err := p.config.Repository.Get(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to get record from repository")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to get record from repository, %s", err))
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
		span.AddEvent("record not found")
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
	ctx, span := p.tracer.Start(r.Context(), "handlePost")
	defer span.End()
	event := &models.Event{}
	if err := p.fromJson(ctx, w, r, event); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to decode json")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to decode json, %s", err))
		return
	}
	if err := event.Validate(); err != nil {
		span.AddEvent(
			"record validation failed",
			trace.WithAttributes(attribute.String("error", err.Error())),
		)
		p.toJson(ctx, w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Invalid Body",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event.Id = uuid.New().String()
	err := p.config.Repository.Save(ctx, event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to save record in repository")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to save record in repository, %s", err))
		p.toJson(ctx, w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	p.toJson(ctx, w, event, http.StatusCreated)
}

// Processa requisições PUT.
func (p *HttpHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	ctx, span := p.tracer.Start(r.Context(), "handlePut")
	defer span.End()
	id := r.PathValue("id")
	if id == "" {
		span.AddEvent("{id} not provided")
		p.toJson(ctx, w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event := &models.Event{}
	if err := p.fromJson(ctx, w, r, event); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to decode json")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to decode json, %s", err))
		return
	}
	if err := event.Validate(); err != nil {
		span.AddEvent(
			"record validation failed",
			trace.WithAttributes(attribute.String("error", err.Error())),
		)
		p.toJson(ctx, w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Invalid Body",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event.Id = id
	err := p.config.Repository.Save(ctx, event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to save record in repository")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to save record in repository, %s", err))
		p.toJson(ctx, w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	p.toJson(ctx, w, event, http.StatusCreated)
}

// Processa requisições DELETE.
func (p *HttpHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx, span := p.tracer.Start(r.Context(), "handleDelete")
	defer span.End()
	id := r.PathValue("id")
	if id == "" {
		span.AddEvent("{id} not provided")
		p.toJson(ctx, w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event, err := p.config.Repository.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to delete record from repository")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to delete record from repository, %s", err))
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
		span.AddEvent("record not found")
		p.toJson(ctx, w, models.ErrorResponse{
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
	ctx, span := p.tracer.Start(r.Context(), "handleFind")
	defer span.End()
	// deve fazer o parser para validar se não há erros no formulario
	err := r.ParseForm()
	if err != nil {
		span.AddEvent(
			"unable to parse form data",
			trace.WithAttributes(attribute.String("error", err.Error())),
		)
		p.toJson(ctx, w, models.ErrorResponse{
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
			span.AddEvent(
				"unable to parse value of {from} to RFC3339",
				trace.WithAttributes(attribute.String("error", err.Error())),
			)
			p.toJson(ctx, w, models.ErrorResponse{
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
			span.AddEvent(
				"unable to parse value of {to} to RFC3339",
				trace.WithAttributes(attribute.String("error", err.Error())),
			)
			p.toJson(ctx, w, models.ErrorResponse{
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
			span.AddEvent(
				"unable to parse value of {statusCode} to interger",
				trace.WithAttributes(attribute.String("error", err.Error())),
			)
			p.toJson(ctx, w, models.ErrorResponse{
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
		span.SetStatus(codes.Error, "unable to find record in repository")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to find record in repository, %s", err))
		p.toJson(ctx, w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	p.toJson(ctx, w, events, http.StatusOK)
}

// Converte o corpo da requisição de JSON para o objeto fornecido.
func (p *HttpHandler) fromJson(ctx context.Context, w http.ResponseWriter, r *http.Request, object interface{}) error {
	ctx, span := p.tracer.Start(ctx, "fromJson")
	defer span.End()
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(object)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to decode json")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to decode json, %s", err))
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
	err := json.NewEncoder(w).Encode(object)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unable to encode json")
		slog.ErrorContext(ctx, fmt.Sprintf("unable to encode json, %s", err))
	}
	return err
}
