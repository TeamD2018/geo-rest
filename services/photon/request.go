package photon

import (
	"net/url"
	"strconv"
)

type SearchQuery struct {
	Query string
	Tags  []string
	Limit int
}

func NewSearchQuery(query string, limit int, tags ...string) *SearchQuery {
	return &SearchQuery{
		Query: query,
		Tags:  tags,
		Limit: limit,
	}
}

func (s *SearchQuery) AddTag(tag string) {
	if s.Tags != nil {
		s.Tags = make([]string, 1)
	}
	s.Tags = append(s.Tags, tag)
}

func (s *SearchQuery) String() string {
	values := url.Values{}
	values.Add("q", s.Query)
	for _, tag := range s.Tags {
		values.Add("osm_tag", tag)
	}
	values.Add("limit", strconv.Itoa(s.Limit))
	return values.Encode()
}
