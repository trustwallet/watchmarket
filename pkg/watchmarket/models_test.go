package watchmarket

import (
	"github.com/go-playground/assert/v2"
	"testing"
	"time"
)

func TestUnixToDuration(t *testing.T) {
	wantedDuration := time.Second * 10

	assert.Equal(t, wantedDuration, UnixToDuration(10))
	assert.Equal(t, time.Second*0, UnixToDuration(0))
	assert.Equal(t, time.Minute, UnixToDuration(60))
}

func TestDurationToUnix(t *testing.T) {
	wantedUnixTime := 10
	assert.Equal(t, uint(wantedUnixTime), DurationToUnix(time.Second*10))
	assert.Equal(t, uint(0), DurationToUnix(time.Second*0))
}
