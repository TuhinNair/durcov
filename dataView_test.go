package durcov

import (
	"os"
	"testing"
)

func TestDataView(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	pool, err := GetPgxPool(dbURL)
	if err != nil {
		t.Error(err)
	}
	defer pool.Close()

	dataStore := CovidDataStore{}
	dataStore.SetDBConnection(pool)

	exampleData, err := ExampleTestData()
	if err != nil {
		t.Fatal(err)
	}
	err = dataStore.StoreData(exampleData)
	if err != nil {
		t.Fatal(err)
	}

	dataView := CovidBotView{}
	dataView.SetDBConnection(pool)

	tests := map[string]func(t *testing.T){
		"Total global deaths view": func(t *testing.T) {
			globalDeaths, err := dataView.LatestGlobalView(Deaths)
			if err != nil {
				t.Error(err)
			}
			if globalDeaths != exampleData.global.stats.totalDeaths {
				t.Errorf("Total global deaths mismatch. Expected=%d Got=%d", exampleData.global.stats.totalDeaths, globalDeaths)
			}
		},
		"Total global active view": func(t *testing.T) {
			globalActive, err := dataView.LatestGlobalView(Active)
			if err != nil {
				t.Error(err)
			}
			expectdActive := calculateActive(
				exampleData.global.stats.totalConfirmed,
				exampleData.global.stats.totalDeaths,
				exampleData.global.stats.totalRecovered,
			)
			if globalActive != expectdActive {
				t.Errorf("Total global active mismatch. Expected=%d Got=%d", expectdActive, globalActive)
			}
		},
		"Countries deaths view": func(t *testing.T) {
			for _, country := range exampleData.countries {
				countryDeaths, err := dataView.LatestCountryView(country.code, Deaths)
				if err != nil {
					t.Error(err)
				}
				if countryDeaths != country.stats.totalDeaths {
					t.Errorf("country deaths mismatch for country=%s. Expected=%d Got=%d", country.name, country.stats.totalDeaths, countryDeaths)
				}
			}
		},
		"Countries active view": func(t *testing.T) {
			for _, country := range exampleData.countries {
				countryActive, err := dataView.LatestCountryView(country.code, Active)
				if err != nil {
					t.Error(err)
				}
				expectedActive := calculateActive(
					country.stats.totalConfirmed,
					country.stats.totalDeaths,
					country.stats.totalRecovered,
				)
				if countryActive != expectedActive {
					t.Errorf("country active mismatch for country=%s. Expected=%d Got=%d", country.name, expectedActive, countryActive)
				}
			}
		},
		"Unmathced country code returns specific error": func(t *testing.T) {
			_, err := dataView.LatestCountryView("--", Active)
			if err, ok := err.(*NoCountryMatchedError); !ok {
				t.Errorf("Unexpected error. Expected=*NoCountryMatchedError Got=%T", err)
			}
			_, err = dataView.LatestCountryView("--", Deaths)
			if err, ok := err.(*NoCountryMatchedError); !ok {
				t.Errorf("Unexpected error. Expected=*NoCountryMatchedError Got=%T", err)
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
	}
}
