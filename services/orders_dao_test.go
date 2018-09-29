package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestOrdersElasticDAO_EnsureMapping(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	c := NewOrdersElasticDAO(client, logger, OrdersIndex)

	err := c.EnsureMapping()
	assert.Nil(t, err)

	exists, err := client.IndexExists(OrdersIndex).Do(context.Background())
	assert.Nil(t, err)
	assert.True(t, exists)
}
