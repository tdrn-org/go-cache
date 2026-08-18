//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package cache_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-cache"
)

func TestNotFound(t *testing.T) {
	load := cache.NotFound[string, string]()
	cached, err := load(t.Context(), t.Name())
	require.ErrorIs(t, err, cache.ErrNotFound)
	require.Equal(t, "", cached)
}
