package bondportfolio

import (
	"errors"
	"time"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
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
	if lot.PricePerUnit <= 0 && lot.PricePerUnitInRUB <= 0 {
		return lot, errors.New("Either the price per unit or price per unit in RUB must be provided")
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

	if lot.Bond.IsRubleBond() && lot.PricePerUnitInRUB > 0.0 {
		lot.PricePerUnit = lot.PricePerUnitInRUB
	}

	if lot.Bond.IsRubleBond() && lot.PricePerUnit > 0.0 {
		lot.PricePerUnitInRUB = lot.PricePerUnit
	}

	if !lot.Bond.IsRubleBond() {
		forexRate, err := forexservice.GetExchangeRate(lot.Bond.NominalCurrency, lot.Bond.Currency, lot.OpeningDate)
		if err != nil {
			return lot, nil
		}
		if lot.PricePerUnit > 0 {
			lot.PricePerUnitInRUB = forexRate.Rate * lot.PricePerUnit
		}
		if lot.PricePerUnitInRUB > 0 {
			lot.PricePerUnit = lot.PricePerUnitInRUB / forexRate.Rate
		}
	}

	return lot, nil
}
