package provider

import (
	"encoding/json"
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

// restCountriesErrorEnvelope is what the retired v3.1 API answers with now: an
// object carrying a deprecation message, served with HTTP 200.
type restCountriesErrorEnvelope struct {
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (p *RestCountriesProvider) fetchCountries() ([]model.RestCountryResponse, error) {
	url := "https://restcountries.com/v3.1/all?fields=name,cca2,cca3,capital,currencies"

	var raw json.RawMessage
	if err := util.HTTPGetJSON(url, &raw); err != nil {
		return nil, fmt.Errorf("rest-countries API error: %w", err)
	}

	var countries []model.RestCountryResponse
	if err := json.Unmarshal(raw, &countries); err == nil {
		return countries, nil
	}

	// The array we expected is gone. Report why rather than letting an opaque
	// "cannot unmarshal object into []RestCountryResponse" reach the client.
	var envelope restCountriesErrorEnvelope
	if json.Unmarshal(raw, &envelope) == nil && len(envelope.Errors) > 0 {
		return nil, fmt.Errorf("rest-countries API is no longer available: %s (set GEONAMES_USERNAME to use GeoNames instead)", envelope.Errors[0].Message)
	}
	return nil, fmt.Errorf("rest-countries API returned an unexpected response shape (set GEONAMES_USERNAME to use GeoNames instead)")
}
