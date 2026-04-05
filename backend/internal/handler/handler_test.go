package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pochkachaiki/parkingspace/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceMock struct {
	store map[string]*model.Record
}

func NewServiceMock() *ServiceMock {
	return &ServiceMock{
		store: map[string]*model.Record{
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

func (m *ServiceMock) StartSession(ctx context.Context, r *model.RecordDto) (model.StatusName, error) {
	if r == nil {
		return model.Failure, nil
	}

	if r.PhoneNumber == "00000000000" {
		return model.Failure, errors.New("internal error")
	}

	// Проверяем, есть ли уже сессия для этого номера телефона
	existing, ok := m.store[r.PhoneNumber]
	if ok || existing != nil {
		return model.Failure, nil
	}

	// Создаем новую запись
	rec := &model.Record{
		ID:           primitive.NewObjectID(),
		ClientName:   r.ClientName,
		PhoneNumber:  r.PhoneNumber,
		LicensePlate: r.LicensePlate,
		SpotNumber:   r.SpotNumber,
		StartTime:    time.Now().UTC(),
		EndTime:      time.Now().UTC().Add(time.Hour),
		Status:       "active",
	}

	// Вычисляем EndTime на основе Duration (если указана)
	if r.Duration != nil && *r.Duration != "" {
		// Парсим duration (например "1h")
		duration, err := time.ParseDuration(*r.Duration)
		if err == nil {
			endTime := rec.StartTime.Add(duration)
			rec.EndTime = endTime
		}
	}

	// Проверяем, свободно ли место
	for _, check := range m.store {
		if check.SpotNumber == r.SpotNumber && check.Status == "active" {
			// Место занято - возвращаем Occupied БЕЗ ошибки
			return model.Occupied, nil
		}
	}

	// Создаем запись в репозитории
	m.store[r.PhoneNumber] = rec

	return model.Success, nil
}

func (m *ServiceMock) GetSession(ctx context.Context, phone string) (*model.RecordDto, error) {
	if phone == "00000000000" {
		return nil, errors.New("internal error")
	}

	rec, ok := m.store[phone]
	if !ok || rec == nil {
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

func (m *ServiceMock) ProlongSession(ctx context.Context, phone string, duration string) (*model.RecordDto, error) {
	if phone == "00000000000" {
		return nil, errors.New("internal error")
	}

	rec, ok := m.store[phone]
	if !ok || rec == nil {
		return nil, nil
	}

	// Если EndTime не установлена, устанавливаем её как StartTime + Duration
	timeToAdd, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}
	rec.EndTime = rec.EndTime.Add(timeToAdd)

	// Обновляем запись в репозитории
	m.store[rec.PhoneNumber].EndTime = rec.EndTime

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
func (m *ServiceMock) StopSession(ctx context.Context, phone string) error {
	if phone == "00000000000" {
		return errors.New("internal error")
	}
	delete(m.store, phone)
	return nil
}

// === Handler Tests ===

func TestHandler_StartSession(t *testing.T) {
	srv := NewServiceMock()
	handler := NewHandler(srv, log.Default())

	tests := []struct {
		name       string
		dto        interface{}
		wantCode   int
		wantStatus model.StatusName
	}{
		{
			name: "creating session success",
			dto: &model.RecordDto{
				ClientName:   "egor",
				PhoneNumber:  "+78888888888",
				LicensePlate: "A321BC321",
				SpotNumber:   100,
			},
			wantCode:   http.StatusCreated,
			wantStatus: model.Success,
		},
		{
			name: "session exist",
			dto: &model.RecordDto{
				ClientName:   "egor",
				PhoneNumber:  "+78888888888",
				LicensePlate: "A321BC321",
				SpotNumber:   101,
			},
			wantCode:   http.StatusOK,
			wantStatus: model.Failure,
		},
		{
			name: "spot occupied",
			dto: &model.RecordDto{
				ClientName:   "egor",
				PhoneNumber:  "+78888888881",
				LicensePlate: "A321BC321",
				SpotNumber:   100,
			},
			wantCode:   http.StatusOK,
			wantStatus: model.Occupied,
		},
		{
			name:       "json invalid",
			dto:        []byte("not json"),
			wantCode:   http.StatusBadRequest,
			wantStatus: model.Failure,
		},
		{
			name: "internal error",
			dto: &model.RecordDto{
				ClientName:   "egor",
				PhoneNumber:  "00000000000",
				LicensePlate: "A321BC321",
				SpotNumber:   100,
			},
			wantCode:   http.StatusInternalServerError,
			wantStatus: model.Failure,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.dto)
			req := httptest.NewRequest("POST", "/api/sessions", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler.StartSession(w, req)

			if w.Code != tt.wantCode {
				t.Fatalf("status want %d, got %d", tt.wantCode, w.Code)
			}

			if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
				var resp model.Response
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if resp.Status != tt.wantStatus {
					t.Fatalf("response status want %s, got %s", tt.wantStatus, resp.Status)
				}
			}
		})
	}
}

func TestHandler_GetSession(t *testing.T) {
	srv := NewServiceMock()
	handler := NewHandler(srv, log.Default())

	tests := []struct {
		name     string
		phone    string
		wantCode int
		wantDto  *model.RecordDto
	}{
		{
			name:  "get session successful",
			phone: "+79999999999",

			wantCode: http.StatusOK,
			wantDto: &model.RecordDto{
				ClientName:   "Egor",
				PhoneNumber:  "+79999999999",
				LicensePlate: "A123BC123",
				SpotNumber:   1,
			},
		},
		{
			name:     "session don't exist",
			phone:    "+78888888888",
			wantCode: http.StatusNotFound,
			wantDto:  nil,
		},
		{
			name:     "internal error",
			phone:    "00000000000",
			wantCode: http.StatusInternalServerError,
			wantDto:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/sessions/%s", tt.phone), nil)
			// req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler.GetSession(w, req)

			if w.Code != tt.wantCode {
				t.Fatalf("status want %d, got %d", tt.wantCode, w.Code)
			}

			if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
				var resp model.RecordDto
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if !(resp.ClientName == tt.wantDto.ClientName &&
					resp.SpotNumber == tt.wantDto.SpotNumber &&
					resp.PhoneNumber == tt.wantDto.PhoneNumber &&
					resp.LicensePlate == tt.wantDto.LicensePlate) {
					t.Fatalf("dto want: %v, got: %v", *tt.wantDto, resp)
				}
			}
		})
	}
}

func TestHandler_ProlongSession(t *testing.T) {
	srv := NewServiceMock()
	handler := NewHandler(srv, log.Default())

	tests := []struct {
		name     string
		phone    string
		dto      *model.ProlongSessionDto
		wantCode int
		// wantDto  *model.RecordDto
	}{
		{
			name:     "prolong session successful",
			phone:    "+79999999999",
			dto:      &model.ProlongSessionDto{Duration: "1h"},
			wantCode: http.StatusOK,
		},
		{
			name:     "session don't exist",
			phone:    "+78888888888",
			dto:      &model.ProlongSessionDto{Duration: "1h"},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "internal error",
			phone:    "00000000000",
			dto:      &model.ProlongSessionDto{Duration: "1h"},
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.dto)
			// record, _ := srv.GetSession(context.Background(), tt.phone)
			// var resEndTime time.Time
			// if record != nil {
			// 	dur, _ := time.ParseDuration(tt.dto.Duration)
			// 	resEndTime = record.EndTime.Add(dur)
			// 	// fmt.Printf("record end time: %v, after addition: %v", record.EndTime, resEndTime)
			// }
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/sessions/%s", tt.phone), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler.ProlongSession(w, req)

			if w.Code != tt.wantCode {
				t.Fatalf("status want %d, got %d", tt.wantCode, w.Code)
			}

			if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
				var resp model.RecordDto
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				// if !resp.EndTime.Equal(resEndTime) {
				// 	t.Fatalf("endTime want: %v, got: %v", resEndTime, *resp.EndTime)
				// }
			}
		})
	}
}

func TestHandler_StopSession(t *testing.T) {
	srv := NewServiceMock()
	handler := NewHandler(srv, log.Default())

	tests := []struct {
		name     string
		phone    string
		wantCode int
		// wantDto  *model.RecordDto
	}{
		{
			name:     "stop session successful",
			phone:    "+79999999999",
			wantCode: http.StatusOK,
		},
		{
			name:     "session don't exist",
			phone:    "+78888888888",
			wantCode: http.StatusOK,
		},
		{
			name:     "internal error",
			phone:    "00000000000",
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/sessions/%s", tt.phone), nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler.StopSession(w, req)

			if w.Code != tt.wantCode {
				t.Fatalf("status want %d, got %d", tt.wantCode, w.Code)
			}

			if w.Code != http.StatusInternalServerError {
				if rec, err := srv.GetSession(context.Background(), tt.phone); rec != nil || err != nil {
					t.Fatalf("record not deleted for %v, err: %v", tt.phone, err)
				}

			}
		})
	}
}
