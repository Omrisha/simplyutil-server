# SimplyUtil Go Backend Server

A unified backend server that consolidates all API calls for the SimplyUtil iOS app.

## Features

- 🌍 **Cities API**: Cities with country and currency info, from GeoNames or REST Countries
- 🏛️ **Landmarks API**: Tourist attractions from Wikipedia (default) or Foursquare
- ☁️ **Weather API**: Forecasts from Open-Meteo
- 💱 **Exchange Rates API**: Real-time currency rates
- ⚡ **Concurrent Fetching**: Parallel API calls for better performance
- 🔒 **Secure**: API keys hidden server-side
- 🚀 **Deployed on Fly.io**: See [.github/workflows/fly-deploy.yml](.github/workflows/fly-deploy.yml)

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Set Environment Variables

Create a `.env` file. Every variable is optional; the defaults are shown below.

```bash
# Port to listen on (default: 8080)
PORT=8080

# Landmarks. Wikipedia is used when LANDMARK_PROVIDER=wikipedia or when
# FOURSQUARE_API_KEY is unset, so no key is needed to get started.
LANDMARK_PROVIDER=wikipedia
FOURSQUARE_API_KEY=YOUR_FOURSQUARE_SERVICE_KEY

# Cities. Without a GeoNames username the server falls back to REST Countries,
# which returns capital cities only and supports no search or pagination.
GEONAMES_USERNAME=YOUR_GEONAMES_USERNAME

# Exchange rates
EXCHANGE_RATE_API_KEY=YOUR_EXCHANGERATE_API_KEY
```

### 3. Run Locally

```bash
go run .
```

Server will start at `http://localhost:8080`

### 4. Test Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Get all cities
curl http://localhost:8080/api/v1/cities

# Search cities (requires GEONAMES_USERNAME)
curl "http://localhost:8080/api/v1/cities?search=lon&country=GB&page=1&pageSize=20"

# Same, using the QUERY method with a JSON body
curl -X QUERY http://localhost:8080/api/v1/cities \
  -H "Content-Type: application/json" \
  -d '{"search": "lon", "country": "GB", "page": 1, "pageSize": 20}'

# Get landmarks
curl "http://localhost:8080/api/v1/landmarks?city=London&country=England"

# Same, using the QUERY method with a JSON body
curl -X QUERY http://localhost:8080/api/v1/landmarks \
  -H "Content-Type: application/json" \
  -d '{"city": "London", "country": "England"}'

# Get weather
curl "http://localhost:8080/api/v1/weather?city=London"

# Get exchange rates
curl http://localhost:8080/api/v1/rates/USD

# Get all city data at once
curl http://localhost:8080/api/v1/cities/London/England
```

## Running Tests

Run all unit tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Integration tests that call live APIs skip themselves unconditionally, so
`go test ./...` never touches the network. Remove the `t.Skip` in the relevant
test to run one against the real API.

Run specific package tests:
```bash
go test ./handler
go test ./model
go test ./service
go test ./util
```

## API Documentation

### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "ok",
  "timestamp": 1707868800
}
```

### GET /api/v1/cities &nbsp;·&nbsp; QUERY /api/v1/cities
Get cities with country and currency info. With GeoNames this returns cities with a
population of at least 50,000; without it, REST Countries capital cities. With the QUERY method, filters are sent as a JSON body instead of query parameters:
```json
{
  "search": "lon",
  "country": "GB",
  "page": 1,
  "pageSize": 100
}
```

Filtering and pagination require the GeoNames provider. Without `GEONAMES_USERNAME`
set, a request carrying any of `search`, `country`, `page`, or `pageSize` returns
`501 Not Implemented` rather than silently returning the unfiltered list. `page`
defaults to 1 and `pageSize` to 100 (maximum 1000); values outside those bounds
fall back to the defaults.

QUERY request bodies are capped at 1 MiB; a larger body returns `413 Request Entity Too Large`.

**Response:**
```json
{
  "cities": [
    {
      "id": 1,
      "name": "London",
      "threeLetterCode": "GBR",
      "currency": "GBP",
      "country": "United Kingdom"
    }
  ],
  "count": 195
}
```

Filtered requests echo the pagination back:
```json
{
  "cities": [...],
  "count": 20,
  "page": 1,
  "pageSize": 20
}
```

### GET /api/v1/landmarks?city={name}&country={name} &nbsp;·&nbsp; QUERY /api/v1/landmarks
Get landmarks for a city. With the QUERY method, send `{"city": "London", "country": "England"}` as a JSON body.
`city` is required; `country` is optional and only helps disambiguate geocoding.

**Response:** `rating` is always 0 and `category` always `"Landmark"` when the
Wikipedia provider is active. `imageUrl` and `category` are omitted when empty.
```json
{
  "landmarks": [
    {
      "name": "Tower Bridge",
      "address": "Tower Bridge Rd",
      "latitude": 51.5055,
      "longitude": -0.0754,
      "rating": 9.5,
      "imageUrl": "https://.../tower-bridge.jpg",
      "category": "Monument"
    }
  ],
  "count": 20
}
```

### GET /api/v1/weather?city={name} &nbsp;·&nbsp; QUERY /api/v1/weather
Get weather forecast for a city (next 24 hours, hourly). With the QUERY method, send `{"city": "London"}` as a JSON body.

**Response:**
```json
{
  "weather": {
    "latitude": 51.5074,
    "longitude": -0.1278,
    "hourly": [
      {
        "time": "2026-02-13T12:00",
        "temperature": 15.5,
        "windSpeed": 12.3,
        "relativeHumidity": 65
      }
    ]
  }
}
```

### GET /api/v1/rates/{currency}
Get exchange rates for a base currency.

**Response:**
```json
{
  "baseCurrency": "USD",
  "rates": {
    "EUR": 0.85,
    "GBP": 0.73,
    "JPY": 110.5
  },
  "timestamp": "2026-02-13T12:00:00Z"
}
```

### GET /api/v1/cities/{city}/{country}
Get all data for a city in one request (landmarks + weather + rates), fetched
concurrently. Rates are always USD-based, regardless of the city's own currency.

Partial failures still return `200`: any section that failed is replaced by a
`<section>_error` key rather than failing the whole request.

**Response:**
```json
{
  "city": "London",
  "country": "England",
  "landmarks": [...],
  "weather": {...},
  "rates_error": "exchange-rate API error: ..."
}
```

**Note:** the QUERY method is not available on this endpoint or on `/api/v1/rates/{currency}`.

## Deployment

### Pre-Deployment Checklist

Before deploying to any platform, always:

```bash
# Run tests
go test ./...

# Build to verify compilation
go build -o simplyutil-server .

# Test the binary
./simplyutil-server &
curl http://localhost:8080/health
pkill simplyutil-server
```

### Deploy to Railway

1. **Run tests**: `go test ./...`
2. Install Railway CLI: `npm install -g @railway/cli`
3. Login: `railway login`
4. Initialize: `railway init`
5. Deploy: `railway up`
6. Set env vars: `railway variables set FOURSQUARE_API_KEY=xxx`

### Deploy to Render

1. **Run tests**: `go test ./...`
2. Connect your GitHub repo
3. Create new Web Service
4. Build command: `go test ./... && go build -o server`
5. Start command: `./server`
6. Add environment variables in dashboard

### Deploy to Fly.io

1. **Run tests**: `go test ./...`
2. Install flyctl: `curl -L https://fly.io/install.sh | sh`
3. Login: `fly auth login`
4. Launch: `fly launch`
5. Deploy: `fly deploy`
6. Set secrets: `fly secrets set FOURSQUARE_API_KEY=xxx`

**Note**: the [Dockerfile](Dockerfile) already runs `go test ./...` before building,
so a failing test fails the deploy. Pushes to `main` deploy automatically via
[.github/workflows/fly-deploy.yml](.github/workflows/fly-deploy.yml).

## Project Structure

This project follows clean architecture principles with clear separation of concerns. For detailed information about the architecture, design patterns, and project organization, see [ARCHITECTURE.md](ARCHITECTURE.md).

## Adding Caching (Optional)

To add Redis caching:

```go
// Add to go.mod
require github.com/go-redis/redis/v8 v8.11.5

// In service/landmark_service.go
var rdb = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

func fetchLandmarksWithCache(city string) {
    // Check cache first
    cached, err := rdb.Get(ctx, "landmarks:"+city).Result()
    if err == nil {
        return parseCached(cached)
    }
    
    // Fetch from API
    data := fetchFromAPI(city)
    
    // Cache for 1 hour
    rdb.Set(ctx, "landmarks:"+city, data, time.Hour)
    
    return data
}
```

## Performance Tips

1. **Enable caching**: Add Redis for frequently accessed data
2. **Rate limiting**: Add a limiter such as `github.com/ulule/limiter/v3` or `github.com/didip/tollbooth`
3. **Connection pooling**: Configure http.Client properly
4. **Gzip compression**: Add middleware for response compression
5. **Monitoring**: Use Prometheus + Grafana

## License

MIT
