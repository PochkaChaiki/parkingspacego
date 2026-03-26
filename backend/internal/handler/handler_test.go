package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/pochkachaiki/parkingspace/internal/model"
	"github.com/pochkachaiki/parkingspace/internal/repository"
	"github.com/pochkachaiki/parkingspace/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// === Minimal Mock Collection for Handler Tests ===
type testMockCollection struct {
	mu      sync.Mutex
	records map[primitive.ObjectID]*model.Record
}

func newTestMockCollection() *testMockCollection {
	return &testMockCollection{records: make(map[primitive.ObjectID]*model.Record)}
}

func (m *testMockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec := document.(*model.Record)
	m.records[rec.ID] = rec
	return &mongo.InsertOneResult{InsertedID: rec.ID}, nil
}

func (m *testMockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var results []*model.Record
	for _, rec := range m.records {
		results = append(results, rec)
	}
	docs := make([]interface{}, len(results))
	for i, r := range results {
		docs[i] = r
	}
	return mongo.NewCursorFromDocuments(docs, nil, nil)
}

func (m *testMockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	m.mu.Lock()
	defer m.mu.Unlock()
	if f, ok := filter.(bson.M); ok {
		if phone, exists := f["phone_number"]; exists {
			for _, rec := range m.records {
				if rec.PhoneNumber == phone {
					return mongo.NewSingleResultFromDocument(rec, nil, nil)
				}
			}
		}
	}
	return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, nil)
}

func (m *testMockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if f, ok := filter.(bson.M); ok {
		if id, exists := f["_id"]; exists {
			if objID, ok := id.(primitive.ObjectID); ok {
				if rec, ok := m.records[objID]; ok {
					if u, ok := update.(bson.M); ok {
						if setMap, ok := u["$set"].(bson.M); ok {
							if endTime, ok := setMap["end_time"].(*time.Time); ok {
								rec.EndTime = endTime
							}
						}
					}
					return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
				}
			}
		}
	}
	return &mongo.UpdateResult{MatchedCount: 0}, nil
}

func (m *testMockCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if f, ok := filter.(bson.M); ok {
		if phone, exists := f["phone_number"]; exists {
			for id, rec := range m.records {
				if rec.PhoneNumber == phone {
					delete(m.records, id)
					return &mongo.DeleteResult{DeletedCount: 1}, nil
				}
			}
		}
	}
	return &mongo.DeleteResult{DeletedCount: 0}, nil
}

// === Handler Tests ===

func TestEnterStartSession(t *testing.T) {
	repo := repository.NewMongoRepository(newTestMockCollection())
	srv := service.NewService(repo)
	handler := NewHandler(srv, log.Default())

	reqBody := &model.RecordDto{
		ClientName:   "Иван",
		PhoneNumber:  "+79991234567",
		LicensePlate: "A123BC140",
		SpotNumber:   42,
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.StartSession(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var resp model.Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// if resp.Token == "" {
	// 	t.Fatalf("expected non-empty token in response")
	// }

	// if resp.PassCode == "" {
	// 	t.Fatalf("expected non-empty pass_code in response")
	// }

	if resp.Status != model.Success {
		t.Fatalf("expected status 'success' in response")
	}
}

func TestHandlerGetSession(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMongoRepository(newTestMockCollection())
	srv := service.NewService(repo)
	handler := NewHandler(srv, log.Default())

	const recordsNum = 3
	// Создаем несколько записей
	for i := 1; i <= recordsNum; i++ {
		_, err := srv.StartSession(ctx, &model.RecordDto{
			ClientName:   fmt.Sprintf("Client%d", i),
			PhoneNumber:  fmt.Sprintf("7999000%04d", i),
			LicensePlate: fmt.Sprintf("A%03dBC", i),
			SpotNumber:   i,
		})
		if err != nil {
			t.Fatalf("StartSession failed: %v", err)
		}
	}

	// Проверяем каждую запись
	for i := 1; i <= recordsNum; i++ {
		phone := fmt.Sprintf("7999000%04d", i)
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/sessions/%s", phone), nil)
		w := httptest.NewRecorder()
		handler.GetSession(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var record model.RecordDto
		if err := json.NewDecoder(w.Body).Decode(&record); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if record.ClientName != fmt.Sprintf("Client%d", i) {
			t.Fatalf("expected name \"Client%d\", got \"%v\"", i, record.ClientName)
		}

		if record.PhoneNumber != phone {
			t.Fatalf("expected phone number \"%s\", got \"%v\"", phone, record.PhoneNumber)
		}

		if record.LicensePlate != fmt.Sprintf("A%03dBC", i) {
			t.Fatalf("expected license plate \"A%03dBC\", got \"%v\"", i, record.LicensePlate)
		}

		if record.SpotNumber != i {
			t.Fatalf("expected spot number \"%d\", got \"%d\"", i, record.SpotNumber)
		}
	}
}

func TestHandlerProlongSession(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMongoRepository(newTestMockCollection())
	srv := service.NewService(repo)
	handler := NewHandler(srv, log.Default())

	// Создаем запись с Duration чтобы EndTime был установлен
	phone := "79990001122"
	duration := "1h"
	srv.StartSession(ctx, &model.RecordDto{
		ClientName:   "Анна",
		PhoneNumber:  phone,
		LicensePlate: "B777BB",
		SpotNumber:   7,
		Duration:     &duration,
	})

	record, _ := srv.GetSession(ctx, phone)

	if record == nil || record.EndTime == nil {
		t.Fatalf("expected record with EndTime")
	}

	initialEndTime := *record.EndTime

	body, _ := json.Marshal(&model.ProlongSessionDto{Duration: "1h"})

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/sessions/%s", phone), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ProlongSession(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res model.RecordDto
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.EndTime == nil {
		t.Fatalf("expected EndTime in response")
	}

	expectedEndTime := initialEndTime.Add(time.Hour)
	if !res.EndTime.Equal(expectedEndTime) {
		t.Fatalf("expected end_time: %v, got: %v", expectedEndTime, *res.EndTime)
	}

}

func TestHandlerStopSession(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMongoRepository(newTestMockCollection())
	srv := service.NewService(repo)
	handler := NewHandler(srv, log.Default())

	// Создаем две записи для одного номера телефона
	for i := 1; i <= 2; i++ {
		srv.StartSession(ctx, &model.RecordDto{
			ClientName:   fmt.Sprintf("Client%d", i),
			PhoneNumber:  fmt.Sprintf("7999000%04d", i),
			LicensePlate: fmt.Sprintf("C%03dA", i),
			SpotNumber:   i + 10,
		})

	}

	// Удаляем по номеру телефона
	phone := fmt.Sprintf("7999000%04d", 2)
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/sessions/%s", phone), nil)
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resps[1].Token))
	w := httptest.NewRecorder()

	handler.StopSession(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}

	// Проверяем, что записи действительно удалены
	rec, err := repo.GetByPhone(ctx, phone)
	if err != nil {
		t.Fatalf("failed to list by phone: %v", err)
	}

	if rec != nil {
		t.Fatalf("failed to delete: %v", rec)
	}
}
