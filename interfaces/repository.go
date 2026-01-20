package interfaces

import (
	"api/models"
	"context"
	"time"
)

// Define a interface do repositório.
type Repository interface {
	Create(ctx context.Context) error
	Save(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, id string) (*models.Event, error)
	Get(ctx context.Context, id string) (*models.Event, error)
	FindByDateAndReturnCode(ctx context.Context, from time.Time, to time.Time, statusCode int) ([]*models.Event, error)
}
