package provider

import (
	"fmt"
	"net/url"
	"simplyutil-server/model"
	"simplyutil-server/util"
)

// GeoNamesProvider implements city fetching using the GeoNames API
type GeoNamesProvider struct {
	username string
}

// GeoNamesQuery holds optional search criteria supported by the GeoNames API
type GeoNamesQuery struct {
	Page        int
	PageSize    int
	Search      string // q: free-text name search
	CountryCode string // country: 2-letter ISO code
}

// NewGeoNamesProvider creates a new GeoNames provider.
func NewGeoNamesProvider(username string) *GeoNamesProvider {
	return &GeoNamesProvider{username: username}
}

// FetchCities returns the first page of cities with population >= 50000
func (p *GeoNamesProvider) FetchCities() ([]model.GeoNamesPlace, error) {
	return p.FetchCitiesWithQuery(GeoNamesQuery{Page: 1, PageSize: 100})
}

// FetchCitiesPaged makes a single GeoNames API call for the requested page
func (p *GeoNamesProvider) FetchCitiesPaged(page, pageSize int) ([]model.GeoNamesPlace, error) {
	return p.FetchCitiesWithQuery(GeoNamesQuery{Page: page, PageSize: pageSize})
}

// FetchCitiesWithQuery makes a single GeoNames API call with optional search criteria
func (p *GeoNamesProvider) FetchCitiesWithQuery(q GeoNamesQuery) ([]model.GeoNamesPlace, error) {
	if p.username == "" {
		return nil, fmt.Errorf("GEONAMES_USERNAME not set")
	}

	params := url.Values{}
	params.Set("featureClass", "P")
	params.Set("minPopulation", "50000")
	params.Set("maxRows", fmt.Sprintf("%d", q.PageSize))
	params.Set("startRow", fmt.Sprintf("%d", (q.Page-1)*q.PageSize))
	params.Set("username", p.username)

	if q.Search != "" {
		params.Set("q", q.Search)
	}
	if q.CountryCode != "" {
		params.Set("country", q.CountryCode)
	}

	apiURL := "http://api.geonames.org/searchJSON?" + params.Encode()

	var response model.GeoNamesResponse
	if err := util.HTTPGetJSON(apiURL, &response); err != nil {
		return nil, fmt.Errorf("geonames API error: %w", err)
	}

	return response.GeoNames, nil
}
