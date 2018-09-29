package models

type Mapper interface {
	GetMapping() (indexName string, mapping string)
	EnsureMapping() error
}
