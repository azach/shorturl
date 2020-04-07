package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache()

	// Test getting value that hasn't been set
	_, exists := cache.Get("foo")
	assert.False(t, exists)

	// Test setting value
	err := cache.Set("foo", "bar")
	assert.NoError(t, err)

	// Test getting back value
	value, exists := cache.Get("foo")
	assert.True(t, exists)
	assert.Equal(t, "bar", value)

	// Test case sensitivity
	_, exists = cache.Get("Foo")
	assert.False(t, exists)
	err = cache.Set("Foo", "baz")
	assert.NoError(t, err)

	value, exists = cache.Get("Foo")
	assert.True(t, exists)
	assert.Equal(t, "baz", value)
}
