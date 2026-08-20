//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

// Package redis provides Redis based implementation for the
// different cache types.
package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tdrn-org/go-cache"
)

// Redis cache name
const Name cache.Name = "redis"

// KeyFunc functions are responsible for converting a key to
// a unique string key usable for Redis.
type KeyFunc[K any] func(K) string

// StringKey is the identitiy [KeyFunc].
func StringKey(key string) string {
	return key
}

type redisKeyValue[K any, V any] struct {
	rdb        *redis.Client
	ttl        time.Duration
	keyFunc    KeyFunc[K]
	serializer cache.Serializer[V]
	noValue    V
	logger     *slog.Logger
}

type Options = redis.Options

// NewKeyValue creates a new [cache.KeyValue] cache.
// Parameter options defines the Redis client options. Parameter ttl defines the
// evict strategy for the cache.
// A positive ttl defines an access ttl (entry is discarded ttl after last access).
// A negative ttl defines a created ttl (entry is discarded ttl after creation).
// Putting both to 0 disables cache evicting.
func NewKeyValue[K any, V any](options *Options, ttl time.Duration, keyFunc KeyFunc[K], serializer cache.Serializer[V]) (cache.KeyValue[K, V], error) {
	rdb := redis.NewClient(options)
	kv := &redisKeyValue[K, V]{
		rdb:        rdb,
		ttl:        ttl,
		keyFunc:    keyFunc,
		serializer: serializer,
		logger:     slog.With(slog.String("redis", fmt.Sprintf("%s://%s#%d", options.Network, options.Addr, options.DB))),
	}
	return kv, nil
}

// [cache.KeyValue.Get] implementation
func (kv *redisKeyValue[K, V]) Get(ctx context.Context, key K) (V, error) {
	keyString := kv.keyFunc(key)
	var encodedValue []byte
	var err error
	if kv.ttl > 0 {
		// A positive ttl defines an access ttl: refresh the entry's expiry
		// on every read, so it is discarded ttl after last access.
		encodedValue, err = kv.rdb.GetEx(ctx, keyString, kv.ttl).Bytes()
	} else {
		encodedValue, err = kv.rdb.Get(ctx, keyString).Bytes()
	}
	if errors.Is(err, redis.Nil) {
		return kv.noValue, cache.ErrNotFound
	} else if err != nil {
		return kv.noValue, fmt.Errorf("failed to get cache value (cause: %w)", err)
	}
	return kv.serializer.Unmarshal(encodedValue)
}

// [cache.KeyValue.Put] implementation
func (kv *redisKeyValue[K, V]) Put(ctx context.Context, key K, value V) {
	encodedValue, err := kv.serializer.Marshal(value)
	if err != nil {
		kv.logger.Warn("failed to marshal cache value (cause: %w)", slog.Any("err", err))
	}
	var ttl time.Duration
	var expireAt time.Time
	if kv.ttl >= 0 {
		ttl = kv.ttl
	} else {
		expireAt = time.Now().Add(-kv.ttl)
	}
	err = kv.rdb.SetArgs(ctx, kv.keyFunc(key), encodedValue, redis.SetArgs{
		TTL:      ttl,
		ExpireAt: expireAt,
	}).Err()
	if err != nil {
		kv.logger.Warn("failed to set cache value (cause: %w)", slog.Any("err", err))
	}
}

// [cache.KeyValue.Delete] implementation
func (kv *redisKeyValue[K, V]) Delete(ctx context.Context, key K) {
	err := kv.rdb.Del(ctx, kv.keyFunc(key)).Err()
	if err != nil {
		kv.logger.Warn("failed to delete cache value (cause: %w)", slog.Any("err", err))
	}
}

// [cache.KeyValue.Close] implementation
func (kv *redisKeyValue[K, V]) Close() error {
	return kv.rdb.Close()
}
