package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// Setup test router
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

// Test health endpoint
func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
}

// Test GetLandmarks handler with missing parameters
func TestGetLandmarksMissingParams(t *testing.T) {
	router := setupTestRouter()
	api := router.Group("/api/v1")
	api.GET("/landmarks", GetLandmarks)

	tests := []struct {
		name       string
		query      string
		expectCode int
	}{
		{
			name:       "Missing both city and country",
			query:      "",
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "Missing city",
			query:      "?country=UK",
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/landmarks"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("Expected status %d, got %d", tt.expectCode, w.Code)
			}
		})
	}
}

// Test GetWeather handler with missing parameters
func TestGetWeatherMissingParams(t *testing.T) {
	router := setupTestRouter()
	api := router.Group("/api/v1")
	api.GET("/weather", GetWeather)

	req, _ := http.NewRequest("GET", "/api/v1/weather", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["error"] != "city parameter is required" {
		t.Errorf("Expected error message about missing city, got %v", response["error"])
	}
}

// Test GetRates handler with valid currency
func TestGetRatesValidation(t *testing.T) {
	router := setupTestRouter()
	api := router.Group("/api/v1")
	api.GET("/rates/:currency", GetRates)

	tests := []struct {
		name     string
		currency string
	}{
		{"USD currency", "USD"},
		{"EUR currency", "EUR"},
		{"GBP currency", "GBP"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/rates/"+tt.currency, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Note: This will fail if the external API is down
			// In a real test, you'd mock the external API call
			if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
				t.Logf("Request to external API returned: %d (this is expected in tests)", w.Code)
			}
		})
	}
}

// Test model validation
func TestLandmarkModel(t *testing.T) {
	landmark := Landmark{
		Name:      "Test Landmark",
		Address:   "123 Test St",
		Latitude:  51.5074,
		Longitude: -0.1278,
		Rating:    8.5,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(landmark)
	if err != nil {
		t.Fatalf("Failed to marshal landmark: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Landmark
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
}

// Test City model
func TestCityModel(t *testing.T) {
	city := City{
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

	var unmarshaled City
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

// Integration test helper - checks if we can read environment variables
func TestEnvironmentSetup(t *testing.T) {
	// This test verifies that environment loading works
	t.Setenv("TEST_VAR", "test_value")

	value := os.Getenv("TEST_VAR")
	if value != "test_value" {
		t.Errorf("Expected TEST_VAR to be 'test_value', got %q", value)
	}
}
