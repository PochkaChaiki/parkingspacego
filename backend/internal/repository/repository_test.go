package repository

import (
	"context"
	"errors"
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
	firstID := primitive.NewObjectID()
	secondID := primitive.NewObjectID()
	thirdID := primitive.NewObjectID()
	return &mockCollection{
		records: map[primitive.ObjectID]*model.Record{
			firstID: {
				ID:           firstID,
				ClientName:   "Egor",
				PhoneNumber:  "+79999999999",
				LicensePlate: "A123BC123",
				SpotNumber:   1,
				StartTime:    time.Now().UTC(),
				EndTime:      time.Now().UTC().Add(time.Hour),
				Status:       "active",
			},
			secondID: {
				ID:           secondID,
				ClientName:   "Arisha",
				PhoneNumber:  "+79999999998",
				LicensePlate: "A123BC122",
				SpotNumber:   2,
				StartTime:    time.Now().UTC(),
				EndTime:      time.Now().UTC().Add(time.Hour),
				Status:       "active",
			},
			thirdID: {
				ID:           thirdID,
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

func (m *mockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec, ok := document.(*model.Record)
	if !ok {
		return nil, ErrNotFound
	}

	if rec.PhoneNumber == "00000000000" {
		return nil, errors.New("internal error")
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

	if phone == "00000000000" {
		return nil, errors.New("internal error")
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

	if phone == "00000000000" {
		return mongo.NewSingleResultFromDocument(nil, errors.New("internal error"), nil)
	}

	// Поиск по _id
	if !objID.IsZero() {
		if rec, ok := m.records[objID]; ok {
			return mongo.NewSingleResultFromDocument(rec, nil, nil)
		}
		return mongo.NewSingleResultFromDocument(&model.Record{}, mongo.ErrNoDocuments, nil)
	}

	// Поиск по phone_number
	if phone != "" {
		for _, rec := range m.records {
			if rec.PhoneNumber == phone {
				return mongo.NewSingleResultFromDocument(rec, nil, nil)
			}
		}
		return mongo.NewSingleResultFromDocument(&model.Record{}, mongo.ErrNoDocuments, nil)
	}

	return mongo.NewSingleResultFromDocument(&model.Record{}, mongo.ErrNoDocuments, nil)
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

	if objID.IsZero() {
		return nil, errors.New("internal error")
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
						rec.EndTime = *tt
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

	if phone == "00000000000" {
		return nil, errors.New("internal error")
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

func TestRepository_Create(t *testing.T) {

	coll := newMockCollection()
	repo := NewMongoRepository(coll)

	tests := []struct {
		name    string
		rec     *model.Record
		wantErr bool
	}{
		{
			name: "creating doc success",
			rec: &model.Record{
				ID:           primitive.NewObjectID(),
				ClientName:   "egor",
				PhoneNumber:  "+78888888888",
				LicensePlate: "A321BC321",
				SpotNumber:   100,
				StartTime:    time.Now().UTC(),
				EndTime:      time.Now().UTC().Add(time.Hour),
				Status:       "active",
			},
			wantErr: false,
		},
		{
			name: "doc without id",
			rec: &model.Record{
				ClientName:   "egor",
				PhoneNumber:  "+78888888887",
				LicensePlate: "A321BC321",
				SpotNumber:   101,
				StartTime:    time.Now().UTC(),
				EndTime:      time.Now().UTC().Add(time.Hour),
			},
			wantErr: false,
		},
		{
			name:    "rec is nil",
			rec:     nil,
			wantErr: true,
		},
		{
			name: "internal error",
			rec: &model.Record{
				ClientName:   "egor",
				PhoneNumber:  "00000000000",
				LicensePlate: "A321BC321",
				SpotNumber:   100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := repo.Create(ctx, tt.rec)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("error was not expected")
				}
			}
		})
	}
}

func TestRepository_GetAll(t *testing.T) {
	coll := newMockCollection()
	repo := NewMongoRepository(coll)
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name:    "get all docs successful",
			size:    len(coll.records),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			recs, err := repo.GetAll(ctx)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("error was not expected")
					return
				}
			}
			if len(recs) != tt.size {
				t.Fatalf("want %d records, but got %d", tt.size, len(recs))
				return
			}
		})
	}
}

func TestRepository_GetByPhone(t *testing.T) {
	coll := newMockCollection()
	repo := NewMongoRepository(coll)
	tests := []struct {
		name    string
		phone   string
		rec     *model.Record
		wantErr bool
	}{
		{
			name:  "get doc by phone successful",
			phone: "+79999999999",
			rec: &model.Record{
				ID:           primitive.NewObjectID(),
				ClientName:   "Egor",
				PhoneNumber:  "+79999999999",
				LicensePlate: "A123BC123",
				SpotNumber:   1,
				StartTime:    time.Now().UTC(),
				EndTime:      time.Now().UTC().Add(time.Hour),
				Status:       "active",
			},
			wantErr: false,
		},
		{
			name:    "get doc by phone failure",
			phone:   "+77777777777",
			rec:     nil,
			wantErr: false,
		},
		{
			name:    "internal error",
			phone:   "00000000000",
			rec:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			rec, err := repo.GetByPhone(ctx, tt.phone)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("error was not expected: %v", err)
				}
				return
			}

			if rec == nil {
				if tt.rec != nil {
					t.Fatalf("expected record but got nil")
				}
				return
			}

			if !(rec.ClientName == tt.rec.ClientName &&
				rec.SpotNumber == tt.rec.SpotNumber &&
				rec.PhoneNumber == tt.rec.PhoneNumber &&
				rec.LicensePlate == tt.rec.LicensePlate) {
				t.Fatalf("rec want: %v, got: %v", *tt.rec, *rec)
				return
			}

		})
	}
}

// TestRepositoryUpdate - тест обновления времени окончания
func TestRepository_Update(t *testing.T) {
	coll := newMockCollection()
	repo := NewMongoRepository(coll)

	tests := []struct {
		name    string
		phone   string
		endTime time.Time
		wantErr bool
	}{
		{
			name:    "update doc by phone successful",
			phone:   "+79999999999",
			endTime: time.Now().UTC().Add(2 * time.Hour),
			wantErr: false,
		},
		{
			name:    "internal error",
			phone:   "",
			endTime: time.Now().UTC().Add(2 * time.Hour),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if !tt.wantErr {
				var rec model.Record
				coll.FindOne(ctx, bson.M{"phone_number": tt.phone}).Decode(&rec)
				err := repo.Update(ctx, rec.ID.Hex(), tt.endTime)
				if err != nil {
					t.Fatalf("error was not expected: %v", err)
					return
				}
			} else {
				err := repo.Update(ctx, primitive.NilObjectID.Hex(), tt.endTime)
				if err == nil {
					t.Fatalf("error was expected, but got nil")
					return
				}
			}

		})
	}
}

// TestRepositoryDeleteByPhone - тест удаления по номеру телефона
func TestRepository_DeleteByPhone(t *testing.T) {
	coll := newMockCollection()
	repo := NewMongoRepository(coll)

	tests := []struct {
		name    string
		phone   string
		rec     *model.Record
		wantErr bool
	}{
		{
			name:    "delete doc by phone successful",
			phone:   "+79999999999",
			wantErr: false,
		},
		{
			name:    "delete doc by phone that not exist",
			phone:   "+77777777777",
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
			err := repo.DeleteByPhone(ctx, tt.phone)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("error was not expected: %v", err)
				}
				return
			}

			var rec model.Record
			err = coll.FindOne(ctx, bson.M{"phone_number": tt.phone}).Decode(&rec)

			if err != nil {
				if err != mongo.ErrNoDocuments {
					t.Fatalf("such error was not expected: %v", err)
				}
				return
			} else {
				if tt.wantErr {
					t.Fatalf(" error was expected but got nil")
				}
			}

		})
	}
}
