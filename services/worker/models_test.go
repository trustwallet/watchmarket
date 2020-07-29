package worker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_isHigherPriority(t *testing.T) {
	priorities := []string{"1", "2", "3", "4"}
	assert.True(t, isHigherPriority(priorities, "3", "2"))
	assert.False(t, isHigherPriority(priorities, "3", "3"))
	assert.False(t, isHigherPriority(priorities, "3", "4"))
}
