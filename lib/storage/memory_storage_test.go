package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	// Test getting value that hasn't been set
	_, exists := storage.Get("foo")
	assert.False(t, exists)

	// Test setting value
	err := storage.Set("foo", "bar")
	assert.NoError(t, err)

	// Test getting back value
	value, exists := storage.Get("foo")
	assert.True(t, exists)
	assert.Equal(t, "bar", value)

	// Test case sensitivity
	_, exists = storage.Get("Foo")
	assert.False(t, exists)
	err = storage.Set("Foo", "baz")
	assert.NoError(t, err)

	value, exists = storage.Get("Foo")
	assert.True(t, exists)
	assert.Equal(t, "baz", value)
}
