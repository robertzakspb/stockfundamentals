package transactionprocessor

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/transactionsdb"
)

func SaveTransactions(transactions []transaction.Transaction) error {
	mappedDbModels := mapTransactionsToDbModel(transactions)

	err := transactionsdb.SaveTransactions(mappedDbModels)

	return err
}
