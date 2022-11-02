package arepo

import (
	"context"
	"errors"
	"time"

	"github.com/gabriellasaro/acache"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrNotUpdated          = errors.New("not updated")
	ErrNotDeleted          = errors.New("not deleted")
	ErrNotSelectOmitFields = errors.New("not select/omit fields")
)

type AbstractRepository[T any] interface {
	GetByID(ctx context.Context, id ID, opts ...*options.FindOneOptions) (*T, error)
	Select(fields ...string) SelectAndOmitFields[T]
	Omit(fields ...string) SelectAndOmitFields[T]
	InsertOne(ctx context.Context, document *T, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []*T, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	UpdateOneByID(ctx context.Context, id ID, update any, opts ...*options.UpdateOptions) error
	DeleteOneByID(ctx context.Context, id ID, opts ...*options.DeleteOptions) error
	WithCache(cache acache.Cache[acache.Key], radicalKey acache.Key, expCache time.Duration) AbstractRepositoryWithCache[T]
}

type AbstractRepositoryWithCache[T any] interface {
	GetByID(ctx context.Context, id ID, opts ...*options.FindOneOptions) (*T, error)
	UpdateOneByID(ctx context.Context, id ID, update any, opts ...*options.UpdateOptions) error
	DeleteOneByID(ctx context.Context, id ID, opts ...*options.DeleteOptions) error
}

type SelectAndOmitFields[T any] interface {
	GetByID(ctx context.Context, id ID) (*T, error)
}
