package portfolio

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/invest-core/quote/quotefetcher"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	"github.com/compoundinvest/stockfundamentals/internal/domain/portfolio/lot"
)

type Lot = lot.Lot

type Portfolio struct {
	Lots []Lot `json:"lots"`
	Cash float64
}

func (portfolio Portfolio) GetPositionByTicker(ticker string) (Lot, error) {
	for _, position := range portfolio.UniquePositions() {
		if position.Security.GetTicker() == ticker {
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
			if lot.Security.GetFigi() == uniquePosition.Security.GetFigi() {
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

func (portfolio Portfolio) PrintAllPositions() {

	positions := portfolio.UniquePositions()

	//Fetching quotes
	securities := []entity.Security{}
	for _, lot := range positions {
		securities = append(securities, entity.Security{
			Figi:   lot.Security.GetFigi(),
			Ticker: lot.Security.GetTicker(),
			MIC:    lot.Security.GetMic(),
		})
	}

	quotes := quotefetcher.FetchQuotesFor(securities)

	//Calculating the total portfolio value
	totalPortfolioValue := 0.0
	for _, lot := range positions {
		for _, quote := range quotes {
			//TODO: Refactor this abomination
			if lot.Security.GetFigi() == "" {
				logger.Log("Position is missing figi. Ticker: "+lot.Security.GetTicker()+". Quantity: "+strconv.FormatFloat(lot.Quantity, 'E', -1, 64), logger.ERROR)
			}
			if lot.Security.GetFigi() == quote.Figi() {
				marketValue, _ := lot.MarketValue(quote)
				totalPortfolioValue += marketValue
			}
		}
	}

	slices.SortFunc(positions, func(a Lot, b Lot) int {

		return 1
	})

	//Displaying the portfolio
	for _, lot := range positions {
		profitOrLoss := 0.0
		var stockQuote entity.SimpleQuote
		didFindQuote := false
		for _, quote := range quotes {
			if lot.Security.GetFigi() == quote.Figi() {
				profitOrLoss = lot.CurrentReturn(quote)
				stockQuote = quote
				didFindQuote = true
			}
		}
		if !didFindQuote {
			fmt.Println("Unable to fetch quotes for ", lot.Security.GetTicker(), "Quantity: ", lot.Quantity, "Spent on position: ", lot.CostBasis())
			continue
		}
		fmt.Printf("%-6s", lot.Security.GetTicker())
		fmt.Printf("Quantity: %.0f | ", lot.Quantity)
		fmt.Printf("Opening Price: %.1f | ", lot.PricePerUnit)
		fmt.Printf("Profit: %.2f | ", profitOrLoss*100)
		mv, _ := lot.MarketValue(stockQuote)
		fmt.Printf("Percentage of portfolio: %.2f %% | ", mv/totalPortfolioValue*100)
		fmt.Printf("Market value: %.0f\n", mv)
	}
}
