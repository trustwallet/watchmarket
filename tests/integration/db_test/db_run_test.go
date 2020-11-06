// +build integration

package db_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/db/postgres"
	"github.com/trustwallet/watchmarket/tests/integration/setup"
	"os"
	"testing"
)

var databaseInstance *postgres.Instance

func TestMain(m *testing.M) {
	databaseInstance = setup.RunPgContainer()
	code := m.Run()
	setup.StopPgContainer()
	os.Exit(code)
}

func TestPgSetup(t *testing.T) {
	assert.NotNil(t, databaseInstance)
	assert.NotNil(t, databaseInstance.Gorm)
}
