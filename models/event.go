package models

import (
	"fmt"
	"time"
)

// Define a estrutura do registro na tabela DynamoDB.
type Event struct {
	Id            string            `json:"id,omitempty" dynamodbav:"id"`
	Date          time.Time         `json:"date" dynamodbav:"date"`
	StatusCode    int               `json:"statusCode" dynamodbav:"statusCode"`
	StatusMessage string            `json:"statusMessage" dynamodbav:"statusMessage"`
	Expiration    int64             `json:"expiration" dynamodbav:"expiration"`
	Metadata      map[string]string `json:"metadata,omitempty" dynamodbav:"metadata,omitempty"`
}

// Valida os campos do registro.
func (e *Event) Validate() error {
	if e.Date.IsZero() {
		return fmt.Errorf("invalid date")
	}
	if e.StatusCode < 0 {
		return fmt.Errorf("invalid status code")
	}
	return nil
}
