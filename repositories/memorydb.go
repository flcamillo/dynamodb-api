package repositories

import (
	"api/models"
	"context"
	"time"
)

// Define a estrutura do repositório de memória.
type MemoryDB struct {
	db  []*models.Event
	ttl time.Duration
}

// Cria uma nova instância do repositório de memória.
func NewMemoryDB(ttl time.Duration) *MemoryDB {
	return &MemoryDB{
		db:  make([]*models.Event, 0),
		ttl: ttl,
	}
}

// Cria o repositório (não faz nada no caso do MemoryDB).
func (p *MemoryDB) Create(ctx context.Context) error {
	return nil
}

// Salva o registro na memória.
// Se já houver registro com o mesmo id, ele será substituído.
func (p *MemoryDB) Save(ctx context.Context, event *models.Event) error {
	if event.Expiration == 0 && p.ttl > 0 {
		event.Expiration = time.Now().Add(p.ttl).Unix()
	}
	for k, v := range p.db {
		if v.Id == event.Id {
			p.db[k] = event
			return nil
		}
	}
	p.db = append(p.db, event)
	return nil
}

// Deleta o registro da memória pelo id.
func (p *MemoryDB) Delete(ctx context.Context, id string) (event *models.Event, err error) {
	for k, v := range p.db {
		if v.Id == id {
			p.db = append(p.db[:k], p.db[k+1:]...)
			return v, nil
		}
	}
	return nil, nil
}

// Recupera o registro da memória pelo id.
func (p *MemoryDB) Get(ctx context.Context, id string) (event *models.Event, err error) {
	for _, v := range p.db {
		if v.Expiration != 0 && time.Unix(v.Expiration, 0).Before(time.Now()) {
			continue
		}
		if v.Id == id {
			return v, nil
		}
	}
	return nil, nil
}

// Encontra registros pela data e código de retorno.
func (p *MemoryDB) FindByDateAndReturnCode(ctx context.Context, from time.Time, to time.Time, statusCode int) (events []*models.Event, err error) {
	for _, v := range p.db {
		if v.Expiration != 0 && time.Unix(v.Expiration, 0).Before(time.Now()) {
			continue
		}
		if (v.Date.After(from) || v.Date.Equal(from)) && (v.Date.Before(to) || v.Date.Equal(to)) && v.StatusCode == statusCode {
			events = append(events, v)
		}
	}
	return events, nil
}
