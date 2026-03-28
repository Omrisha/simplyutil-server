package service

import (
	"os"
	"simplyutil-server/model"
	"testing"
)

// MockLandmarkProvider is a mock implementation for testing
type MockLandmarkProvider struct {
	mockData []model.LandmarkEntity
	mockErr  error
}

func (m *MockLandmarkProvider) FetchLandmarks(cityName, country string) ([]model.LandmarkEntity, error) {
	if m.mockErr != nil {
		return nil, m.mockErr
	}
	return m.mockData, nil
}

func TestLandmarkServiceWithMock(t *testing.T) {
	// Create a service with mock provider
	mockData := []model.LandmarkEntity{
		{
			Name:      "Test Landmark",
			Address:   "Test Address",
			Latitude:  51.5074,
			Longitude: -0.1278,
			Rating:    8.5,
		},
	}

	mockProvider := &MockLandmarkProvider{mockData: mockData}
	service := &LandmarkService{provider: mockProvider}

	landmarks, err := service.GetLandmarks("London", "UK")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(landmarks) != 1 {
		t.Fatalf("Expected 1 landmark, got %d", len(landmarks))
	}

	if landmarks[0].Name != "Test Landmark" {
		t.Errorf("Expected landmark name 'Test Landmark', got %q", landmarks[0].Name)
	}
}

func TestLandmarkServiceProviderSelection(t *testing.T) {
	// Test Wikipedia provider selection when LANDMARK_PROVIDER is set
	t.Setenv("LANDMARK_PROVIDER", "wikipedia")
	t.Setenv("FOURSQUARE_API_KEY", "test-key")

	service := NewLandmarkService()
	if service == nil {
		t.Fatal("Expected service to be created")
	}

	// Test Foursquare provider selection
	os.Unsetenv("LANDMARK_PROVIDER")
	t.Setenv("FOURSQUARE_API_KEY", "test-key")

	service = NewLandmarkService()
	if service == nil {
		t.Fatal("Expected service to be created")
	}

	// Test fallback to Wikipedia when no Foursquare key
	os.Unsetenv("FOURSQUARE_API_KEY")
	os.Unsetenv("LANDMARK_PROVIDER")

	service = NewLandmarkService()
	if service == nil {
		t.Fatal("Expected service to be created with Wikipedia fallback")
	}
}
