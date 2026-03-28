package model

// FoursquareV3Response for Foursquare Places API (new structure)
type FoursquareV3Response struct {
	Results []struct {
		FsqPlaceID string  `json:"fsq_place_id"`
		Name       string  `json:"name"`
		Latitude   float64 `json:"latitude"`
		Longitude  float64 `json:"longitude"`
		Location   struct {
			Address          string `json:"address"`
			FormattedAddress string `json:"formatted_address"`
		} `json:"location"`
		Categories []struct {
			Name string `json:"name"`
		} `json:"categories"`
		Photos []struct {
			Prefix string `json:"prefix"`
			Suffix string `json:"suffix"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"photos"`
		Rating float64 `json:"rating"`
	} `json:"results"`
}

// OpenMeteoResponse for Open-Meteo API
type OpenMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Hourly    struct {
		Time             []string  `json:"time"`
		Temperature      []float64 `json:"temperature_2m"`
		WindSpeed        []float64 `json:"wind_speed_10m"`
		RelativeHumidity []int     `json:"relative_humidity_2m"`
	} `json:"hourly"`
}

// ExchangeRateResponse for Exchange Rate API
type ExchangeRateResponse struct {
	Result             string             `json:"result"`
	BaseCode           string             `json:"base_code"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
}

// RestCountryResponse for REST Countries API
type RestCountryResponse struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	CCA3       string   `json:"cca3"`
	Capital    []string `json:"capital"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

// NominatimResponse for geocoding
type NominatimResponse struct {
	Lat float64 `json:"lat,string"`
	Lon float64 `json:"lon,string"`
}

// WikipediaGeoSearchResponse for Wikipedia geosearch API
type WikipediaGeoSearchResponse struct {
	Query struct {
		GeoSearch []struct {
			PageID int     `json:"pageid"`
			Title  string  `json:"title"`
			Lat    float64 `json:"lat"`
			Lon    float64 `json:"lon"`
			Dist   float64 `json:"dist"`
		} `json:"geosearch"`
	} `json:"query"`
}

// WikipediaPageResponse for Wikipedia page details API
type WikipediaPageResponse struct {
	Query struct {
		Pages map[string]struct {
			PageID    int    `json:"pageid"`
			Title     string `json:"title"`
			Thumbnail struct {
				Source string `json:"source"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnail"`
			Coordinates []struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"coordinates"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}
