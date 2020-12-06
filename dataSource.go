package durcov

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// DataSource describes an API for fetching and decoding expected data
type DataSource interface {
	UseURL(url string) error
	FetchData() (*Data, error)
}

// CovidAPI represents expected covid data from a chosen source
type CovidAPI struct {
	url *url.URL
}

// UseURL sets the endpoint to be used while fetching data.
// Must be set before calling FetchData
func (c *CovidAPI) UseURL(covidURL string) error {
	u, err := url.ParseRequestURI(covidURL)
	if err != nil {
		return err
	}
	c.url = u
	return nil
}

// FetchData makes a get request on the set endpoint and decodes data into the expected format
func (c *CovidAPI) FetchData() (*Data, error) {
	if c.url == nil {
		return nil, errors.New("No URL set to fetch data")
	}
	resp, err := http.Get(c.url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return unmarshalData(body)
}

func unmarshalData(rawData []byte) (*Data, error) {
	data := Data{}

	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
