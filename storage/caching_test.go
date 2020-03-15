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
	err := subject.Set(responseKey, CacheData{})

	assert.Equal(t, nil, err)
}

func TestResponseSetDbFails(t *testing.T) {
	mockDb := &mocks.DB{}
	addHMErr := errors.New("boom")
	mockDb.On("AddHM", EntityCache, responseKey, mock.AnythingOfType("*storage.CacheData")).Return(addHMErr)

	subject := &Storage{mockDb}
	err := subject.Set(responseKey, CacheData{})

	assert.Equal(t, addHMErr, err)

}

func TestResponseSetExistingKey(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("AddHM", EntityCache, responseKey, mock.AnythingOfType("*storage.CacheData")).Return(nil)

	subject := &Storage{mockDb}

	err := subject.Set(responseKey, CacheData{})
	assert.Equal(t, nil, err)

	err = subject.Set(responseKey, CacheData{})
	assert.Equal(t, nil, err)
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
	err := subject.Delete(responseKey)

	assert.Equal(t, nil, err)
}

func TestResponseDeleteDbFails(t *testing.T) {
	mockDb := &mocks.DB{}
	addHMErr := errors.New("boom")
	mockDb.On("DeleteHM", EntityCache, responseKey).Return(addHMErr)

	subject := &Storage{mockDb}
	err := subject.Delete(responseKey)

	assert.Equal(t, addHMErr, err)
}

func TestResponseGetAndExpiredWithDeletion(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("GetHMValue", EntityCache, responseKey, mock.AnythingOfType("*storage.CacheData")).Return(nil)
	mockDb.On("DeleteHM", EntityCache, responseKey).Return(nil)

	subject := &Storage{mockDb}
	res, err := subject.Get(responseKey)

	assert.Equal(t, nil, err)
	assert.Equal(t, CacheData{}, res)
	assert.True(t, res.Validate(1, 10))

	if res.Validate(1, 10) {
		err = subject.Delete(responseKey)
	}

	assert.Equal(t, nil, err)
}
