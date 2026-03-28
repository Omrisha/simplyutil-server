package model

import "time"

// CityEntity represents a city/country with currency info
type CityEntity struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	ThreeLetterCode string `json:"threeLetterCode"`
	Currency        string `json:"currency"`
	Country         string `json:"country"`
}

// LandmarkEntity represents a place of interest
type LandmarkEntity struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Rating    float64 `json:"rating"`
	ImageURL  string  `json:"imageUrl,omitempty"`
	Category  string  `json:"category,omitempty"`
}

// WeatherEntity represents weather forecast
type WeatherEntity struct {
	Latitude  float64                 `json:"latitude"`
	Longitude float64                 `json:"longitude"`
	Hourly    []HourlyForecastEntity  `json:"hourly"`
}

// HourlyForecastEntity represents hourly weather data
type HourlyForecastEntity struct {
	Time             string  `json:"time"`
	Temperature      float64 `json:"temperature"`
	WindSpeed        float64 `json:"windSpeed"`
	RelativeHumidity int     `json:"relativeHumidity"`
}

// RatesEntity represents exchange rates
type RatesEntity struct {
	BaseCurrency string             `json:"baseCurrency"`
	Rates        map[string]float64 `json:"rates"`
	Timestamp    time.Time          `json:"timestamp"`
}
