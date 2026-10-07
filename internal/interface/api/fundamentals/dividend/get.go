package apidividend

import (
	"net/http"

	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	divcalapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/dividend-calendar"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetDividendForecasts(c *gin.Context) {
	forecasts, err := appdividend.GetDividendForecasts()

	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	dtos := mapDividendForecastDomainToDto(forecasts)

	c.JSON(http.StatusOK, dtos)
}

func GetDividendForecastsGroupedBySecurity(c *gin.Context) {
	forecasts, err := appdividend.GetDivForecastsGroupedBySecurity()

	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	dtos := mapSecurityDivForecastToDto(forecasts)

	c.JSON(http.StatusOK, dtos)
}

func GetDividendForecastsForAccount(c *gin.Context) {
	accountPayouts, err := appdividend.GetDividendForecastsForAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	dtos := []divcalapi.PayoutDto{}

	for i := range accountPayouts {
		dto := divcalapi.MapPayoutToDto(dividend.Payout(accountPayouts[i]))
		dtos = append(dtos, dto)
	}

	c.JSON(http.StatusOK, dtos)
}

func GetFutureDividendPayoutsForAccount(c *gin.Context) {
	accountId, err := shared.GetFromQueryParams("accountId", c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	accountUuid, err := uuid.Parse(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	payouts, err := appdividend.GetDividendPayoutsForAccount(accountUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	dtos := []divcalapi.PayoutDto{}
	for i := range payouts {
		dto := divcalapi.MapPayoutToDto(dividend.Payout(payouts[i]))
		dtos = append(dtos, dto)
	}

	c.JSON(http.StatusOK, dtos)
}
