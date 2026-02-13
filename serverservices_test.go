package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test parseV2APIKey function
func TestParseV2APIKey(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedID     string
		expectedSecret string
	}{
		{
			name:           "Valid key with separator",
			input:          "CLIENT_ID+CLIENT_SECRET",
			expectedID:     "CLIENT_ID",
			expectedSecret: "CLIENT_SECRET",
		},
		{
			name:           "Key without separator",
			input:          "SINGLE_KEY",
			expectedID:     "SINGLE_KEY",
			expectedSecret: "SINGLE_KEY",
		},
		{
			name:           "Key with multiple plus signs",
			input:          "ID+SECRET+EXTRA",
			expectedID:     "ID",
			expectedSecret: "SECRET+EXTRA",
		},
		{
			name:           "Empty key",
			input:          "",
			expectedID:     "",
			expectedSecret: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, secret := parseV2APIKey(tt.input)
			if id != tt.expectedID {
				t.Errorf("Expected ID %q, got %q", tt.expectedID, id)
			}
			if secret != tt.expectedSecret {
				t.Errorf("Expected secret %q, got %q", tt.expectedSecret, secret)
			}
		})
	}
}

// Test geocodeCity with mock server
func TestGeocodeCity(t *testing.T) {
	// Create mock server
	mockResponse := []NominatimResult{
		{Lat: 51.5074, Lon: -0.1278},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") == "" {
			t.Error("Expected query parameter 'q' to be set")
		}
		if r.URL.Query().Get("format") != "json" {
			t.Error("Expected format=json")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Note: This test would need to modify geocodeCity to accept a custom URL
	// For now, this is a template for how to test with mocking
	t.Skip("Skipping: requires refactoring geocodeCity to accept custom URL for testing")
}

// Test FoursquareV2Response JSON parsing
func TestFoursquareV2ResponseParsing(t *testing.T) {
	jsonData := `{
		"response": {
			"groups": [{
				"items": [{
					"venue": {
						"name": "Big Ben",
						"rating": 9.5,
						"location": {
							"lat": 51.5007,
							"lng": -0.1246,
							"address": "Westminster",
							"formattedAddress": ["Westminster", "London"]
						}
					}
				}]
			}]
		}
	}`

	var response FoursquareV2Response
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(response.Response.Groups) == 0 {
		t.Fatal("Expected at least one group")
	}

	if len(response.Response.Groups[0].Items) == 0 {
		t.Fatal("Expected at least one item")
	}

	venue := response.Response.Groups[0].Items[0].Venue
	if venue.Name != "Big Ben" {
		t.Errorf("Expected name 'Big Ben', got %q", venue.Name)
	}
	if venue.Rating != 9.5 {
		t.Errorf("Expected rating 9.5, got %f", venue.Rating)
	}
	if venue.Location.Lat != 51.5007 {
		t.Errorf("Expected lat 51.5007, got %f", venue.Location.Lat)
	}
}

// Test OpenMeteoResponse JSON parsing
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

// Test ExchangeRateResponse JSON parsing
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

// Test RestCountry JSON parsing
func TestRestCountryParsing(t *testing.T) {
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

	var country RestCountry
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

// Test NominatimResult JSON parsing
func TestNominatimResultParsing(t *testing.T) {
	jsonData := `{
		"lat": "51.5074",
		"lon": "-0.1278"
	}`

	var result NominatimResult
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
