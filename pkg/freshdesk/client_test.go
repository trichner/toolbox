package freshdesk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient_BadConfig(t *testing.T) {
	_, err := NewClient()

	assert.NotNil(t, err)
}

func TestNewClient(t *testing.T) {
	client, err := NewClient(BaseUrl("https://dundermifflin.freshdesk.com"), Token("123"))

	assert.NoError(t, err)
	assert.NotNil(t, client)
}
