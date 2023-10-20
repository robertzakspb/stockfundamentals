package securityinfo

import (
	"time"
	"github.com/google/uuid"
)

type PolygonSecurityDTO struct {
	Results []struct {
		Ticker          string    `json:"ticker"`
		Name            string    `json:"name"`
		Market          string    `json:"market"`
		Locale          string    `json:"locale"`
		PrimaryExchange string    `json:"primary_exchange"`
		Type            string    `json:"type"`
		Active          bool      `json:"active"`
		CurrencyName    string    `json:"currency_name"`
		Cik             string    `json:"cik,omitempty"`
		CompositeFigi   string    `json:"composite_figi,omitempty"`
		ShareClassFigi  string    `json:"share_class_figi,omitempty"`
		LastUpdatedUtc  time.Time `json:"last_updated_utc"`
	} `json:"results"`
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Count     int    `json:"count"`
	NextURL   string `json:"next_url"`
	Error string `json:"error"`
}

// Converts the API DTO into a BLL entity
func (dto PolygonSecurityDTO) convertToPolygonSecurity() []PolygonSecurity {
	parsedSecurities := []PolygonSecurity{}
	for _, security := range dto.Results {
		uuid, err := uuid.NewRandom()
		if err != nil {
			continue
		}
		parsedSecurity := PolygonSecurity {
			uuid.String(),
			security.Ticker,
			security.Name,
			security.PrimaryExchange,
			security.CurrencyName,
		}
		parsedSecurities = append(parsedSecurities, parsedSecurity)
	}
	return parsedSecurities
}
