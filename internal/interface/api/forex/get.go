package forexapi

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func GetFilteredExchangeRates(c *gin.Context) {
	filters := ydbfilter.MapQueryFiltersToYdb(c.Request.URL.Query(), forexservice.ForexRate{})

	rates, err := forexservice.GetFilteredExchangeRates(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.StringResponse{Message: "Failed to fetch forex rates due to an error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}
