// +build integration

package integration

import (
	"github.com/trustwallet/watchmarket/tests/integration/setup"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup.RunRedisContainer()
	code := m.Run()
	setup.StopRedisContainer()
	os.Exit(code)
}
