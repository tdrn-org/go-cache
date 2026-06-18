/*
 * Copyright 2026 Holger de Carne
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package cache provides differnt cache type (e.g. Key/Value) with
// plugable backends.
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

// LoadFunc defines the signature for functions responsible for
// loading cache values based on the value's key.
type LoadFunc[K comparable, V any] func(ctx context.Context, key K) (V, error)

// NotFound is a [LoadFunc] not loading any value at all.
// This is suitable for test scenarios or to disable caching.
func NotFound[K comparable, V any](value V) LoadFunc[K, V] {
	return func(_ context.Context, _ K) (V, error) {
		return value, ErrNotFound
	}
}

// Cache defines the basic cache interface used
// to retrieve a cached value for a key.
type Cache[K comparable, V any] interface {
	// Get gets the cached value for the given key.
	// 2nd argument indicates whether the value
	// was found (hit) or not (no-hit).
	Get(ctx context.Context, key K) (V, bool)
}

// NoCache provides a cache implementation
// without any caching at all.
// This is suitable for test scenarios or to disable caching.
type NoCache[K comparable, V any] struct {
	Value V
	Found bool
}

// Get imlements [Cache.Get].
func (c *NoCache[K, V]) Get(ctx context.Context, key K) (V, bool) {
	return c.Value, c.Found
}

// KeyValue defines the key-value cache type.
type KeyValue[K comparable, V any] interface {
	Cache[K, V]
	// Put adds key-value association to the cache.
	// Any previously established key-value association
	// is overwritten.
	Put(ctx context.Context, key K, value V)
	io.Closer
}
