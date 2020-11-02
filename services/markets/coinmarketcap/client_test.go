package coinmarketcap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("pro.api", "web.api", "widget.api", "key")
	assert.NotNil(t, c)
	assert.Equal(t, "pro.api", c.proApiURL)
	assert.Equal(t, "web.api", c.webApiURL)
	assert.Equal(t, "widget.api", c.widgetApiURL)
}
