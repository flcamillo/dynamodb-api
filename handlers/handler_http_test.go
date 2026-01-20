package handlers

import (
	"api/interfaces"
	"api/models"
	"api/repositories"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createTestMemoryDB() interfaces.Repository {
	return repositories.NewMemoryDB(1 * time.Hour)
}

func TestNewHttpHandler(t *testing.T) {
	repo := createTestMemoryDB()
	handler := NewHttpHandler(repo)

	if handler == nil {
		t.Fatalf("NewHttpHandler returned nil")
	}
	if handler.repository != repo {
		t.Errorf("Repository not set correctly")
	}
}

func TestHandleHealth(t *testing.T) {
	repo := createTestMemoryDB()
	handler := NewHttpHandler(repo)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code: expected %d, got %d", http.StatusOK, w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	if string(body) != "OK" {
		t.Errorf("Body: expected 'OK', got '%s'", string(body))
	}
}

func TestHandlePost(t *testing.T) {
	repo := createTestMemoryDB()
	repo.Create(context.Background())
	handler := NewHttpHandler(repo)

	event := models.Event{
		Date:          time.Now(),
		StatusCode:    201,
		StatusMessage: "Created",
		Metadata:      map[string]string{"key": "value"},
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/eventos", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.handlePost(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status code: expected %d, got %d", http.StatusCreated, w.Code)
	}

	var response models.Event
	json.NewDecoder(w.Body).Decode(&response)

	if response.Id == "" {
		t.Errorf("Response should have ID")
	}
	if response.StatusCode != 201 {
		t.Errorf("StatusCode mismatch")
	}
}

func TestHandlePostInvalidJSON(t *testing.T) {
	repo := createTestMemoryDB()
	repo.Create(context.Background())
	handler := NewHttpHandler(repo)

	req := httptest.NewRequest("POST", "/eventos", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.handlePost(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code: expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlePostInvalidDate(t *testing.T) {
	repo := createTestMemoryDB()
	repo.Create(context.Background())
	handler := NewHttpHandler(repo)

	event := models.Event{
		StatusCode:    200,
		StatusMessage: "OK",
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/eventos", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.handlePost(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code: expected %d, got %d", http.StatusBadRequest, w.Code)
	}

	var errorResp models.ErrorResponse
	json.NewDecoder(w.Body).Decode(&errorResp)

	if errorResp.Title != "Invalid Body" {
		t.Errorf("Error title mismatch")
	}
}

func TestHandleGet(t *testing.T) {
	repo := createTestMemoryDB()
	repo.Create(context.Background())
	handler := NewHttpHandler(repo)

	// Save an event first
	event := &models.Event{
		Id:            "123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
	}
	repo.Save(context.Background(), event)

	// Create request
	req := httptest.NewRequest("GET", "/eventos/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.handleGet(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code: expected %d, got %d", http.StatusOK, w.Code)
	}

	var response models.Event
	json.NewDecoder(w.Body).Decode(&response)

	if response.Id != "123" {
		t.Errorf("Event ID mismatch")
	}
}

func TestHandleGetNotFound(t *testing.T) {
	repo := createTestMemoryDB()
	repo.Create(context.Background())
	handler := NewHttpHandler(repo)

	req := httptest.NewRequest("GET", "/eventos/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()

	handler.handleGet(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status code: expected %d, got %d", http.StatusNotFound, w.Code)
	}

	var errorResp models.ErrorResponse
	json.NewDecoder(w.Body).Decode(&errorResp)

	if errorResp.Title != "Not Found" {
		t.Errorf("Error title mismatch")
	}
}

func TestHandleGetMissingID(t *testing.T) {
	repo := createTestMemoryDB()
	handler := NewHttpHandler(repo)

	req := httptest.NewRequest("GET", "/eventos/", nil)
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()

	handler.handleGet(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code: expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlePut(t *testing.T) {
	repo := createTestMemoryDB()
	repo.Create(context.Background())
	handler := NewHttpHandler(repo)

	event := models.Event{
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "Updated",
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("PUT", "/eventos/123", bytes.NewReader(body))
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.handlePut(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status code: expected %d, got %d", http.StatusCreated, w.Code)
	}

	var response models.Event
	json.NewDecoder(w.Body).Decode(&response)

	if response.Id != "123" {
		t.Errorf("Event ID should be 123, got %s", response.Id)
	}
}

func TestHandlePutMissingID(t *testing.T) {
	repo := createTestMemoryDB()
	handler := NewHttpHandler(repo)

	event := models.Event{
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "Updated",
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("PUT", "/eventos/", bytes.NewReader(body))
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()

	handler.handlePut(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code: expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlePutInvalidJSON(t *testing.T) {
	repo := createTestMemoryDB()
	handler := NewHttpHandler(repo)

	req := httptest.NewRequest("PUT", "/eventos/123", bytes.NewReader([]byte("invalid")))
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.handlePut(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code: expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleDelete(t *testing.T) {
	repo := createTestMemoryDB()
	repo.Create(context.Background())
	handler := NewHttpHandler(repo)

	// Save an event first
	event := &models.Event{
		Id:            "123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
	}
	repo.Save(context.Background(), event)

	req := httptest.NewRequest("DELETE", "/eventos/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.handleDelete(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Status code: expected %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestHandleDeleteNotFound(t *testing.T) {
	repo := createTestMemoryDB()
	handler := NewHttpHandler(repo)

	req := httptest.NewRequest("DELETE", "/eventos/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()

	handler.handleDelete(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status code: expected %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandleDeleteMissingID(t *testing.T) {
	repo := createTestMemoryDB()
	handler := NewHttpHandler(repo)

	req := httptest.NewRequest("DELETE", "/eventos/", nil)
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()

	handler.handleDelete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code: expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleRequestRegistration(t *testing.T) {
	repo := createTestMemoryDB()
	handler := NewHttpHandler(repo)

	router := http.NewServeMux()
	handler.HandleRequest(router)

	// Verify that routes are registered
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Route not registered correctly")
	}
}

func TestFromJsonError(t *testing.T) {
	repo := createTestMemoryDB()
	handler := NewHttpHandler(repo)

	req := httptest.NewRequest("POST", "/eventos", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	event := &models.Event{}
	handler.fromJson(w, req, event)

	// fromJson calls toJson which writes to response
	// In case of decode error, toJson is called which returns nil or error from encoder
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code: expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestToJsonSuccess(t *testing.T) {
	repo := createTestMemoryDB()
	handler := NewHttpHandler(repo)

	w := httptest.NewRecorder()
	response := models.ErrorResponse{
		Type:   "test",
		Status: 200,
		Title:  "Test",
	}

	err := handler.toJson(w, response, http.StatusOK)

	if err != nil {
		t.Errorf("toJson returned error: %v", err)
	}
	if w.Code != http.StatusOK {
		t.Errorf("Status code mismatch")
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header not set correctly")
	}
}

func TestHandleGetError(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewHttpHandler(repo)

	req := httptest.NewRequest("GET", "/eventos/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.handleGet(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status code: expected %d, got %d", http.StatusNotFound, w.Code)
	}

	var errorResp models.ErrorResponse
	json.NewDecoder(w.Body).Decode(&errorResp)

	if errorResp.Title != "Not Found" {
		t.Errorf("Error title mismatch")
	}
}

func TestHandlePostError(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewHttpHandler(repo)

	event := models.Event{
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/eventos", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.handlePost(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status code: expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestHandlePostInvalidStatusCode(t *testing.T) {
	repo := createTestMemoryDB()
	repo.Create(context.Background())
	handler := NewHttpHandler(repo)

	event := models.Event{
		Date:          time.Now(),
		StatusCode:    -1,
		StatusMessage: "Invalid",
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/eventos", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.handlePost(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code: expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlePutInvalidStatusCode(t *testing.T) {
	repo := createTestMemoryDB()
	repo.Create(context.Background())
	handler := NewHttpHandler(repo)

	event := models.Event{
		Date:          time.Now(),
		StatusCode:    -1,
		StatusMessage: "Invalid",
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("PUT", "/eventos/123", bytes.NewReader(body))
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.handlePut(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code: expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandlePutError(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewHttpHandler(repo)

	event := models.Event{
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "Updated",
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("PUT", "/eventos/123", bytes.NewReader(body))
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.handlePut(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status code: expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestHandleDeleteError(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewHttpHandler(repo)

	req := httptest.NewRequest("DELETE", "/eventos/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.handleDelete(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status code: expected %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandlePutInvalidDate(t *testing.T) {
	repo := createTestMemoryDB()
	repo.Create(context.Background())
	handler := NewHttpHandler(repo)

	event := models.Event{
		StatusCode:    200,
		StatusMessage: "Updated",
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("PUT", "/eventos/123", bytes.NewReader(body))
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.handlePut(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code: expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
