package accountsapi

import (
	"time"

	"github.com/google/uuid"
)

type AccountDto struct {
	Id              uuid.UUID `json:"id"`
	OpeningDate     time.Time `json:"openingDate"`
	Type            string    `json:"type"`
	Broker          string    `json:"broker"`
	Holder          string    `json:"holder"`
	PrimaryCurrency string    `json:"primaryCurrency"`
	CashBalance     float64   `json:"cashBalance"`
}
