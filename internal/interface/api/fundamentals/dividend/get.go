package apidividend

import (
	"net/http"

	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
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
	accountIdString, _ := shared.GetFromQueryParams("accountId", c.Request.URL.Query())
	accountId, err := uuid.Parse(accountIdString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	accountPayouts, err := appdividend.GetDividendForecastsForAccount(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
		return
	}

	dtos := []divcalapi.PayoutDto{}

	for i := range accountPayouts {
		dto := divcalapi.MapPayoutToDto(divcalapi.Payout(accountPayouts[i]))
		dtos = append(dtos, dto)
	}

	calendar := divcalapi.DividendCalendarDto{
		AccountIds:    []uuid.UUID{accountId},
		FuturePayouts: dtos,
	}

	c.JSON(http.StatusOK, calendar)
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
		dto := divcalapi.MapPayoutToDto(divcalapi.Payout(payouts[i]))
		dtos = append(dtos, dto)
	}

	c.JSON(http.StatusOK, dtos)
}
