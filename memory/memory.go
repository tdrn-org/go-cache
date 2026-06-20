//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

// Package memory provides memory based implementation for the
// different cache types.
package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/maypok86/otter/v2"
	"github.com/tdrn-org/go-cache"
)

// Memory cache name
const Name cache.Name = "memory"

type memoryKeyValue[K comparable, V any] struct {
	cache *otter.Cache[K, V]
	load  cache.LoadFunc[K, V]
}

// NewKeyValue creates a new [cache.KeyValue] cache.
// Parameters size and ttl define the evict strategy for the cache.
// Putting both to 0 disables cache evicting.
func NewKeyValue[K comparable, V any](size int, ttl time.Duration, load cache.LoadFunc[K, V]) (cache.KeyValue[K, V], error) {
	var expiryCalculator otter.ExpiryCalculator[K, V]
	if ttl > 0 {
		expiryCalculator = otter.ExpiryAccessing[K, V](ttl)
	}
	options := &otter.Options[K, V]{
		MaximumSize:      size,
		ExpiryCalculator: expiryCalculator,
	}
	cache, err := otter.New(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory cache (cause: %w)", err)
	}
	return &memoryKeyValue[K, V]{cache: cache, load: load}, nil
}

// [cache.KeyValue.Get] implementation
func (kv *memoryKeyValue[K, V]) Get(ctx context.Context, key K) (V, error) {
	value, err := kv.cache.Get(ctx, key, otter.LoaderFunc[K, V](kv.load))
	return value, err
}

// [cache.KeyValue.Put] implementation
func (kv *memoryKeyValue[K, V]) Put(ctx context.Context, key K, value V) {
	kv.cache.Set(key, value)
}

// [cache.KeyValue.Close] implementation
func (kv *memoryKeyValue[K, V]) Close() error {
	return nil
}
