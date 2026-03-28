package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() (*gin.Engine, *Handler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler()
	return router, handler
}

func TestNewHandler(t *testing.T) {
	handler := NewHandler()
	if handler == nil {
		t.Fatal("Expected handler to be created")
	}

	if handler.landmarkService == nil {
		t.Error("Expected landmark service to be initialized")
	}
	if handler.weatherService == nil {
		t.Error("Expected weather service to be initialized")
	}
	if handler.ratesService == nil {
		t.Error("Expected rates service to be initialized")
	}
	if handler.citiesService == nil {
		t.Error("Expected cities service to be initialized")
	}
}

func TestGetCitiesHandler(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/cities", handler.GetCities)

	// This is an integration test, skip if you don't want to call real APIs
	t.Skip("Skipping integration test - calls real REST Countries API")

	req, _ := http.NewRequest("GET", "/cities", nil)
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

	if _, exists := response["cities"]; !exists {
		t.Error("Expected 'cities' in response")
	}
}

func TestGetLandmarksHandlerMissingParam(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/landmarks", handler.GetLandmarks)

	req, _ := http.NewRequest("GET", "/landmarks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["error"] != "city parameter is required" {
		t.Errorf("Expected error message about missing city parameter, got %v", response["error"])
	}
}

func TestGetWeatherHandlerMissingParam(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/weather", handler.GetWeather)

	req, _ := http.NewRequest("GET", "/weather", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["error"] != "city parameter is required" {
		t.Errorf("Expected error message about missing city parameter, got %v", response["error"])
	}
}

func TestGetRatesHandlerMissingParam(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/rates/:currency", handler.GetRates)

	req, _ := http.NewRequest("GET", "/rates/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// When param is missing, Gin returns 404
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for missing param, got %d", w.Code)
	}
}

func TestGetCityDataHandler(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/cities/:name/:country", handler.GetCityData)

	// This is an integration test
	t.Skip("Skipping integration test - calls multiple real APIs")

	req, _ := http.NewRequest("GET", "/cities/London/UK", nil)
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

	if response["city"] != "London" {
		t.Errorf("Expected city 'London', got %v", response["city"])
	}

	if response["country"] != "UK" {
		t.Errorf("Expected country 'UK', got %v", response["country"])
	}
}
