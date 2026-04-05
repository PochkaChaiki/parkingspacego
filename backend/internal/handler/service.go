package handler

import (
	"context"

	"github.com/pochkachaiki/parkingspace/internal/model"
)

type Service interface {
	StartSession(ctx context.Context, r *model.RecordDto) (model.StatusName, error)
	GetSession(ctx context.Context, phone string) (*model.RecordDto, error)
	ProlongSession(ctx context.Context, phone string, duration string) (*model.RecordDto, error)
	StopSession(ctx context.Context, phone string) error
}
