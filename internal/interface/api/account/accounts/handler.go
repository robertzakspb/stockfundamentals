package accountsapi

import (
	"encoding/json"
	"io"
	"net/http"

	accountservice "github.com/compoundinvest/stockfundamentals/internal/application/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func GetAllAccounts(c *gin.Context) {
	accounts, err := accountservice.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	dtos := MapAccountsToDtos(accounts)

	c.JSON(http.StatusOK, dtos)
}

func CreateAccount(c *gin.Context) {
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
}
