package main

import "time"

// City represents a city/country with currency info
type City struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	ThreeLetterCode string `json:"threeLetterCode"`
	Currency        string `json:"currency"`
	Country         string `json:"country"`
}

// Landmark represents a place of interest
type Landmark struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Rating    float64 `json:"rating"`
}

// WeatherData represents weather forecast
type WeatherData struct {
	Latitude  float64           `json:"latitude"`
	Longitude float64           `json:"longitude"`
	Hourly    []HourlyForecast  `json:"hourly"`
}

type HourlyForecast struct {
	Time             string  `json:"time"`
	Temperature      float64 `json:"temperature"`
	WindSpeed        float64 `json:"windSpeed"`
	RelativeHumidity int     `json:"relativeHumidity"`
}

// RatesData represents exchange rates
type RatesData struct {
	BaseCurrency string             `json:"baseCurrency"`
	Rates        map[string]float64 `json:"rates"`
	Timestamp    time.Time          `json:"timestamp"`
}

// External API response types

// FoursquareV3Response for Foursquare Places API v3
type FoursquareV3Response struct {
	Results []struct {
		Name     string `json:"name"`
		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Address   string  `json:"address"`
		} `json:"location"`
		Rating float64 `json:"rating"`
	} `json:"results"`
}

// FoursquareV2Response for Foursquare API v2 (legacy)
type FoursquareV2Response struct {
	Response struct {
		Groups []struct {
			Items []struct {
				Venue struct {
					Name     string  `json:"name"`
					Rating   float64 `json:"rating"`
					Location struct {
						Lat               float64  `json:"lat"`
						Lng               float64  `json:"lng"`
						Address           string   `json:"address"`
						FormattedAddress  []string `json:"formattedAddress"`
					} `json:"location"`
				} `json:"venue"`
			} `json:"items"`
		} `json:"groups"`
	} `json:"response"`
}

// OpenMeteoResponse for Open-Meteo API
type OpenMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Hourly    struct {
		Time             []string  `json:"time"`
		Temperature      []float64 `json:"temperature_2m"`
		WindSpeed        []float64 `json:"wind_speed_10m"`
		RelativeHumidity []int     `json:"relative_humidity_2m"`
	} `json:"hourly"`
}

// ExchangeRateResponse for Exchange Rate API
type ExchangeRateResponse struct {
	Result              string             `json:"result"`
	BaseCode            string             `json:"base_code"`
	ConversionRates     map[string]float64 `json:"conversion_rates"`
	TimeLastUpdateUnix  int64              `json:"time_last_update_unix"`
}

// RestCountry for REST Countries API
type RestCountry struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	CCA3       string   `json:"cca3"`
	Capital    []string `json:"capital"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

// NominatimResult for geocoding
type NominatimResult struct {
	Lat float64 `json:"lat,string"`
	Lon float64 `json:"lon,string"`
}
