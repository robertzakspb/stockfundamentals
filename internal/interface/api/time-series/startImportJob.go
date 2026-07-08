package timeseries

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/quoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/timeseries"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func StartTimeSeriesImportJob(c *gin.Context) {
	go timeseries.FetchAndSaveHistoricalQuotes()
	c.JSON(http.StatusOK, shared.StringResponse{Message: "Successfully started the time series import job"})
}

func StartQuoteSnapShotJob(c *gin.Context) {
	go quoteservice.CreateQuoteSnapshot()
	c.JSON(http.StatusOK, shared.StringResponse{Message: "Successfully started the quote snapshot job"})
}
