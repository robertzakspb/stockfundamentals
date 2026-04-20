package accountreturnapi

import (
	"net/http"

	accountmvservice "github.com/compoundinvest/stockfundamentals/internal/application/account/market-value"
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/gin-gonic/gin"
)

func GetAccountReturn(c *gin.Context) {
	parsedFilters := ydbfilter.MapQueryFiltersToYdb(c.Request.URL.Query(), accountmvdomain.Return{})
	accountReturn, err := accountmvservice.GetAccountReturn(parsedFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	dto := mapDomainToDto(accountReturn)

	c.JSON(http.StatusOK, dto)
}

func GetPortfolioOverview(c *gin.Context) {
	// parsedFilters := ydbfilter.MapQueryFiltersToYdb(c.Request.URL.Query(), accountmvdomain.Return{})
	// accountReturn, err := accountmvservice.GetAccountReturn(parsedFilters)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, err.Error())
	// 	return
	// }
}

func StartMarketValueSnapshotJob(c *gin.Context) {
	go accountmvservice.SaveAccountMarketValueSnapshots()

	c.JSON(http.StatusOK, "The account market value snapshot job has been successfully started")
}
