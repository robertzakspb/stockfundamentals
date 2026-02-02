package timeseries

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/timeseries"
	"github.com/gin-gonic/gin"
)

func StartTimeSeriesImportJob(c *gin.Context) {
	go timeseries.FetchAndSaveHistoricalQuotes()
	c.JSON(http.StatusOK, "Successfully started the time series import job")
}