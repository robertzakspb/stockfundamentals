package main

import (
	"net/http"

	api_security "github.com/compoundinvest/stockfundamentals/internal/interface/api/security"
	timeseries "github.com/compoundinvest/stockfundamentals/internal/interface/api/time-series"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/dataseed"

	"github.com/compoundinvest/stockfundamentals/internal/interface/api/account/portfolio"
	dividend "github.com/compoundinvest/stockfundamentals/internal/interface/api/fundamentals/dividend"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//TODO: Add the middleware code here
	addEndpoints(router)

	router.Run("localhost:8080")
}

func addEndpoints(router *gin.Engine) {
	router.POST("/migration/initial-seed", dataseed.InitialSeed)

	router.GET("/health-check", healthCheck)

	router.GET("/portfolio", portfolio.GetPortfolio)

	router.POST("/fetch/dividends", dividend.StartDividendFetchingJob)
	router.POST("/fetch/securities", api_security.StartSecurityMasterImportJob)
	router.POST("/fetch/time-series", timeseries.StartTimeSeriesImportJob)
}

func healthCheck(c *gin.Context) {
	//TODO: Extract into a middleware method
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	c.JSON(http.StatusOK, "The service is running without any issues")
}
