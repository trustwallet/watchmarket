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

func TestParseID(t *testing.T) {
	testStruct := []struct {
		givenID     string
		wantedCoin  uint
		wantedToken string
		wantedType  CoinType
		wantedError error
	}{
		{"714_TWT-8C2",
			714,
			"TWT-8C2",
			Token,
			nil,
		},
	}

	for _, tt := range testStruct {
		coin, token, givenType, err := ParseID(tt.givenID)
		assert.Equal(t, tt.wantedCoin, coin)
		assert.Equal(t, tt.wantedToken, token)
		assert.Equal(t, tt.wantedType, givenType)
		assert.Equal(t, tt.wantedError, err)
	}
}

func TestBuildID(t *testing.T) {
	testStruct := []struct {
		wantedID   string
		givenCoin  uint
		givenToken string
	}{
		{"714_TWT-8C2",
			714,
			"TWT-8C2",
		},
		{"60",
			60,
			"",
		},
		{"0",
			0,
			"",
		},
		{"0_:fnfjunwpiucU#*0! 02",
			0,
			":fnfjunwpiucU#*0! 02",
		},
	}

	for _, tt := range testStruct {
		id := BuildID(tt.givenCoin, tt.givenToken)
		assert.Equal(t, tt.wantedID, id)
	}
}
