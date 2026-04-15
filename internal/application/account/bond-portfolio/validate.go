package bondportfolio

import (
	"errors"
	"time"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bondservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
)

func validateLot(lot bonds.BondLot) (bonds.BondLot, error) {
	if lot.Figi == "" && lot.Isin == "" {
		return lot, errors.New("Missing both figi and ISIN in the bond")
	}
	if lot.Quantity <= 0 {
		return lot, errors.New("Invalid quantity")
	}
	if lot.OpeningDate.After(time.Now()) {
		return lot, errors.New("Invalid opening date")
	}
	if lot.ModificationDate.After(time.Now()) {
		return lot, errors.New("Invalid modification date")
	}
	if lot.PricePerUnitPercentage <= 0 && lot.PricePerUnitPercentage > 2000 {
		return lot, errors.New("The price per unit as percentage should be greater than 0 and no greater than 2000 ")
	}

	return lot, nil
}

func addMissingInformationToLot(lot bonds.BondLot) (bonds.BondLot, error) {
	if lot.Figi == "" {
		bond, err := bondservice.GetBondByIsin(lot.Isin)
		lot.Bond = bond
		lot.Figi = bond.Figi

		if err != nil {
			return lot, err
		}
	}

	if lot.Isin == "" {
		lot.Isin = lot.Bond.Isin
	}

	return lot, nil
}
