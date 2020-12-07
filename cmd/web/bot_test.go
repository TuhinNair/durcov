package main

import "testing"

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
