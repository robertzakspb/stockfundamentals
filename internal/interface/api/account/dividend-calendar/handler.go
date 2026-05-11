package divcalapi

import (
	"net/http"
	"strings"

	divcalendarservice "github.com/compoundinvest/stockfundamentals/internal/application/account/dividend-calendar"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAccountDividendCalendar(c *gin.Context) {
	query, err := shared.GetFromQueryParams("accountIds", c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusNotAcceptable, shared.ResponseError{Errors: []string{err.Error()}})
		return
	}
	accountIds := strings.Split(query, ",")

	uuids := uuid.UUIDs{}
	for _, accountId := range accountIds[1:] {
		uuid, err := uuid.Parse(accountId)
		if err != nil {
			logger.Log("Failed to parse a UUID from "+accountId, logger.ERROR)
			continue
		}
		uuids = append(uuids, uuid)
	}

	divCalendar, err := divcalendarservice.GetAccountDividendCalendar(uuids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ResponseError{Errors: []string{err.Error()}})
		return
	}

	dto := mapDivCalToDto(DividendCalendar(divCalendar))

	c.JSON(http.StatusOK, dto)
}
