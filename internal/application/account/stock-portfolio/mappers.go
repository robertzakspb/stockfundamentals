package portfolio

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/portfoliodb"
)

func mapDbLotsToLots(dbLots []portfoliodb.LotDb) []lot.Lot {
	mappedLots := make([]lot.Lot, len(dbLots))
	for i, dbLot := range dbLots {
		mappedLot := lot.Lot{
			Id:           dbLot.Id,
			CreatedAt:    dbLot.CreatedAt,
			UpdatedAt:    dbLot.UpdatedAt,
			Quantity:     dbLot.Quantity,
			PricePerUnit: dbLot.PricePerUnit,
			Currency:     dbLot.Currency,
			AccountId:    dbLot.AccountId,
			Figi:         dbLot.Figi,
			IsClosed:     dbLot.IsClosed,
		}
		mappedLots[i] = mappedLot
	}
	return mappedLots
}

func mapLotToDbLot(lots []lot.Lot) []portfoliodb.LotDb {
	dbLots := make([]portfoliodb.LotDb, len(lots))
	for i, lot := range lots {
		dbLot := portfoliodb.LotDb{
			Id:           lot.Id,
			CreatedAt:    lot.CreatedAt,
			UpdatedAt:    lot.UpdatedAt,
			Quantity:     lot.Quantity,
			PricePerUnit: lot.PricePerUnit,
			Currency:     lot.Currency,
			AccountId:    lot.AccountId,
			Figi:         lot.Figi,
			IsClosed:     lot.IsClosed,
		}
		dbLots[i] = dbLot
	}
	return dbLots
}
