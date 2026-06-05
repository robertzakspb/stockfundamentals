package tranlotrelationservice

import (
	tranlotrelation "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/tran-lot-relation"
	tranlotrelationdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/tran-lot-relation-db"
)

func mapTranLotRelationsToDbModel(relations []tranlotrelation.TransactionLotRelation) []tranlotrelationdb.TransactionLotRelationDb {
	dbModels := make([]tranlotrelationdb.TransactionLotRelationDb, len(relations))

	for i := range relations {
		dbModel := tranlotrelationdb.TransactionLotRelationDb{
			Id:            relations[i].Id,
			TransactionId: relations[i].TransactionId,
			StockLotId:    relations[i].StockLotId,
			BondLotId:     relations[i].BondLotId,
			Date:          relations[i].Date,
			Quantity:      relations[i].Quantity,
		}
		dbModels[i] = dbModel
	}

	return dbModels
}
