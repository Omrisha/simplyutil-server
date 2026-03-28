package service

import (
	"testing"
)

func TestNewCitiesService(t *testing.T) {
	service := NewCitiesService()
	if service == nil {
		t.Fatal("Expected cities service to be created")
	}
}

func TestCitiesServiceGetCities(t *testing.T) {
	// This is an integration test that calls the real API
	// Skip if you don't want to make real API calls
	t.Skip("Skipping integration test - calls real REST Countries API")

	service := NewCitiesService()
	cities, err := service.GetCities()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(cities) == 0 {
		t.Error("Expected cities data")
	}

	// Check first city has required fields
	if len(cities) > 0 {
		city := cities[0]
		if city.Name == "" {
			t.Error("Expected city to have a name")
		}
		if city.Country == "" {
			t.Error("Expected city to have a country")
		}
		if city.Currency == "" {
			t.Error("Expected city to have a currency")
		}
	}
}
