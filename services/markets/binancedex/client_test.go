package binancedex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("demo.api")
	assert.NotNil(t, client)
	assert.Equal(t, "demo.api", client.BaseUrl)
}
