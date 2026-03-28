package provider

import (
	"fmt"
	"net/url"
	"simplyutil-server/model"
	"simplyutil-server/util"
)

// WikipediaProvider implements landmark fetching using Wikipedia/Wikimedia APIs
type WikipediaProvider struct{}

// NewWikipediaProvider creates a new Wikipedia provider
func NewWikipediaProvider() *WikipediaProvider {
	return &WikipediaProvider{}
}

// FetchLandmarks fetches landmarks using Wikipedia and Wikimedia Commons
func (w *WikipediaProvider) FetchLandmarks(cityName, country string) ([]model.LandmarkEntity, error) {
	// Geocode the city to get coordinates
	lat, lon, err := util.GeocodeCity(cityName, country)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	// Use Wikipedia API to search for nearby landmarks
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

	var geoResponse model.WikipediaGeoSearchResponse
	if err := util.HTTPGetJSONWithHeaders(fullURL, headers, &geoResponse); err != nil {
		return nil, fmt.Errorf("wikipedia geosearch error: %w", err)
	}

	if len(geoResponse.Query.GeoSearch) == 0 {
		return []model.LandmarkEntity{}, nil
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

	var pageResponse model.WikipediaPageResponse
	if err := util.HTTPGetJSONWithHeaders(fullURL, headers, &pageResponse); err != nil {
		return nil, fmt.Errorf("wikipedia page details error: %w", err)
	}

	// Build landmarks from geo search and page details
	landmarks := make([]model.LandmarkEntity, 0)
	for _, geoResult := range geoResponse.Query.GeoSearch {
		pageID := fmt.Sprintf("%d", geoResult.PageID)
		pageDetail, exists := pageResponse.Query.Pages[pageID]

		var imageURL string
		if exists && pageDetail.Thumbnail.Source != "" {
			imageURL = pageDetail.Thumbnail.Source
		}

		landmark := model.LandmarkEntity{
			Name:      geoResult.Title,
			Address:   fmt.Sprintf("%s, %s", cityName, country),
			Latitude:  geoResult.Lat,
			Longitude: geoResult.Lon,
			Rating:    0, // Wikipedia doesn't provide ratings
			ImageURL:  imageURL,
			Category:  "Landmark",
		}
		landmarks = append(landmarks, landmark)
	}

	return landmarks, nil
}
