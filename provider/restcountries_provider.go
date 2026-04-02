package provider

import (
	"fmt"
	"simplyutil-server/model"
	"simplyutil-server/util"
)

// CountryLookup holds country metadata keyed by 2-letter ISO code
type CountryLookup struct {
	CCA3     string
	Currency string
}

// RestCountriesProvider implements city fetching using the REST Countries API
type RestCountriesProvider struct{}

// NewRestCountriesProvider creates a new REST Countries provider
func NewRestCountriesProvider() *RestCountriesProvider {
	return &RestCountriesProvider{}
}

// FetchCities returns capital cities for all countries
func (p *RestCountriesProvider) FetchCities() ([]model.CityEntity, error) {
	countries, err := p.fetchCountries()
	if err != nil {
		return nil, err
	}

	cities := make([]model.CityEntity, 0, len(countries))
	id := 1
	for _, country := range countries {
		if len(country.Capital) == 0 || country.Currencies == nil {
			continue
		}
		var currencyCode string
		for code := range country.Currencies {
			currencyCode = code
			break
		}
		cities = append(cities, model.CityEntity{
			ID:              id,
			Name:            country.Capital[0],
			ThreeLetterCode: country.CCA3,
			Currency:        currencyCode,
			Country:         country.Name.Common,
		})
		id++
	}
	return cities, nil
}

// FetchCountryLookup returns a map from 2-letter ISO code to country metadata
func (p *RestCountriesProvider) FetchCountryLookup() (map[string]CountryLookup, error) {
	countries, err := p.fetchCountries()
	if err != nil {
		return nil, err
	}

	lookup := make(map[string]CountryLookup, len(countries))
	for _, country := range countries {
		if country.CCA2 == "" || country.Currencies == nil {
			continue
		}
		var currencyCode string
		for code := range country.Currencies {
			currencyCode = code
			break
		}
		lookup[country.CCA2] = CountryLookup{
			CCA3:     country.CCA3,
			Currency: currencyCode,
		}
	}
	return lookup, nil
}

func (p *RestCountriesProvider) fetchCountries() ([]model.RestCountryResponse, error) {
	url := "https://restcountries.com/v3.1/all?fields=name,cca2,cca3,capital,currencies"
	var countries []model.RestCountryResponse
	if err := util.HTTPGetJSON(url, &countries); err != nil {
		return nil, fmt.Errorf("rest-countries API error: %w", err)
	}
	return countries, nil
}
