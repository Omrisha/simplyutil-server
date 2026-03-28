package model

import (
	"encoding/json"
	"testing"
)

func TestFoursquareV3ResponseParsing(t *testing.T) {
	jsonData := `{
		"results": [{
			"fsq_place_id": "test-id-123",
			"name": "Big Ben",
			"latitude": 51.5007,
			"longitude": -0.1246,
			"location": {
				"address": "Westminster",
				"formatted_address": "Westminster, London"
			},
			"categories": [{
				"name": "Historic Site"
			}],
			"photos": [{
				"prefix": "https://example.com/photo_",
				"suffix": ".jpg",
				"width": 1920,
				"height": 1080
			}],
			"rating": 9.5
		}]
	}`

	var response FoursquareV3Response
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(response.Results) == 0 {
		t.Fatal("Expected at least one result")
	}

	place := response.Results[0]
	if place.Name != "Big Ben" {
		t.Errorf("Expected name 'Big Ben', got %q", place.Name)
	}
	if place.Rating != 9.5 {
		t.Errorf("Expected rating 9.5, got %f", place.Rating)
	}
	if place.Latitude != 51.5007 {
		t.Errorf("Expected latitude 51.5007, got %f", place.Latitude)
	}
}

func TestOpenMeteoResponseParsing(t *testing.T) {
	jsonData := `{
		"latitude": 51.5,
		"longitude": -0.12,
		"hourly": {
			"time": ["2024-01-01T00:00", "2024-01-01T01:00"],
			"temperature_2m": [10.5, 11.2],
			"wind_speed_10m": [5.3, 6.1],
			"relative_humidity_2m": [75, 73]
		}
	}`

	var response OpenMeteoResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if response.Latitude != 51.5 {
		t.Errorf("Expected latitude 51.5, got %f", response.Latitude)
	}
	if len(response.Hourly.Time) != 2 {
		t.Errorf("Expected 2 time entries, got %d", len(response.Hourly.Time))
	}
	if response.Hourly.Temperature[0] != 10.5 {
		t.Errorf("Expected first temperature 10.5, got %f", response.Hourly.Temperature[0])
	}
}

func TestExchangeRateResponseParsing(t *testing.T) {
	jsonData := `{
		"result": "success",
		"base_code": "USD",
		"conversion_rates": {
			"EUR": 0.85,
			"GBP": 0.73,
			"JPY": 110.0
		},
		"time_last_update_unix": 1609459200
	}`

	var response ExchangeRateResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if response.BaseCode != "USD" {
		t.Errorf("Expected base code 'USD', got %q", response.BaseCode)
	}
	if response.ConversionRates["EUR"] != 0.85 {
		t.Errorf("Expected EUR rate 0.85, got %f", response.ConversionRates["EUR"])
	}
	if response.TimeLastUpdateUnix != 1609459200 {
		t.Errorf("Expected timestamp 1609459200, got %d", response.TimeLastUpdateUnix)
	}
}

func TestRestCountryResponseParsing(t *testing.T) {
	jsonData := `{
		"name": {
			"common": "United Kingdom",
			"official": "United Kingdom of Great Britain and Northern Ireland"
		},
		"cca3": "GBR",
		"capital": ["London"],
		"currencies": {
			"GBP": {
				"name": "British pound",
				"symbol": "£"
			}
		}
	}`

	var country RestCountryResponse
	err := json.Unmarshal([]byte(jsonData), &country)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if country.Name.Common != "United Kingdom" {
		t.Errorf("Expected 'United Kingdom', got %q", country.Name.Common)
	}
	if country.CCA3 != "GBR" {
		t.Errorf("Expected CCA3 'GBR', got %q", country.CCA3)
	}
	if len(country.Capital) == 0 || country.Capital[0] != "London" {
		t.Errorf("Expected capital 'London', got %v", country.Capital)
	}
}

func TestNominatimResponseParsing(t *testing.T) {
	jsonData := `{
		"lat": "51.5074",
		"lon": "-0.1278"
	}`

	var result NominatimResponse
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result.Lat != 51.5074 {
		t.Errorf("Expected lat 51.5074, got %f", result.Lat)
	}
	if result.Lon != -0.1278 {
		t.Errorf("Expected lon -0.1278, got %f", result.Lon)
	}
}
