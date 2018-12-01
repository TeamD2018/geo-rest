package photon

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL *url.URL
}

func NewPhotonClient(baseurl string) *Client {
	url, err := url.Parse(baseurl)
	if err != nil {
		panic(err)
	}
	return &Client{
		baseURL: url,
	}
}

func (p *Client) Search(query *SearchQuery) ([]byte, error) {
	url := *p.baseURL
	url.RawQuery = query.String()
	queryString := url.String()
	response, err := http.Get(queryString)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
