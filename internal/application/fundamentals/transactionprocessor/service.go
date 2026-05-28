package transactionprocessor

import (
	"errors"

	accountservice "github.com/compoundinvest/stockfundamentals/internal/application/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/transaction"
)

func ProcessOrderExecutions(transactions []transaction.Transaction) error {
	if len(transactions) == 0 {
		return errors.New("Provided zero transactions")
	}
	for i := range transactions {
		if transactions[i].Type != transaction.OrderExecution {
			return errors.New("Encountered a transaction of type " + string(transactions[i].Type) + " while processing order executions")
		}
	}

	accountIds := ExtractAccountsFrom(transactions)
	_, err := accountservice.GetAccountsById(accountIds)
	if err != nil {
		return err
	}

	return nil
}

func AdjustPositionLots() {

}
