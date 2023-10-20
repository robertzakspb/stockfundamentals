package lot

import (
	"math"
	"time"

	quote "github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/invest-core/quote/moexapi"
)

type Lot struct {
	Quantity         float64   `json:"quantity"`
	CostBasis        float64   `json:"costBasis"`
	OpeningDate      time.Time `json:"openingDate"`
	Currency         string    `json:"currency"`
	CompanyName      string    `json:"companyName"`
	Ticker           string    `json:"ticker"`
	BrokerName       string    `json:"brokerName"`
	AnnualizedReturn float64   `json:"anualizedReturn"`
}

// Returns the current profit on the lot given a quote (expressed as a percentage)
func (lot Lot) currentReturn(quote quote.SimpleQuote) float64 {
	if lot.CostBasis == 0 {
		return 0
	}
	return (quote.Quote() - lot.CostBasis) / lot.CostBasis
}

// Returns the current profit on the lot given a quote (in currency)
func (lot Lot) currentProfit(quote quote.SimpleQuote) float64 {
	return quote.Quote() - lot.CostBasis
}

// Returns the tax liabilities the investor would incur if the lot were sold
func (lot Lot) TaxAmountIfGainsRealized(quote quote.SimpleQuote) float64 {
	threeYearsAgo := time.Now().AddDate(-3, 0, 0)
	threeYearsHavePassedSincePositonOpened := lot.OpeningDate.Before(threeYearsAgo)
	if threeYearsHavePassedSincePositonOpened {
		return 0
	}

	currentProfit := lot.currentProfit(quote)
	if currentProfit <= 0 {
		return 0
	}

	potentialTaxOwed := currentProfit * 0.13
	return potentialTaxOwed
}

func CalculateLotsAnnualizedReturns(lots []Lot) {
	allTickers := []string{}
	for _, lot := range lots {
		allTickers = append(allTickers, lot.Ticker)
	}

	quotes := moexapi.FetchQuotes(allTickers)

	for i := 0; i < len(lots); i++ {
		for _, quote := range quotes {
			if quote.Ticker() == lots[i].Ticker {
				lots[i].calculateAnnualizedProfit(quote)
				break
			}
		}
	}
}

func (lot *Lot) calculateAnnualizedProfit(quote quote.SimpleQuote) {
	durationInYears := time.Since(lot.OpeningDate).Hours() / 24 / 365
	currentReturn := lot.currentReturn(quote)

	annualizedReturn := math.Pow((1+currentReturn), 1/durationInYears) - 1
	annualizedReturnAsPercentage := annualizedReturn * 100

	lot.AnnualizedReturn = annualizedReturnAsPercentage
}
