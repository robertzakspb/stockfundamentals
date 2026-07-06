package transactionprocessor

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/transactionsdb"
)

func GetAllTransactions() ([]transaction.Transaction, error) {
	dbModels, err := transactionsdb.GetAllTransactions()
	if err != nil {
		return []transaction.Transaction{}, err
	}

	mappedTransactions, err := mapTransactionDbModelsToTransactions(dbModels)
	if err != nil {
		return []transaction.Transaction{}, err
	}

	return mappedTransactions, nil
}
