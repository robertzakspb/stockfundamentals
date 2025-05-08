package main

import (
	"fmt"
	"github.com/compoundinvest/stockfundamentals/dataseed"
	"github.com/compoundinvest/stockfundamentals/features/fundamentaldata/dividend"
	"github.com/compoundinvest/stockfundamentals/features/fundamentaldata/security"
	"github.com/compoundinvest/stockfundamentals/features/portfolio"
)

func main() {
	myPorfolio := portfolio.GeMyPortfolio()
	myPorfolio.PrintAllPositions()

	totalSum := 0.0
	for _, dividend := range myPorfolio.UpcomingDividends() {
		fmt.Println("Ticker:", dividend.Ticker, "| Payout:", dividend.GrossPayout())
		totalSum += dividend.GrossPayout()
	}
	fmt.Printf("Total projected payout: %.1f\n", totalSum)
	dataseed.InitialSeed()

	fetchExternalData()
	security.FetchSecuritiesFromDB()
	// security.FetchAndSaveSecurities()
	// stocks, err := security.FetchSecuritiesFromDB()
	// if err != nil {
	// fmt.Println(err.Error())
	// }

	// security.FetchSecuritiesFromDB()
	// dividend.FetchAndSaveDividendsForAllStocks()
	// fetchExternalData()
}

func fetchExternalData() {
	security.FetchAndSaveSecurities()
	dividend.FetchAndSaveAllDividends()
}
