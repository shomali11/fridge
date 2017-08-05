package fridge

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetrievalDetails_New(t *testing.T) {
	retrievalDetails := &RetrievalDetails{}

	retrievalOption := WithRestock(func() (string, error) {
		return "Hi", nil
	})

	retrievalOption(retrievalDetails)

	assert.NotNil(t, retrievalDetails.Restock)

	value, err := retrievalDetails.Restock()

	assert.Equal(t, value, "Hi")
	assert.Nil(t, err)
}

func TestRetrievalDetails_Defaults(t *testing.T) {
	retrievalDetails := newRetrievalDetails()

	assert.Nil(t, retrievalDetails.Restock)
}

func TestRetrievalDetails_Override(t *testing.T) {
	retrievalDetails := newRetrievalDetails(WithRestock(func() (string, error) {
		return "Hi", nil
	}))

	assert.NotNil(t, retrievalDetails.Restock)

	value, err := retrievalDetails.Restock()

	assert.Equal(t, value, "Hi")
	assert.Nil(t, err)
}
