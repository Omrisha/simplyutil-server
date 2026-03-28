package service

import (
	"testing"
)

func TestNewRatesService(t *testing.T) {
	service := NewRatesService()
	if service == nil {
		t.Fatal("Expected rates service to be created")
	}

	if service.apiKey == "" {
		t.Error("Expected API key to be set (either from env or default)")
	}
}

func TestRatesServiceWithEnvKey(t *testing.T) {
	t.Setenv("EXCHANGE_RATE_API_KEY", "test-key-123")

	service := NewRatesService()
	if service.apiKey != "test-key-123" {
		t.Errorf("Expected API key 'test-key-123', got %q", service.apiKey)
	}
}

func TestRatesServiceGetRates(t *testing.T) {
	// This is an integration test that calls the real API
	// Skip if you don't want to make real API calls
	t.Skip("Skipping integration test - calls real Exchange Rate API")

	service := NewRatesService()
	rates, err := service.GetRates("USD")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rates.BaseCurrency != "USD" {
		t.Errorf("Expected base currency 'USD', got %q", rates.BaseCurrency)
	}

	if len(rates.Rates) == 0 {
		t.Error("Expected rates data")
	}
}
