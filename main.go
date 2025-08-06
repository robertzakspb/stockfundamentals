package main

import (

	"github.com/compoundinvest/stockfundamentals/dataseed"
	// "fmt"
	// "github.com/compoundinvest/stockfundamentals/features/fundamentaldata/dividend"
	"github.com/compoundinvest/stockfundamentals/features/fundamentaldata/security"
	"github.com/compoundinvest/stockfundamentals/internal/domain/portfolio"
)

func main() {
	dataseed.InitialSeed()
	fetchExternalData()
	portfolio.GeMyPortfolio().PrintAllPositions()
}

func fetchExternalData() {
	security.FetchAndSaveSecurities()
	// dividend.FetchAndSaveAllDividends()
	// timeseries.FetchAndSaveHistoricalQuotes()
}
