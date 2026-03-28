package provider

import (
	"fmt"
	"net/url"
	"os"
	"simplyutil-server/model"
	"simplyutil-server/util"
)

// FoursquareProvider implements landmark fetching using Foursquare Places API
type FoursquareProvider struct {
	apiKey string
}

// NewFoursquareProvider creates a new Foursquare provider
func NewFoursquareProvider() *FoursquareProvider {
	return &FoursquareProvider{
		apiKey: os.Getenv("FOURSQUARE_API_KEY"),
	}
}

// FetchLandmarks fetches landmarks using Foursquare Places API
func (f *FoursquareProvider) FetchLandmarks(cityName, country string) ([]model.LandmarkEntity, error) {
	if f.apiKey == "" {
		return nil, fmt.Errorf("FOURSQUARE_API_KEY not set")
	}

	// Geocode the city to get coordinates
	lat, lon, err := util.GeocodeCity(cityName, country)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	// Build Foursquare Places API request
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

	var fsResponse model.FoursquareV3Response
	if err := util.HTTPGetJSONWithHeaders(fullURL, headers, &fsResponse); err != nil {
		return nil, fmt.Errorf("foursquare API error: %w", err)
	}

	// Convert to LandmarkEntity format
	landmarks := make([]model.LandmarkEntity, 0)
	for _, place := range fsResponse.Results {
		// Build image URL if photos exist
		var imageURL string
		if len(place.Photos) > 0 {
			photo := place.Photos[0]
			imageURL = fmt.Sprintf("%s300x300%s", photo.Prefix, photo.Suffix)
		}

		// Get primary category if available
		var category string
		if len(place.Categories) > 0 {
			category = place.Categories[0].Name
		}

		landmark := model.LandmarkEntity{
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
