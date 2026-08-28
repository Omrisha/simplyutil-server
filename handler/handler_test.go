package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestQueryLandmarksHandlerMissingBody(t *testing.T) {
	router, handler := setupTestRouter()
	router.Handle("QUERY", "/landmarks", handler.QueryLandmarks)

	req, _ := http.NewRequest("QUERY", "/landmarks", nil)
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

func TestQueryLandmarksHandlerInvalidJSON(t *testing.T) {
	router, handler := setupTestRouter()
	router.Handle("QUERY", "/landmarks", handler.QueryLandmarks)

	req, _ := http.NewRequest("QUERY", "/landmarks", strings.NewReader("{not json"))
	req.Header.Set("Content-Type", "application/json")
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

	if response["error"] != "invalid JSON body" {
		t.Errorf("Expected invalid JSON body error, got %v", response["error"])
	}
}

func TestQueryWeatherHandlerMissingBody(t *testing.T) {
	router, handler := setupTestRouter()
	router.Handle("QUERY", "/weather", handler.QueryWeather)

	req, _ := http.NewRequest("QUERY", "/weather", nil)
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

func TestQueryCitiesHandlerInvalidJSON(t *testing.T) {
	router, handler := setupTestRouter()
	router.Handle("QUERY", "/cities", handler.QueryCities)

	req, _ := http.NewRequest("QUERY", "/cities", strings.NewReader("{not json"))
	req.Header.Set("Content-Type", "application/json")
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

	if response["error"] != "invalid JSON body" {
		t.Errorf("Expected invalid JSON body error, got %v", response["error"])
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

// queryRequest builds a QUERY request. A negative length simulates a chunked
// body, which is how a real server reports an unknown Content-Length.
func queryRequest(t *testing.T, path, body string, contentLength int64) *http.Request {
	t.Helper()
	req, err := http.NewRequest("QUERY", path, io.NopCloser(strings.NewReader(body)))
	if err != nil {
		t.Fatalf("Failed to build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = contentLength
	return req
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	return response
}

// An empty chunked body carries no fields, so it must produce the same
// "city required" error as an empty body with a known length of zero.
func TestQueryLandmarksHandlerEmptyChunkedBody(t *testing.T) {
	router, handler := setupTestRouter()
	router.Handle("QUERY", "/landmarks", handler.QueryLandmarks)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, queryRequest(t, "/landmarks", "", -1))

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
	if got := decodeBody(t, w)["error"]; got != "city parameter is required" {
		t.Errorf("Expected error message about missing city parameter, got %v", got)
	}
}

func TestQueryLandmarksHandlerOversizedBody(t *testing.T) {
	router, handler := setupTestRouter()
	router.Handle("QUERY", "/landmarks", handler.QueryLandmarks)

	body := `{"city": "` + strings.Repeat("a", maxRequestBodyBytes+1) + `"}`
	w := httptest.NewRecorder()
	router.ServeHTTP(w, queryRequest(t, "/landmarks", body, int64(len(body))))

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected status 413, got %d", w.Code)
	}
	if got := decodeBody(t, w)["error"]; got != "request body too large" {
		t.Errorf("Expected body-too-large error, got %v", got)
	}
}

// An explicit zero is a supplied filter, not an omitted one, and filtering is
// unavailable without GeoNames - so it must be reported, not silently dropped.
func TestCitiesFilteringUnsupportedWithoutGeoNames(t *testing.T) {
	t.Setenv("GEONAMES_USERNAME", "")

	tests := []struct {
		name string
		req  func(t *testing.T) *http.Request
	}{
		{
			name: "GET with explicit zero page",
			req: func(t *testing.T) *http.Request {
				req, _ := http.NewRequest("GET", "/cities?page=0", nil)
				return req
			},
		},
		{
			name: "QUERY with explicit zero page",
			req: func(t *testing.T) *http.Request {
				return queryRequest(t, "/cities", `{"page": 0}`, -1)
			},
		},
		{
			name: "QUERY with search over a chunked body",
			req: func(t *testing.T) *http.Request {
				return queryRequest(t, "/cities", `{"search": "lon"}`, -1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler := setupTestRouter()
			router.GET("/cities", handler.GetCities)
			router.Handle("QUERY", "/cities", handler.QueryCities)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, tt.req(t))

			if w.Code != http.StatusNotImplemented {
				t.Errorf("Expected status 501, got %d", w.Code)
			}
			if got := decodeBody(t, w)["error"]; got == nil {
				t.Error("Expected an error explaining that filtering is unsupported")
			}
		})
	}
}

func TestCitiesQueryParamsPagination(t *testing.T) {
	intPtr := func(v int) *int { return &v }

	tests := []struct {
		name           string
		params         citiesQueryParams
		filtered       bool
		page, pageSize int
	}{
		{"omitted", citiesQueryParams{}, false, defaultPage, defaultPageSize},
		{"explicit zero page", citiesQueryParams{Page: intPtr(0)}, true, defaultPage, defaultPageSize},
		{"explicit zero page size", citiesQueryParams{PageSize: intPtr(0)}, true, defaultPage, defaultPageSize},
		{"page size over the maximum", citiesQueryParams{PageSize: intPtr(maxPageSize + 1)}, true, defaultPage, defaultPageSize},
		{"search only", citiesQueryParams{Search: "lon"}, true, defaultPage, defaultPageSize},
		{"valid values", citiesQueryParams{Page: intPtr(3), PageSize: intPtr(25)}, true, 3, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.filtered(); got != tt.filtered {
				t.Errorf("filtered() = %v, want %v", got, tt.filtered)
			}
			page, pageSize := tt.params.pagination()
			if page != tt.page || pageSize != tt.pageSize {
				t.Errorf("pagination() = (%d, %d), want (%d, %d)", page, pageSize, tt.page, tt.pageSize)
			}
		})
	}
}
