package apis

import (
	"api/handlers"
	"api/interfaces"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Estrutura da API para servidor HTTP.
type HttpApi struct {
	router     *http.ServeMux
	server     *http.Server
	repository interfaces.Repository
}

// Cria uma nova instância da API para servidor HTTP.
func NewHttpApi(port int, repository interfaces.Repository) *HttpApi {
	router := http.NewServeMux()
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        router,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1024 * 1024,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
	}
	return &HttpApi{
		router:     router,
		server:     server,
		repository: repository,
	}
}

// Inicia a API para servidor HTTP.
func (p *HttpApi) Run() {
	// deve inicializar um context com cancelamento para receber sinais de término
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// inicia o servidor em uma goroutine
	handler := handlers.NewHttpHandler(p.repository)
	handler.HandleRequest(p.router)
	errChan := make(chan error, 1)
	go func() {
		log.Printf("starting server on :%s", p.server.Addr)
		errChan <- p.server.ListenAndServe()
	}()
	// aguarda o sinal de término
	select {
	case err := <-errChan:
		log.Printf("server error: %v", err)
	case <-ctx.Done():
		stop()
	}
	p.server.Shutdown(context.Background())
	log.Println("server exiting")
}
