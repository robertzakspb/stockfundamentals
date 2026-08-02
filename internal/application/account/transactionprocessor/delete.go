package transactionprocessor

import (
	"github.com/compoundinvest/stockfundamentals/internal/application/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/transactionsdb"
	"github.com/google/uuid"
)

func DeleteTbankTransactions() error {
	tBankAccountIds := []uuid.UUID{uuid.MustParse(shared.TINKOFF_IIS_ACCOUNT_ID)} //Quick and dirty solution

	err := transactionsdb.DeleteTbankTransactions(tBankAccountIds)

	return err
}
