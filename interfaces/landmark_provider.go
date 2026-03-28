package interfaces

import "simplyutil-server/model"

// LandmarkProvider interface for different landmark data sources
type LandmarkProvider interface {
	FetchLandmarks(cityName, country string) ([]model.LandmarkEntity, error)
}
