//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

// Package cache provides different cache type (e.g. Key/Value) with
// pluggable backends.
package cache

import (
	"context"
	"errors"
	"io"
)

// Name defines the cache name type.
type Name string

// Stringer interface
func (n Name) String() string {
	return string(n)
}

// ErrNotFound is returned by [LoadFunc] to indicate the requested
// key was not found.
var ErrNotFound error = errors.New("not found")

// Cache defines the basic cache interface used
// to retrieve a cached value for a key.
type Cache[K any, V any] interface {
	// Get gets the cached value for the given key.
	// Error ErrNotFound indicates a cache miss.
	Get(ctx context.Context, key K) (V, error)
}

// NoCache provides a cache implementation
// without any caching at all.
// This is suitable for test scenarios or to disable caching.
type NoCache[K any, V any] struct {
	Value V
	Found bool
}

// Get implements [Cache.Get].
func (c *NoCache[K, V]) Get(ctx context.Context, key K) (V, bool) {
	return c.Value, c.Found
}

// KeyValue defines the key-value cache type.
type KeyValue[K any, V any] interface {
	Cache[K, V]
	// Put adds key-value association to the cache.
	// Any previously established key-value association
	// is overwritten.
	Put(ctx context.Context, key K, value V)
	// Delete removes the given key from the cache.
	Delete(ctx context.Context, key K)
	io.Closer
}
