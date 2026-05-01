package portfolio

import (
	"net/http"
	"sync"

	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/gin-gonic/gin"
)

func GetAccountPortfolio(c *gin.Context) {
	filters := ydbfilter.MapQueryFiltersToYdb(c.Request.URL.Query(), lot.Lot{})

	accountPortfolio, err := portfolio.GetAccountPortfolio(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	accountPortfolio.Lots, err = portfolio.PopulateLotSecurities(accountPortfolio.Lots)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	accountPortfolio, err = portfolio.PopulateLotsWithQuotes(accountPortfolio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
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
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, "The portfolio has been successfully updated")
	}
}
