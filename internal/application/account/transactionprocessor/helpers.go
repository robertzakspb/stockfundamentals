package transactionprocessor

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
	"github.com/google/uuid"
)

func ExtractAccountsFrom(transactions []transaction.Transaction) (accountIds uuid.UUIDs) {
	for i := range transactions {
		accountIds = append(accountIds, transactions[i].AccountId)
	}
	return accountIds
}

func GroupByAccount(transactions []transaction.Transaction) map[uuid.UUID][]transaction.Transaction {
	var grouped = map[uuid.UUID][]transaction.Transaction{}

	for _, t := range transactions {
		grouped[t.AccountId] = append(grouped[t.AccountId], t)
	}

	return grouped
}
