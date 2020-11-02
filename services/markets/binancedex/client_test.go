package binancedex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("demo.api")
	assert.NotNil(t, client)
	assert.Equal(t, "demo.api", client.baseURL)
}
