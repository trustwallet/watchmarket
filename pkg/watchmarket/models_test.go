package watchmarket

import (
	"github.com/stretchr/testify/assert"
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

func TestIsFiatRate(t *testing.T) {
	assert.False(t, IsFiatRate("BTC"))
	assert.True(t, IsFiatRate("USD"))
}

func TestChart_IsEmpty(t *testing.T) {
	var emptyChart Chart
	assert.True(t, emptyChart.IsEmpty())
	emptyChart.Prices = []ChartPrice{{Price: 1}}
	assert.False(t, emptyChart.IsEmpty())
}

func TestIsRespectableValue(t *testing.T) {
	assert.True(t, IsRespectableValue(1, 0))
	assert.False(t, IsRespectableValue(0, 100))
	assert.True(t, IsRespectableValue(100, 100))
}

func TestIsSuitableUpdateTime(t *testing.T) {
	assert.True(t, IsSuitableUpdateTime(time.Now(), time.Second*10))
	assert.False(t, IsSuitableUpdateTime(time.Date(1999, 1, 1, 1, 1, 1, 1, time.Local), time.Nanosecond))
}

func TestTruncateWithPrecision(t *testing.T) {
	assert.Equal(t, 6.11, TruncateWithPrecision(6.111, 2))
	assert.Equal(t, float64(6), TruncateWithPrecision(6.111, 0))
}

func TestCoinDetails_IsEmpty(t *testing.T) {
	emptyDetails := CoinDetails{}
	assert.True(t, emptyDetails.IsEmpty())
	var i Info
	i.Name = "1"
	emptyDetails.Info = &i
	assert.False(t, emptyDetails.IsEmpty())
}
