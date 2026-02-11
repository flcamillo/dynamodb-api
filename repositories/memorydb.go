package repositories

import (
	"api/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Define a configuração do repositório de memória.
type MemoryDBConfig struct {
	// tempo de expiração dos registros
	TTL time.Duration
}

// Define a estrutura do repositório de memória.
type MemoryDB struct {
	// banco de dados em memória
	db map[string]*models.Event
	// configuração do repositório
	config *MemoryDBConfig
	// configura o tracer
	tracer trace.Tracer
}

// Cria uma nova instância do repositório de memória.
func NewMemoryDB(config *MemoryDBConfig) *MemoryDB {
	return &MemoryDB{
		db:     make(map[string]*models.Event),
		config: config,
		tracer: otel.Tracer("memorydb.repository"),
	}
}

// Cria um span contextualizado para o banco de dados de memória.
func (p *MemoryDB) newSpan(ctx context.Context, operation string, statement string) (context.Context, trace.Span) {
	ctx, span := p.tracer.Start(
		ctx,
		operation,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "memorydb"),
			attribute.String("db.name", "events"),
			attribute.String("db.operation", operation),
		),
	)
	if statement != "" {
		span.SetAttributes(attribute.String("db.statement", statement))
	}
	return ctx, span
}

// Cria o repositório (não faz nada no caso do MemoryDB).
func (p *MemoryDB) Create(ctx context.Context) error {
	return nil
}

// Salva o registro na memória.
// Se já houver registro com o mesmo id, ele será substituído.
func (p *MemoryDB) Save(ctx context.Context, event *models.Event) error {
	ctx, span := p.newSpan(ctx, "save", "")
	defer span.End()
	if event.Expiration == 0 && p.config.TTL > 0 {
		event.Expiration = time.Now().Add(p.config.TTL).Unix()
	}
	if event.Id == "" {
		event.Id = uuid.New().String()
	}
	p.db[event.Id] = event
	return nil
}

// Deleta o registro da memória pelo id.
func (p *MemoryDB) Delete(ctx context.Context, id string) (event *models.Event, err error) {
	ctx, span := p.newSpan(ctx, "delete", "id = "+id)
	defer span.End()
	event, ok := p.db[id]
	if !ok {
		span.AddEvent("record not found")
		return nil, nil
	}
	delete(p.db, id)
	return event, nil
}

// Recupera o registro da memória pelo id.
func (p *MemoryDB) Get(ctx context.Context, id string) (event *models.Event, err error) {
	ctx, span := p.newSpan(ctx, "get", "id = "+id)
	defer span.End()
	event, ok := p.db[id]
	if !ok {
		span.AddEvent("record not found")
		return nil, nil
	}
	if event.Expiration != 0 && time.Unix(event.Expiration, 0).Before(time.Now()) {
		delete(p.db, id)
		return nil, nil
	}
	return event, nil
}

// Encontra registros pela data e código de retorno.
func (p *MemoryDB) FindByDateAndReturnCode(ctx context.Context, from time.Time, to time.Time, statusCode int) (events []*models.Event, err error) {
	ctx, span := p.newSpan(ctx, "query", fmt.Sprintf("from = %s to = %s statusCode = %d", from.Format(time.RFC3339), to.Format(time.RFC3339), statusCode))
	defer span.End()
	expired := make([]string, 0)
	for k, v := range p.db {
		if v.Expiration != 0 && time.Unix(v.Expiration, 0).Before(time.Now()) {
			expired = append(expired, k)
			continue
		}
		if (v.Date.After(from) || v.Date.Equal(from)) && (v.Date.Before(to) || v.Date.Equal(to)) && v.StatusCode == statusCode {
			events = append(events, v)
		}
	}
	for _, id := range expired {
		delete(p.db, id)
	}
	return events, nil
}
