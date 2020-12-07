package main

import (
	"log"
	"os"

	"github.com/jackc/pgx"

	durcov "github.com/TuhinNair/durcov"
)

func main() {
	var err error

	dbURL := os.Getenv("DATABASE_URL")
	covidEndpoint := os.Getenv("COVID_API_ENDPOINT")

	pgxpool, err := durcov.GetPgxPool(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pgxpool.Close()

	err = fetchAndStoreData(pgxpool, covidEndpoint)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchAndStoreData(pgxpool *pgx.ConnPool, covidEndpoint string) error {
	data, err := fetchData(covidEndpoint)
	if err != nil {
		return err
	}
	return storeData(pgxpool, data)
}

func fetchData(covidEndpoint string) (*durcov.Data, error) {
	dataSource := durcov.CovidAPI{}
	dataSource.UseURL(covidEndpoint)
	data, err := dataSource.FetchData()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func storeData(pgxpool *pgx.ConnPool, data *durcov.Data) error {
	dataStore := durcov.CovidDataStore{}
	dataStore.SetDBConnection(pgxpool)
	err := dataStore.StoreData(data)
	return err
}
