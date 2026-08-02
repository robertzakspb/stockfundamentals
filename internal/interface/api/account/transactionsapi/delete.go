package transactionsapi

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/account/transactionprocessor"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func DeleteTbankTransactions(c *gin.Context) {
	err := transactionprocessor.DeleteTbankTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, shared.StringResponse{Message: "The T Bank transactions have been successfully deleted"})
}
