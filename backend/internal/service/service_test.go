package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pochkachaiki/parkingspace/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockRepository имитирует Repository для тестирования сервиса
type mockRepository struct {
	records map[string]*model.Record // по phone_number
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		records: map[string]*model.Record{
			"+79999999999": {
				ID:           primitive.NewObjectID(),
				ClientName:   "Egor",
				PhoneNumber:  "+79999999999",
				LicensePlate: "A123BC123",
				SpotNumber:   1,
				StartTime:    time.Now().UTC(),
				EndTime:      time.Now().UTC().Add(time.Hour),
				Status:       "active",
			},
			"+79999999998": {
				ID:           primitive.NewObjectID(),
				ClientName:   "Arisha",
				PhoneNumber:  "+79999999998",
				LicensePlate: "A123BC122",
				SpotNumber:   2,
				StartTime:    time.Now().UTC(),
				EndTime:      time.Now().UTC().Add(time.Hour),
				Status:       "active",
			},
			"+79999999997": {
				ID:           primitive.NewObjectID(),
				ClientName:   "Kirill",
				PhoneNumber:  "+79999999997",
				LicensePlate: "A123BC121",
				SpotNumber:   3,
				StartTime:    time.Now().UTC(),
				EndTime:      time.Now().UTC().Add(time.Hour),
				Status:       "active",
			},
		},
	}
}

func (m *mockRepository) Create(ctx context.Context, rec *model.Record) (*model.Record, error) {
	if rec.PhoneNumber == "00000000000" {
		return nil, errors.New("internal error")
	}

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
	if phone == "00000000000" {
		return nil, errors.New("internal error")
	}

	if rec, exists := m.records[phone]; exists {
		return rec, nil
	}
	return nil, nil
}

func (m *mockRepository) Update(ctx context.Context, id string, endTime time.Time) error {
	if id == "" {
		return errors.New("internal error")
	}
	for _, rec := range m.records {
		if rec.ID.Hex() == id {
			rec.EndTime = endTime
			rec.Status = "active"
			return nil
		}
	}
	return nil
}

func (m *mockRepository) DeleteByPhone(ctx context.Context, phone string) error {
	if phone == "00000000000" {
		return errors.New("internal error")
	}
	delete(m.records, phone)
	return nil
}

// ============= Tests (RED/GREEN TDD) =============

func TestService_StartSession(t *testing.T) {
	repo := newMockRepository()
	srv := NewService(repo)

	tests := []struct {
		name       string
		dto        *model.RecordDto
		wantCode   int
		wantStatus model.StatusName
		wantErr    bool
	}{
		{
			name: "creating session success",
			dto: &model.RecordDto{
				ClientName:   "egor",
				PhoneNumber:  "+78888888888",
				LicensePlate: "A321BC321",
				SpotNumber:   100,
			},
			wantStatus: model.Success,
			wantErr:    false,
		},
		{
			name: "session exist",
			dto: &model.RecordDto{
				ClientName:   "egor",
				PhoneNumber:  "+78888888888",
				LicensePlate: "A321BC321",
				SpotNumber:   101,
			},
			wantStatus: model.Failure,
			wantErr:    false,
		},
		{
			name: "spot occupied",
			dto: &model.RecordDto{
				ClientName:   "egor",
				PhoneNumber:  "+78888888881",
				LicensePlate: "A321BC321",
				SpotNumber:   100,
			},
			wantStatus: model.Occupied,
			wantErr:    false,
		},
		{
			name:       "json was invalid",
			dto:        nil,
			wantStatus: model.Failure,
			wantErr:    true,
		},
		{
			name: "internal error",
			dto: &model.RecordDto{
				ClientName:   "egor",
				PhoneNumber:  "00000000000",
				LicensePlate: "A321BC321",
				SpotNumber:   100,
			},
			wantStatus: model.Failure,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			status, err := srv.StartSession(ctx, tt.dto)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("error was not expected")
					return
				}
			}

			if status != tt.wantStatus {
				t.Fatalf("response status want %s, got %s", tt.wantStatus, status)
				return
			}

			if status == model.Success {
				rec, _ := repo.GetByPhone(ctx, tt.dto.PhoneNumber)
				if rec == nil {
					t.Fatalf("session did not started")
					return
				}
			}
		})
	}
}

func TestService_GetSession(t *testing.T) {
	repo := newMockRepository()
	srv := NewService(repo)

	tests := []struct {
		name    string
		phone   string
		wantDto *model.RecordDto
		wantErr bool
	}{
		{
			name:  "get session successful",
			phone: "+79999999999",
			wantDto: &model.RecordDto{
				ClientName:   "Egor",
				PhoneNumber:  "+79999999999",
				LicensePlate: "A123BC123",
				SpotNumber:   1,
			},
			wantErr: false,
		},
		{
			name:    "session don't exist",
			phone:   "+78888888888",
			wantDto: nil,
			wantErr: false,
		},
		{
			name:    "internal error",
			phone:   "00000000000",
			wantDto: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			resp, err := srv.GetSession(ctx, tt.phone)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("error was not expected")
					return
				}
			}

			if tt.wantDto == nil {
				if resp != nil {
					t.Fatalf("expected nil got not nil")
				}
				return
			} else {
				if resp == nil {
					t.Fatalf("expected not nil got nil")
					return
				}
			}

			if !(resp.ClientName == tt.wantDto.ClientName &&
				resp.SpotNumber == tt.wantDto.SpotNumber &&
				resp.PhoneNumber == tt.wantDto.PhoneNumber &&
				resp.LicensePlate == tt.wantDto.LicensePlate) {
				t.Fatalf("dto want: %v, got: %v", *tt.wantDto, resp)
				return
			}

		})
	}
}

func TestService_ProlongSession(t *testing.T) {
	repo := newMockRepository()
	srv := NewService(repo)

	tests := []struct {
		name     string
		phone    string
		duration string
		wantDto  bool
		wantErr  bool
	}{
		{
			name:     "prolong session successful",
			phone:    "+79999999999",
			duration: "1h",
			wantDto:  true,
			wantErr:  false,
		},
		{
			name:     "session don't exist",
			phone:    "+78888888888",
			duration: "1h",
			wantDto:  false,
			wantErr:  false,
		},
		{
			name:     "internal error",
			phone:    "00000000000",
			duration: "1h",
			wantDto:  false,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			record, err := repo.GetByPhone(ctx, tt.phone)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("error occured: %v", err)
				}
				return
			}
			var resEndTime time.Time
			if record != nil {
				dur, _ := time.ParseDuration(tt.duration)
				resEndTime = record.EndTime.Add(dur)
			}
			resp, err := srv.ProlongSession(ctx, tt.phone, tt.duration)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("error was not expected: %v", err)
					return
				}
			}

			if resp == nil {
				if tt.wantDto {
					t.Fatalf("expected resp but got nil")
					return
				}
				return
			}

			if !resp.EndTime.Equal(resEndTime) {
				t.Fatalf("endTime want: %v, got: %v", resEndTime, *resp.EndTime)
				return
			}

		})
	}
}

func TestService_StopSession(t *testing.T) {
	repo := newMockRepository()
	srv := NewService(repo)

	tests := []struct {
		name    string
		phone   string
		wantErr bool
		// wantDto  *model.RecordDto
	}{
		{
			name:    "stop session successful",
			phone:   "+79999999999",
			wantErr: false,
		},
		{
			name:    "session don't exist",
			phone:   "+78888888888",
			wantErr: false,
		},
		{
			name:    "internal error",
			phone:   "00000000000",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := srv.StopSession(ctx, tt.phone)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("error was not expected but got: %v", err)
				}
				return
			}

			if rec, err := repo.GetByPhone(ctx, tt.phone); rec != nil || err != nil {
				t.Fatalf("record not deleted for %v, err: %v", tt.phone, err)
			}

		})
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
