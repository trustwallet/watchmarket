package coingecko

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("api", "USD", 500)
	assert.NotNil(t, c)
	assert.Equal(t, "api", c.BaseUrl)
	assert.Equal(t, "USD", c.currency)
	assert.Equal(t, 500, c.bucketSize)
}
