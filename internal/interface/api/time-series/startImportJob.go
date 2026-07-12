package timeseries

import (
	"net/http"

	jobservice "github.com/compoundinvest/stockfundamentals/internal/application/jobs"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/timeseries"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func StartTimeSeriesImportJob(c *gin.Context) {
	go timeseries.FetchAndSaveHistoricalQuotes()
	c.JSON(http.StatusOK, shared.StringResponse{Message: "Successfully started the time series import job"})
}

func StartQuoteSnapShotJob(c *gin.Context) {
	go jobservice.ExecuteQuoteSnapshotJob()
	c.JSON(http.StatusOK, shared.StringResponse{Message: "Successfully started the quote snapshot job"})
}
