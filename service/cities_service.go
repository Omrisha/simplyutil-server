package service

import (
	"fmt"
	"os"
	"simplyutil-server/interfaces"
	"simplyutil-server/model"
	"simplyutil-server/provider"
)

// CitiesService handles cities/countries data operations
type CitiesService struct {
	provider         interfaces.CityProvider
	geoNamesProvider *provider.GeoNamesProvider
}

// NewCitiesService creates a new cities service with the appropriate provider
func NewCitiesService() *CitiesService {
	geoNamesUsername := os.Getenv("GEONAMES_USERNAME")
	if geoNamesUsername != "" {
		return &CitiesService{geoNamesProvider: provider.NewGeoNamesProvider(geoNamesUsername)}
	}
	return &CitiesService{provider: provider.NewRestCountriesProvider()}
}

// GetCities fetches cities from the configured provider
func (s *CitiesService) GetCities() ([]model.CityEntity, error) {
	if s.geoNamesProvider != nil {
		return s.enrichGeoNames(s.geoNamesProvider.FetchCities())
	}
	return s.provider.FetchCities()
}

// GetCitiesWithQuery fetches cities using GeoNames search criteria — only supported with GeoNames
func (s *CitiesService) GetCitiesWithQuery(q provider.GeoNamesQuery) ([]model.CityEntity, error) {
	if s.geoNamesProvider == nil {
		return nil, fmt.Errorf("search/pagination is only supported with the GeoNames provider")
	}
	return s.enrichGeoNames(s.geoNamesProvider.FetchCitiesWithQuery(q))
}

func (s *CitiesService) enrichGeoNames(places []model.GeoNamesPlace, err error) ([]model.CityEntity, error) {
	if err != nil {
		return nil, err
	}

	countryLookup, err := provider.NewRestCountriesProvider().FetchCountryLookup()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch country lookup: %w", err)
	}

	cities := make([]model.CityEntity, 0, len(places))
	for i, place := range places {
		info, ok := countryLookup[place.CountryCode]
		if !ok {
			continue
		}
		cities = append(cities, model.CityEntity{
			ID:              i + 1,
			Name:            place.Name,
			ThreeLetterCode: info.CCA3,
			Currency:        info.Currency,
			Country:         place.CountryName,
		})
	}
	return cities, nil
}

// IsGeoNames returns true if the GeoNames provider is active
func (s *CitiesService) IsGeoNames() bool {
	return s.geoNamesProvider != nil
}
