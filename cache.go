package arepo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gabriellasaro/acache"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type abstractRepoWithCache[T any] struct {
	repo      *abstractRepo[T]
	cache     acache.Cache[acache.Key]
	rKey      acache.Key
	rKeyForID acache.Key
	expCache  time.Duration
}

func (a *abstractRepo[T]) WithCache(
	cache acache.Cache[acache.Key],
	radicalKey acache.Key,
	expCache time.Duration,
) AbstractRepositoryWithCache[T] {
	return &abstractRepoWithCache[T]{
		repo:      a,
		cache:     cache,
		rKey:      radicalKey,
		rKeyForID: radicalKey.Add("_id"),
		expCache:  expCache,
	}
}

func (a *abstractRepoWithCache[T]) GetByID(ctx context.Context, id ID, opts ...*options.FindOneOptions) (*T, error) {
	data := new(T)

	key := a.rKeyForID.Add(id.Hex())
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

func (a *abstractRepoWithCache[T]) UpdateOneByID(ctx context.Context, id ID, update any) error {
	a.deleteCacheByID(id)

	return a.repo.UpdateOneByID(ctx, id, update)
}

func (a *abstractRepoWithCache[T]) DeleteOneByID(ctx context.Context, id ID, opts ...*options.DeleteOptions) error {
	a.deleteCacheByID(id)

	return a.repo.DeleteOneByID(ctx, id, opts...)
}

func (a *abstractRepoWithCache[T]) deleteCacheByID(id ID) {
	go func() {
		_ = a.cache.Delete(context.Background(), a.rKeyForID.Add(id.Hex()))
	}()
}

func (a *abstractRepoWithCache[T]) WithCustomFilter() CacheWithCustomFilter[T] {
	return &cacheWithCustomFilter[T]{
		repo:     a.repo,
		cache:    a.cache,
		rKey:     a.rKey.Add("custom"),
		expCache: a.expCache,
	}
}

type cacheWithCustomFilter[T any] struct {
	repo     *abstractRepo[T]
	cache    acache.Cache[acache.Key]
	rKey     acache.Key
	expCache time.Duration
}

func (c *cacheWithCustomFilter[T]) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) (*T, error) {
	key, err := c.customKey(filter)
	if err != nil {
		return nil, err
	}

	data := new(T)

	if err := c.cache.GetJSON(ctx, key, data); err == nil {
		return data, nil
	}

	data, err = c.repo.FindOne(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	go func() {
		_ = c.cache.SetJSON(context.Background(), key, data, c.expCache)
	}()

	return data, nil
}

func (c *cacheWithCustomFilter[T]) Find(ctx context.Context, filter any, opts ...*options.FindOptions) ([]*T, error) {
	key, err := c.customKey(filter)
	if err != nil {
		return nil, err
	}

	var data []*T

	if err := c.cache.GetJSON(ctx, key, data); err == nil {
		return data, nil
	}

	data, err = c.repo.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	go func() {
		_ = c.cache.SetJSON(context.Background(), key, data, c.expCache)
	}()

	return data, nil
}

func (c *cacheWithCustomFilter[T]) customKey(filter any) (acache.Key, error) {
	hash, err := filterHash(filter)
	if err != nil {
		return "", err
	}

	return c.rKey.Add(hash), nil
}

func filterHash(filter any) (string, error) {
	content, err := json.Marshal(filter)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(content)

	return hex.EncodeToString(sum[0:]), nil
}
