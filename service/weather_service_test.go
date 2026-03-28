package service

import (
	"testing"
)

func TestNewWeatherService(t *testing.T) {
	service := NewWeatherService()
	if service == nil {
		t.Fatal("Expected weather service to be created")
	}
}

func TestWeatherServiceGetWeather(t *testing.T) {
	// This is an integration test that calls the real API
	// Skip if you don't want to make real API calls
	t.Skip("Skipping integration test - calls real Open-Meteo API")

	service := NewWeatherService()
	weather, err := service.GetWeather("London")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if weather.Latitude == 0 || weather.Longitude == 0 {
		t.Error("Expected non-zero coordinates")
	}

	if len(weather.Hourly) == 0 {
		t.Error("Expected hourly weather data")
	}
}
