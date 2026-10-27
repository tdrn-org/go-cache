//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package memory

import (
	"github.com/maypok86/otter/v2"
	"github.com/tdrn-org/go-cache"
)

func mapErrNotFound[V any](value V, err error) (V, error) {
	switch err {
	case cache.ErrNotFound:
		return value, otter.ErrNotFound
	case otter.ErrNotFound:
		return value, cache.ErrNotFound
	default:
		return value, err
	}
}

func mapHit[V any](value V, hit bool) (V, error) {
	if !hit {
		return value, cache.ErrNotFound
	}
	return value, nil
}
