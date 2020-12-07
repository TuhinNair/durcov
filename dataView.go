package durcov

import (
	"fmt"
	"log"

	"github.com/jackc/pgx"
	"gopkg.in/errgo.v2/fmt/errors"
)

// Datum represents a datapoint
type Datum int

// Represt covid datapoints
const (
	Confirmed Datum = iota + 1
	Deaths
	Recovered
	Active
)

// DataView describes functions for obtaining view data
type DataView interface {
	SetDBConnection(pgxpool *pgx.ConnPool)
	LatestGlobalView(datapoint Datum) (int64, error)
	LatestCountryView(countryCode string, datapoint Datum) (int64, error)
}

// NoCountryMatchedError when data for a given country code is not found in the database
type NoCountryMatchedError struct {
	attemptedCode string
}

func (n *NoCountryMatchedError) Error() string {
	errMsg := fmt.Sprintf("No country matched with code %s", n.attemptedCode)
	return errMsg
}

// CovidBotView represents a view for covid data
type CovidBotView struct {
	pgxpool *pgx.ConnPool
}

// SetDBConnection sets a database connectioon to be used to generate view data.
// Must be set before trying to generate views.
func (c *CovidBotView) SetDBConnection(pool *pgx.ConnPool) {
	c.pgxpool = pool
}

// LatestGlobalView returns the latest (available) global data for the given datapoint
// Returns error if the given datapoint does not have a view implemented.
func (c *CovidBotView) LatestGlobalView(datapoint Datum) (int64, error) {
	if c.pgxpool == nil {
		return 0, errors.New("DB Connection not set in data view")
	}
	switch datapoint {
	case Active:
		log.Println("About to view latestGlobalActive")
		return c.latestGlobalActive()
	case Deaths:
		log.Println("About to view latestGlobalDeaths")
		return c.latestGlobalDeaths()
	default:
		return 0, errors.Newf("Unsupported Op for Global View. Datum Enum %d", datapoint)
	}
}

func (c *CovidBotView) latestGlobalActive() (int64, error) {
	var confirmed int64
	var deaths int64
	var recovered int64

	err := c.pgxpool.QueryRow("SELECT confirmed, deaths, recovered FROM covid_stats WHERE id='GLOBAL';").Scan(&confirmed, &deaths, &recovered)
	if err != nil {
		return 0, err
	}
	globalActive := calculateActive(confirmed, deaths, recovered)
	return globalActive, nil
}

func (c *CovidBotView) latestGlobalDeaths() (int64, error) {
	var deaths int64
	err := c.pgxpool.QueryRow("SELECT deaths FROM covid_stats WHERE id='GLOBAL';").Scan(&deaths)
	if err != nil {
		return 0, err
	}
	return deaths, nil
}

// LatestCountryView returns the latest (available) covid data for the given country code and datapoint.
// Returns err if no match found for the country code (or) if no view implemented for the datapoint.
func (c *CovidBotView) LatestCountryView(countryCode string, datapoint Datum) (int64, error) {
	if c.pgxpool == nil {
		return 0, errors.New("DB Connection not set in data view")
	}
	switch datapoint {
	case Active:
		return c.latestCountryActive(countryCode)
	case Deaths:
		return c.latestCountryDeaths(countryCode)
	default:
		return 0, errors.Newf("Unsupported Op for Country View. Datum Enum %d", datapoint)
	}
}

func (c *CovidBotView) latestCountryActive(countryCode string) (int64, error) {
	var confirmed int64
	var deaths int64
	var recovered int64

	err := c.pgxpool.QueryRow("SELECT confirmed, deaths, recovered FROM covid_stats WHERE id=$1;", countryCode).Scan(&confirmed, &deaths, &recovered)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, &NoCountryMatchedError{countryCode}
		}
		return 0, err
	}
	countryActive := calculateActive(confirmed, deaths, recovered)
	return countryActive, nil
}

func (c *CovidBotView) latestCountryDeaths(countryCode string) (int64, error) {
	var deaths int64
	err := c.pgxpool.QueryRow("SELECT deaths FROM covid_stats WHERE id=$1;", countryCode).Scan(&deaths)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, &NoCountryMatchedError{countryCode}
		}
		return 0, err
	}
	return deaths, nil
}

func calculateActive(confirmed int64, deaths int64, recovered int64) int64 {
	return confirmed - (deaths + recovered)
}
