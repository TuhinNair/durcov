package durcov

import (
	"errors"

	"github.com/jackc/pgx"
)

// DataStore describes an API to store expected data
type DataStore interface {
	SetDBConnection(pgxpool *pgx.ConnPool)
	StoreData(data *Data) error
}

// CovidDataStore represents a connection API to store expected Covid data
type CovidDataStore struct {
	pgxpool *pgx.ConnPool
}

// SetDBConnection sets the connection to the backing database.
// Must be set before calling StoreData
func (c *CovidDataStore) SetDBConnection(pgxpool *pgx.ConnPool) {
	c.pgxpool = pgxpool
}

// StoreData stores given data in the database,
// Note: StoreData overwrites existing data in the database.
func (c *CovidDataStore) StoreData(data *Data) error {
	if c.pgxpool == nil {
		return errors.New("Database connection not set on data store")
	}
	tx, err := c.pgxpool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("TRUNCATE covid_stats")
	if err != nil {
		return err
	}

	err = storeBatchCovidData(tx, data)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func storeBatchCovidData(tx *pgx.Tx, data *Data) error {
	source := [][]interface{}{}

	globalData := []interface{}{
		"GLOBAL",
		"Global",
		"global",
		data.global.stats.totalConfirmed,
		data.global.stats.totalDeaths,
		data.global.stats.totalRecovered,
		data.global.stats.date,
	}
	source = append(source, globalData)

	for _, country := range data.countries {
		countryData := []interface{}{
			country.code,
			country.name,
			country.slug,
			country.stats.totalConfirmed,
			country.stats.totalDeaths,
			country.stats.totalRecovered,
			country.stats.date,
		}

		source = append(source, countryData)
	}

	tableName := pgx.Identifier{"covid_stats"}
	columns := []string{
		"id",
		"name",
		"slug",
		"confirmed",
		"deaths",
		"recovered",
		"collected_at",
	}

	_, err := tx.CopyFrom(tableName, columns, pgx.CopyFromRows(source))

	if err != nil {
		return err
	}
	return nil
}
