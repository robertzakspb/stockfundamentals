package accountreturnapi

import (
	"net/http"

	bondportfolioanalysis "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio-analysis"
	accountmvservice "github.com/compoundinvest/stockfundamentals/internal/application/account/market-value"
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAccountReturn(c *gin.Context) {
	parsedFilters, err := ydbfilter.MapQueryFiltersToYdb[accountmvdomain.Return](c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}
	accountReturn, err := accountmvservice.GetAccountReturn(parsedFilters, "RUB")
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	dto := mapDomainToDto(accountReturn)

	c.JSON(http.StatusOK, dto)
}

func GetPortfolioOverview(c *gin.Context) {
	parsedFilters, err := ydbfilter.MapQueryFiltersToYdb[accountmvdomain.Return](c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}
	portfolioOverview, err := bondportfolioanalysis.GeneratePortfolioOverview(parsedFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, shared.StringResponse{Message: portfolioOverview})
}

func StartMarketValueSnapshotJob(c *gin.Context) {
	go accountmvservice.SaveAccountMarketValueSnapshots()

	c.JSON(http.StatusOK, shared.StringResponse{Message: "The account market value snapshot job has been successfully started"})
}

func GetAccountMarketValueInRUB(c *gin.Context) {
	accountId, err := shared.GetFromQueryParams("accountId", c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}, Message: "Failed to fetch the account market value because accountId was not found in the query parameters"})
		return
	}
	accountUuid, err := uuid.Parse(accountId)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}, Message: "Provided accountId is not a valid UUID: " + accountId})
		return
	}

	accountMarketValue, err := accountmvservice.GetCurrentAccountMarketValue(accountUuid, "RUB")
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorResponse{Errors: []string{err.Error()}, Message: "Failed to fetch the account market value"})
		return
	}
	dto := mapAccountMarketValueToDto(accountMarketValue)

	c.JSON(http.StatusOK, dto)
}
