package service

import (
	"fmt"
	"simplyutil-server/model"
	"simplyutil-server/util"
)

// CitiesService handles cities/countries data operations
type CitiesService struct{}

// NewCitiesService creates a new cities service
func NewCitiesService() *CitiesService {
	return &CitiesService{}
}

// GetCities fetches list of countries/cities from REST Countries API
func (s *CitiesService) GetCities() ([]model.CityEntity, error) {
	url := "https://restcountries.com/v3.1/all?fields=name,cca3,capital,currencies"

	var countries []model.RestCountryResponse
	if err := util.HTTPGetJSON(url, &countries); err != nil {
		return nil, fmt.Errorf("rest-countries API error: %w", err)
	}

	// Convert to CityEntity format
	cities := make([]model.CityEntity, 0)
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

		city := model.CityEntity{
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
