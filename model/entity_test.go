package model

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLandmarkEntityMarshaling(t *testing.T) {
	landmark := LandmarkEntity{
		Name:      "Test Landmark",
		Address:   "123 Test St",
		Latitude:  51.5074,
		Longitude: -0.1278,
		Rating:    8.5,
		ImageURL:  "https://example.com/image.jpg",
		Category:  "Historical",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(landmark)
	if err != nil {
		t.Fatalf("Failed to marshal landmark: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled LandmarkEntity
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal landmark: %v", err)
	}

	if unmarshaled.Name != landmark.Name {
		t.Errorf("Expected name %q, got %q", landmark.Name, unmarshaled.Name)
	}
	if unmarshaled.Latitude != landmark.Latitude {
		t.Errorf("Expected latitude %f, got %f", landmark.Latitude, unmarshaled.Latitude)
	}
	if unmarshaled.Rating != landmark.Rating {
		t.Errorf("Expected rating %f, got %f", landmark.Rating, unmarshaled.Rating)
	}
}

func TestCityEntityMarshaling(t *testing.T) {
	city := CityEntity{
		ID:              1,
		Name:            "London",
		ThreeLetterCode: "GBR",
		Currency:        "GBP",
		Country:         "United Kingdom",
	}

	jsonData, err := json.Marshal(city)
	if err != nil {
		t.Fatalf("Failed to marshal city: %v", err)
	}

	var unmarshaled CityEntity
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal city: %v", err)
	}

	if unmarshaled.Name != city.Name {
		t.Errorf("Expected name %q, got %q", city.Name, unmarshaled.Name)
	}
	if unmarshaled.Currency != city.Currency {
		t.Errorf("Expected currency %q, got %q", city.Currency, unmarshaled.Currency)
	}
}

func TestWeatherEntityMarshaling(t *testing.T) {
	weather := WeatherEntity{
		Latitude:  51.5074,
		Longitude: -0.1278,
		Hourly: []HourlyForecastEntity{
			{
				Time:             "2024-01-01T00:00",
				Temperature:      10.5,
				WindSpeed:        5.3,
				RelativeHumidity: 75,
			},
		},
	}

	jsonData, err := json.Marshal(weather)
	if err != nil {
		t.Fatalf("Failed to marshal weather: %v", err)
	}

	var unmarshaled WeatherEntity
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal weather: %v", err)
	}

	if unmarshaled.Latitude != weather.Latitude {
		t.Errorf("Expected latitude %f, got %f", weather.Latitude, unmarshaled.Latitude)
	}
	if len(unmarshaled.Hourly) != 1 {
		t.Errorf("Expected 1 hourly entry, got %d", len(unmarshaled.Hourly))
	}
}

func TestEnvironmentSetup(t *testing.T) {
	// This test verifies that environment loading works
	t.Setenv("TEST_VAR", "test_value")

	value := os.Getenv("TEST_VAR")
	if value != "test_value" {
		t.Errorf("Expected TEST_VAR to be 'test_value', got %q", value)
	}
}
