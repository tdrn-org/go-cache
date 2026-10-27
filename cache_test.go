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
	// Put, Get, Delete tests
	kv := memory.MustNewKeyValue(0, 0, cache.NotFound[string, string]())
	runKeyValuePutGetDeleteTest(t, kv)
	require.NoError(t, kv.Close())

	// Access expiry test
	accessExpiryTTL := 500 * time.Millisecond
	kv = memory.MustNewKeyValue(0, accessExpiryTTL, cache.NotFound[string, string]())
	runKeyValueAccessExpiryTest(t, kv, accessExpiryTTL)
	require.NoError(t, kv.Close())

	// Create expiry test
	createExpiryTTL := -500 * time.Millisecond
	kv = memory.MustNewKeyValue(0, createExpiryTTL, cache.NotFound[string, string]())
	runKeyValueCreateExpiryTest(t, kv, accessExpiryTTL)
	require.NoError(t, kv.Close())
}

func TestRedisKeyValue(t *testing.T) {
	options := redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	}
	if options.Addr == "" {
		t.Skip("No Redis Addr set; skipping tests")
	}

	// Put, Get, Delete tests
	kv := redis.MustNewKeyValue(&options, 0, redis.StringKey, cache.StringSerializer())
	runKeyValuePutGetDeleteTest(t, kv)
	require.NoError(t, kv.Close())

	// Access expiry test
	accessExpiryTTL := 500 * time.Millisecond
	kv = redis.MustNewKeyValue(&options, accessExpiryTTL, redis.StringKey, cache.StringSerializer())
	runKeyValueAccessExpiryTest(t, kv, accessExpiryTTL)
	require.NoError(t, kv.Close())

	// Create expiry test
	createExpiryTTL := -500 * time.Millisecond
	kv = redis.MustNewKeyValue(&options, createExpiryTTL, redis.StringKey, cache.StringSerializer())
	runKeyValueCreateExpiryTest(t, kv, accessExpiryTTL)
	require.NoError(t, kv.Close())
}

func runKeyValuePutGetDeleteTest(t *testing.T, kv cache.KeyValue[string, string]) {
	for keyValue := range 100 {
		key := fmt.Sprintf("%d", keyValue)
		kv.Put(t.Context(), key, key)
		value, err := kv.Get(t.Context(), key)
		require.NoError(t, err)
		require.Equal(t, key, value)
		kv.Delete(t.Context(), key)
	}
}

func runKeyValueAccessExpiryTest(t *testing.T, kv cache.KeyValue[string, string], ttl time.Duration) {
	key := t.Name() + "Key"
	value := t.Name() + "Value"
	kv.Put(t.Context(), key, value)
	for range 10 {
		time.Sleep(ttl / 2)
		cached, err := kv.Get(t.Context(), key)
		require.NoError(t, err)
		require.Equal(t, value, cached)
	}
	time.Sleep(2 * ttl)
	cached, err := kv.Get(t.Context(), key)
	require.ErrorIs(t, err, cache.ErrNotFound)
	require.Empty(t, cached)
}

func runKeyValueCreateExpiryTest(t *testing.T, kv cache.KeyValue[string, string], ttl time.Duration) {
	key := t.Name() + "Key"
	value := t.Name() + "Value"
	kv.Put(t.Context(), key, value)
	time.Sleep(ttl / 2)
	cached, err := kv.Get(t.Context(), key)
	require.NoError(t, err)
	require.Equal(t, value, cached)
	time.Sleep(ttl)
	cached, err = kv.Get(t.Context(), key)
	require.ErrorIs(t, err, cache.ErrNotFound)
	require.Empty(t, cached)
}
