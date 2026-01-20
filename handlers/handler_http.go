package handlers

import (
	"api/interfaces"
	"api/models"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// Estrutura do HttpHandler.
type HttpHandler struct {
	repository interfaces.Repository
}

// Cria uma nova instância do HttpHandler.
func NewHttpHandler(repository interfaces.Repository) *HttpHandler {
	return &HttpHandler{
		repository: repository,
	}
}

// Registra os handlers HTTP no roteador fornecido.
func (p *HttpHandler) HandleRequest(router *http.ServeMux) {
	router.HandleFunc("GET /health", p.handleHealth)
	router.HandleFunc("GET /eventos/{id}", p.handleGet)
	router.HandleFunc("POST /eventos", p.handlePost)
	router.HandleFunc("PUT /eventos/{id}", p.handlePut)
	router.HandleFunc("DELETE /eventos/{id}", p.handleDelete)
}

// Processa requisições para checagem de saúde da aplicação.
func (p *HttpHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Processa requisições GET.
func (p *HttpHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		p.toJson(w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event, err := p.repository.Get(r.Context(), id)
	if err != nil {
		p.toJson(w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	if event == nil {
		p.toJson(w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Not Found",
			Status:   http.StatusNotFound,
			Detail:   "Event not found",
			Instance: r.URL.String(),
		}, http.StatusNotFound)
		return
	}
	p.toJson(w, event, http.StatusOK)
}

// Processa requisições POST.
func (p *HttpHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	event := &models.Event{}
	if err := p.fromJson(w, r, event); err != nil {
		return
	}
	if err := event.Validate(); err != nil {
		p.toJson(w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Invalid Body",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event.Id = uuid.New().String()
	err := p.repository.Save(r.Context(), event)
	if err != nil {
		p.toJson(w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	p.toJson(w, event, http.StatusCreated)
}

// Processa requisições PUT.
func (p *HttpHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		p.toJson(w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event := &models.Event{}
	if err := p.fromJson(w, r, event); err != nil {
		return
	}
	if err := event.Validate(); err != nil {
		p.toJson(w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Invalid Body",
			Status:   http.StatusBadRequest,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event.Id = id
	err := p.repository.Save(r.Context(), event)
	if err != nil {
		p.toJson(w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	p.toJson(w, event, http.StatusCreated)
}

// Processa requisições DELETE.
func (p *HttpHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		p.toJson(w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   "Missing event ID in URL",
			Instance: r.URL.String(),
		}, http.StatusBadRequest)
		return
	}
	event, err := p.repository.Delete(r.Context(), id)
	if err != nil {
		p.toJson(w, models.ErrorResponse{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: r.URL.String(),
		}, http.StatusInternalServerError)
		return
	}
	if event == nil {
		p.toJson(w, models.ErrorResponse{
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

// Converte o corpo da requisição de JSON para o objeto fornecido.
func (p *HttpHandler) fromJson(w http.ResponseWriter, r *http.Request, object interface{}) error {
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(object)
	if err != nil {
		return p.toJson(w, models.ErrorResponse{
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
func (p *HttpHandler) toJson(w http.ResponseWriter, object interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(object)
}
