package timeseries

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/timeseries"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func StartTimeSeriesImportJob(c *gin.Context) {
	go timeseries.FetchAndSaveHistoricalQuotes()
	c.JSON(http.StatusOK, shared.StringResponse{Message: "Successfully started the time series import job"})
}
