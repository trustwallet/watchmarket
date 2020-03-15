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
	mockDb.On("AddHM", EntityCache, responseKey, &[]byte{}).Return(nil)

	subject := &Storage{mockDb}
	err := subject.Set(responseKey, []byte{})

	assert.Equal(t, nil, err)
}

func TestResponseSetDbFails(t *testing.T) {
	mockDb := &mocks.DB{}
	addHMErr := errors.New("boom")
	mockDb.On("AddHM", EntityCache, responseKey, &[]byte{}).Return(addHMErr)

	subject := &Storage{mockDb}
	err := subject.Set(responseKey, []byte{})

	assert.Equal(t, addHMErr, err)

}

func TestResponseGet(t *testing.T) {
	mockDb := &mocks.DB{}

	mockDb.On("GetHMValue", EntityCache, responseKey, mock.Anything).Return(nil)

	subject := &Storage{mockDb}
	res, err := subject.Get(responseKey)

	assert.Equal(t, nil, err)
	assert.Equal(t, string([]byte{}), string(res))
}

func TestResponseGetDbFails(t *testing.T) {
	mockDb := &mocks.DB{}
	addHMErr := errors.New("boom")
	mockDb.On("GetHMValue", EntityCache, responseKey, mock.Anything).Return(addHMErr)

	subject := &Storage{mockDb}
	res, err := subject.Get(responseKey)

	assert.Equal(t, addHMErr, err)
	assert.Equal(t, string([]byte{}), string(res))
}

func TestResponseGetExistingKey(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("GetHMValue", EntityCache, responseKey, mock.Anything).Return(nil)

	subject := &Storage{mockDb}
	res, err := subject.Get(responseKey)

	assert.Equal(t, nil, err)
	assert.Equal(t, string([]byte{}), string(res))

	resTwo, err := subject.Get(responseKey)

	assert.Equal(t, nil, err)
	assert.Equal(t, string([]byte{}), string(resTwo))
}

func TestResponseSetExistingKey(t *testing.T) {
	mockDb := &mocks.DB{}
	mockDb.On("AddHM", EntityCache, responseKey, &[]byte{}).Return(nil)

	subject := &Storage{mockDb}

	err := subject.Set(responseKey, []byte{})
	assert.Equal(t, nil, err)

	err = subject.Set(responseKey, []byte{})
	assert.Equal(t, nil, err)
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
