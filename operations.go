package arepo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindOne[T any](ctx context.Context, collection *mongo.Collection, filter any, opts ...*options.FindOneOptions) (*T, error) {
	result := collection.FindOne(ctx, filter, opts...)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, result.Err()
	}

	var data T

	if err := result.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func FindOneByID[T any](ctx context.Context, collection *mongo.Collection, id ID, opts ...*options.FindOneOptions) (*T, error) {
	return FindOne[T](ctx, collection, bson.M{
		"_id": primitive.ObjectID(id),
	})
}

func Find[T any](ctx context.Context, collection *mongo.Collection, filter any, opts ...*options.FindOptions) ([]T, error) {
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var list []T

	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func UpdateOne(ctx context.Context, collection *mongo.Collection, filter, update any, opts ...*options.UpdateOptions) error {
	r, err := collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	if r.MatchedCount == 0 {
		return ErrNotFound
	}

	if r.ModifiedCount == 0 {
		return ErrNotUpdated
	}

	return nil
}

func UpdateOneByID(ctx context.Context, collection *mongo.Collection, id ID, update any, opts ...*options.UpdateOptions) error {
	return UpdateOne(
		ctx,
		collection,
		bson.M{
			"_id": primitive.ObjectID(id),
		},
		update,
		opts...,
	)
}

func DeleteOne(ctx context.Context, collection *mongo.Collection, filter any, opts ...*options.DeleteOptions) error {
	r, err := collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return err
	}

	if r.DeletedCount == 0 {
		return ErrNotDeleted
	}

	return nil
}

func DeleteOneByID(ctx context.Context, collection *mongo.Collection, id ID, opts ...*options.DeleteOptions) error {
	return DeleteOne(
		ctx,
		collection,
		bson.M{
			"_id": primitive.ObjectID(id),
		},
		opts...,
	)
}
