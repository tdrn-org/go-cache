//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package cache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-cache"
	"github.com/tdrn-org/go-cache/memory"
)

func TestMemoryKeyValue(t *testing.T) {
	kv, err := memory.NewKeyValue(0, time.Second, cache.NotFound[string](""))
	require.NoError(t, err)

	runKeyValueTest(t, kv)

	err = kv.Close()
	require.NoError(t, err)
}

func runKeyValueTest(t *testing.T, kv cache.KeyValue[string, string]) {
	const count = 1000
	for keyValue := range count {
		key := fmt.Sprintf("%d", keyValue)
		kv.Put(t.Context(), key, key)
		value, err := kv.Get(t.Context(), key)
		require.NoError(t, err)
		require.Equal(t, key, value)
		kv.Delete(t.Context(), key)
	}
}
