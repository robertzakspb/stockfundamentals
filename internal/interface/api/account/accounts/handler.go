package accountsapi

import (
	"encoding/json"
	"io"
	"net/http"

	accountservice "github.com/compoundinvest/stockfundamentals/internal/application/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
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
		logger.Log("Failed to read the account json from the POST payload: "+err.Error(), logger.ERROR)
		return
	}

	dto := NewAccountDto{}
	err = json.Unmarshal(jsonData, &dto)
	if err != nil {
		logger.Log("Failed to unmarshal the account json in the POST payload: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	mappedAccount, err := mapDtoToAccount(dto)
	if err != nil {
		logger.Log("Failed to map the provided account: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	err = accountservice.SaveAccounts([]account.Account{mappedAccount})
	if err != nil {
		logger.Log("Failed to create the account. Reason: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, shared.StringResponse{Message: "The account has been successfuly saved!"})
}
