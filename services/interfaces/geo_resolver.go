package interfaces

import "github.com/TeamD2018/geo-rest/models"

type GeoResolver interface {
	Resolve(location *models.Location) error
}
