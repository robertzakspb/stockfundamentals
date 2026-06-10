package transactionsapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/account/transactionprocessor"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func ProcessOrderExecutions(c *gin.Context) {
	bodyReader := c.Request.Body
	defer bodyReader.Close()

	jsonData, err := io.ReadAll(bodyReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ResponseError{Errors: []string{err.Error()}})
		logger.Log("Failed to read the transaction json from the POST payload: "+err.Error(), logger.ERROR)
		return
	}

	dtos := []TransactionDto{}
	err = json.Unmarshal(jsonData, &dtos)
	if err != nil {
		logger.Log("Failed to unmarshal the transactions json in the POST payload: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, shared.ResponseError{Errors: []string{err.Error()}})
		return
	}

	transactions, err := mapTransactionDtosToTransactions(dtos)
	if err != nil {
		logger.Log("Invalid transaction data was provided: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, shared.ResponseError{Errors: []string{err.Error()}})
	}

	err = transactionprocessor.ProcessStockOrderExecutions(transactions)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ResponseError{Errors: []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, shared.StringResponse{Message: "The transactions have been successfully processed"})
}
