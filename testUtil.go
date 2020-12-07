package durcov

import "time"

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
