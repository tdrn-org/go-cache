//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package cache_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-cache"
	"github.com/tdrn-org/go-cache/memory"
	"github.com/tdrn-org/go-cache/redis"
)

func TestMemoryKeyValue(t *testing.T) {
	kv, err := memory.NewKeyValue(0, time.Second, cache.NotFound[string, string]())
	require.NoError(t, err)

	runKeyValueTest(t, kv)

	err = kv.Close()
	require.NoError(t, err)
}

func TestRedisKeyValue(t *testing.T) {
	options := redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	}
	if options.Addr == "" {
		t.Skip("No Redis Addr set; skipping tests")
	}

	kv, err := redis.NewKeyValue(&options, 0, redis.StringKey, cache.StringSerializer())
	require.NoError(t, err)

	runKeyValueTest(t, kv)

	err = kv.Close()
	require.NoError(t, err)
}

func TestRedisAccessTTL(t *testing.T) {
	options := redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	}
	if options.Addr == "" {
		t.Skip("No Redis Addr set; skipping tests")
	}

	const ttl = 500 * time.Millisecond
	kv, err := redis.NewKeyValue(&options, ttl, redis.StringKey, cache.StringSerializer())
	require.NoError(t, err)
	defer kv.Close()

	ctx := t.Context()
	kv.Put(ctx, "access-ttl", "access-ttl")

	// A positive ttl is an access ttl: every read refreshes the expiry, so the
	// entry must survive well beyond ttl as long as it keeps being accessed.
	for range 4 {
		time.Sleep(ttl / 2)
		value, err := kv.Get(ctx, "access-ttl")
		require.NoError(t, err)
		require.Equal(t, "access-ttl", value)
	}

	// Once it is no longer accessed, the entry must be discarded after ttl.
	time.Sleep(ttl + 100*time.Millisecond)
	_, err = kv.Get(ctx, "access-ttl")
	require.ErrorIs(t, err, cache.ErrNotFound)
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
