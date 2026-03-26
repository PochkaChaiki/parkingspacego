package service

import (
	"context"
	"testing"
	"time"

	"github.com/pochkachaiki/parkingspace/internal/model"
)

// mockRepository имитирует Repository для тестирования сервиса
type mockRepository struct {
	records map[string]*model.Record // по phone_number
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		records: make(map[string]*model.Record),
	}
}

func (m *mockRepository) Create(ctx context.Context, rec *model.Record) (*model.Record, error) {
	m.records[rec.PhoneNumber] = rec
	return rec, nil
}

func (m *mockRepository) GetAll(ctx context.Context) ([]*model.Record, error) {
	var all []*model.Record
	for _, rec := range m.records {
		all = append(all, rec)
	}
	return all, nil
}

func (m *mockRepository) GetByPhone(ctx context.Context, phone string) (*model.Record, error) {
	if rec, exists := m.records[phone]; exists {
		return rec, nil
	}
	return nil, nil
}

func (m *mockRepository) Update(ctx context.Context, id string, endTime *time.Time) error {
	for _, rec := range m.records {
		if rec.ID.Hex() == id {
			rec.EndTime = endTime
			rec.Status = "completed"
			return nil
		}
	}
	return nil
}

func (m *mockRepository) DeleteByPhone(ctx context.Context, phone string) error {
	delete(m.records, phone)
	return nil
}

// ============= Tests (RED/GREEN TDD) =============

// TestServiceStartSession - тест создания новой сессии парковки
func TestServiceStartSession(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()
	srv := NewService(repo)

	req := &model.RecordDto{
		ClientName:   "Иван",
		PhoneNumber:  "+79991234567",
		LicensePlate: "A123BC140",
		SpotNumber:   42,
	}

	resp, err := srv.StartSession(ctx, req)
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	if resp == nil {
		t.Fatalf("expected response, got nil")
	}

	if resp.Status != model.Success {
		t.Fatalf("expected status 'success', got %q", resp.Status)
	}

	// Проверяем что запись создана в репо
	rec, err := repo.GetByPhone(ctx, "+79991234567")
	if err != nil {
		t.Fatalf("GetByPhone failed: %v", err)
	}

	if rec == nil {
		t.Fatalf("expected record in repository")
	}

	if rec.ClientName != "Иван" {
		t.Fatalf("expected ClientName 'Иван', got %q", rec.ClientName)
	}
}

// TestServiceStartSessionDuplicate - тест когда пользователь уже имеет активную сессию
func TestServiceStartSessionDuplicate(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()
	srv := NewService(repo)

	phone := "+79991234567"

	// Первая сессия
	req1 := &model.RecordDto{
		ClientName:   "Иван",
		PhoneNumber:  phone,
		LicensePlate: "A123BC140",
		SpotNumber:   42,
	}

	_, err := srv.StartSession(ctx, req1)
	if err != nil {
		t.Fatalf("first StartSession failed: %v", err)
	}

	// Вторая сессия с тем же номером
	req2 := &model.RecordDto{
		ClientName:   "Иван",
		PhoneNumber:  phone,
		LicensePlate: "A123BC141",
		SpotNumber:   43,
	}

	resp, err := srv.StartSession(ctx, req2)
	if err != nil {
		t.Fatalf("second StartSession failed: %v", err)
	}

	if resp.Status != model.Failure {
		t.Fatalf("expected status 'failure' for duplicate, got %q", resp.Status)
	}
}

// TestServiceGetSession - тест получения информации о сессии
func TestServiceGetSession(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()
	srv := NewService(repo)

	phone := "+79991234567"

	// Создаем сессию
	req := &model.RecordDto{
		ClientName:   "Иван",
		PhoneNumber:  phone,
		LicensePlate: "A123BC140",
		SpotNumber:   42,
	}

	_, err := srv.StartSession(ctx, req)
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	// Получаем информацию
	dto, err := srv.GetSession(ctx, phone)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}

	if dto == nil {
		t.Fatalf("expected RecordDto, got nil")
	}

	if dto.ClientName != "Иван" {
		t.Fatalf("expected ClientName 'Иван', got %q", dto.ClientName)
	}

	if dto.PhoneNumber != phone {
		t.Fatalf("expected phone %q, got %q", phone, dto.PhoneNumber)
	}
}

// TestServiceProlongSession - тест продления сессии парковки
func TestServiceProlongSession(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()
	srv := NewService(repo)

	phone := "+79991234567"

	// Создаем сессию
	req := &model.RecordDto{
		ClientName:   "Иван",
		PhoneNumber:  phone,
		LicensePlate: "A123BC140",
		SpotNumber:   42,
	}

	_, err := srv.StartSession(ctx, req)
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	// Получаем начальную сессию
	initialDto, err := srv.GetSession(ctx, phone)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}

	initialEndTime := initialDto.EndTime

	// Продляем на 1 час
	duration := time.Hour
	updatedDto, err := srv.ProlongSession(ctx, phone, duration)
	if err != nil {
		t.Fatalf("ProlongSession failed: %v", err)
	}

	if updatedDto == nil {
		t.Fatalf("expected RecordDto, got nil")
	}

	if updatedDto.EndTime == nil {
		t.Fatalf("expected EndTime to be set")
	}

	// Если было начальное EndTime, проверяем что оно увеличено
	if initialEndTime != nil {
		expectedEndTime := initialEndTime.Add(duration)
		if !updatedDto.EndTime.Equal(expectedEndTime) {
			t.Fatalf("expected EndTime ~%v, got %v", expectedEndTime, *updatedDto.EndTime)
		}
	}
}

// TestServiceStopSession - тест завершения сессии парковки
func TestServiceStopSession(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()
	srv := NewService(repo)

	phone := "+79991234567"

	// Создаем сессию
	req := &model.RecordDto{
		ClientName:   "Иван",
		PhoneNumber:  phone,
		LicensePlate: "A123BC140",
		SpotNumber:   42,
	}

	_, err := srv.StartSession(ctx, req)
	if err != nil {
		t.Fatalf("StartSession failed: %v", err)
	}

	// Проверяем что она существует
	dto, err := srv.GetSession(ctx, phone)
	if err != nil || dto == nil {
		t.Fatalf("session should exist before stop")
	}

	// Завершаем сессию
	err = srv.StopSession(ctx, phone)
	if err != nil {
		t.Fatalf("StopSession failed: %v", err)
	}

	// Проверяем что она удалена
	dto, err = srv.GetSession(ctx, phone)
	if err != nil {
		t.Fatalf("GetSession after stop should not error: %v", err)
	}

	if dto != nil {
		t.Fatalf("expected session to be deleted, got %v", dto)
	}
}
