package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mocks "github.com/trustwallet/watchmarket/mocks/storage"
	"testing"
)

const (
	responseKey = "TEST_KEY"
)

func TestResponseSet(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("AddHM", EntityCache, responseKey, mock.AnythingOfType("*storage.CacheData")).Return(nil)

	subject := &Storage{mockDb}
	res, err := subject.Set(responseKey, CacheData{})

	assert.Equal(t, nil, err)
	assert.Equal(t, SaveResultSuccess, res)
}

func TestResponseSetDbFails(t *testing.T) {
	mockDb := &mocks.DB{}
	addHMErr := errors.New("boom")
	mockDb.On("AddHM", EntityCache, responseKey, mock.AnythingOfType("*storage.CacheData")).Return(addHMErr)

	subject := &Storage{mockDb}
	res, err := subject.Set(responseKey, CacheData{})

	assert.Equal(t, addHMErr, err)
	assert.Equal(t, SaveResultStorageFailure, res)
}

func TestResponseSetExistingKey(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("AddHM", EntityCache, responseKey, mock.AnythingOfType("*storage.CacheData")).Return(nil)

	subject := &Storage{mockDb}
	res, err := subject.Set(responseKey, CacheData{})

	assert.Equal(t, nil, err)
	assert.Equal(t, SaveResultSuccess, res)

	resTwo, err := subject.Set(responseKey, CacheData{})

	assert.Equal(t, nil, err)
	assert.Equal(t, SaveResultSuccess, resTwo)

}

func TestResponseGet(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("GetHMValue", EntityCache, responseKey, mock.AnythingOfType("*storage.CacheData")).Return(nil)

	subject := &Storage{mockDb}
	res, err := subject.Get(responseKey)

	assert.Equal(t, nil, err)
	assert.Equal(t, CacheData{}, res)
}

func TestResponseGetDbFails(t *testing.T) {
	mockDb := &mocks.DB{}
	addHMErr := errors.New("boom")
	mockDb.On("GetHMValue", EntityCache, responseKey, mock.AnythingOfType("*storage.CacheData")).Return(addHMErr)

	subject := &Storage{mockDb}
	res, err := subject.Get(responseKey)

	assert.Equal(t, addHMErr, err)
	assert.Equal(t, CacheData{}, res)

}

func TestResponseGetExistingKey(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("GetHMValue", EntityCache, responseKey, mock.AnythingOfType("*storage.CacheData")).Return(nil)

	subject := &Storage{mockDb}
	res, err := subject.Get(responseKey)

	assert.Equal(t, nil, err)
	assert.Equal(t, CacheData{}, res)

	resTwo, err := subject.Get(responseKey)

	assert.Equal(t, nil, err)
	assert.Equal(t, CacheData{}, resTwo)
}

func TestResponseDelete(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("DeleteHM", EntityCache, responseKey).Return(nil)

	subject := &Storage{mockDb}
	result, err := subject.Delete(responseKey)

	assert.Equal(t, nil, err)
	assert.Equal(t, SaveResultSuccess, result)
}

func TestResponseDeleteDbFails(t *testing.T) {
	mockDb := &mocks.DB{}
	addHMErr := errors.New("boom")
	mockDb.On("DeleteHM", EntityCache, responseKey).Return(addHMErr)

	subject := &Storage{mockDb}
	result, err := subject.Delete(responseKey)

	assert.Equal(t, addHMErr, err)
	assert.Equal(t, SaveResultStorageFailure, result)
}

func TestResponseGetAndExpiredWithDeletion(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("GetHMValue", EntityCache, responseKey, mock.AnythingOfType("*storage.CacheData")).Return(nil)
	mockDb.On("DeleteHM", EntityCache, responseKey).Return(nil)

	subject := &Storage{mockDb}
	res, err := subject.Get(responseKey)

	assert.Equal(t, nil, err)
	assert.Equal(t, CacheData{}, res)
	assert.True(t, res.IsExpired())

	var result SaveResult
	if res.IsExpired() {
		result, err = subject.Delete(responseKey)
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, SaveResultSuccess, result)
}
