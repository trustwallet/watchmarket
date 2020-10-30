package watchmarket

import (
	"errors"
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

func TestParseID(t *testing.T) {
	testStruct := []struct {
		givenID     string
		wantedCoin  uint
		wantedToken string
		wantedType  CoinType
		wantedError error
	}{
		{"c714_tTWT-8C2",
			714,
			"TWT-8C2",
			Token,
			nil,
		},
		{"tTWT-8C2_c714",
			714,
			"TWT-8C2",
			Token,
			nil,
		},
		{"c714",
			714,
			"",
			Coin,
			nil,
		},
		{"tTWT-8C2",
			0,
			"",
			Coin,
			errors.New("bad ID"),
		},
		{"c714_TWT-8C2",
			714,
			"",
			Coin,
			nil,
		},
	}

	for _, tt := range testStruct {
		coin, token, err := ParseID(tt.givenID)
		assert.Equal(t, tt.wantedCoin, coin)
		assert.Equal(t, tt.wantedToken, token)
		assert.Equal(t, tt.wantedError, err)
	}
}

func TestBuildID(t *testing.T) {
	testStruct := []struct {
		wantedID   string
		givenCoin  uint
		givenToken string
	}{
		{"c714_tTWT-8C2",
			714,
			"TWT-8C2",
		},
		{"c60",
			60,
			"",
		},
		{"c0",
			0,
			"",
		},
		{"c0_t:fnfjunwpiucU#*0! 02",
			0,
			":fnfjunwpiucU#*0! 02",
		},
	}

	for _, tt := range testStruct {
		id := BuildID(tt.givenCoin, tt.givenToken)
		assert.Equal(t, tt.wantedID, id)
	}
}

func Test_removeFirstChar(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal case", "Bob", "ob"},
		{"Empty String", "", ""},
		{"One Char String Test", "A", ""},
		{"Another normaal", "abcdef", "bcdef"},
	}

	for _, tt := range tests {
		var got = removeFirstChar(tt.input)
		if got != tt.expected {
			t.Fatalf("Got %v, Expected %v.", got, tt.expected)
		}
	}
}

func Test_findCoinID(t *testing.T) {
	tests := []struct {
		name        string
		words       []string
		expected    uint
		expectedErr error
	}{
		{"Normal case", []string{"c100", "t60", "e30"}, 100, nil},
		{"Empty coin", []string{"d100", "t60", "e30"}, 0, errors.New("no coin")},
		{"Empty words", []string{}, 0, errors.New("no coin")},
		{"Bad coin", []string{"cd100", "t60", "e30"}, 0, errors.New("bad coin")},
		{"Bad coin #2", []string{"c", "t60", "e30"}, 0, errors.New("bad coin")},
	}

	for _, tt := range tests {
		got, err := findCoinID(tt.words)
		assert.Equal(t, tt.expected, got)
		assert.Equal(t, tt.expectedErr, err)
	}
}

func Test_findTokenID(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		expected string
	}{
		{"Normal case", []string{"c100", "t60", "e30"}, "60"},
		{"Empty token", []string{"d100", "a", "e30"}, ""},
		{"Empty words", []string{}, ""},
		{"Bad token", []string{"cd100", "t", "e30"}, ""},
	}

	for _, tt := range tests {
		got := findTokenID(tt.words)
		assert.Equal(t, tt.expected, got)
	}
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
