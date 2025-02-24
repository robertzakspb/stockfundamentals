package main

import (
	// "path"

	// "github.com/compoundinvest/stockfundamentals/Features/fundamentaldata/moexapi/securityinfo"
	"fmt"

	"github.com/compoundinvest/stockfundamentals/Features/portfolio"
	"github.com/compoundinvest/stockfundamentals/dataseed"
)

type Stock struct {
	Id            string `sql:"id"`
	Company_name  string `sql:"company_name"`
	Is_public     bool   `sql:"is_public"`
	Isin          string `sql:"isin"`
	Security_type string `sql:"security_type"`
	Country_iso2  string `sql:"country_iso2"`
	Ticker        string `sql:"ticker"`
	Share_count   uint64 `sql:"share_count"`
	Sector        string `sql:"sector"`
}

func main() {
	// portfolio.GeMyPortfolio().PrintAllPositions()
	totalSum := 0.0
	for _, dividend := range portfolio.GeMyPortfolio().UpcomingDividends() {
		fmt.Println("Ticker:", dividend.Ticker, "| Payout:", dividend.GrossPayout())
		totalSum += dividend.GrossPayout()
	}
	fmt.Printf("Total projected payout: %.1f", totalSum)
	dataseed.InitialSeed()
}
