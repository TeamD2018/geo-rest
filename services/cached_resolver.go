package services

import (
	"context"
	"github.com/TeamD2018/geo-rest/models"
)

type CachedResolver struct {
	tntResolver   *TntResolver
	gmapsResolver *GMapsResolver
}

func NewCachedResolver(tntresolver *TntResolver, gmapsresolver *GMapsResolver) *CachedResolver {
	return &CachedResolver{
		tntResolver:   tntresolver,
		gmapsResolver: gmapsresolver,
	}
}

func (c *CachedResolver) Resolve(location *models.Location, ctx context.Context) error {
	if err := c.tntResolver.Resolve(location, ctx); err != nil {
		if err := c.gmapsResolver.Resolve(location, ctx); err != nil {
			return err
		}
		if err := c.tntResolver.SaveToCache(location); err != nil {
			return err
		}
	}
	return nil
}
