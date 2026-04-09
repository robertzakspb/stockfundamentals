package accountreturnapi

import (
	"net/http"

	accountmv "github.com/compoundinvest/stockfundamentals/internal/application/account/market-value"
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/dividend-calendar/market-value"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/gin-gonic/gin"
)

func GetAccountReturn(c *gin.Context) {
	parsedFilters := ydbfilter.MapQueryFiltersToYdb(c.Request.URL.Query(), accountmvdomain.Return{})
	accountReturn, err := accountmv.GetAccountReturn(parsedFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	dto := mapDomainToDto(accountReturn)

	c.JSON(http.StatusOK, dto)
}
