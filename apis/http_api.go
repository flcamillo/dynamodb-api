package apis

import (
	"api/handlers"
	"api/interfaces"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

// Configuração da API para servidor HTTP.
type HttpApiConfig struct {
	// endereço do servidor
	Address string
	// porta do servidor
	Port int
	// repositório de dados
	Repository interfaces.Repository
}

// Estrutura da API para servidor HTTP.
type HttpApi struct {
	// servidor HTTP
	server *http.Server
	// configuração da API
	config *HttpApiConfig
}

// Cria uma nova instância da API para servidor HTTP.
func NewHttpApi(config *HttpApiConfig) *HttpApi {
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Address, config.Port),
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1024 * 1024,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
	}
	h := &HttpApi{
		server: server,
		config: config,
	}
	return h
}

// Middleware para logar todas as requests, inclusive 404 e coletar metricas.
func (p *HttpApi) basicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		slog.InfoContext(
			r.Context(),
			fmt.Sprintf("request duration {%dms} status code {%d} method {%s} path {%s} remote address {%s} agent {%s}",
				duration.Milliseconds(),
				rw.statusCode,
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				r.UserAgent(),
			),
		)
	})
}

// Inicia a API para servidor HTTP.
func (p *HttpApi) Run() {
	// deve inicializar um context com cancelamento para receber sinais de término
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// configura o handler
	router := http.NewServeMux()
	handler := handlers.NewHttpHandler(&handlers.HttpHandlerConfig{
		Repository: p.config.Repository,
	})
	handler.HandleRequest(router)
	p.server.Handler = p.basicMiddleware(router)
	// inicia o servidor em uma goroutine
	errChan := make(chan error, 1)
	go func() {
		slog.Info(fmt.Sprintf("starting server on: %s", p.server.Addr))
		errChan <- p.server.ListenAndServe()
	}()
	// aguarda o sinal de término
	select {
	case err := <-errChan:
		slog.Error(fmt.Sprintf("server error: %s", err))
	case <-ctx.Done():
		stop()
	}
	p.server.Shutdown(context.Background())
	slog.Info("server exiting")
}

// Encerra a API para servidor HTTP.
func (p *HttpApi) Shutdown() {
	p.server.Shutdown(context.Background())
}
