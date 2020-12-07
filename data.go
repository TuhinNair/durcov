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

type expectedJSONShape struct {
	Global    expectedJSONGlobalShape     `json:"Global"`
	Countries []*expectedJSONCountryShape `json:"Countries"`
	Date      time.Time                   `json:"Date"`
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

type expectedJSONGlobalShape struct {
	TotalConfirmed int64
	TotalDeaths    int64
	TotalRecovered int64
}

// UnmarshalJSON complies with the json package for custom decoding
func (d *Data) UnmarshalJSON(data []byte) error {
	var ingress expectedJSONShape
	if err := json.Unmarshal(data, &ingress); err != nil {
		return err
	}

	globalStats := &statistics{
		// When decoding the interface{} type go will cast all numeric types to float64 by default
		totalConfirmed: ingress.Global.TotalConfirmed,
		totalDeaths:    ingress.Global.TotalDeaths,
		totalRecovered: ingress.Global.TotalRecovered,
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
