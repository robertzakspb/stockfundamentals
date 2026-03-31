package divcalapi

import (
	"net/http"
	"strings"

	divcalendarservice "github.com/compoundinvest/stockfundamentals/internal/application/account/dividend-calendar"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAccountDividendCalendar(c *gin.Context) {
	query := c.Request.URL.Query()["accountIds"]
	accountIds := strings.Split(query[0], ",")

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
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, divCalendar)
}
