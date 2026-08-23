package transactionsapi

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/account/transactionprocessor"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func GetFilteredTransactions(c *gin.Context) {
	parsedQuery, err := shared.ParseFiltrationPaginationAndSorting[TransactionDto](c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	transactions, err := transactionprocessor.GetFilteredTransactions(parsedQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	dtos := mapTransactionsToDto(transactions)

	c.JSON(http.StatusOK, dtos)
}
