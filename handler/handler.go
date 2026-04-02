package handler

import (
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

// GetCities returns cities. When using GeoNames, supports ?page, ?pageSize, ?search, ?country.
func (h *Handler) GetCities(c *gin.Context) {
	_, hasPage := c.GetQuery("page")
	_, hasPageSize := c.GetQuery("pageSize")
	search := c.Query("search")
	country := c.Query("country")

	if h.citiesService.IsGeoNames() && (hasPage || hasPageSize || search != "" || country != "") {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 1000 {
			pageSize = 100
		}
		Í
		cities, err := h.citiesService.GetCitiesWithQuery(provider.GeoNamesQuery{
			Page:        page,
			PageSize:    pageSize,
			Search:      search,
			CountryCode: country,
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

// GetLandmarks returns landmarks for a city
func (h *Handler) GetLandmarks(c *gin.Context) {
	cityName := c.Query("city")
	country := c.Query("country")

	if cityName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

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
