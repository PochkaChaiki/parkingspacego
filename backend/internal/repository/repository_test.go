package repository

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/pochkachaiki/parkingspace/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mockCollection имитирует mongo.Collection для тестирования
type mockCollection struct {
	mu      sync.Mutex
	records map[primitive.ObjectID]*model.Record
}

func newMockCollection() *mockCollection {
	return &mockCollection{records: make(map[primitive.ObjectID]*model.Record)}
}

func (m *mockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := document.(*model.Record)
	if !ok {
		return nil, ErrNotFound
	}
	m.records[rec.ID] = rec
	return &mongo.InsertOneResult{InsertedID: rec.ID}, nil
}

func (m *mockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	phone := ""
	if f, ok := filter.(bson.M); ok {
		if v, exists := f["phone_number"]; exists {
			phone, _ = v.(string)
		}
	}

	var results []*model.Record
	for _, rec := range m.records {
		if phone == "" || rec.PhoneNumber == phone {
			results = append(results, rec)
		}
	}

	docs := make([]interface{}, len(results))
	for i, r := range results {
		docs[i] = r
	}
	return mongo.NewCursorFromDocuments(docs, nil, nil)
}

func (m *mockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	phone := ""
	objID := primitive.NilObjectID

	if f, ok := filter.(bson.M); ok {
		if v, exists := f["phone_number"]; exists {
			phone, _ = v.(string)
		}
		if v, exists := f["_id"]; exists {
			if id, ok2 := v.(primitive.ObjectID); ok2 {
				objID = id
			}
		}
	}

	// Поиск по _id
	if !objID.IsZero() {
		if rec, ok := m.records[objID]; ok {
			return mongo.NewSingleResultFromDocument(rec, nil, nil)
		}
		return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, nil)
	}

	// Поиск по phone_number
	if phone != "" {
		for _, rec := range m.records {
			if rec.PhoneNumber == phone {
				return mongo.NewSingleResultFromDocument(rec, nil, nil)
			}
		}
		return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, nil)
	}

	return mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, nil)
}

func (m *mockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	objID := primitive.NilObjectID
	if f, ok := filter.(bson.M); ok {
		if v, exists := f["_id"]; exists {
			if id, ok2 := v.(primitive.ObjectID); ok2 {
				objID = id
			}
		}
	}

	rec, ok := m.records[objID]
	if !ok {
		return &mongo.UpdateResult{MatchedCount: 0}, nil
	}

	if u, ok := update.(bson.M); ok {
		if setMap, ok := u["$set"].(bson.M); ok {
			for k, v := range setMap {
				if k == "end_time" {
					if tt, ok := v.(*time.Time); ok {
						rec.EndTime = tt
					}
				}
				if k == "status" {
					if s, ok := v.(string); ok {
						rec.Status = s
					}
				}
			}
		}
	}

	return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
}

func (m *mockCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	phone := ""
	if f, ok := filter.(bson.M); ok {
		if v, exists := f["phone_number"]; exists {
			phone, _ = v.(string)
		}
	}

	count := int64(0)
	for id, rec := range m.records {
		if rec.PhoneNumber == phone {
			delete(m.records, id)
			count++
		}
	}
	return &mongo.DeleteResult{DeletedCount: count}, nil
}

// ============= TESTS (RED/GREEN TDD) =============

// TestRepositoryCreate - тест создания записи
func TestRepositoryCreate(t *testing.T) {
	ctx := context.Background()
	coll := newMockCollection()
	repo := NewMongoRepository(coll)

	rec := &model.Record{
		ClientName:   "Иван",
		PhoneNumber:  "+79991234567",
		LicensePlate: "A123BC140",
		SpotNumber:   42,
		StartTime:    time.Now().UTC(),
		Status:       "active",
	}

	created, err := repo.Create(ctx, rec)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if created == nil {
		t.Fatalf("expected record, got nil")
	}

	if created.ID.IsZero() {
		t.Fatalf("expected non-zero ID")
	}

	if created.Status != "active" {
		t.Fatalf("expected status 'active', got %q", created.Status)
	}
}

// TestRepositoryGetByPhone - тест получения записи по номеру телефона
func TestRepositoryGetByPhone(t *testing.T) {
	ctx := context.Background()
	coll := newMockCollection()
	repo := NewMongoRepository(coll)

	phone := "+79991234567"
	rec := &model.Record{
		ClientName:   "Иван",
		PhoneNumber:  phone,
		LicensePlate: "A123BC140",
		SpotNumber:   42,
		StartTime:    time.Now().UTC(),
		Status:       "active",
	}

	created, err := repo.Create(ctx, rec)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := repo.GetByPhone(ctx, phone)
	if err != nil {
		t.Fatalf("GetByPhone failed: %v", err)
	}

	if retrieved == nil {
		t.Fatalf("expected record, got nil")
	}

	if retrieved.ID != created.ID {
		t.Fatalf("expected ID %v, got %v", created.ID, retrieved.ID)
	}
}

// TestRepositoryUpdate - тест обновления времени окончания
func TestRepositoryUpdate(t *testing.T) {
	ctx := context.Background()
	coll := newMockCollection()
	repo := NewMongoRepository(coll)

	rec := &model.Record{
		ClientName:   "Иван",
		PhoneNumber:  "+79991234567",
		LicensePlate: "A123BC140",
		SpotNumber:   42,
		StartTime:    time.Now().UTC(),
		Status:       "active",
	}

	created, err := repo.Create(ctx, rec)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	endTime := time.Now().UTC().Add(time.Hour)
	err = repo.Update(ctx, created.ID.Hex(), &endTime)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	retrieved, err := repo.GetByPhone(ctx, "+79991234567")
	if err != nil {
		t.Fatalf("GetByPhone failed: %v", err)
	}

	if retrieved.EndTime == nil {
		t.Fatalf("expected EndTime to be set")
	}
}

// TestRepositoryDeleteByPhone - тест удаления по номеру телефона
func TestRepositoryDeleteByPhone(t *testing.T) {
	ctx := context.Background()
	coll := newMockCollection()
	repo := NewMongoRepository(coll)

	phone := "+79991234567"
	rec := &model.Record{
		ClientName:   "Иван",
		PhoneNumber:  phone,
		LicensePlate: "A123BC140",
		SpotNumber:   42,
		StartTime:    time.Now().UTC(),
		Status:       "active",
	}

	_, err := repo.Create(ctx, rec)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = repo.DeleteByPhone(ctx, phone)
	if err != nil {
		t.Fatalf("DeleteByPhone failed: %v", err)
	}

	retrieved, err := repo.GetByPhone(ctx, phone)
	if err != nil {
		t.Fatalf("GetByPhone failed: %v", err)
	}

	if retrieved != nil {
		t.Fatalf("expected nil after delete, got %v", retrieved)
	}
}
