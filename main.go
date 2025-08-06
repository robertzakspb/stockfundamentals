package main

import (
	"github.com/compoundinvest/stockfundamentals/dataseed"
	"github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/application/portfolio"
)

func main() {
	dataseed.InitialSeed()
	fetchExternalData()
	portfolio.GeMyPortfolio().PrintAllPositions()
}

func fetchExternalData() {
	security_master.FetchAndSaveSecurities()
	// dividend.FetchAndSaveAllDividends()
	// timeseries.FetchAndSaveHistoricalQuotes()
}
