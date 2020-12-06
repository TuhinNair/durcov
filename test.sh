#! /bin/bash

export TEST_DATABASE_URL=$4 
echo "Using DB URL: $TEST_DATABASE_URL"
PGPASSWORD=$1 psql -U $2 -d $3 -f ./schema/covid_stats_schema_07_12_2020.sql
go test -v ./...