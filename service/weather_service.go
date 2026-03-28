package service

import (
	"fmt"
	"net/url"
	"simplyutil-server/model"
	"simplyutil-server/util"
)

// WeatherService handles weather-related operations
type WeatherService struct{}

// NewWeatherService creates a new weather service
func NewWeatherService() *WeatherService {
	return &WeatherService{}
}

// GetWeather fetches weather data from Open-Meteo API
func (s *WeatherService) GetWeather(cityName string) (model.WeatherEntity, error) {
	// Geocode the city
	lat, lon, err := util.GeocodeCity(cityName, "")
	if err != nil {
		return model.WeatherEntity{}, fmt.Errorf("geocoding failed: %w", err)
	}

	// Build Open-Meteo API request
	baseURL := "https://api.open-meteo.com/v1/forecast"
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%f", lat))
	params.Add("longitude", fmt.Sprintf("%f", lon))
	params.Add("hourly", "temperature_2m,relative_humidity_2m,wind_speed_10m")
	params.Add("forecast_days", "1")

	fullURL := baseURL + "?" + params.Encode()

	var weatherResponse model.OpenMeteoResponse
	if err := util.HTTPGetJSON(fullURL, &weatherResponse); err != nil {
		return model.WeatherEntity{}, fmt.Errorf("open-meteo API error: %w", err)
	}

	// Convert to WeatherEntity format
	hourly := make([]model.HourlyForecastEntity, 0)
	for i := range weatherResponse.Hourly.Time {
		hourly = append(hourly, model.HourlyForecastEntity{
			Time:             weatherResponse.Hourly.Time[i],
			Temperature:      weatherResponse.Hourly.Temperature[i],
			WindSpeed:        weatherResponse.Hourly.WindSpeed[i],
			RelativeHumidity: weatherResponse.Hourly.RelativeHumidity[i],
		})
	}

	return model.WeatherEntity{
		Latitude:  lat,
		Longitude: lon,
		Hourly:    hourly,
	}, nil
}
