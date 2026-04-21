package portfolio

import (
	"net/http"
	"sync"

	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPortfolio(c *gin.Context) {
	userPortfolio, err := portfolio.GeStockPortfolio()
	if err != nil {
		c.JSON(http.StatusNoContent, "Failed to fetch user positions")
	} else {
		portfolioWithPLs := userPortfolio.WithPLs()
		c.JSON(http.StatusOK, portfolioWithPLs)
	}
}

func GetAccountPortfolio(c *gin.Context) {
	accountIDs := uuid.UUIDs{} //FIXME: Parse from the query parameters
	portfolio, err := portfolio.GetAccountPortfolio(accountIDs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, portfolio)
	}
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
