package durcov

import (
	"encoding/json"
	"time"
)

// Data represents a combination of global and country based statistics
type Data struct {
	global    *global
	countries []*country
}

type country struct {
	name  string
	slug  string
	code  string
	stats *statistics
}

type global struct {
	stats *statistics
}

type statistics struct {
	totalConfirmed int64
	totalDeaths    int64
	totalRecovered int64
	date           time.Time
}

type expectedJSONCountryShape struct {
	Name           string    `json:"Country"`
	Slug           string    `json:"Slug"`
	Code           string    `json:"CountryCode"`
	TotalConfirmed int64     `json:"TotalConfirmed"`
	TotalDeaths    int64     `json:"TotalDeaths"`
	TotalRecovered int64     `json:"TotalRecovered"`
	Date           time.Time `json:"Date"`
}

type expectedJSONShape struct {
	Global    map[string]interface{}      `json:"Global"`
	Countries []*expectedJSONCountryShape `json:"Countries"`
	Date      time.Time                   `json:"Date"`
}

// UnmarshalJSON complies with the json package for custom decoding
func (d *Data) UnmarshalJSON(data []byte) error {
	var ingress expectedJSONShape
	if err := json.Unmarshal(data, &ingress); err != nil {
		return err
	}

	globalStats := &statistics{
		totalConfirmed: int64(ingress.Global["TotalConfirmed"].(float64)),
		totalDeaths:    int64(ingress.Global["TotalDeaths"].(float64)),
		totalRecovered: int64(ingress.Global["TotalRecovered"].(float64)),
		date:           ingress.Date,
	}

	countries := []*country{}
	for _, countryData := range ingress.Countries {
		countryStats := &statistics{
			totalConfirmed: countryData.TotalConfirmed,
			totalDeaths:    countryData.TotalDeaths,
			totalRecovered: countryData.TotalRecovered,
			date:           countryData.Date,
		}
		country := &country{
			name:  countryData.Name,
			slug:  countryData.Slug,
			code:  countryData.Code,
			stats: countryStats,
		}
		countries = append(countries, country)
	}

	d.global = &global{globalStats}
	d.countries = countries
	return nil
}
