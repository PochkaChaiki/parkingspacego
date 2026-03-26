package handler

import (
	"context"
	"time"

	"github.com/pochkachaiki/parkingspace/internal/model"
)

type Service interface {
	StartSession(ctx context.Context, r *model.RecordDto) (*model.Response, error)
	GetSession(ctx context.Context, phone string) (*model.RecordDto, error)
	ProlongSession(ctx context.Context, phone string, duration time.Duration) (*model.RecordDto, error)
	StopSession(ctx context.Context, phone string) error
}
