package accountsapi

import (
	"net/http"

	accountservice "github.com/compoundinvest/stockfundamentals/internal/application/account/account"
	"github.com/gin-gonic/gin"
)

func GetAllAccounts(c *gin.Context) {
	accounts, err := accountservice.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	dtos := mapAccountsToDtos(accounts)

	c.JSON(http.StatusOK, dtos)
}
