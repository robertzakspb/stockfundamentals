package main

import (
	"github.com/compoundinvest/stockfundamentals/internal/application/portfolio"
	dividend "github.com/compoundinvest/stockfundamentals/internal/interface/api/fundamentals/dividend"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/dataseed"
)

func main() {
	dataseed.InitialSeed()
	fetchExternalData()
	portfolio.GeMyPortfolio().PrintAllPositions()
}

func fetchExternalData() {
	security_master.FetchAndSaveSecurities()
	dividend.FetchAndSaveAllDividends()
	// timeseries.FetchAndSaveHistoricalQuotes()
}
