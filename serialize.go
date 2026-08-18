//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package cache

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Serializer interface is used to marshal and unmarshal objects
// for byte based cache backends.
type Serializer[V any] interface {
	// Marshal converts a value to a byte array.
	Marshal(v V) ([]byte, error)
	// Unmarshal restores a marshaled value from a byte array.
	Unmarshal(data []byte) (V, error)
}

type bytesSerializer struct{}

func (*bytesSerializer) Marshal(v []byte) ([]byte, error) {
	return v, nil
}

func (*bytesSerializer) Unmarshal(data []byte) ([]byte, error) {
	return data, nil
}

// BytesSerializer creates an identity [Serializer] instance.
func BytesSerializer() Serializer[[]byte] {
	return &bytesSerializer{}
}

type stringSerializer struct{}

func (*stringSerializer) Marshal(v string) ([]byte, error) {
	return []byte(v), nil
}

func (*stringSerializer) Unmarshal(data []byte) (string, error) {
	return string(data), nil
}

// StringSerializer creates a [Serializer] instance serializing
// strings to byte arrays and vice versa.
func StringSerializer() Serializer[string] {
	return &stringSerializer{}
}

type jsonSerializer[V any] struct{}

func (*jsonSerializer[V]) Marshal(v V) ([]byte, error) {
	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(v)
	if err != nil {
		return nil, fmt.Errorf("failed to encode JSON (cause: %w)", err)
	}
	return buffer.Bytes(), nil
}

func (*jsonSerializer[V]) Unmarshal(data []byte) (V, error) {
	var v V
	err := json.NewDecoder(bytes.NewReader(data)).Decode(&v)
	if err != nil {
		return v, fmt.Errorf("failed to decode JSON (cause: %w)", err)
	}
	return v, nil
}

// JSONSerializer create a [Serializer] instance serialzing
// JSON objects to byte arrays and vice versa.
func JSONSerializer[V any]() Serializer[V] {
	return &jsonSerializer[V]{}
}
