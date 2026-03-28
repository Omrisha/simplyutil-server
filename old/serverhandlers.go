package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// getLandmarkProvider returns the appropriate landmark provider based on configuration
func getLandmarkProvider() LandmarkProvider {
	// Check if we should use Wikipedia (if LANDMARK_PROVIDER env var is set to "wikipedia")
	// or if Foursquare API key is not available
	provider := os.Getenv("LANDMARK_PROVIDER")
	foursquareKey := os.Getenv("FOURSQUARE_API_KEY")

	if provider == "wikipedia" || foursquareKey == "" {
		return NewWikipediaLandmarkProvider()
	}

	// Default to Foursquare, but fall back to Wikipedia on error
	return NewFoursquareLandmarkProvider()
}

// GetCities returns a list of all cities/countries with currencies
func GetCities(c *gin.Context) {
	// Fetch from REST Countries API
	cities, err := fetchCitiesFromRestCountries()
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
func GetCityData(c *gin.Context) {
	cityName := c.Param("name")
	country := c.Param("country")

	// Get the configured landmark provider
	landmarkProvider := getLandmarkProvider()

	// Fetch all data concurrently
	landmarksCh := make(chan landmarksResult)
	weatherCh := make(chan weatherResult)
	ratesCh := make(chan ratesResult)

	go func() {
		landmarks, err := landmarkProvider.FetchLandmarks(cityName, country)
		landmarksCh <- landmarksResult{landmarks: landmarks, err: err}
	}()

	go func() {
		weather, err := fetchWeatherFromOpenMeteo(cityName)
		weatherCh <- weatherResult{weather: weather, err: err}
	}()

	go func() {
		// For now, just return USD rates - could get city's currency from DB
		rates, err := fetchRatesFromExchangeAPI("USD")
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
func GetLandmarks(c *gin.Context) {
	cityName := c.Query("city")
	country := c.Query("country")

	if cityName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

	// Get the configured landmark provider
	landmarkProvider := getLandmarkProvider()

	landmarks, err := landmarkProvider.FetchLandmarks(cityName, country)
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
func GetWeather(c *gin.Context) {
	cityName := c.Query("city")

	if cityName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

	weather, err := fetchWeatherFromOpenMeteo(cityName)
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
func GetRates(c *gin.Context) {
	currency := c.Param("currency")

	if currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currency parameter is required"})
		return
	}

	rates, err := fetchRatesFromExchangeAPI(currency)
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
	landmarks []LandmarkEntity
	err       error
}

type weatherResult struct {
	weather WeatherEntity
	err     error
}

type ratesResult struct {
	rates RatesEntity
	err   error
}
