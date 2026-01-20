package repositories

import (
	"api/models"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// MockDynamoDBClient is a mock implementation of DynamoDBClient
type MockDynamoDBClient struct {
	CreateTableFunc      func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	DescribeTableFunc    func(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
	UpdateTimeToLiveFunc func(ctx context.Context, params *dynamodb.UpdateTimeToLiveInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateTimeToLiveOutput, error)
	PutItemFunc          func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItemFunc       func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	GetItemFunc          func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	QueryFunc            func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

func (m *MockDynamoDBClient) CreateTable(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	if m.CreateTableFunc != nil {
		return m.CreateTableFunc(ctx, params, optFns...)
	}
	return nil, nil
}

func (m *MockDynamoDBClient) DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	if m.DescribeTableFunc != nil {
		return m.DescribeTableFunc(ctx, params, optFns...)
	}
	return nil, nil
}

func (m *MockDynamoDBClient) UpdateTimeToLive(ctx context.Context, params *dynamodb.UpdateTimeToLiveInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateTimeToLiveOutput, error) {
	if m.UpdateTimeToLiveFunc != nil {
		return m.UpdateTimeToLiveFunc(ctx, params, optFns...)
	}
	return nil, nil
}

func (m *MockDynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if m.PutItemFunc != nil {
		return m.PutItemFunc(ctx, params, optFns...)
	}
	return nil, nil
}

func (m *MockDynamoDBClient) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	if m.DeleteItemFunc != nil {
		return m.DeleteItemFunc(ctx, params, optFns...)
	}
	return nil, nil
}

func (m *MockDynamoDBClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if m.GetItemFunc != nil {
		return m.GetItemFunc(ctx, params, optFns...)
	}
	return nil, nil
}

func (m *MockDynamoDBClient) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, params, optFns...)
	}
	return nil, nil
}

// Create tests
func TestDynamoDBCreateSuccess(t *testing.T) {
	callCount := 0
	client := &MockDynamoDBClient{
		CreateTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
			callCount++
			return &dynamodb.CreateTableOutput{}, nil
		},
		DescribeTableFunc: func(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
			return &dynamodb.DescribeTableOutput{
				Table: &types.TableDescription{
					TableStatus: types.TableStatusActive,
				},
			}, nil
		},
		UpdateTimeToLiveFunc: func(ctx context.Context, params *dynamodb.UpdateTimeToLiveInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateTimeToLiveOutput, error) {
			return &dynamodb.UpdateTimeToLiveOutput{}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	err := repo.Create(context.Background())

	if err != nil {
		t.Errorf("Create() returned unexpected error: %v", err)
	}
	if callCount == 0 {
		t.Errorf("CreateTable was not called")
	}
}

func TestDynamoDBCreateTableAlreadyExists(t *testing.T) {
	client := &MockDynamoDBClient{
		CreateTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
			return nil, fmt.Errorf("already exists")
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	err := repo.Create(context.Background())

	if err != nil {
		t.Errorf("Create() should not return error for 'already exists', got: %v", err)
	}
}

func TestDynamoDBCreateTableError(t *testing.T) {
	client := &MockDynamoDBClient{
		CreateTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
			return nil, fmt.Errorf("connection error")
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	err := repo.Create(context.Background())

	if err == nil {
		t.Errorf("Create() should return error for connection failure")
	}
}

func TestDynamoDBCreateUpdateTimeToLiveError(t *testing.T) {
	client := &MockDynamoDBClient{
		CreateTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
			return &dynamodb.CreateTableOutput{}, nil
		},
		DescribeTableFunc: func(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
			return &dynamodb.DescribeTableOutput{
				Table: &types.TableDescription{
					TableStatus: types.TableStatusActive,
				},
			}, nil
		},
		UpdateTimeToLiveFunc: func(ctx context.Context, params *dynamodb.UpdateTimeToLiveInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateTimeToLiveOutput, error) {
			return nil, fmt.Errorf("TTL error")
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	err := repo.Create(context.Background())

	if err == nil {
		t.Errorf("Create() should return error when UpdateTimeToLive fails")
	}
}

// Save tests
func TestDynamoDBSaveSuccess(t *testing.T) {
	client := &MockDynamoDBClient{
		PutItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			if params.TableName == nil || *params.TableName != "events" {
				t.Errorf("Table name mismatch")
			}
			return &dynamodb.PutItemOutput{}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	event := &models.Event{
		Id:            "test-123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
	}

	err := repo.Save(context.Background(), event)

	if err != nil {
		t.Errorf("Save() returned error: %v", err)
	}
	if event.Expiration == 0 {
		t.Errorf("Save() should set expiration")
	}
}

func TestDynamoDBSaveWithExistingExpiration(t *testing.T) {
	client := &MockDynamoDBClient{
		PutItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	futureExpiration := time.Now().Add(2 * time.Hour).Unix()
	event := &models.Event{
		Id:            "test-123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
		Expiration:    futureExpiration,
	}

	repo.Save(context.Background(), event)

	if event.Expiration != futureExpiration {
		t.Errorf("Save() should preserve existing expiration")
	}
}

func TestDynamoDBSavePutItemError(t *testing.T) {
	client := &MockDynamoDBClient{
		PutItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, fmt.Errorf("put item failed")
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	event := &models.Event{
		Id:            "test-123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
	}

	err := repo.Save(context.Background(), event)

	if err == nil {
		t.Errorf("Save() should return error when PutItem fails")
	}
}

// Delete tests
func TestDynamoDBDeleteSuccess(t *testing.T) {
	deletedEvent := &models.Event{
		Id:            "test-123",
		Date:          time.Now(),
		StatusCode:    200,
		StatusMessage: "OK",
	}

	item, _ := attributevalue.MarshalMap(deletedEvent)

	client := &MockDynamoDBClient{
		DeleteItemFunc: func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
			if params.ReturnValues != types.ReturnValueAllOld {
				t.Errorf("ReturnValues should be ReturnValueAllOld")
			}
			return &dynamodb.DeleteItemOutput{
				Attributes: item,
			}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	event, err := repo.Delete(context.Background(), "test-123")

	if err != nil {
		t.Errorf("Delete() returned error: %v", err)
	}
	if event == nil {
		t.Errorf("Delete() should return the deleted event")
	}
	if event.Id != "test-123" {
		t.Errorf("Deleted event ID mismatch")
	}
}

func TestDynamoDBDeleteNotFound(t *testing.T) {
	client := &MockDynamoDBClient{
		DeleteItemFunc: func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
			return &dynamodb.DeleteItemOutput{
				Attributes: nil,
			}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	event, err := repo.Delete(context.Background(), "nonexistent")

	if err != nil {
		t.Errorf("Delete() returned error: %v", err)
	}
	if event != nil {
		t.Errorf("Delete() should return nil for nonexistent item")
	}
}

func TestDynamoDBDeleteError(t *testing.T) {
	client := &MockDynamoDBClient{
		DeleteItemFunc: func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
			return nil, fmt.Errorf("delete failed")
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	_, err := repo.Delete(context.Background(), "test-123")

	if err == nil {
		t.Errorf("Delete() should return error")
	}
}

// Get tests
func TestDynamoDBGetSuccess(t *testing.T) {
	now := time.Now()
	getEvent := &models.Event{
		Id:            "test-123",
		Date:          now,
		StatusCode:    200,
		StatusMessage: "OK",
	}

	item, _ := attributevalue.MarshalMap(getEvent)

	client := &MockDynamoDBClient{
		GetItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{
				Item: item,
			}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	event, err := repo.Get(context.Background(), "test-123")

	if err != nil {
		t.Errorf("Get() returned error: %v", err)
	}
	if event == nil {
		t.Errorf("Get() should return event")
	}
	if event.Id != "test-123" {
		t.Errorf("Event ID mismatch")
	}
}

func TestDynamoDBGetNotFound(t *testing.T) {
	client := &MockDynamoDBClient{
		GetItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	event, err := repo.Get(context.Background(), "nonexistent")

	if err != nil {
		t.Errorf("Get() returned error: %v", err)
	}
	if event != nil {
		t.Errorf("Get() should return nil for nonexistent item")
	}
}

func TestDynamoDBGetError(t *testing.T) {
	client := &MockDynamoDBClient{
		GetItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return nil, fmt.Errorf("get failed")
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	_, err := repo.Get(context.Background(), "test-123")

	if err == nil {
		t.Errorf("Get() should return error")
	}
}

// FindByDateAndReturnCode tests
func TestDynamoDBFindByDateAndReturnCodeSuccess(t *testing.T) {
	now := time.Now()
	event1 := &models.Event{
		Id:            "1",
		Date:          now,
		StatusCode:    200,
		StatusMessage: "OK",
	}
	event2 := &models.Event{
		Id:            "2",
		Date:          now.Add(10 * time.Minute),
		StatusCode:    200,
		StatusMessage: "OK",
	}

	item1, _ := attributevalue.MarshalMap(event1)
	item2, _ := attributevalue.MarshalMap(event2)

	client := &MockDynamoDBClient{
		QueryFunc: func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{item1, item2},
			}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	from := now.Add(-1 * time.Hour)
	to := now.Add(1 * time.Hour)

	events, err := repo.FindByDateAndReturnCode(context.Background(), from, to, 200)

	if err != nil {
		t.Errorf("FindByDateAndReturnCode() returned error: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
}

func TestDynamoDBFindByDateAndReturnCodeEmpty(t *testing.T) {
	client := &MockDynamoDBClient{
		QueryFunc: func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{},
			}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	now := time.Now()

	events, err := repo.FindByDateAndReturnCode(context.Background(), now.Add(-1*time.Hour), now.Add(1*time.Hour), 200)

	if err != nil {
		t.Errorf("FindByDateAndReturnCode() returned error: %v", err)
	}
	if len(events) > 0 {
		t.Errorf("Expected empty results")
	}
}

func TestDynamoDBFindByDateAndReturnCodeError(t *testing.T) {
	client := &MockDynamoDBClient{
		QueryFunc: func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return nil, fmt.Errorf("query failed")
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	now := time.Now()

	_, err := repo.FindByDateAndReturnCode(context.Background(), now.Add(-1*time.Hour), now.Add(1*time.Hour), 200)

	if err == nil {
		t.Errorf("FindByDateAndReturnCode() should return error")
	}
}

func TestDynamoDBFindByDateAndReturnCodeUnmarshalError(t *testing.T) {
	client := &MockDynamoDBClient{
		QueryFunc: func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{
					{
						"id":   &types.AttributeValueMemberS{Value: "1"},
						"date": &types.AttributeValueMemberN{Value: "invalid"},
					},
				},
			}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	now := time.Now()

	_, err := repo.FindByDateAndReturnCode(context.Background(), now.Add(-1*time.Hour), now.Add(1*time.Hour), 200)

	if err == nil {
		t.Errorf("FindByDateAndReturnCode() should return error for unmarshal failure")
	}
}

// Edge cases and additional coverage
func TestDynamoDBRepositoryStructure(t *testing.T) {
	client := &MockDynamoDBClient{}
	repo := NewDynamoDBRepository(client, "table", 1*time.Hour)

	if repo.client == nil {
		t.Errorf("Client should not be nil")
	}
	if repo.name == "" {
		t.Errorf("Name should not be empty")
	}
	if repo.ttl <= 0 {
		t.Errorf("TTL should be positive")
	}
}

func TestDynamoDBSaveWithZeroTTL(t *testing.T) {
	client := &MockDynamoDBClient{
		PutItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 0)
	event := &models.Event{
		Id:         "test-123",
		Date:       time.Now(),
		StatusCode: 200,
	}

	err := repo.Save(context.Background(), event)

	if err != nil {
		t.Errorf("Save() should work with zero TTL")
	}
	if event.Expiration == 0 {
		t.Errorf("Save() should set expiration even with zero TTL")
	}
}

func TestDynamoDBFindByDateAndReturnCodeParameterValidation(t *testing.T) {
	queryUsed := false
	client := &MockDynamoDBClient{
		QueryFunc: func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			queryUsed = true
			if params.IndexName == nil || *params.IndexName != "date-statusCode-index" {
				t.Errorf("IndexName mismatch")
			}
			if params.KeyConditionExpression == nil {
				t.Errorf("KeyConditionExpression should not be nil")
			}
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, nil
		},
	}

	repo := NewDynamoDBRepository(client, "events", 1*time.Hour)
	now := time.Now()
	repo.FindByDateAndReturnCode(context.Background(), now.Add(-1*time.Hour), now.Add(1*time.Hour), 404)

	if !queryUsed {
		t.Errorf("Query was not called with proper parameters")
	}
}
