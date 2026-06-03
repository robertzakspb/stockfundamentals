package transactionprocessor

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/transactionsdb"
)

func mapTransactionsToDbModel(transactions []transaction.Transaction) []transactionsdb.TransactionDbModel {
	dbModels := make([]transactionsdb.TransactionDbModel, len(transactions))

	for i := range transactions {
		dbModel := transactionsdb.TransactionDbModel{
			Id:           transactions[i].Id,
			AccountId:    transactions[i].AccountId,
			Figi:         transactions[i].Figi,
			Type:         string(transactions[i].Type),
			Timestamp:    transactions[i].Timestamp,
			Side:         string(transactions[i].Side),
			Quantity:     transactions[i].Quantity,
			PricePerUnit: transactions[i].PricePerUnit,
			Currency:     transactions[i].Currency,
			Description:  transactions[i].Description,
		}
		dbModels[i] = dbModel
	}

	return dbModels
}
