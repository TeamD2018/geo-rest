package photon

type IPhotonClient interface {
	Search(query *SearchQuery) ([]byte, error)
}
