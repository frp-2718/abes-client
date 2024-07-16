package abes

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAbesClient(t *testing.T) {
	endpoint := "https://www.sudoc.fr/services/multiwhere/"
	ac := NewAbesClient(nil)
	assert := assert.New(t)
	assert.NotNil(ac)
	assert.NotNil(ac.Multiwhere.client)
	assert.Equal(ac.Multiwhere.endpoint, endpoint, "Wrong endpoint")
	assert.Equal(ac.Multiwhere.max_ppns, MAX_MULTIWHERE_PPNS)

	client := &http.Client{
		Timeout: 4 * time.Second,
	}
	ac = NewAbesClient(client)
	assert.NotNil(ac)
	assert.NotNil(ac.Multiwhere.client)
	assert.Same(ac.Multiwhere.client, client, "Custom HTTP client not registered")
}
