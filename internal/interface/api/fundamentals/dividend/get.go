package apidividend

import (
	"net/http"

	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	"github.com/gin-gonic/gin"
)

func GetDividendForecasts(c *gin.Context) {
	forecasts, err := appdividend.GetDividendForecasts()

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} else {
		dtos := mapDividendForecastDomainToDto(forecasts)
		c.JSON(http.StatusOK, dtos)
	}
}

func GetDividendForecastsGroupedBySecurity(c *gin.Context) {
	forecasts, err := appdividend.GetDivForecastsGroupedBySecurity()

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} else {
		dtos := mapSecurityDivForecastToDto(forecasts)
		c.JSON(http.StatusOK, dtos)
	}
}
