package coingecko

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("api", 500)
	assert.NotNil(t, c)
	assert.Equal(t, "api", c.baseURL)
	assert.Equal(t, 500, c.bucketSize)
}
