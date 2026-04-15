package bondportfolio

import (
	"strings"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
)

func mapBondLotToDbModel(lot bonds.BondLot) bondsdb.BondPositionLotDb {
	return bondsdb.BondPositionLotDb{
		Id:                lot.Id,
		Figi:              strings.ToUpper(lot.Figi),
		Isin:              strings.ToUpper(lot.Isin),
		OpeningDate:       lot.OpeningDate,
		ModificationDate:  lot.ModificationDate,
		AccountId:         lot.AccountId,
		Quantity:          lot.Quantity,
		PricePerUnitPercentage:      lot.PricePerUnitPercentage,
		AccruedInterest:   lot.AccruedInterest,
	}
}

func mapDbBondToDomain(lot bondsdb.BondPositionLotDb) bonds.BondLot {
	return bonds.BondLot{
		Id:                lot.Id,
		Figi:              lot.Figi,
		Isin:              lot.Isin,
		OpeningDate:       lot.OpeningDate,
		ModificationDate:  lot.ModificationDate,
		AccountId:         lot.AccountId,
		Quantity:          lot.Quantity,
		PricePerUnitPercentage:      lot.PricePerUnitPercentage,
		AccruedInterest:   lot.AccruedInterest,
	}
}
