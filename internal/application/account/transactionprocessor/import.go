package transactionprocessor

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"opensource.tbank.ru/invest/invest-go/investgo"
	investapi "opensource.tbank.ru/invest/invest-go/proto"
)

func ImportTBankTransactions() error {
	config, err := investgo.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	client, err := investgo.NewClient(ctx, config, nil)
	if err != nil {
		logger.Log("Failed to initialize the investgo API client: ", logger.ALERT)
		return err
	}
	operationService := client.NewOperationsServiceClient()

	accountService := client.NewUsersServiceClient()
	status := investapi.AccountStatus_ACCOUNT_STATUS_ALL
	accsResp, err := accountService.GetAccounts(&status)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}

	accounts := accsResp.GetAccounts()
	if len(accounts) == 0 {
		logger.Log("No accounts found in T Bank API", logger.ALERT)
		return err
	}

	transactions := []*investapi.OperationItem{}

	for i := range accounts {
		if accounts[i].Id != "2236996963" {
			continue //We only need to import transactions for this one account
		}
		cursor := ""
		for {
			request := &investgo.GetOperationsByCursorRequest{
				AccountId:          accounts[i].Id,
				Cursor:             cursor,
				WithoutCommissions: false,
				WithoutTrades:      false,
				WithoutOvernights:  false,
			}
			response, err := operationService.GetOperationsByCursor(request)
			if response == nil {
				logger.Log("Unexpectedly received a nil response from investgo API", logger.ALERT)
			}
			if err != nil {
				logger.LogError(err, logger.ERROR)
				return err
			}

			for _, transaction := range response.GetItems() {
				transactions = append(transactions, transaction)
			}
			if !response.HasNext {
				break
			}
			cursor = response.NextCursor
		}
	}

	mappedTransactions := mapTBankTransactionsToTransaction(transactions)

	err = SaveTransactions(mappedTransactions)
	if err != nil {
		logger.LogError(err, logger.ERROR)
	}
	logger.Log("The T Bank Transactions import job has been successfully executed", logger.ERROR)

	return nil
}
