package tranlotrelationservice

import (
	tranlotrelation "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/tran-lot-relation"
	tranlotrelationdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/tran-lot-relation-db"
	"github.com/google/uuid"
)

func mapTranLotRelationsToDbModel(relations []tranlotrelation.TransactionLotRelation) []tranlotrelationdb.TransactionLotRelationDb {
	dbModels := make([]tranlotrelationdb.TransactionLotRelationDb, len(relations))

	for i := range relations {
		dbModel := tranlotrelationdb.TransactionLotRelationDb{
			Id:         uuid.New(),
			StockLotId: relations[i].StockLotId,
			BondLotId:  relations[i].BondLotId,
			Date:       relations[i].Date,
			Quantity:   relations[i].Quantity,
		}
		dbModels[i] = dbModel
	}

	return dbModels
}
