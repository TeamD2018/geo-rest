package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCouriersElasticDAO_EnsureMapping(t *testing.T) {
	c := NewCouriersDAO(client, CourierIndex)

	err := c.EnsureMapping()
	assert.Nil(t, err)

	exists, err := client.IndexExists(CourierIndex).Do(context.Background())
	assert.Nil(t, err)
	assert.True(t, exists)
}
