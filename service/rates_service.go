package service

import (
	"fmt"
	"os"
	"simplyutil-server/model"
	"simplyutil-server/util"
	"time"
)

// RatesService handles exchange rate operations
type RatesService struct {
	apiKey string
}

// NewRatesService creates a new rates service
func NewRatesService() *RatesService {
	apiKey := os.Getenv("EXCHANGE_RATE_API_KEY")
	if apiKey == "" {
		// Fallback to hardcoded key (consider moving to env var)
		apiKey = "c04b66e4d1f1f147c60834b3"
	}
	return &RatesService{
		apiKey: apiKey,
	}
}

// GetRates fetches exchange rates for a given base currency
func (s *RatesService) GetRates(baseCurrency string) (model.RatesEntity, error) {
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", s.apiKey, baseCurrency)

	var ratesResponse model.ExchangeRateResponse
	if err := util.HTTPGetJSON(url, &ratesResponse); err != nil {
		return model.RatesEntity{}, fmt.Errorf("exchange-rate API error: %w", err)
	}

	return model.RatesEntity{
		BaseCurrency: ratesResponse.BaseCode,
		Rates:        ratesResponse.ConversionRates,
		Timestamp:    time.Unix(ratesResponse.TimeLastUpdateUnix, 0),
	}, nil
}
