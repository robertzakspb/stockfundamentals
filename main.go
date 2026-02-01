package main

import (
	"net/http"

	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"

	"github.com/compoundinvest/stockfundamentals/internal/application/timeseries"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/dataseed"

	"github.com/compoundinvest/stockfundamentals/internal/interface/api/account/portfolio"
	dividend "github.com/compoundinvest/stockfundamentals/internal/interface/api/fundamentals/dividend"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	addEndpoints(router)

	router.Run("localhost:8080")
}

func addEndpoints(router *gin.Engine) {
	router.GET("/health-check", healthCheck)
	router.GET("/portfolio", portfolio.GetPortfolio)
	router.POST("/migration/initial-seed", dataseed.InitialSeed)
	router.POST("/fetch/dividends", dividend.FetchAndSaveAllDividends)
	router.POST("/fetch/securities", security_master.FetchAndSaveSecurities)   //TODO: User the handler from the api layer
	router.POST("/fetch/time-series", timeseries.FetchAndSaveHistoricalQuotes) //TODO: User the handler from the api layer
}

func healthCheck(c *gin.Context) {
	//TODO: Extract into a middleware method
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	c.JSON(http.StatusOK, "The service is running without any issues")
}
