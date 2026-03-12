package bondportfolio

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
)

func mapBondLotToDbModel(lot bonds.BondLot) bondsdb.BondPositionLotDb {
	return bondsdb.BondPositionLotDb{
		Id:               lot.Id,
		Figi:             lot.Figi,
		OpeningDate:      lot.OpeningDate,
		ModificationDate: lot.ModificationDate,
		AccountId:        lot.AccountId,
		Quantity:         lot.Quantity,
		PricePerUnit:     lot.PricePerUnit,
	}
}
