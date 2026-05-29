package transactionprocessor

import (
	"github.com/google/uuid"
)

func ExtractAccountsFrom(transactions []Transaction) (accountIds uuid.UUIDs) {
	for i := range transactions {
		accountIds = append(accountIds, transactions[i].AccountId)
	}
	return accountIds
}

func GroupByAccount(transactions []Transaction) map[uuid.UUID][]Transaction {
	var grouped = map[uuid.UUID][]Transaction{}

	for _, t := range transactions {
		grouped[t.AccountId] = append(grouped[t.AccountId], t)
	}

	return grouped
}
