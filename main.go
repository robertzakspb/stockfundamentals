package main

import (
	// "path"

	// "github.com/compoundinvest/stockfundamentals/Features/fundamentaldata/moexapi/securityinfo"
	"fmt"

	"github.com/compoundinvest/stockfundamentals/Features/portfolio"
	"github.com/compoundinvest/stockfundamentals/dataseed"
)

func main() {
	myPorfolio := portfolio.GeMyPortfolio()
	myPorfolio.PrintAllPositions()

	totalSum := 0.0
	for _, dividend := range myPorfolio.UpcomingDividends() {
		fmt.Println("Ticker:", dividend.Ticker, "| Payout:", dividend.GrossPayout())
		totalSum += dividend.GrossPayout()
	}
	fmt.Printf("Total projected payout: %.1f", totalSum)
	dataseed.InitialSeed()
}