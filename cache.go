package arepo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cache[K ~string] interface {
	GetJSON(ctx context.Context, key K, dest any) error
	SetJSON(ctx context.Context, key K, value any, exp time.Duration) error
	Delete(ctx context.Context, key K) error
}

type AbstractRepoWithCache[T any, K ~string] struct {
	repo      *AbstractRepo[T]
	cache     cache[K]
	rKey      K
	rKeyForID K
	expCache  time.Duration
}

func NewRepositoryWithCache[T any, K ~string](
	repo *AbstractRepo[T],
	cache cache[K],
	radicalKey K,
	expCache time.Duration,
) *AbstractRepoWithCache[T, K] {
	return &AbstractRepoWithCache[T, K]{
		repo:      repo,
		cache:     cache,
		rKey:      radicalKey,
		rKeyForID: radicalKey + ":_id",
		expCache:  expCache,
	}
}

func (a *AbstractRepoWithCache[T, k]) GetByID(ctx context.Context, id primitive.ObjectID, opts ...*options.FindOneOptions) (*T, error) {
	data := new(T)

	key := a.rKeyForID + ":" + k(id.Hex())
	if err := a.cache.GetJSON(ctx, key, data); err == nil {
		return data, nil
	}

	data, err := a.repo.GetByID(ctx, id, opts...)
	if err != nil {
		return nil, err
	}

	go func() {
		_ = a.cache.SetJSON(context.Background(), key, data, a.expCache)
	}()

	return data, nil
}

func (a *AbstractRepoWithCache[T, k]) UpdateOneByID(ctx context.Context, id primitive.ObjectID, update any) error {
	a.deleteCacheByID(id)

	return a.repo.UpdateOneByID(ctx, id, update)
}

func (a *AbstractRepoWithCache[T, k]) DeleteOneByID(ctx context.Context, id primitive.ObjectID, opts ...*options.DeleteOptions) error {
	a.deleteCacheByID(id)

	return a.repo.DeleteOneByID(ctx, id, opts...)
}

func (a *AbstractRepoWithCache[T, k]) deleteCacheByID(id primitive.ObjectID) {
	go func() {
		_ = a.cache.Delete(context.Background(), a.rKeyForID+":"+k(id.Hex()))
	}()
}
