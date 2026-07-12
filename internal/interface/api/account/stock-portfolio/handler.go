package portfolioapi

import (
	"net/http"
	"sync"

	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func GetAccountPortfolio(c *gin.Context) {
	filters, err := ydbfilter.MapQueryFiltersToYdb[LotDto](c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	accountPortfolio, err := portfolio.GetAccountPortfolio(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	accountPortfolio.Lots, err = portfolio.PopulateLotSecurities(accountPortfolio.Lots)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	accountPortfolio, err = portfolio.PopulateLotsWithQuotes(accountPortfolio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	dto := mapPortfolioToDto(accountPortfolio)
	c.JSON(http.StatusOK, dto)
}

func UpdatePortfolio(c *gin.Context) {
	var wg sync.WaitGroup
	var err error
	wg.Go(func() { err = portfolio.UpdatePortfolio() })
	wg.Go(func() { err = bondportfolio.ImportTinkoffBondLots() })

	wg.Wait()

	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, shared.StringResponse{Message: "The portfolio has been successfully updated"})
}
