package arepo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotSelectOmitFields = errors.New("not select/omit fields")

type AbstractRepo[T any] struct {
	collection *mongo.Collection
}

func NewAbstractRepository[T any](collection *mongo.Collection) *AbstractRepo[T] {
	return &AbstractRepo[T]{
		collection: collection,
	}
}

func (a *AbstractRepo[T]) GetByID(ctx context.Context, id primitive.ObjectID, opts ...*options.FindOneOptions) (*T, error) {
	return FindOneByID[T](ctx, a.collection, id, opts...)
}

func (a *AbstractRepo[T]) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) (*T, error) {
	return FindOne[T](ctx, a.collection, filter, opts...)
}

func (a *AbstractRepo[T]) FindOneAndUpdate(ctx context.Context, filter, update any, opts ...*options.FindOneAndUpdateOptions) (*T, error) {
	return FindOneAndUpdate[T](ctx, a.collection, filter, update, opts...)
}

func (a *AbstractRepo[T]) Find(ctx context.Context, filter any, opts ...*options.FindOptions) ([]*T, error) {
	return Find[T](ctx, a.collection, filter, opts...)
}

func (a *AbstractRepo[T]) InsertOne(ctx context.Context, document *T, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return a.collection.InsertOne(ctx, document, opts...)
}

func (a *AbstractRepo[T]) InsertMany(ctx context.Context, documents []*T, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	docs := make([]any, 0, len(documents))

	for _, doc := range documents {
		docs = append(docs, doc)
	}

	return a.collection.InsertMany(ctx, docs, opts...)
}

func (a *AbstractRepo[T]) UpdateOneByID(ctx context.Context, id primitive.ObjectID, update any) error {
	return UpdateOneByID(ctx, a.collection, id, update)
}

func (a *AbstractRepo[T]) UpdateOne(ctx context.Context, filter, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return a.collection.UpdateOne(ctx, filter, update, opts...)
}

func (a *AbstractRepo[T]) DeleteOneByID(ctx context.Context, id primitive.ObjectID, opts ...*options.DeleteOptions) error {
	return DeleteOneByID(ctx, a.collection, id, opts...)
}

func (a *AbstractRepo[T]) DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) error {
	return DeleteOne(ctx, a.collection, filter, opts...)
}

func (a *AbstractRepo[T]) Select(fields ...string) *SelectAndOmitFields[T] {
	setProjection := bson.D{}

	for i := range fields {
		setProjection = append(setProjection, bson.E{Key: fields[i], Value: 1})
	}

	return &SelectAndOmitFields[T]{
		collection:    a.collection,
		setProjection: setProjection,
	}
}

func (a *AbstractRepo[T]) Omit(fields ...string) *SelectAndOmitFields[T] {
	setProjection := bson.D{}

	for i := range fields {
		setProjection = append(setProjection, bson.E{Key: fields[i], Value: 0})
	}

	return &SelectAndOmitFields[T]{
		collection:    a.collection,
		setProjection: setProjection,
	}
}

type SelectAndOmitFields[T any] struct {
	collection    *mongo.Collection
	setProjection bson.D
}

func (s *SelectAndOmitFields[T]) Select(fields ...string) *SelectAndOmitFields[T] {
	for i := range fields {
		s.setProjection = append(s.setProjection, bson.E{Key: fields[i], Value: 1})
	}

	return s
}

func (s *SelectAndOmitFields[T]) Omit(fields ...string) *SelectAndOmitFields[T] {
	for i := range fields {
		s.setProjection = append(s.setProjection, bson.E{Key: fields[i], Value: 0})
	}

	return s
}

func (s *SelectAndOmitFields[T]) GetByID(ctx context.Context, id primitive.ObjectID) (*T, error) {
	if len(s.setProjection) == 0 {
		return nil, ErrNotSelectOmitFields
	}

	return FindOneByID[T](ctx, s.collection, id, options.FindOne().SetProjection(s.setProjection))
}

func (s *SelectAndOmitFields[T]) FindOne(ctx context.Context, filter any) (*T, error) {
	return FindOne[T](ctx, s.collection, filter, options.FindOne().SetProjection(s.setProjection))
}

func (s *SelectAndOmitFields[T]) FindOneAndUpdate(ctx context.Context, filter, update any) (*T, error) {
	return FindOneAndUpdate[T](ctx, s.collection, filter, update, options.FindOneAndUpdate().SetProjection(s.setProjection))
}

func (s *SelectAndOmitFields[T]) Find(ctx context.Context, filter any) ([]*T, error) {
	return Find[T](ctx, s.collection, filter, options.Find().SetProjection(s.setProjection))
}
