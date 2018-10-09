package interfaces

import (
	"context"
	"github.com/TeamD2018/geo-rest/models"
)

type GeoResolver interface {
	Resolve(location *models.Location, ctx context.Context) error
}
