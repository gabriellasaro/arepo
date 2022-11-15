package arepo

import (
	"context"
	"time"

	"github.com/gabriellasaro/acache"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type abstractRepoWithCache[T any] struct {
	repo     *abstractRepo[T]
	cache    acache.Cache[acache.Key]
	rKey     acache.Key
	expCache time.Duration
}

func (a *abstractRepo[T]) WithCache(
	cache acache.Cache[acache.Key],
	radicalKey acache.Key,
	expCache time.Duration,
) AbstractRepositoryWithCache[T] {
	return &abstractRepoWithCache[T]{
		repo:     a,
		cache:    cache,
		rKey:     radicalKey,
		expCache: expCache,
	}
}

func (a *abstractRepoWithCache[T]) GetByID(ctx context.Context, id ID, opts ...*options.FindOneOptions) (*T, error) {
	data := new(T)

	key := a.rKey.Add(id.Hex())
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

func (a *abstractRepoWithCache[T]) UpdateOneByID(ctx context.Context, id ID, update any, opts ...*options.UpdateOptions) error {
	a.deleteCache(id)

	return a.repo.UpdateOneByID(ctx, id, update, opts...)
}

func (a *abstractRepoWithCache[T]) DeleteOneByID(ctx context.Context, id ID, opts ...*options.DeleteOptions) error {
	a.deleteCache(id)

	return a.repo.DeleteOneByID(ctx, id, opts...)
}

func (a *abstractRepoWithCache[T]) deleteCache(id ID) {
	go func() {
		_ = a.cache.Delete(context.Background(), a.rKey.Add(id.Hex()))
	}()
}
