package tranlotrelationservice

import (
	tranlotrelation "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/tran-lot-relation"
	tranlotrelationdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/tran-lot-relation-db"
)

func SaveTranLotRelations(relations []tranlotrelation.TransactionLotRelation) error {
	dbModels := mapTranLotRelationsToDbModel(relations)

	err := tranlotrelationdb.SaveTranLotRelations(dbModels)

	return err
}
