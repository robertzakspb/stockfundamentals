package bondportfolio

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
)

// TODO: Optimize to remove the loop
func SaveBondPositionLots(lots []bonds.BondLot) error {
	for _, lot := range lots {
		err := SaveBondPositionLot(lot)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveBondPositionLot(lot bonds.BondLot) error {
	lot, err := validateLot(lot)
	if err != nil {
		return err
	}

	lot, err = addMissingInformationToLot(lot)
	if err != nil {
		return err
	}

	mappedLot := mapBondLotToDbModel(lot)

	err = bondsdb.SaveBondPositionLots([]bondsdb.BondPositionLotDb{mappedLot})
	if err != nil {
		return err
	}

	return nil
}
