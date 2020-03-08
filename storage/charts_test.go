package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mocks "github.com/trustwallet/watchmarket/mocks/storage"
	"testing"
)

const keyCharts = "KeyCharts"

func TestSaveChartsDbFails(t *testing.T) {
	mockDb := &mocks.DB{}
	dbFailure := errors.New("boom")
	mockDb.On("AddHM", EntityCharts, keyCharts, mock.AnythingOfType("**storage.ChartData")).Return(dbFailure)
	mockRates := ChartData{}
	subject := &Storage{mockDb}
	result, err := subject.SaveCharts(keyCharts, &mockRates)
	assert.Equal(t, result, SaveResultAddHMFailure)
	assert.Equal(t, err, dbFailure)
}

func TestSaveChartsDbSuccess(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("AddHM", EntityCharts, keyCharts, mock.AnythingOfType("**storage.ChartData")).Return(nil)
	mockRates := ChartData{}
	subject := &Storage{mockDb}
	result, err := subject.SaveCharts(keyCharts, &mockRates)
	assert.Equal(t, result, SaveResultSuccess)
	assert.Equal(t, err, nil)
}

func TestGetChartsDbFails(t *testing.T) {
	mockDb := &mocks.DB{}
	dbFailure := errors.New("boom")
	mockDb.On("GetHMValue", EntityCharts, keyCharts, mock.AnythingOfType("*storage.ChartData")).Return(dbFailure)
	subject := &Storage{mockDb}
	_, err := subject.GetCharts(keyCharts)
	assert.Equal(t, dbFailure, err)
}

func TestGetChartsSuccess(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("GetHMValue", EntityCharts, keyCharts, mock.AnythingOfType("*storage.ChartData")).Return(nil)
	subject := &Storage{mockDb}
	res, err := subject.GetCharts(keyCharts)
	assert.Nil(t, err)
	assert.Equal(t, &ChartData{}, res)
}
