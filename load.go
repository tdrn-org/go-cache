//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package cache

import "context"

// LoadFunc defines the signature for functions responsible for
// loading cache values based on the value's key.
type LoadFunc[K any, V any] func(ctx context.Context, key K) (V, error)

// NotFound is a [LoadFunc] not loading any value at all.
// This is suitable for test scenarios or to disable caching.
func NotFound[K any, V any]() LoadFunc[K, V] {
	return func(_ context.Context, _ K) (V, error) {
		var value V
		return value, ErrNotFound
	}
}
