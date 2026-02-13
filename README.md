# SimplyUtil Go Backend Server

A unified backend server that consolidates all API calls for the SimplyUtil iOS app.

## Features

- 🌍 **Cities API**: List of countries with currencies
- 🏛️ **Landmarks API**: Tourist attractions from Foursquare
- ☁️ **Weather API**: Forecasts from Open-Meteo
- 💱 **Exchange Rates API**: Real-time currency rates
- ⚡ **Concurrent Fetching**: Parallel API calls for better performance
- 🔒 **Secure**: API keys hidden server-side
- 🚀 **Fast**: Deployed on Railway/Render/Fly.io

## Setup

### 1. Install Dependencies

```bash
cd server
go mod download
```

### 2. Set Environment Variables

Create a `.env` file:

```bash
PORT=8080
FOURSQUARE_API_KEY=YOUR_CLIENT_ID+YOUR_CLIENT_SECRET
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

# Get landmarks
curl "http://localhost:8080/api/v1/landmarks?city=London&country=England"

# Get weather
curl "http://localhost:8080/api/v1/weather?city=London"

# Get exchange rates
curl http://localhost:8080/api/v1/rates/USD

# Get all city data at once
curl http://localhost:8080/api/v1/cities/London/England
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

### GET /api/v1/cities
Get list of all cities with currencies.

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

### GET /api/v1/landmarks?city={name}&country={name}
Get landmarks for a city.

**Response:**
```json
{
  "landmarks": [
    {
      "name": "Tower Bridge",
      "address": "Tower Bridge Rd",
      "latitude": 51.5055,
      "longitude": -0.0754,
      "rating": 9.5
    }
  ],
  "count": 20
}
```

### GET /api/v1/weather?city={name}
Get weather forecast for a city.

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
Get all data for a city in one request (landmarks + weather + rates).

**Response:**
```json
{
  "city": "London",
  "country": "England",
  "landmarks": [...],
  "weather": {...},
  "rates": {...}
}
```

## Deployment

### Deploy to Railway

1. Install Railway CLI: `npm install -g @railway/cli`
2. Login: `railway login`
3. Initialize: `railway init`
4. Deploy: `railway up`
5. Set env vars: `railway variables set FOURSQUARE_API_KEY=xxx`

### Deploy to Render

1. Connect your GitHub repo
2. Create new Web Service
3. Build command: `go build -o server`
4. Start command: `./server`
5. Add environment variables in dashboard

### Deploy to Fly.io

1. Install flyctl: `curl -L https://fly.io/install.sh | sh`
2. Login: `fly auth login`
3. Launch: `fly launch`
4. Deploy: `fly deploy`
5. Set secrets: `fly secrets set FOURSQUARE_API_KEY=xxx`

## Project Structure

```
server/
├── main.go          # Entry point, router setup
├── handlers.go      # HTTP request handlers
├── services.go      # External API integrations
├── models.go        # Data structures
├── go.mod           # Go dependencies
└── README.md        # This file
```

## Adding Caching (Optional)

To add Redis caching:

```go
// Add to go.mod
require github.com/go-redis/redis/v8 v8.11.5

// In services.go
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
2. **Rate limiting**: Use `gin-contrib/rate` package
3. **Connection pooling**: Configure http.Client properly
4. **Gzip compression**: Add middleware for response compression
5. **Monitoring**: Use Prometheus + Grafana

## License

MIT
