package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// httpGetJSON performs an HTTP GET request and decodes the JSON response
func httpGetJSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return err
	}

	return nil
}

// httpGetJSONWithHeaders performs an HTTP GET request with custom headers and decodes the JSON response
func httpGetJSONWithHeaders(url string, headers map[string]string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return err
	}

	return nil
}

// LandmarkProvider interface for different landmark data sources
type LandmarkProvider interface {
	FetchLandmarks(cityName, country string) ([]LandmarkEntity, error)
}

// FoursquareLandmarkProvider implements LandmarkProvider using Foursquare Places API
type FoursquareLandmarkProvider struct {
	apiKey string
}

// NewFoursquareLandmarkProvider creates a new Foursquare provider
func NewFoursquareLandmarkProvider() *FoursquareLandmarkProvider {
	return &FoursquareLandmarkProvider{
		apiKey: os.Getenv("FOURSQUARE_API_KEY"),
	}
}

// FetchLandmarks fetches landmarks using Foursquare Places API
func (f *FoursquareLandmarkProvider) FetchLandmarks(cityName, country string) ([]LandmarkEntity, error) {
	if f.apiKey == "" {
		return nil, fmt.Errorf("FOURSQUARE_API_KEY not set")
	}

	// First, geocode the city to get coordinates
	lat, lon, err := geocodeCity(cityName, country)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	// Build Foursquare Places API request (new endpoint structure)
	baseURL := "https://places-api.foursquare.com/places/search"
	params := url.Values{}
	params.Add("ll", fmt.Sprintf("%f,%f", lat, lon))
	params.Add("radius", "5000")
	params.Add("categories", "16000") // Landmarks & Outdoors category
	params.Add("limit", "20")
	params.Add("fields", "fsq_place_id,name,latitude,longitude,location,categories,photos,rating")

	fullURL := baseURL + "?" + params.Encode()

	// Set up headers with Bearer token and version
	headers := map[string]string{
		"Authorization":        "Bearer " + f.apiKey,
		"Accept":               "application/json",
		"X-Places-Api-Version": "2025-06-17",
	}

	var fsResponse FoursquareV3Response
	if err := httpGetJSONWithHeaders(fullURL, headers, &fsResponse); err != nil {
		return nil, fmt.Errorf("foursquare API error: %w", err)
	}

	// Convert to our LandmarkEntity format
	landmarks := make([]LandmarkEntity, 0)
	for _, place := range fsResponse.Results {
		// Build image URL if photos exist
		var imageURL string
		if len(place.Photos) > 0 {
			photo := place.Photos[0]
			// Use original size or specify dimensions (e.g., 300x300)
			imageURL = fmt.Sprintf("%s300x300%s", photo.Prefix, photo.Suffix)
		}

		// Get primary category if available
		var category string
		if len(place.Categories) > 0 {
			category = place.Categories[0].Name
		}

		landmark := LandmarkEntity{
			Name:      place.Name,
			Address:   place.Location.FormattedAddress,
			Latitude:  place.Latitude,
			Longitude: place.Longitude,
			Rating:    place.Rating,
			ImageURL:  imageURL,
			Category:  category,
		}
		landmarks = append(landmarks, landmark)
	}

	return landmarks, nil
}

// WikipediaLandmarkProvider implements LandmarkProvider using Wikipedia/Wikimedia APIs
type WikipediaLandmarkProvider struct{}

// NewWikipediaLandmarkProvider creates a new Wikipedia provider
func NewWikipediaLandmarkProvider() *WikipediaLandmarkProvider {
	return &WikipediaLandmarkProvider{}
}

// FetchLandmarks fetches landmarks using Wikipedia and Wikimedia Commons
func (w *WikipediaLandmarkProvider) FetchLandmarks(cityName, country string) ([]LandmarkEntity, error) {
	// First, geocode the city to get coordinates
	lat, lon, err := geocodeCity(cityName, country)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	// Use Wikipedia API to search for nearby landmarks
	// Using geosearch to find articles near the coordinates
	baseURL := "https://en.wikipedia.org/w/api.php"
	params := url.Values{}
	params.Add("action", "query")
	params.Add("list", "geosearch")
	params.Add("gscoord", fmt.Sprintf("%f|%f", lat, lon))
	params.Add("gsradius", "5000") // 5km radius
	params.Add("gslimit", "20")
	params.Add("format", "json")

	fullURL := baseURL + "?" + params.Encode()

	// Wikipedia requires a User-Agent header
	headers := map[string]string{
		"User-Agent": "SimplyUtil/1.0 (Travel App; omrishapira@example.com)",
	}

	var geoResponse WikipediaGeoSearchResponse
	if err := httpGetJSONWithHeaders(fullURL, headers, &geoResponse); err != nil {
		return nil, fmt.Errorf("wikipedia geosearch error: %w", err)
	}

	if len(geoResponse.Query.GeoSearch) == 0 {
		return []LandmarkEntity{}, nil
	}

	// Get page IDs to fetch images and details
	pageIDs := make([]string, 0)
	for _, result := range geoResponse.Query.GeoSearch {
		pageIDs = append(pageIDs, fmt.Sprintf("%d", result.PageID))
	}

	// Fetch page details including images
	params = url.Values{}
	params.Add("action", "query")
	params.Add("pageids", fmt.Sprintf("%s", pageIDs[0])) // Start with first page
	for i := 1; i < len(pageIDs) && i < 20; i++ {
		params.Set("pageids", params.Get("pageids")+"|"+pageIDs[i])
	}
	params.Add("prop", "pageimages|coordinates|extracts")
	params.Add("piprop", "thumbnail")
	params.Add("pithumbsize", "300")
	params.Add("exintro", "true")
	params.Add("explaintext", "true")
	params.Add("exsentences", "1")
	params.Add("format", "json")

	fullURL = baseURL + "?" + params.Encode()

	var pageResponse WikipediaPageResponse
	if err := httpGetJSONWithHeaders(fullURL, headers, &pageResponse); err != nil {
		return nil, fmt.Errorf("wikipedia page details error: %w", err)
	}

	// Build landmarks from geo search and page details
	landmarks := make([]LandmarkEntity, 0)
	for _, geoResult := range geoResponse.Query.GeoSearch {
		pageID := fmt.Sprintf("%d", geoResult.PageID)
		pageDetail, exists := pageResponse.Query.Pages[pageID]

		var imageURL string
		if exists && pageDetail.Thumbnail.Source != "" {
			imageURL = pageDetail.Thumbnail.Source
		}

		landmark := LandmarkEntity{
			Name:      geoResult.Title,
			Address:   fmt.Sprintf("%s, %s", cityName, country),
			Latitude:  geoResult.Lat,
			Longitude: geoResult.Lon,
			Rating:    0, // Wikipedia doesn't provide ratings
			ImageURL:  imageURL,
			Category:  "Landmark", // Could be enhanced with Wikipedia categories
		}
		landmarks = append(landmarks, landmark)
	}

	return landmarks, nil
}

// fetchLandmarksFromFoursquare is kept for backward compatibility
func fetchLandmarksFromFoursquare(cityName, country string) ([]LandmarkEntity, error) {
	provider := NewFoursquareLandmarkProvider()
	return provider.FetchLandmarks(cityName, country)
}

// fetchWeatherFromOpenMeteo fetches weather data from Open-Meteo API
func fetchWeatherFromOpenMeteo(cityName string) (WeatherEntity, error) {
	// Geocode the city
	lat, lon, err := geocodeCity(cityName, "")
	if err != nil {
		return WeatherEntity{}, fmt.Errorf("geocoding failed: %w", err)
	}

	// Build Open-Meteo API request
	baseURL := "https://api.open-meteo.com/v1/forecast"
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%f", lat))
	params.Add("longitude", fmt.Sprintf("%f", lon))
	params.Add("hourly", "temperature_2m,relative_humidity_2m,wind_speed_10m")
	params.Add("forecast_days", "1")

	fullURL := baseURL + "?" + params.Encode()

	var weatherResponse OpenMeteoResponse
	if err := httpGetJSON(fullURL, &weatherResponse); err != nil {
		return WeatherEntity{}, fmt.Errorf("open-meteo API error: %w", err)
	}

	// Convert to our format
	hourly := make([]HourlyForecastEntity, 0)
	for i := range weatherResponse.Hourly.Time {
		hourly = append(hourly, HourlyForecastEntity{
			Time:             weatherResponse.Hourly.Time[i],
			Temperature:      weatherResponse.Hourly.Temperature[i],
			WindSpeed:        weatherResponse.Hourly.WindSpeed[i],
			RelativeHumidity: weatherResponse.Hourly.RelativeHumidity[i],
		})
	}

	return WeatherEntity{
		Latitude:  lat,
		Longitude: lon,
		Hourly:    hourly,
	}, nil
}

// fetchRatesFromExchangeAPI fetches exchange rates
func fetchRatesFromExchangeAPI(baseCurrency string) (RatesEntity, error) {
	apiKey := "c04b66e4d1f1f147c60834b3" // Consider moving to env var
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", apiKey, baseCurrency)

	var ratesResponse ExchangeRateResponse
	if err := httpGetJSON(url, &ratesResponse); err != nil {
		return RatesEntity{}, fmt.Errorf("exchange-rate API error: %w", err)
	}

	return RatesEntity{
		BaseCurrency: ratesResponse.BaseCode,
		Rates:        ratesResponse.ConversionRates,
		Timestamp:    time.Unix(ratesResponse.TimeLastUpdateUnix, 0),
	}, nil
}

// fetchCitiesFromRestCountries fetches list of countries/cities
func fetchCitiesFromRestCountries() ([]CityEntity, error) {
	url := "https://restcountries.com/v3.1/all?fields=name,cca3,capital,currencies"

	var countries []RestCountryResponse
	if err := httpGetJSON(url, &countries); err != nil {
		return nil, fmt.Errorf("rest-countries API error: %w", err)
	}

	// Convert to our CityEntity format
	cities := make([]CityEntity, 0)
	id := 1
	for _, country := range countries {
		if len(country.Capital) == 0 || country.Currencies == nil {
			continue
		}

		// Get first currency
		var currencyCode string
		for code := range country.Currencies {
			currencyCode = code
			break
		}

		city := CityEntity{
			ID:              id,
			Name:            country.Capital[0],
			ThreeLetterCode: country.CCA3,
			Currency:        currencyCode,
			Country:         country.Name.Common,
		}
		cities = append(cities, city)
		id++
	}

	return cities, nil
}

// geocodeCity converts city name to coordinates using a simple geocoding service
func geocodeCity(cityName, country string) (float64, float64, error) {
	// Using Nominatim (OpenStreetMap) for free geocoding
	query := cityName
	if country != "" {
		query = fmt.Sprintf("%s, %s", cityName, country)
	}

	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", query)
	params.Add("format", "json")
	params.Add("limit", "1")

	fullURL := baseURL + "?" + params.Encode()

	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("User-Agent", "SimplyUtil-iOS-App")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var results []NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, 0, err
	}

	if len(results) == 0 {
		return 0, 0, fmt.Errorf("location not found: %s", query)
	}

	return results[0].Lat, results[0].Lon, nil
}
