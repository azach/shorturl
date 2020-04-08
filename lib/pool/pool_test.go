package pool

import (
	"testing"

	"github.com/azach/shorturl/lib/storage"
	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	memStorage := storage.NewMemoryStorage()
	pool := NewPool(memStorage, &Options{minPoolSize: 2, minPoolGenerationSize: 0})

	// Empty pool should still generate a value
	assert.NotEqual(t, "", pool.Get())

	// Deplete pool
	pool.Get()
	assert.Equal(t, 0, len(pool.queue))

	// Generate should replenish pool to minPoolSize
	pool.Generate()
	assert.Equal(t, 2, len(pool.queue))

	// Future addition: mock out ID generator to test
	// that existing keys in storage are skipped
}
