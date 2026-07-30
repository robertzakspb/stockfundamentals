package transactionsapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/account/transactionprocessor"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	accountsapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/accounts"
	portfolioapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/stock-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func ProcessOrderExecutions(c *gin.Context) {
	bodyReader := c.Request.Body
	defer bodyReader.Close()

	jsonData, err := io.ReadAll(bodyReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
		logger.Log("Failed to read the transaction json from the POST payload: "+err.Error(), logger.ERROR)
		return
	}

	dtos := []TransactionDto{}
	err = json.Unmarshal(jsonData, &dtos)
	if err != nil {
		logger.Log("Failed to unmarshal the transactions json in the POST payload: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	transactions, err := mapTransactionDtosToTransactions(dtos)
	if err != nil {
		logger.Log("Invalid transaction data was provided: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
	}

	err = transactionprocessor.ProcessStockOrderExecutions(transactions)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, shared.StringResponse{Message: "The transactions have been successfully processed"})
}

func PreviewTransactions(c *gin.Context) {
	bodyReader := c.Request.Body
	defer bodyReader.Close()

	jsonData, err := io.ReadAll(bodyReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
		logger.Log("Failed to read the transaction json from the POST payload: "+err.Error(), logger.ERROR)
		return
	}

	dtos := []TransactionDto{}
	err = json.Unmarshal(jsonData, &dtos)
	if err != nil {
		logger.Log("Failed to unmarshal the transactions json in the POST payload: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	transactions, err := mapTransactionDtosToTransactions(dtos)
	if err != nil {
		logger.Log("Invalid transaction data was provided: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
	}

	accounts, transactions, lot, relations, err := transactionprocessor.PreviewTransactionProcessing(transactions)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	accountDtos := accountsapi.MapAccountsToDtos(accounts)
	transactionDtos := mapTransactionsToDto(transactions)
	lotDtos := portfolioapi.MapLotsToDtos(lot)
	relationDtos := mapTranLotRelationsToDtos(relations)

	dto := AccountsTransactionsLotsRelationsDto{
		AdjustedAccounts: accountDtos,
		AdjustedLots:     lotDtos,
		Relations:        relationDtos,
		Transactions:     transactionDtos,
	}

	c.JSON(http.StatusOK, dto)
}

func StartTBankTransactionsImportJob(c *gin.Context) {
	go transactionprocessor.ImportTBankTransactions()

	c.JSON(http.StatusOK, shared.StringResponse{Message: "The T Bank transactions import job has been started"})
}
