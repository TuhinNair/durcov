package durcov

import (
	"os"
	"testing"
	"time"
)

func exampleTestData() (*Data, error) {
	exampleTime, err := time.Parse(time.RFC3339, "2020-12-04T03:49:29Z")
	if err != nil {
		return nil, err
	}

	exampleData := Data{
		&global{
			&statistics{
				totalConfirmed: 10000000,
				totalDeaths:    500000,
				totalRecovered: 500000,
				date:           exampleTime,
			},
		},
		[]*country{
			&country{
				name: "Afghanistan",
				code: "AF",
				slug: "afghanistan",
				stats: &statistics{
					totalConfirmed: 46980,
					totalDeaths:    1822,
					totalRecovered: 37026,
					date:           exampleTime,
				},
			},
			&country{
				name: "Singapore",
				code: "SG",
				slug: "singapore",
				stats: &statistics{
					totalConfirmed: 46980,
					totalDeaths:    1822,
					totalRecovered: 37026,
					date:           exampleTime,
				},
			},
		},
	}

	return &exampleData, nil
}

func TestDataStore(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	pool, err := GetPgxPool(dbURL)
	if err != nil {
		t.Error(err)
	}

	dataStore := CovidDataStore{}
	dataStore.SetDBConnection(pool)

	exampleData, err := exampleTestData()
	if err != nil {
		t.Error(err)
	}
	dataStore.StoreData(exampleData)

	tests := map[string]func(t *testing.T){
		"Global data in database matches input": func(t *testing.T) {
			var confirmed int64
			var deaths int64
			var recovered int64
			var time time.Time

			err := pool.QueryRow("SELECT confirmed, deaths, recovered, collected_at FROM covid_stats WHERE id='GLOBAL';").Scan(&confirmed, &deaths, &recovered, &time)
			if err != nil {
				t.Error(err)
			}
			if confirmed != exampleData.global.stats.totalConfirmed {
				t.Errorf("global confirmed mismatch. Expected=%d Got=%d", exampleData.global.stats.totalConfirmed, confirmed)
			}
			if deaths != exampleData.global.stats.totalDeaths {
				t.Errorf("global deaths mismatch. Expected=%d Got=%d", exampleData.global.stats.totalDeaths, deaths)
			}
			if recovered != exampleData.global.stats.totalRecovered {
				t.Errorf("global recovered mismatch. Expected=%d Got=%d", exampleData.global.stats.totalRecovered, recovered)
			}
			if time != exampleData.global.stats.date {
				t.Errorf("global collected at time mismatch. Expected=%v Got=%v", exampleData.global.stats.date, time)
			}
		},
		"Countries data in database matches input": func(t *testing.T) {
			for _, country := range exampleData.countries {
				var name string
				var slug string
				var id string
				var confirmed int64
				var deaths int64
				var recovered int64
				var time time.Time

				countryCode := country.code

				err := pool.QueryRow("SELECT id, slug, name, confirmed, deaths, recovered, collected_at FROM covid_stats WHERE id=$1;", countryCode).Scan(&id, &slug, &name, &confirmed, &deaths, &recovered, &time)
				if err != nil {
					t.Error(err)
				}
				if id != countryCode {
					t.Errorf("Country code mismatch. Expected=%s Got=%s", countryCode, id)
				}
				if name != country.name {
					t.Errorf("Country name mismatch. Expected=%s Got=%s", country.name, name)
				}
				if slug != country.slug {
					t.Errorf("Country slug mismatch. Expected=%s Got=%s", country.slug, slug)
				}
				if confirmed != country.stats.totalConfirmed {
					t.Errorf("Country confirmed mismatch. Expected=%d Got=%d", country.stats.totalConfirmed, confirmed)
				}
				if deaths != country.stats.totalDeaths {
					t.Errorf("Country deaths mismatch. Expected=%d Got=%d", country.stats.totalDeaths, deaths)
				}
				if recovered != country.stats.totalRecovered {
					t.Errorf("Country recovered mismatch. Expected=%d Got=%d", country.stats.totalRecovered, recovered)
				}
				if time != country.stats.date {
					t.Errorf("Country collected at time mismatch. Expected=%v Got=%v", country.stats.date, time)
				}
			}
		},
	}

	for name, test := range tests {
		t.Run(name, test)
	}
}
