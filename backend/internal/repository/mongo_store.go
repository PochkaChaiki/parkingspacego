package repository

import (
	"context"
	"errors"
	"time"

	"github.com/pochkachaiki/parkingspace/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNotFound = errors.New("record not found")
)

// mongoCollection позволяет мокировать mongo.Collection в тестах.
type mongoCollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

// NewMongoRepository создает репозиторий, работающий с MongoDB.
func NewMongoRepository(coll mongoCollection) Repository {
	return &mongoRepository{coll: coll}
}

type mongoRepository struct {
	coll mongoCollection
}

func (r *mongoRepository) Create(ctx context.Context, rec *model.Record) (*model.Record, error) {
	if rec == nil {
		return nil, errors.New("record is nil")
	}
	if rec.Status == "" {
		rec.Status = "active"
	}
	if rec.ID.IsZero() {
		rec.ID = primitive.NewObjectID()
	}
	_, err := r.coll.InsertOne(ctx, rec)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// GetAll возвращает все паркинг-сессии из базы.
func (r *mongoRepository) GetAll(ctx context.Context) ([]*model.Record, error) {
	cursor, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var results []*model.Record
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetByPhone возвращает первую сессию по номеру телефона.
// Если сессия не найдена, возвращает nil без ошибки.
func (r *mongoRepository) GetByPhone(ctx context.Context, phone string) (*model.Record, error) {
	var rec model.Record
	err := r.coll.FindOne(ctx, bson.M{"phone_number": phone}).Decode(&rec)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		// Если ошибка при декодировании документа, это означает что документ не найден
		// (например, ошибка "document is nil")
		return nil, err
	}
	return &rec, nil
}

// Update обновляет время окончания сессии парковки по ID.
func (r *mongoRepository) Update(ctx context.Context, id string, endTime time.Time) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": bson.M{"end_time": endTime},
	}
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *mongoRepository) DeleteByPhone(ctx context.Context, phone string) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{"phone_number": phone})
	return err
}
