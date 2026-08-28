package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"simplyutil-server/model"
	"simplyutil-server/provider"
	"simplyutil-server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler holds all service dependencies
type Handler struct {
	landmarkService *service.LandmarkService
	weatherService  *service.WeatherService
	ratesService    *service.RatesService
	citiesService   *service.CitiesService
}

// NewHandler creates a new handler with all services
func NewHandler() *Handler {
	return &Handler{
		landmarkService: service.NewLandmarkService(),
		weatherService:  service.NewWeatherService(),
		ratesService:    service.NewRatesService(),
		citiesService:   service.NewCitiesService(),
	}
}

// maxRequestBodyBytes caps how much of a QUERY body we are willing to read.
const maxRequestBodyBytes = 1 << 20 // 1 MiB

// bindOptionalJSON decodes the request body into target. An absent or empty
// body is not an error - it means "no fields supplied", and callers apply their
// own required-field checks. It writes the error response and returns false
// when the body is present but malformed or oversized.
func bindOptionalJSON(c *gin.Context, target any) bool {
	if c.Request.Body == nil {
		return true
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestBodyBytes)

	// Decode rather than inspect ContentLength: a chunked body reports a length
	// of -1, so an empty one is only recognisable by the EOF it decodes to.
	err := json.NewDecoder(c.Request.Body).Decode(target)
	switch {
	case err == nil, errors.Is(err, io.EOF):
		return true
	case errors.As(err, new(*http.MaxBytesError)):
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error":   "request body too large",
			"message": fmt.Sprintf("body exceeds the %d byte limit", maxRequestBodyBytes),
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid JSON body",
			"message": err.Error(),
		})
	}
	return false
}

// Pagination bounds for the GeoNames-backed cities endpoints.
const (
	defaultPage     = 1
	defaultPageSize = 100
	maxPageSize     = 1000
)

// citiesQueryParams are the filter/pagination options accepted by both the
// GET query string and the QUERY JSON body. Page and PageSize are pointers so
// an explicitly supplied zero is distinguishable from an omitted field.
type citiesQueryParams struct {
	Page     *int   `json:"page"`
	PageSize *int   `json:"pageSize"`
	Search   string `json:"search"`
	Country  string `json:"country"`
}

// filtered reports whether the caller asked for any filtering or pagination.
func (p citiesQueryParams) filtered() bool {
	return p.Page != nil || p.PageSize != nil || p.Search != "" || p.Country != ""
}

// pagination returns the page and page size to request, falling back to the
// defaults for anything omitted or outside the supported bounds.
func (p citiesQueryParams) pagination() (page, pageSize int) {
	page, pageSize = defaultPage, defaultPageSize
	if p.Page != nil && *p.Page >= 1 {
		page = *p.Page
	}
	if p.PageSize != nil && *p.PageSize >= 1 && *p.PageSize <= maxPageSize {
		pageSize = *p.PageSize
	}
	return page, pageSize
}

// GetCities returns cities. When using GeoNames, supports ?page, ?pageSize, ?search, ?country.
func (h *Handler) GetCities(c *gin.Context) {
	params := citiesQueryParams{
		Search:  c.Query("search"),
		Country: c.Query("country"),
	}
	if raw, ok := c.GetQuery("page"); ok {
		page, _ := strconv.Atoi(raw)
		params.Page = &page
	}
	if raw, ok := c.GetQuery("pageSize"); ok {
		pageSize, _ := strconv.Atoi(raw)
		params.PageSize = &pageSize
	}

	h.respondCities(c, params)
}

// QueryCities is the QUERY-method variant of GetCities; filters arrive as a JSON body.
func (h *Handler) QueryCities(c *gin.Context) {
	var params citiesQueryParams
	if !bindOptionalJSON(c, &params) {
		return
	}
	h.respondCities(c, params)
}

func (h *Handler) respondCities(c *gin.Context, params citiesQueryParams) {
	if params.filtered() {
		// Only GeoNames can filter; say so rather than quietly serving the
		// unfiltered list as though the criteria had been applied.
		if !h.citiesService.IsGeoNames() {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error":   "search and pagination are not supported by the configured city provider",
				"message": "set GEONAMES_USERNAME to enable filtering",
			})
			return
		}

		page, pageSize := params.pagination()
		cities, err := h.citiesService.GetCitiesWithQuery(provider.GeoNamesQuery{
			Page:        page,
			PageSize:    pageSize,
			Search:      params.Search,
			CountryCode: params.Country,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch cities",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"cities":   cities,
			"count":    len(cities),
			"page":     page,
			"pageSize": pageSize,
		})
		return
	}

	cities, err := h.citiesService.GetCities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch cities",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"cities": cities,
		"count":  len(cities),
	})
}

// GetCityData returns all data for a specific city in one request
func (h *Handler) GetCityData(c *gin.Context) {
	cityName := c.Param("name")
	country := c.Param("country")

	// Fetch all data concurrently
	landmarksCh := make(chan landmarksResult)
	weatherCh := make(chan weatherResult)
	ratesCh := make(chan ratesResult)

	go func() {
		landmarks, err := h.landmarkService.GetLandmarks(cityName, country)
		landmarksCh <- landmarksResult{landmarks: landmarks, err: err}
	}()

	go func() {
		weather, err := h.weatherService.GetWeather(cityName)
		weatherCh <- weatherResult{weather: weather, err: err}
	}()

	go func() {
		// For now, just return USD rates
		rates, err := h.ratesService.GetRates("USD")
		ratesCh <- ratesResult{rates: rates, err: err}
	}()

	// Collect results
	landmarksRes := <-landmarksCh
	weatherRes := <-weatherCh
	ratesRes := <-ratesCh

	// Build response (include partial data even if some calls fail)
	response := gin.H{
		"city":    cityName,
		"country": country,
	}

	if landmarksRes.err == nil {
		response["landmarks"] = landmarksRes.landmarks
	} else {
		response["landmarks_error"] = landmarksRes.err.Error()
	}

	if weatherRes.err == nil {
		response["weather"] = weatherRes.weather
	} else {
		response["weather_error"] = weatherRes.err.Error()
	}

	if ratesRes.err == nil {
		response["rates"] = ratesRes.rates
	} else {
		response["rates_error"] = ratesRes.err.Error()
	}

	c.JSON(http.StatusOK, response)
}

// cityQueryParams is the QUERY JSON body for the landmarks and weather endpoints.
type cityQueryParams struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

// bindCityQuery parses the QUERY body and enforces the required city field.
// It writes the error response and returns false when the request is invalid.
func bindCityQuery(c *gin.Context, params *cityQueryParams) bool {
	if !bindOptionalJSON(c, params) {
		return false
	}
	if params.City == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return false
	}
	return true
}

// GetLandmarks returns landmarks for a city
func (h *Handler) GetLandmarks(c *gin.Context) {
	cityName := c.Query("city")
	country := c.Query("country")

	if cityName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

	h.respondLandmarks(c, cityName, country)
}

// QueryLandmarks is the QUERY-method variant of GetLandmarks; the city arrives as a JSON body.
func (h *Handler) QueryLandmarks(c *gin.Context) {
	var params cityQueryParams
	if !bindCityQuery(c, &params) {
		return
	}
	h.respondLandmarks(c, params.City, params.Country)
}

func (h *Handler) respondLandmarks(c *gin.Context, cityName, country string) {
	landmarks, err := h.landmarkService.GetLandmarks(cityName, country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch landmarks",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"landmarks": landmarks,
		"count":     len(landmarks),
	})
}

// GetWeather returns weather forecast for a city
func (h *Handler) GetWeather(c *gin.Context) {
	cityName := c.Query("city")

	if cityName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

	h.respondWeather(c, cityName)
}

// QueryWeather is the QUERY-method variant of GetWeather; the city arrives as a JSON body.
func (h *Handler) QueryWeather(c *gin.Context) {
	var params cityQueryParams
	if !bindCityQuery(c, &params) {
		return
	}
	h.respondWeather(c, params.City)
}

func (h *Handler) respondWeather(c *gin.Context, cityName string) {
	weather, err := h.weatherService.GetWeather(cityName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch weather",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"weather": weather,
	})
}

// GetRates returns exchange rates for a currency
func (h *Handler) GetRates(c *gin.Context) {
	currency := c.Param("currency")

	if currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currency parameter is required"})
		return
	}

	rates, err := h.ratesService.GetRates(currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch rates",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rates)
}

// Result types for concurrent fetching
type landmarksResult struct {
	landmarks []model.LandmarkEntity
	err       error
}

type weatherResult struct {
	weather model.WeatherEntity
	err     error
}

type ratesResult struct {
	rates model.RatesEntity
	err   error
}
