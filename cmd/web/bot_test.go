package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/TuhinNair/durcov"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	exitVal := m.Run()
	os.Exit(exitVal)
}
func TestBotMatchRequest(t *testing.T) {
	tests := []struct {
		input                 string
		expectedParsedRequest *parsedRequest
		expectError           bool
	}{
		{
			"CASES TOTAL",
			&parsedRequest{
				_Cases,
				"TOTAL",
			},
			false,
		},
		{
			"CASES in",
			&parsedRequest{
				_Cases,
				"IN",
			},
			false,
		},
		{
			"deaths TOTAL",
			&parsedRequest{
				_Deaths,
				"TOTAL",
			},
			false,
		},
		{
			"DEATHS uS",
			&parsedRequest{
				_Deaths,
				"US",
			},
			false,
		},
		{
			"cases       AL",
			&parsedRequest{
				_Cases,
				"AL",
			},
			false,
		},
		{
			"DEATHS    TOTAL",
			&parsedRequest{
				_Deaths,
				"TOTAL",
			},
			false,
		},
		{
			"asiodjoai",
			nil,
			true,
		},
		{
			"sfsf sadoija asfdijas",
			nil,
			true,
		},
		{
			"Cases total cases total",
			nil,
			true,
		},
		{
			"Case IN",
			nil,
			true,
		},
		{
			"DEATHS TOT",
			nil,
			true,
		},
		{
			"Cases I",
			nil,
			true,
		},
	}

	testBot := Bot{nil}

	for _, test := range tests {
		parsedReq, botErr := testBot.matchRequest(test.input)
		if botErr != nil {
			if !test.expectError {
				t.Fatalf("Didn't Expect error. Input: %s, Error: %v", test.input, botErr.err)
			}
		} else if test.expectError {
			t.Fatalf("Expected error. Input: %s, Parsed Request: %+v", test.input, parsedReq)
		} else {
			if parsedReq.Code != test.expectedParsedRequest.Code {
				t.Fatalf("Parsed Request Code Error. Expected=%s Got=%s", test.expectedParsedRequest.Code, parsedReq.Code)
			}

			if parsedReq.Type != test.expectedParsedRequest.Type {
				t.Fatalf("Parsed Request Type error. Expected=%d Got=%d", test.expectedParsedRequest.Type, parsedReq.Type)
			}
		}
	}
}

func TestBotResponseGeneration(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	pool, err := durcov.GetPgxPool(dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	dataStore := &durcov.CovidDataStore{}
	dataStore.SetDBConnection(pool)

	exampleData, err := durcov.ExampleTestData()
	if err != nil {
		t.Fatal(err)
	}
	err = dataStore.StoreData(exampleData)
	if err != nil {
		t.Fatal(err)
	}

	dataView := &durcov.CovidBotView{}
	dataView.SetDBConnection(pool)

	testBot := &Bot{dataView}

	tests := []struct {
		input    string
		expected string
	}{
		{
			"CASES TOTAL",
			"Total Active Cases: 9,000,000",
		},
		{
			"CASES AF",
			"[AF] Afghanistan Active Cases: 8,132",
		},
		{
			"DEATHS TOTAL",
			"Total Deaths: 500,000",
		},
		{
			"DEATHS SG",
			"[SG] Singapore Deaths: 1,822",
		},
		{
			"abcdefghijklmnopqrstuvwxyz",
			"Sorry, I'm not sure how to respond to that.",
		},
		{
			"DEATH SG",
			"Sorry, I'm not sure how to respond to that.",
		},
		{
			"CASES                                                TOTAL",
			"Sorry, that message is too long for me.",
		},
		{
			"cases IN",
			"Sorry, that code doesn't match any countries I know.",
		},
	}

	for _, test := range tests {
		response := testBot.respond(test.input)
		if response != test.expected {
			t.Errorf("Response mismatch. Expected=%s Got=%s", test.expected, response)
		}
	}
}
