package services

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"sync"
)

type ConcurrentLookupService struct {
	LookupService interfaces.LookupInterface
	Concurrency   int
}

func (ls *ConcurrentLookupService) LookupAll(
	source <-chan *models.OSMEntity,
	destination chan<- *models.OSMPolygonSuggestion,
	errc chan<- error) func() {
	var wg sync.WaitGroup
	wg.Add(ls.Concurrency)
	for i := 0; i < ls.Concurrency; i++ {
		go ls.lookup(source, destination, errc, &wg)
	}
	return func() {
		wg.Wait()
		close(destination)
		close(errc)
	}
}

func (ls *ConcurrentLookupService) lookup(
	source <-chan *models.OSMEntity,
	destination chan<- *models.OSMPolygonSuggestion,
	errc chan<- error,
	wg *sync.WaitGroup) {
	defer wg.Done()
	for entity := range source {
		resolvedName, err := ls.LookupService.Lookup(entity)
		if err != nil {
			select {
			case errc <- err:
			default:
			}
		} else {
			destination <- &models.OSMPolygonSuggestion{OSMID: int64(entity.OSMID), OSMType: entity.OSMType, Name: resolvedName}
		}
	}
}
