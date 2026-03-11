package bonds

import (
	"time"

	"github.com/google/uuid"
)

type BondPosition struct {
	Id               uuid.UUID
	Figi             string
	OpeningDate      time.Time
	ModificationDate time.Time
	AccountId        uuid.UUID
	Quantity         float64
	PricePerUnit     float64
}

func (b BondPosition) PricePerUnitPercentage(nominalValue float64) float64 {
	if nominalValue != 0 {
		return b.PricePerUnit / nominalValue
	}
	return -1
}
