package durcov

import (
	"testing"
	"time"
)

var exampleJSON = []byte(`{
	"Global":{
	   "NewConfirmed":646333,
	   "TotalConfirmed":64520350,
	   "NewDeaths":12444,
	   "TotalDeaths":1493624,
	   "NewRecovered":461134,
	   "TotalRecovered":41488406
	},
	"Countries":[
	   {
		  "Country":"Afghanistan",
		  "CountryCode":"AF",
		  "Slug":"afghanistan",
		  "NewConfirmed":263,
		  "TotalConfirmed":46980,
		  "NewDeaths":25,
		  "TotalDeaths":1822,
		  "NewRecovered":119,
		  "TotalRecovered":37026,
		  "Date":"2020-12-04T03:49:29Z",
		  "Premium":{
			 
		  }
	   },
	   {
		  "Country":"Albania",
		  "CountryCode":"AL",
		  "Slug":"albania",
		  "NewConfirmed":705,
		  "TotalConfirmed":39719,
		  "NewDeaths":17,
		  "TotalDeaths":839,
		  "NewRecovered":528,
		  "TotalRecovered":19912,
		  "Date":"2020-12-04T03:49:29Z",
		  "Premium":{
			 
		  }
	   },
	   {
		  "Country":"Algeria",
		  "CountryCode":"DZ",
		  "Slug":"algeria",
		  "NewConfirmed":932,
		  "TotalConfirmed":85084,
		  "NewDeaths":17,
		  "TotalDeaths":2464,
		  "NewRecovered":585,
		  "TotalRecovered":54990,
		  "Date":"2020-12-04T03:49:29Z",
		  "Premium":{
			 
		  }
	   }
	],
	"Date":"2020-12-04T03:49:29Z"
 }`)

func TestJSONDecoding(t *testing.T) {
	data, err := unmarshalData(exampleJSON)
	if err != nil {
		t.Error(err)
	}

	globalData := data.global.stats

	if globalData.totalConfirmed != 64520350 {
		t.Errorf("global.totalConfirmed mismatch. Got=%d Expected=%d", globalData.totalConfirmed, 64520350)
	}
	if globalData.totalDeaths != 1493624 {
		t.Errorf("global.totalDeaths mismatch. Got=%d Expected=%d", globalData.totalDeaths, 1493624)
	}
	if globalData.totalRecovered != 41488406 {
		t.Errorf("global.totalRecovered mismatch. Got=%d Expected=%d", globalData.totalRecovered, 41488406)
	}
	tm, err := time.Parse(time.RFC3339, "2020-12-04T03:49:29Z")
	if err != nil {
		t.Error(err)
	}
	if globalData.date != tm {
		t.Errorf("global.Date mismatch. Got=%v Expected=%v", globalData.date, tm)
	}

	countryData := data.countries[2]
	if countryData == nil {
		t.Error("Expected existing country data. Should not be nil")
	}

	if countryData.name != "Algeria" {
		t.Errorf("country.Name mismatch. Got=%s Expected=%s", countryData.name, "Algeria")
	}
	if countryData.code != "DZ" {
		t.Errorf("country.CountryCode mismatch. Got=%s Expected=%s", countryData.code, "DZ")
	}
	if countryData.slug != "algeria" {
		t.Errorf("country.Slug mismatch. Got=%s Expected=%s", countryData.slug, "algeria")
	}
	if countryData.stats.totalConfirmed != 85084 {
		t.Errorf("country.totalConfirmed mismatch. Got=%d Expected=%d", countryData.stats.totalConfirmed, 85084)
	}
	if countryData.stats.totalDeaths != 2464 {
		t.Errorf("country.totalDeaths mismatch. Got=%d Expected=%d", countryData.stats.totalDeaths, 2464)
	}
	if countryData.stats.totalRecovered != 54990 {
		t.Errorf("country.totalRecovered mismatch. Got=%d Expected=%d", countryData.stats.totalRecovered, 54990)
	}
	tm, err = time.Parse(time.RFC3339, "2020-12-04T03:49:29Z")
	if err != nil {
		t.Error(err)
	}
	if countryData.stats.date != tm {
		t.Errorf("country.Date mismatch. Got=%v Expected=%v", countryData.stats.date, tm)
	}

}
