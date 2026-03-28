package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"simplyutil-server/model"
	"time"
)

// GeocodeCity converts city name to coordinates using Nominatim (OpenStreetMap)
func GeocodeCity(cityName, country string) (float64, float64, error) {
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

	var results []model.NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, 0, err
	}

	if len(results) == 0 {
		return 0, 0, fmt.Errorf("location not found: %s", query)
	}

	return results[0].Lat, results[0].Lon, nil
}
