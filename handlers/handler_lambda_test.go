package handlers

import (
	"api/models"
	"api/repositories"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func TestNewLambdaHandler(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	if handler == nil {
		t.Fatalf("NewLambdaHandler returned nil")
	}
	if handler.repository != repo {
		t.Errorf("Repository not set correctly")
	}
}

func TestHandleRequestUnsupportedMethod(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	request := events.APIGatewayV2HTTPRequest{
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "PATCH",
				Path:   "/eventos",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 405 {
		t.Errorf("Status code: expected 405, got %d", response.StatusCode)
	}
}

func TestLambdaHandleGetSuccess(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	repo.Create(context.Background())
	handler := NewLambdaHandler(repo)

	event := &models.Event{
		Id:            "123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
	}
	repo.Save(context.Background(), event)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "123"},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/eventos/123",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Status code: expected 200, got %d", response.StatusCode)
	}

	var respEvent models.Event
	json.Unmarshal([]byte(response.Body), &respEvent)

	if respEvent.Id != "123" {
		t.Errorf("Event ID mismatch")
	}
}

func TestLambdaHandleGetNotFound(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "nonexistent"},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/eventos/nonexistent",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 404 {
		t.Errorf("Status code: expected 404, got %d", response.StatusCode)
	}
}

func TestLambdaHandleGetMissingID(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": ""},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/eventos/",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 400 {
		t.Errorf("Status code: expected 400, got %d", response.StatusCode)
	}
}

func TestLambdaHandlePostSuccess(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	repo.Create(context.Background())
	handler := NewLambdaHandler(repo)

	event := models.Event{
		Date:          time.Now(),
		StatusCode:    201,
		StatusMessage: "Created",
	}

	body, _ := json.Marshal(event)

	request := events.APIGatewayV2HTTPRequest{
		Body: string(body),
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/eventos",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 201 {
		t.Errorf("Status code: expected 201, got %d", response.StatusCode)
	}

	var respEvent models.Event
	json.Unmarshal([]byte(response.Body), &respEvent)

	if respEvent.Id == "" {
		t.Errorf("Response should have ID")
	}
}

func TestLambdaHandlePostInvalidJSON(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	request := events.APIGatewayV2HTTPRequest{
		Body: "invalid json",
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/eventos",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 400 {
		t.Errorf("Status code: expected 400, got %d", response.StatusCode)
	}
}

func TestLambdaHandlePutSuccess(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	repo.Create(context.Background())
	handler := NewLambdaHandler(repo)

	event := models.Event{
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "Updated",
	}

	body, _ := json.Marshal(event)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "123"},
		Body:           string(body),
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "PUT",
				Path:   "/eventos/123",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 201 {
		t.Errorf("Status code: expected 201, got %d", response.StatusCode)
	}

	var respEvent models.Event
	json.Unmarshal([]byte(response.Body), &respEvent)

	if respEvent.Id != "123" {
		t.Errorf("Event ID should be 123, got %s", respEvent.Id)
	}
}

func TestLambdaHandlePutMissingID(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	event := models.Event{
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "Updated",
	}

	body, _ := json.Marshal(event)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": ""},
		Body:           string(body),
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "PUT",
				Path:   "/eventos/",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 400 {
		t.Errorf("Status code: expected 400, got %d", response.StatusCode)
	}
}

func TestLambdaHandleDeleteSuccess(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	repo.Create(context.Background())
	handler := NewLambdaHandler(repo)

	event := &models.Event{
		Id:            "123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
	}
	repo.Save(context.Background(), event)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "123"},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "DELETE",
				Path:   "/eventos/123",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 204 {
		t.Errorf("Status code: expected 204, got %d", response.StatusCode)
	}
}

func TestLambdaHandleDeleteNotFound(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "nonexistent"},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "DELETE",
				Path:   "/eventos/nonexistent",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 404 {
		t.Errorf("Status code: expected 404, got %d", response.StatusCode)
	}
}

func TestLambdaHandleDeleteMissingID(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": ""},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "DELETE",
				Path:   "/eventos/",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 400 {
		t.Errorf("Status code: expected 400, got %d", response.StatusCode)
	}
}

func TestLambdaToJsonSuccess(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	data := models.Event{
		Id:            "123",
		StatusCode:    200,
		StatusMessage: "OK",
	}

	response, err := handler.toJson(data, 200)

	if err != nil {
		t.Errorf("toJson returned error: %v", err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Status code mismatch")
	}

	var result models.Event
	json.Unmarshal([]byte(response.Body), &result)

	if result.Id != "123" {
		t.Errorf("Event ID mismatch")
	}
}

func TestLambdaHandleAllMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		repo := repositories.NewMemoryDB(1 * time.Hour)
		handler := NewLambdaHandler(repo)

		request := events.APIGatewayV2HTTPRequest{
			PathParameters: map[string]string{"id": "test"},
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: method,
					Path:   "/eventos/test",
				},
			},
		}

		_, err := handler.HandleRequest(context.Background(), request)
		if err != nil {
			t.Errorf("HandleRequest for %s returned error: %v", method, err)
		}
	}
}

func TestLambdaHandleGetError(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "123"},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/eventos/123",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 404 {
		t.Errorf("Status code: expected 404, got %d", response.StatusCode)
	}
}

func TestLambdaHandlePostInvalidData(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	event := models.Event{
		StatusCode:    200,
		StatusMessage: "OK",
	}

	body, _ := json.Marshal(event)

	request := events.APIGatewayV2HTTPRequest{
		Body: string(body),
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/eventos",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 400 {
		t.Errorf("Status code: expected 400, got %d", response.StatusCode)
	}
}

func TestLambdaHandlePutInvalidJSON(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "123"},
		Body:           "invalid json",
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "PUT",
				Path:   "/eventos/123",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 400 {
		t.Errorf("Status code: expected 400, got %d", response.StatusCode)
	}
}

func TestLambdaHandlePutInvalidData(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	event := models.Event{
		StatusCode:    200,
		StatusMessage: "Updated",
	}

	body, _ := json.Marshal(event)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "123"},
		Body:           string(body),
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "PUT",
				Path:   "/eventos/123",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 400 {
		t.Errorf("Status code: expected 400, got %d", response.StatusCode)
	}
}

func TestLambdaHandlePutError(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	event := models.Event{
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "Updated",
	}

	body, _ := json.Marshal(event)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "123"},
		Body:           string(body),
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "PUT",
				Path:   "/eventos/123",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 201 {
		t.Errorf("Status code: expected 201, got %d", response.StatusCode)
	}
}

func TestLambdaHandleDeleteError(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	request := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"id": "123"},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "DELETE",
				Path:   "/eventos/123",
			},
		},
	}

	response, err := handler.HandleRequest(context.Background(), request)

	if err != nil {
		t.Errorf("HandleRequest returned error: %v", err)
	}
	if response.StatusCode != 404 {
		t.Errorf("Status code: expected 404, got %d", response.StatusCode)
	}
}

func TestLambdaToJsonMarshalError(t *testing.T) {
	repo := repositories.NewMemoryDB(1 * time.Hour)
	handler := NewLambdaHandler(repo)

	data := make(chan int)

	response, err := handler.toJson(data, 200)

	if err == nil {
		t.Errorf("toJson should return error for unmarshalable data")
	}
	if response.StatusCode != 500 {
		t.Errorf("Status code: expected 500, got %d", response.StatusCode)
	}
}
