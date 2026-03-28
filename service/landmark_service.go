package service

import (
	"os"
	"simplyutil-server/interfaces"
	"simplyutil-server/model"
	"simplyutil-server/provider"
)

// LandmarkService handles landmark-related operations
type LandmarkService struct {
	provider interfaces.LandmarkProvider
}

// NewLandmarkService creates a new landmark service with the appropriate provider
func NewLandmarkService() *LandmarkService {
	var landmarkProvider interfaces.LandmarkProvider

	// Check configuration to determine which provider to use
	providerType := os.Getenv("LANDMARK_PROVIDER")
	foursquareKey := os.Getenv("FOURSQUARE_API_KEY")

	if providerType == "wikipedia" || foursquareKey == "" {
		landmarkProvider = provider.NewWikipediaProvider()
	} else {
		landmarkProvider = provider.NewFoursquareProvider()
	}

	return &LandmarkService{
		provider: landmarkProvider,
	}
}

// GetLandmarks fetches landmarks for a given city
func (s *LandmarkService) GetLandmarks(cityName, country string) ([]model.LandmarkEntity, error) {
	return s.provider.FetchLandmarks(cityName, country)
}
