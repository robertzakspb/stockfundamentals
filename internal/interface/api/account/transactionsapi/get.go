package transactionsapi

import (
	"net/http"
	"slices"

	"github.com/compoundinvest/stockfundamentals/internal/application/account/transactionprocessor"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func GetFilteredTransactions(c *gin.Context) {
	parsedQuery, err := shared.ParseFiltrationPaginationAndSorting[TransactionDto](c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	transactions, err := transactionprocessor.GetFilteredTransactions(parsedQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	dtos := mapTransactionsToDto(transactions)

	slices.SortFunc(dtos, func(t1, t2 TransactionDto) int {
		if t1.Timestamp.Before(t2.Timestamp) {
			return 1
		} else {
			return -1
		}
	})

	c.JSON(http.StatusOK, dtos)
}
