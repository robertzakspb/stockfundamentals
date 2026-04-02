package bondportfolio

import (
	"errors"
	"time"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
)

func validateLot(lot bonds.BondLot) (bonds.BondLot, error) {
	if lot.Figi != "" {
		bond, err := bondservice.GetBondByFigi(lot.Figi)
		lot.Figi = bond.Figi
		if err != nil {
			return lot, err
		}
	} else if lot.Isin != "" {
		bond, err := bondservice.GetBondByIsin(lot.Isin)
		lot.Isin = bond.Isin
		if err != nil {
			return lot, err
		}
	} else {
		return lot, errors.New("Missing both figi and ISIN in the bond")
	}

	if lot.Quantity < 0 {
		return lot, errors.New("Invalid quantity")
	}
	if lot.OpeningDate.After(time.Now()) {
		return lot, errors.New("Invalid opening date")
	}
	if lot.ModificationDate.After(time.Now()) {
		return lot, errors.New("Invalid modification date")
	}
	if lot.PricePerUnit < 0 {
		return lot, errors.New("Invalid price per unit")
	}

	return lot, nil
}
