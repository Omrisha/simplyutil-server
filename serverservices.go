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

// fetchLandmarksFromFoursquare fetches landmarks using Foursquare API v2
func fetchLandmarksFromFoursquare(cityName, country string) ([]Landmark, error) {
	// First, geocode the city to get coordinates
	lat, lon, err := geocodeCity(cityName, country)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	// Load API key from environment (after .env is loaded)
	apiKey := os.Getenv("FOURSQUARE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("FOURSQUARE_API_KEY not set")
	}

	// Parse API key (format: CLIENT_ID+CLIENT_SECRET)
	clientID, clientSecret := parseV2APIKey(apiKey)

	// Build Foursquare v2 API request
	baseURL := "https://api.foursquare.com/v2/venues/explore"
	params := url.Values{}
	params.Add("ll", fmt.Sprintf("%f,%f", lat, lon))
	params.Add("client_id", clientID)
	params.Add("client_secret", clientSecret)
	params.Add("v", "20240101")
	params.Add("radius", "5000")
	params.Add("section", "sights")
	params.Add("limit", "20")

	fullURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("foursquare API error: %d - %s", resp.StatusCode, string(body))
	}

	var fsResponse FoursquareV2Response
	if err := json.NewDecoder(resp.Body).Decode(&fsResponse); err != nil {
		return nil, err
	}

	// Convert to our Landmark format
	landmarks := make([]Landmark, 0)
	if len(fsResponse.Response.Groups) > 0 {
		for _, item := range fsResponse.Response.Groups[0].Items {
			venue := item.Venue
			landmark := Landmark{
				Name:      venue.Name,
				Address:   venue.Location.Address,
				Latitude:  venue.Location.Lat,
				Longitude: venue.Location.Lng,
				Rating:    venue.Rating,
			}
			landmarks = append(landmarks, landmark)
		}
	}

	return landmarks, nil
}

// fetchWeatherFromOpenMeteo fetches weather data from Open-Meteo API
func fetchWeatherFromOpenMeteo(cityName string) (WeatherData, error) {
	// Geocode the city
	lat, lon, err := geocodeCity(cityName, "")
	if err != nil {
		return WeatherData{}, fmt.Errorf("geocoding failed: %w", err)
	}

	// Build Open-Meteo API request
	baseURL := "https://api.open-meteo.com/v1/forecast"
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%f", lat))
	params.Add("longitude", fmt.Sprintf("%f", lon))
	params.Add("hourly", "temperature_2m,relative_humidity_2m,wind_speed_10m")
	params.Add("forecast_days", "1")

	fullURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return WeatherData{}, fmt.Errorf("open-meteo API error: %d - %s", resp.StatusCode, string(body))
	}

	var weatherResponse OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return WeatherData{}, err
	}

	// Convert to our format
	hourly := make([]HourlyForecast, 0)
	for i := range weatherResponse.Hourly.Time {
		hourly = append(hourly, HourlyForecast{
			Time:              weatherResponse.Hourly.Time[i],
			Temperature:       weatherResponse.Hourly.Temperature[i],
			WindSpeed:         weatherResponse.Hourly.WindSpeed[i],
			RelativeHumidity:  weatherResponse.Hourly.RelativeHumidity[i],
		})
	}

	return WeatherData{
		Latitude:  lat,
		Longitude: lon,
		Hourly:    hourly,
	}, nil
}

// fetchRatesFromExchangeAPI fetches exchange rates
func fetchRatesFromExchangeAPI(baseCurrency string) (RatesData, error) {
	apiKey := "c04b66e4d1f1f147c60834b3" // Consider moving to env var
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", apiKey, baseCurrency)

	resp, err := http.Get(url)
	if err != nil {
		return RatesData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return RatesData{}, fmt.Errorf("exchange-rate API error: %d - %s", resp.StatusCode, string(body))
	}

	var ratesResponse ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&ratesResponse); err != nil {
		return RatesData{}, err
	}

	return RatesData{
		BaseCurrency: ratesResponse.BaseCode,
		Rates:        ratesResponse.ConversionRates,
		Timestamp:    time.Unix(ratesResponse.TimeLastUpdateUnix, 0),
	}, nil
}

// fetchCitiesFromRestCountries fetches list of countries/cities
func fetchCitiesFromRestCountries() ([]City, error) {
	url := "https://restcountries.com/v3.1/all?fields=name,cca3,capital,currencies"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("rest-countries API error: %d - %s", resp.StatusCode, string(body))
	}

	var countries []RestCountry
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return nil, err
	}

	// Convert to our City format
	cities := make([]City, 0)
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

		city := City{
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

	var results []NominatimResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, 0, err
	}

	if len(results) == 0 {
		return 0, 0, fmt.Errorf("location not found: %s", query)
	}

	return results[0].Lat, results[0].Lon, nil
}

// parseV2APIKey parses Foursquare v2 API key format (CLIENT_ID+CLIENT_SECRET)
func parseV2APIKey(apiKey string) (string, string) {
	// If contains '+', split it
	for i, ch := range apiKey {
		if ch == '+' {
			return apiKey[:i], apiKey[i+1:]
		}
	}
	// If no '+', return the whole key as both (fallback)
	return apiKey, apiKey
}

