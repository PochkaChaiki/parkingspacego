package service

import (
	"context"
	"time"

	"github.com/pochkachaiki/parkingspace/internal/model"
)

type Repository interface {
	Create(ctx context.Context, rec *model.Record) (*model.Record, error)
	GetAll(ctx context.Context) ([]*model.Record, error)
	GetByPhone(ctx context.Context, phone string) (*model.Record, error)
	Update(ctx context.Context, id string, endTime time.Time) error
	DeleteByPhone(ctx context.Context, phone string) error
}

type ParkingService struct {
	repo Repository
}

func NewService(repo Repository) *ParkingService {
	return &ParkingService{repo: repo}
}

// StartSession создает новую сессию парковки.
func (p *ParkingService) StartSession(ctx context.Context, r *model.RecordDto) (model.StatusName, error) {
	if r == nil {
		return model.Failure, nil
	}

	// Проверяем, есть ли уже сессия для этого номера телефона
	existing, err := p.repo.GetByPhone(ctx, r.PhoneNumber)
	if err != nil {
		return model.Failure, err
	}
	if existing != nil {
		// Пользователь уже имеет активную сессию - возвращаем Failure БЕЗ ошибки
		return model.Failure, nil
	}

	duration := "1h"
	if r.Duration != nil {
		duration = *r.Duration
	}

	timeToAdd, err := time.ParseDuration(duration)
	if err != nil {
		return model.Failure, err
	}

	// Создаем новую запись
	rec := &model.Record{
		ClientName:   r.ClientName,
		PhoneNumber:  r.PhoneNumber,
		LicensePlate: r.LicensePlate,
		SpotNumber:   r.SpotNumber,
		StartTime:    time.Now().UTC(),
		EndTime:      time.Now().UTC().Add(timeToAdd),
		Status:       "active",
	}

	// // Вычисляем EndTime на основе Duration (если указана)
	// if r.Duration != nil && *r.Duration != "" {
	// 	// Парсим duration (например "1h")
	// 	duration, err := time.ParseDuration(*r.Duration)
	// 	if err == nil {
	// 		endTime := rec.StartTime.Add(duration)
	// 		rec.EndTime = endTime
	// 	}
	// }

	// Проверяем, свободно ли место
	all, err := p.repo.GetAll(ctx)
	if err != nil {
		return model.Failure, err
	}

	for _, check := range all {
		if check.SpotNumber == r.SpotNumber && check.Status == "active" {
			// Место занято - возвращаем Occupied БЕЗ ошибки
			return model.Occupied, nil
		}
	}

	// Создаем запись в репозитории
	_, err = p.repo.Create(ctx, rec)
	if err != nil {
		return model.Failure, err
	}

	return model.Success, nil
}

// GetSession возвращает информацию о сессии парковки по номеру телефона.
func (p *ParkingService) GetSession(ctx context.Context, phone string) (*model.RecordDto, error) {
	rec, err := p.repo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		// Не найдена - возвращаем nil БЕЗ ошибки
		return nil, nil
	}

	// Преобразуем Model.Record в RecordDto
	dto := &model.RecordDto{
		ClientName:   rec.ClientName,
		PhoneNumber:  rec.PhoneNumber,
		LicensePlate: rec.LicensePlate,
		SpotNumber:   rec.SpotNumber,
		StartTime:    &rec.StartTime,
		EndTime:      &rec.EndTime,
	}

	return dto, nil
}

// ProlongSession продлевает сессию парковки.
// Red: Согласно тестам, метод должен увеличить EndTime на указанную Duration
func (p *ParkingService) ProlongSession(ctx context.Context, phone string, duration string) (*model.RecordDto, error) {
	rec, err := p.repo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, nil
	}

	// Иначе добавляем Duration к существующей EndTime
	timeToAdd, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	rec.EndTime = rec.EndTime.Add(timeToAdd)

	// Обновляем запись в репозитории
	err = p.repo.Update(ctx, rec.ID.Hex(), rec.EndTime)
	if err != nil {
		return nil, err
	}

	// Преобразуем обновленный Model.Record в RecordDto
	dto := &model.RecordDto{
		ClientName:   rec.ClientName,
		PhoneNumber:  rec.PhoneNumber,
		LicensePlate: rec.LicensePlate,
		SpotNumber:   rec.SpotNumber,
		StartTime:    &rec.StartTime,
		EndTime:      &rec.EndTime,
	}

	return dto, nil
}

// StopSession завершает сессию парковки и удаляет её.
// Red: Согласно тестам, метод должен удалить сессию по номеру телефона
func (p *ParkingService) StopSession(ctx context.Context, phone string) error {
	return p.repo.DeleteByPhone(ctx, phone)
}
