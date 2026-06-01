package transactionprocessor

import (
	"errors"
	"sort"
	"strconv"

	accountservice "github.com/compoundinvest/stockfundamentals/internal/application/account/account"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/transaction"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
)

type Transaction transaction.Transaction
type Account account.Account

func ProcessStockOrderExecutions(transactions []Transaction) error {
	if len(transactions) == 0 {
		return errors.New("Provided zero transactions")
	}
	for i := range transactions {
		if transactions[i].Type != transaction.OrderExecution {
			return errors.New("Encountered a transaction of type " + string(transactions[i].Type) + " while processing order executions. Aborted.")
		}
	}

	err := adjustStockLotsAndCashBalances(transactions)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	return nil
}

func adjustStockLotsAndCashBalances(transactions []Transaction) error {
	//Grouping transactions by account, as they are applied to each account separately
	groupedTransactions := GroupByAccount(transactions)

	//Fetching accounts that contains cash balances
	accoundIds := ExtractAccountsFrom(transactions)
	accounts, err := accountservice.GetAccountsById(accoundIds)

	if err != nil {
		return err
	}
	if len(groupedTransactions) != len(accounts) {
		return errors.New("The account count in grouped transactions is " + strconv.Itoa(len(groupedTransactions)) + " whilte the DB account count is " + strconv.Itoa(len(accounts)))
	}

	//Fetching the current stock portfolios to adjust them
	accountfilter := ydbfilter.YdbFilter{
		YqlColumnName:  "account_id",
		Condition:      ydbfilter.Contains,
		ConditionValue: ydbhelper.ConvertUUIDsToYdbList(accoundIds),
	}
	//FIXME: We need to add the closed flag and remove closed lots from the selection
	lots, err := portfolio.GetFilteredLots([]ydbfilter.YdbFilter{accountfilter})
	if err != nil {
		return err
	}
	groupedLots := portfolio.GroupLotsByAccount(lots)
	if len(groupedTransactions) != len(groupedLots) {
		return errors.New("The account count in grouped transactions is " + strconv.Itoa(len(groupedTransactions)) + " whilte the portfolio count is " + strconv.Itoa(len(groupedLots)))
	}

	err = adjustAccountStockLotsAndBalances(accounts, groupedTransactions, groupedLots)
	if err != nil {
		return err
	}

	return nil
}

// Recalculates and saves the adjusted stock lots and balances after the transactions have been applied
func adjustAccountStockLotsAndBalances(accounts []account.Account, transactions map[uuid.UUID][]Transaction, lots map[uuid.UUID][]lot.Lot) error {
	adjustedLots := []lot.Lot{}
	for accountId, accountTransactions := range transactions {
		account, err := accountservice.FindAccountById(accountId, accounts)
		if err != nil {
			return errors.New("Failed to find account " + accountId.String() + " in the list, abandoning the order execution processing")
		}
		lots, found := lots[accountId]
		if !found {
			return errors.New("Failed to find lots for account " + accountId.String() + " in grouped lots")
		}

		updatedAccount, newLots, err := recalculateLotsAndCashBalances(account, accountTransactions, lots)
		if err != nil {
			return err
		}

		adjustedLots = append(adjustedLots, newLots...)
		foundAccount := false
		for i := range accounts {
			if accounts[i].Id == updatedAccount.Id {
				//FIXME: The cash updating logic is to be here
				foundAccount = true
			}
		}
		if !foundAccount {
			return errors.New("Failed to find the account " + accountId.String() + " while updating the cash position")
		}
	}

	flattenedTransactions := []Transaction{}
	for _, t := range transactions {
		flattenedTransactions = append(flattenedTransactions, t...)
	}

	err := saveLotsAndAccountsAndTransactions(accounts, flattenedTransactions, adjustedLots)
	if err != nil {
		return err
	}

	return nil
}

func recalculateLotsAndCashBalances(account account.Account, transactions []Transaction, lots []lot.Lot) (account.Account, []lot.Lot, error) {
	//Transactions must be processed in chronological order, reenacting the user's behavior in the OMS
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Timestamp.After(transactions[j].Timestamp)
	})
	//Lots must be sorted by creation date because sales are applied according to the FIFO principle
	sort.Slice(lots, func(i, j int) bool {
		return lots[i].CreatedAt.After(lots[j].CreatedAt)
	})





	return account, lots, nil
}

func saveLotsAndAccountsAndTransactions(accounts []account.Account, transactions []Transaction, lots []lot.Lot) error {
	//TODO: Save the updated lots
	//TODO: Save the updated accounts (cash balances)
	//TODO: Save the transactions
	//TODO: Establish a relationship between lots and transactions
	return nil
}
