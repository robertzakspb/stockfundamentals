package main

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/portfolio"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"

	// "github.com/compoundinvest/stockfundamentals/internal/application/timeseries"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/dataseed"

	// dividend "github.com/compoundinvest/stockfundamentals/internal/interface/api/fundamentals/dividend"
	"github.com/gin-gonic/gin"
)

func main() {
	dataseed.InitialSeed()

	router := gin.Default()
	router.GET("/health-check", healthCheck)
	router.GET("/portfolio", getUserPortfolio)
	router.Run("localhost:8080")
	fetchExternalData()
}

func fetchExternalData() {
	security_master.FetchAndSaveSecurities()
	// dividend.FetchAndSaveAllDividends()
	// timeseries.FetchAndSaveHistoricalQuotes()
}

func healthCheck(c *gin.Context) {
	//TODO: Extract into a common method
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	c.JSON(http.StatusOK, "The service is running without any issues")
}

func getUserPortfolio(c *gin.Context) {
	userPortfolio := portfolio.GeMyPortfolio().WithPLs()
	c.JSON(http.StatusOK, userPortfolio)
}
