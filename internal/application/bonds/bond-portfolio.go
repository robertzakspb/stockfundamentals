package bondservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	
)

func SaveBondPositionLot(lot bonds.BondPosition) error {
	_, err := GetBondByFigi(lot.Figi)
	if err != nil {
		return err
	}

	return nil
}