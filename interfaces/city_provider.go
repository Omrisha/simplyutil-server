package interfaces

import "simplyutil-server/model"

// CityProvider interface for different city data sources
type CityProvider interface {
	FetchCities() ([]model.CityEntity, error)
}
