package apidividend

import (
	"net/http"

	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	"github.com/gin-gonic/gin"
)

func StartDividendFetchingJob(c *gin.Context) {
	go appdividend.FetchAndSaveAllDividends()

	c.JSON(http.StatusOK, "Successfully started the dividend fetching job")
}

func GetAllDividends(c *gin.Context) {
	dividends, err := dbdividend.GetAllDividends() //FIXME: Refactor the use the service
	dtos := convertDividendToDTO(dividends)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, dtos)
	}
}

func GetUpcomingDividends(c *gin.Context) {
	upcomingDividends, err := dbdividend.GetUpcomingDividends() //FIXME: Refactor the use the service
	dtos := convertDividendToDTO(upcomingDividends)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, dtos)
	}
}
