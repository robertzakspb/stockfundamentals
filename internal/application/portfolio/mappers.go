package portfolio

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/portfoliodb"
)

func mapLotDbToLot(lotDb portfoliodb.LotDb) lot.Lot {
	return lot.Lot{
		Id:           lotDb.Id,
		CreatedAt:    lotDb.CreatedAt,
		UpdatedAt:    lotDb.UpdatedAt,
		Quantity:     lotDb.Quantity,
		PricePerUnit: lotDb.PricePerUnit,
		Currency:     lotDb.Currency,
		AccountId:    lotDb.AccountId,
		SecurityId:   lotDb.Figi,
		CurrentPL:    0,
	}
}

func mapLotToDbLot(lots []lot.Lot) []portfoliodb.LotDb {
	dbLots := []portfoliodb.LotDb{}
	for _, lot := range lots {
		dbLot := portfoliodb.LotDb{
			Id:           lot.Id,
			CreatedAt:    lot.CreatedAt,
			UpdatedAt:    lot.UpdatedAt,
			Quantity:     lot.Quantity,
			PricePerUnit: lot.PricePerUnit,
			Currency:     lot.Currency,
			AccountId:    lot.AccountId,
			Figi:         lot.SecurityId,
		}
		dbLots = append(dbLots, dbLot)

	}
	return dbLots
}
