# Project Architecture

This project follows clean architecture principles with clear separation of concerns.

## Directory Structure

```
simplyutil-server/
├── interfaces/          # Interface definitions
│   └── landmark_provider.go
├── provider/           # Provider implementations
│   ├── foursquare_provider.go
│   └── wikipedia_provider.go
├── service/            # Business logic layer
│   ├── cities_service.go
│   ├── landmark_service.go
│   ├── rates_service.go
│   └── weather_service.go
├── model/              # Data models
│   ├── entity.go       # Domain entities
│   └── response.go     # API response types
├── handler/            # HTTP handlers
│   └── handler.go
├── util/               # Utility functions
│   ├── http_client.go
│   └── geocoding.go
└── main.go             # Application entry point
```

## Layer Responsibilities

### 1. Interfaces Layer (`interfaces/`)
- Defines contracts for providers
- Example: `LandmarkProvider` interface

### 2. Provider Layer (`provider/`)
- Implements external API integrations
- **FoursquareProvider**: Fetches landmarks from Foursquare Places API
- **WikipediaProvider**: Fetches landmarks from Wikipedia/Wikimedia

### 3. Service Layer (`service/`)
- Contains business logic
- Orchestrates providers and data transformations
- Services:
  - **LandmarkService**: Manages landmark fetching with provider selection
  - **WeatherService**: Handles weather data
  - **RatesService**: Manages exchange rates
  - **CitiesService**: Fetches city/country data

### 4. Model Layer (`model/`)
- **entity.go**: Domain models (LandmarkEntity, WeatherEntity, etc.)
- **response.go**: External API response structures

### 5. Handler Layer (`handler/`)
- HTTP request handlers using Gin framework
- Coordinates service calls
- Handles concurrent requests for city data

### 6. Util Layer (`util/`)
- Shared utility functions
- HTTP client helpers
- Geocoding functionality

## Design Patterns Used

### 1. **Strategy Pattern**
The `LandmarkProvider` interface allows switching between Foursquare and Wikipedia providers:
```go
type LandmarkProvider interface {
    FetchLandmarks(cityName, country string) ([]model.LandmarkEntity, error)
}
```

### 2. **Dependency Injection**
Services are injected into handlers:
```go
type Handler struct {
    landmarkService *service.LandmarkService
    weatherService  *service.WeatherService
    ratesService    *service.RatesService
    citiesService   *service.CitiesService
}
```

### 3. **Factory Pattern**
Providers and services are created through constructor functions:
```go
func NewLandmarkService() *LandmarkService
func NewFoursquareProvider() *FoursquareProvider
func NewWikipediaProvider() *WikipediaProvider
```

## Configuration

### Environment Variables
- `PORT`: Server port (default: 8080)
- `LANDMARK_PROVIDER`: Choose "foursquare" or "wikipedia" (default: foursquare)
- `FOURSQUARE_API_KEY`: API key for Foursquare (if using Foursquare provider)
- `EXCHANGE_RATE_API_KEY`: API key for exchange rates (optional)

### Provider Selection Logic
The `LandmarkService` automatically selects the appropriate provider:
1. If `LANDMARK_PROVIDER=wikipedia`, uses Wikipedia
2. If `FOURSQUARE_API_KEY` is not set, falls back to Wikipedia
3. Otherwise, uses Foursquare

## API Endpoints

### Cities
- `GET /api/v1/cities` - List all cities
- `GET /api/v1/cities/:name/:country` - Get all data for a city (landmarks, weather, rates)

### Individual Resources
- `GET /api/v1/landmarks?city=X&country=Y` - Get landmarks
- `GET /api/v1/weather?city=X` - Get weather
- `GET /api/v1/rates/:currency` - Get exchange rates

### Health Check
- `GET /health` - Server health status

## Benefits of This Architecture

1. **Testability**: Each layer can be tested independently
2. **Maintainability**: Clear separation makes code easier to understand
3. **Flexibility**: Easy to add new providers or services
4. **Scalability**: Layers can be scaled independently
5. **Reusability**: Services and utilities can be reused across handlers

## Adding New Features

### Adding a New Provider
1. Create a new file in `provider/`
2. Implement the `LandmarkProvider` interface
3. Update `LandmarkService` to include the new provider option

### Adding a New Service
1. Create a new file in `service/`
2. Define the service struct and constructor
3. Inject the service into `Handler`
4. Add handler methods to expose the service

### Adding a New Endpoint
1. Create a handler method in `handler/handler.go`
2. Register the route in `main.go`
