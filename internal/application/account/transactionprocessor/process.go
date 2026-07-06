package transactionprocessor

import (
	"errors"

	"strconv"

	accountservice "github.com/compoundinvest/stockfundamentals/internal/application/account/account"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	tranlotrelationservice "github.com/compoundinvest/stockfundamentals/internal/application/account/tran-lot-relation-service"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	tranlotrelation "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/tran-lot-relation"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

type Account account.Account

// Recalculates and readjusts accounts, stock lots, and creates lot-transaction relationships
func ProcessStockOrderExecutions(transactions []transaction.Transaction) error {
	if err := validateTransactions(transactions); err != nil {
		return err
	}

	_, _, _, _, err := adjustStockLotsAndCashBalances(transactions, false)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	return nil
}

// Unlike ProcessStockOrderExecutions, the function does not save any entities and simply returns the recalculated positions, transactions, accounts, and relationships
func PreviewTransactionProcessing(transactions []transaction.Transaction) ([]account.Account, []transaction.Transaction, []lot.Lot, []tranlotrelation.TransactionLotRelation, error) {
	if err := validateTransactions(transactions); err != nil {
		return []account.Account{}, []transaction.Transaction{}, []lot.Lot{}, []tranlotrelation.TransactionLotRelation{}, err
	}

	adjustedAccounts, flattenedTransactions, adjustedLots, relations, err := adjustStockLotsAndCashBalances(transactions, true)
	if err != nil {
		return []account.Account{}, []transaction.Transaction{}, []lot.Lot{}, []tranlotrelation.TransactionLotRelation{}, err
	}

	return adjustedAccounts, flattenedTransactions, adjustedLots, relations, nil
}

func adjustStockLotsAndCashBalances(transactions []transaction.Transaction, isPreview bool) ([]account.Account, []transaction.Transaction, []lot.Lot, []tranlotrelation.TransactionLotRelation, error) {
	//Grouping transactions by account, as they are applied to each account separately
	groupedTransactions := GroupByAccount(transactions)

	//Fetching accounts that contains cash balances
	accoundIds := ExtractAccountsFrom(transactions)
	accounts, err := accountservice.GetAccountsById(accoundIds)

	if err != nil {
		return []account.Account{}, []transaction.Transaction{}, []lot.Lot{}, []tranlotrelation.TransactionLotRelation{}, err
	}
	if len(groupedTransactions) != len(accounts) {
		return []account.Account{}, []transaction.Transaction{}, []lot.Lot{}, []tranlotrelation.TransactionLotRelation{}, errors.New("The account count in grouped transactions is " + strconv.Itoa(len(groupedTransactions)) + " while the DB account count is " + strconv.Itoa(len(accounts)))
	}

	//Fetching the current stock portfolios to adjust them
	accountfilter := ydbfilter.YdbFilter{
		YqlColumnName:  "account_id",
		Condition:      ydbfilter.Contains,
		ConditionValue: ydbhelper.ConvertUUIDsToYdbList(accoundIds),
	}
	closedFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "is_closed",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.BoolValue(false),
	}
	lots, err := portfolio.GetFilteredLots([]ydbfilter.YdbFilter{accountfilter, closedFilter})
	if err != nil {
		return []account.Account{}, []transaction.Transaction{}, []lot.Lot{}, []tranlotrelation.TransactionLotRelation{}, err
	}
	groupedLots := portfolio.GroupLotsByAccount(lots)
	if len(groupedTransactions) != len(groupedLots) {
		err := errors.New("The account count in grouped transactions is " + strconv.Itoa(len(groupedTransactions)) + " whilte the portfolio count is " + strconv.Itoa(len(groupedLots)))
		return []account.Account{}, []transaction.Transaction{}, []lot.Lot{}, []tranlotrelation.TransactionLotRelation{}, err
	}

	adjustedAccounts, flattenedTransactions, adjustedLots, relations, err := adjustAccountStockLotsAndBalances(accounts, groupedTransactions, groupedLots)
	if err != nil {
		return []account.Account{}, []transaction.Transaction{}, []lot.Lot{}, []tranlotrelation.TransactionLotRelation{}, err
	}

	if isPreview {
		return adjustedAccounts, flattenedTransactions, adjustedLots, relations, err
	}

	err = saveAllEntities(adjustedAccounts, flattenedTransactions, adjustedLots, relations)
	if err != nil {
		return adjustedAccounts, []transaction.Transaction{}, adjustedLots, relations, err
	}

	return adjustedAccounts, flattenedTransactions, adjustedLots, relations, nil
}

// Recalculates the adjusted stock lots and balances after the transactions have been applied
func adjustAccountStockLotsAndBalances(accounts []account.Account, transactions map[uuid.UUID][]transaction.Transaction, lots map[uuid.UUID][]lot.Lot) ([]account.Account, []transaction.Transaction, []lot.Lot, []tranlotrelation.TransactionLotRelation, error) {
	adjustedLots := []lot.Lot{}
	adjustedAccounts := []account.Account{}
	relations := []tranlotrelation.TransactionLotRelation{}
	for accountId, accountTransactions := range transactions {
		account, err := accountservice.FindAccountById(accountId, accounts)
		if err != nil {
			return adjustedAccounts, []transaction.Transaction{}, adjustedLots, relations, errors.New("Failed to find account " + accountId.String() + " in the list, abandoning the order execution processing")
		}
		lots, found := lots[accountId]
		if !found {
			return adjustedAccounts, []transaction.Transaction{}, adjustedLots, relations, errors.New("Failed to find lots for account " + accountId.String() + " in grouped lots")
		}

		updatedAccount, newLots, newRelations, err := recalculateLotsAndCashBalances(account, accountTransactions, lots)
		if err != nil {
			return adjustedAccounts, []transaction.Transaction{}, adjustedLots, relations, err
		}

		adjustedAccounts = append(adjustedAccounts, updatedAccount)
		adjustedLots = append(adjustedLots, newLots...)
		relations = append(relations, newRelations...)
	}

	flattenedTransactions := []transaction.Transaction{}
	for _, t := range transactions {
		flattenedTransactions = append(flattenedTransactions, t...)
	}

	return adjustedAccounts, flattenedTransactions, adjustedLots, relations, nil
}

func saveAllEntities(accounts []account.Account, transactions []transaction.Transaction, lots []lot.Lot, relations []tranlotrelation.TransactionLotRelation) error {
	err := accountservice.SaveAccounts(accounts)
	if err != nil {
		return err
	}

	err = SaveTransactions(transactions)
	if err != nil {
		return err
	}

	err = portfolio.SaveLots(lots)
	if err != nil {
		return err
	}

	err = tranlotrelationservice.SaveTranLotRelations(relations)
	if err != nil {
		return err
	}

	return nil
}
