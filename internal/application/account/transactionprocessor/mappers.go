package transactionprocessor

import (
	"errors"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
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

func mapTransactionDbModelsToTransactions(dbModels []transactionsdb.TransactionDbModel) ([]transaction.Transaction, error) {
	transactions := make([]transaction.Transaction, len(dbModels))

	for i := range dbModels {
		tType, found := transaction.TypeLookup[dbModels[i].Type]
		if !found {
			return transactions, errors.New("Unable to map the provided value – " + dbModels[i].Type + "- to a transaction type")
		}
		side, found := transaction.OrderSideLookup[dbModels[i].Side]
		if !found {
			return transactions, errors.New("Unable to map the provided value – " + dbModels[i].Type + "- to a transaction side")
		}

		transaction := transaction.Transaction{
			Id:           dbModels[i].Id,
			AccountId:    dbModels[i].AccountId,
			Figi:         dbModels[i].Figi,
			Type:         tType,
			Timestamp:    dbModels[i].Timestamp,
			Side:         side,
			Quantity:     dbModels[i].Quantity,
			PricePerUnit: dbModels[i].PricePerUnit,
			Currency:     dbModels[i].Currency,
			Description:  dbModels[i].Description,
		}
		transactions[i] = transaction
	}

	return transactions, nil
}
