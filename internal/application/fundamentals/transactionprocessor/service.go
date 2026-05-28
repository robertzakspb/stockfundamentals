package transactionprocessor

import (
	"errors"
	"sort"

	accountservice "github.com/compoundinvest/stockfundamentals/internal/application/account/account"
	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
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

	accounts, err := accountservice.GetAccountsById(ExtractAccountsFrom(transactions))
	if err != nil {
		return err
	}
	groupedTransactions := GroupByAccount(transactions)

	for accountId, accountTransactions := range groupedTransactions {
		account, err := accountservice.FindAccountById(accountId, accounts)
		if err != nil {
			return errors.New("Failed to find account " + accountId.String() + " in the list, abandoning the order execution processing")
		}
	}

	return nil
}

func AdjustStockLotsAndCashBalances(transactions map[account.Account][]transaction.Transaction, lots []lot.Lot) {
	sort.Slice()

	//TODO: Adjust the account's cash balance
}

func AdjustBondLotsAndCashBalances(transactions map[account.Account][]transaction.Transaction, lots []bonds.BondLot)
