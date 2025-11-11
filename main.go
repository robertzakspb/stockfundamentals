package main

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/portfolio"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/dataseed"
	dividend "github.com/compoundinvest/stockfundamentals/internal/interface/api/fundamentals/dividend"
	"github.com/gin-gonic/gin"
)

func main() {
	dataseed.InitialSeed()
	// fetchExternalData()
	portfolio.GeMyPortfolio().PrintAllPositions()
	router := gin.Default()
	router.GET("/health-check", healthCheck)
	router.Run("localhost:8080")
}

func fetchExternalData() {
	security_master.FetchAndSaveSecurities()
	dividend.FetchAndSaveAllDividends()
	// timeseries.FetchAndSaveHistoricalQuotes()
}

func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "The service is running without any issues")
}
