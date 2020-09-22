package coinmarketcap

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MappingTest []struct {
	Coin    int    `json:"coin"`
	Type    string `json:"type"`
	ID      int    `json:"id"`
	TokenID string `json:"token_id,omitempty"`
}

func TestMapping(t *testing.T) {
	var r MappingTest
	err := json.Unmarshal([]byte(Mapping), &r)
	assert.Nil(t, err)
	assert.NotNil(t, r)
}
