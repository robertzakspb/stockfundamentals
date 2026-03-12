package bondportfolio

import (
	"errors"
	"time"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
)

func SaveBondPositionLot(lot bonds.BondLot) error {
	err := validateLot(lot)
	if err != nil {
		return err
	}

	mappedLot := mapBondLotToDbModel(lot)

	err = bondsdb.SaveBondPositionLots([]bondsdb.BondPositionLotDb{mappedLot})

	return nil
}

func validateLot(lot bonds.BondLot) error {
	_, err := bondservice.GetBondByFigi(lot.Figi)
	if err != nil {
		return err
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
