package transactionsapi

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/account/transactionprocessor"
	"github.com/gin-gonic/gin"
)

func GetAllTransactions(c *gin.Context) {
	transactions, err := transactionprocessor.GetAllTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, transactions)
}
