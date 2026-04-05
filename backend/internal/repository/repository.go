package repository

import (
	"context"
	"time"

	"github.com/pochkachaiki/parkingspace/internal/model"
)

// Repository описывает операции с данными парковочных записей.
// Интерфейс позволяет мокировать хранилище данных в тестах.
type Repository interface {
	Create(ctx context.Context, rec *model.Record) (*model.Record, error)
	GetAll(ctx context.Context) ([]*model.Record, error)
	GetByPhone(ctx context.Context, phone string) (*model.Record, error)
	// ListByPhone(ctx context.Context, phone string) ([]*model.Record, error)
	Update(ctx context.Context, id string, endTime time.Time) error
	DeleteByPhone(ctx context.Context, phone string) error
}
