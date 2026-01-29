package apis

import (
	"api/handlers"
	"api/interfaces"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Configuração da API para servidor HTTP.
type HttpApiConfig struct {
	// log de aplicação
	Log interfaces.Log
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
	return &HttpApi{
		server: server,
		config: config,
	}
}

// Inicia a API para servidor HTTP.
func (p *HttpApi) Run() {
	// deve inicializar um context com cancelamento para receber sinais de término
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// inicia o servidor em uma goroutine
	router := http.NewServeMux()
	handler := handlers.NewHttpHandler(&handlers.HttpHandlerConfig{
		Log:        p.config.Log,
		Repository: p.config.Repository,
	})
	handler.HandleRequest(router)
	p.server.Handler = otelhttp.NewHandler(router, "http.handler") // configura telemetria no handler
	errChan := make(chan error, 1)
	go func() {
		p.config.Log.Info("starting server on: %s", p.server.Addr)
		errChan <- p.server.ListenAndServe()
	}()
	// aguarda o sinal de término
	select {
	case err := <-errChan:
		p.config.Log.Error("server error: %v", err)
	case <-ctx.Done():
		stop()
	}
	p.server.Shutdown(context.Background())
	p.config.Log.Info("server exiting")
}

// Encerra a API para servidor HTTP.
func (p *HttpApi) Shutdown() {
	p.server.Shutdown(context.Background())
}
