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

func TestBytesSerializer(t *testing.T) {
	data := []byte(t.Name())
	serializer := cache.BytesSerializer()

	encoded, err := serializer.Marshal(data)
	require.NoError(t, err)
	require.Equal(t, data, encoded)

	decoded, err := serializer.Unmarshal(encoded)
	require.NoError(t, err)
	require.Equal(t, data, decoded)
}

func TestStringSerializer(t *testing.T) {
	data := t.Name()
	serializer := cache.StringSerializer()

	encoded, err := serializer.Marshal(data)
	require.NoError(t, err)
	decoded, err := serializer.Unmarshal(encoded)
	require.NoError(t, err)
	require.Equal(t, data, decoded)
}

type testJSON struct {
	S string `json:"s"`
	I int    `json:"i"`
}

func TestJSONSerializer(t *testing.T) {
	data := &testJSON{
		S: t.Name(),
		I: len(t.Name()),
	}
	serializer := cache.JSONSerializer[*testJSON]()

	encoded, err := serializer.Marshal(data)
	require.NoError(t, err)
	decoded, err := serializer.Unmarshal(encoded)
	require.NoError(t, err)
	require.Equal(t, data, decoded)
}
