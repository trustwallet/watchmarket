package fixer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("demo.api", "key", "USD")
	assert.NotNil(t, client)
	assert.Equal(t, "demo.api", client.api)
	assert.Equal(t, "key", client.key)
	assert.Equal(t, "USD", client.currency)
}
