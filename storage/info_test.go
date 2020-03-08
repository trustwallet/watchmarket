package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mocks "github.com/trustwallet/watchmarket/mocks/storage"
	"testing"
)

const keyInfo = "KeyInfo"

func TestSaveInfoDbFails(t *testing.T) {
	mockDb := &mocks.DB{}
	dbFailure := errors.New("boom")
	mockDb.On("AddHM", EntityInfo, keyInfo, mock.AnythingOfType("**storage.CoinInfo")).Return(dbFailure)
	mockRates := CoinInfo{}
	subject := &Storage{mockDb}
	result, err := subject.SaveInfo(keyInfo, &mockRates)
	assert.Equal(t, result, SaveResultAddHMFailure)
	assert.Equal(t, err, dbFailure)
}

func TestSaveInfoDbSuccess(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("AddHM", EntityInfo, keyInfo, mock.AnythingOfType("**storage.CoinInfo")).Return(nil)
	mockRates := CoinInfo{}
	subject := &Storage{mockDb}
	result, err := subject.SaveInfo(keyInfo, &mockRates)
	assert.Equal(t, result, SaveResultSuccess)
	assert.Equal(t, err, nil)
}

func TestGetInfoDbFails(t *testing.T) {
	mockDb := &mocks.DB{}
	dbFailure := errors.New("boom")
	mockDb.On("GetHMValue", EntityInfo, keyInfo, mock.AnythingOfType("*storage.CoinInfo")).Return(dbFailure)
	subject := &Storage{mockDb}
	_, err := subject.GetInfo(keyInfo)
	assert.Equal(t, dbFailure, err)
}

func TestGetInfoSuccess(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("GetHMValue", EntityInfo, keyInfo, mock.AnythingOfType("*storage.CoinInfo")).Return(nil)
	subject := &Storage{mockDb}
	res, err := subject.GetInfo(keyInfo)
	assert.Nil(t, err)
	assert.Equal(t, &CoinInfo{}, res)
}
