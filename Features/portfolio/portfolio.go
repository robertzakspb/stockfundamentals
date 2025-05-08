package portfolio

import (
	"fmt"

	"github.com/compoundinvest/invest-core/quote/quotefetcher"
	"github.com/compoundinvest/stockfundamentals/features/portfolio/lot"
)

type Lot = lot.Lot

type Portfolio struct {
	Lots []Lot `json:"lots"`
}

func (portfolio Portfolio) GetPositionByTicker(ticker string) (Lot, error) {
	for _, position := range portfolio.UniquePositions() {
		if position.Ticker == ticker {
			return position, nil
		}
	}

	return Lot{}, fmt.Errorf("didn't find a position with ticker %s", ticker)
}

func (portfolio Portfolio) UniquePositions() []Lot {
	uniquePositions := []Lot{}
	for _, lot := range portfolio.Lots {
		foundLotWithSameTicker := false
		lotWithSameTickerIndex := 0
		for i, uniquePosition := range uniquePositions {
			if lot.Ticker == uniquePosition.Ticker {
				foundLotWithSameTicker = true
				lotWithSameTickerIndex = i
			}
		}

		if foundLotWithSameTicker {
			mergedLot, err := uniquePositions[lotWithSameTickerIndex].MergeWith(lot)
			if err != nil {
				//If there was an error, add both positions
				uniquePositions = append(uniquePositions, lot)
				uniquePositions = append(uniquePositions, uniquePositions[lotWithSameTickerIndex])
			}
			uniquePositions[lotWithSameTickerIndex] = mergedLot
		} else {
			uniquePositions = append(uniquePositions, lot)
		}
	}

	return uniquePositions
}

// Lists all tickers present in the portfolio
func (portfolio Portfolio) PrintAllPositions() {

	positions := portfolio.UniquePositions()

	//Fetching quotes
	tickersWithExchanges := []quotefetcher.TickerWithMarket{}
	for _, lot := range positions {
		tickersWithExchanges = append(tickersWithExchanges, quotefetcher.TickerWithMarket{Ticker: lot.Ticker, MIC: lot.MIC})
	}

	quotes := quotefetcher.FetchQuotesFor(tickersWithExchanges)

	//Calculating the total portfolio value
	totalPortfolioValue := 0.0
	for _, lot := range positions {
		for _, quote := range quotes {
			if lot.Ticker == quote.Ticker() {
				totalPortfolioValue += lot.MarketValue(quote.Quote())
			}
		}
	}

	//Displaying the portfolio
	for _, lot := range positions {
		profitOrLoss := 0.0
		stockQuote := 0.0
		for _, quote := range quotes {
			if lot.Ticker == quote.Ticker() {
				profitOrLoss = lot.CurrentReturn(quote.Quote())
				stockQuote = quote.Quote()
			}
		}
		fmt.Printf("%-6s", lot.Ticker)
		fmt.Printf("Quantity: %.0f | ", lot.Quantity)
		fmt.Printf("Opening Price: %.1f | ", lot.OpeningPrice)
		fmt.Printf("Profit: %.2f | ", profitOrLoss*100)
		fmt.Printf("Percentage of portfolio: %.2f %% | ", lot.MarketValue(stockQuote)/totalPortfolioValue*100)
		fmt.Printf("Market value: %.0f\n", lot.MarketValue(stockQuote))
	}
}
