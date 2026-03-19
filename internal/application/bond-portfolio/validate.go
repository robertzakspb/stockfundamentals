package bondportfolio

import (
	"errors"
	"time"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
)

func validateLot(lot bonds.BondLot) error {
	if lot.Figi != "" {
		_, err := bondservice.GetBondByFigi(lot.Figi)
		if err != nil {
			return err
		}
	} else if lot.Isin != "" {
		_, err := bondservice.GetBondByIsin(lot.Isin)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Missing both figi and ISIN in the bond")
	}

	if lot.Quantity < 0 {
		return errors.New("Invalid quantity")
	}
	if lot.OpeningDate.After(time.Now()) {
		return errors.New("Invalid opening date")
	}
	if lot.ModificationDate.After(time.Now()) {
		return errors.New("Invalid modification date")
	}
	if lot.PricePerUnit < 0 {
		return errors.New("Invalid price per unit")
	}

	return nil
}
