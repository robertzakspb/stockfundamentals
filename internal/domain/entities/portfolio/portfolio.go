package portfolio

import (
	"fmt"
	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/invest-core/quote/quotefetcher"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
)

type Portfolio struct {
	Lots []lot.Lot `json:"lots"`
	Cash float64
}

func (portfolio Portfolio) UniquePositions() []lot.Lot {
	uniquePositions := []lot.Lot{}
	for _, lot := range portfolio.Lots {
		foundLotWithSameTicker := false
		lotWithSameTickerIndex := 0
		for i, uniquePosition := range uniquePositions {
			if lot.SecurityId.String() == uniquePosition.SecurityId.String() {
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

type LotWithSecurity struct {
	lot   lot.Lot
	stock security.Stock
}

func (portfolio Portfolio) PrintAllPositions() {
	positions := portfolio.UniquePositions()
	lotsWithSecurities := []LotWithSecurity{}

	ids := uuid.UUIDs{}
	for _, p := range positions {
		ids = append(ids, p.SecurityId)
	}
	securities, err := security_master.GetSecuritiesById(ids)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}

	entitySecurities := []entity.Security{}
	for _, s := range securities {
		if s.GetId() == uuid.Nil {
			continue
		}
		entitySecurities = append(entitySecurities, entity.Security{
			Figi:   s.GetFigi(),
			ISIN:   s.GetIsin(),
			Ticker: s.GetTicker(),
			MIC:    s.GetMic(),
		})
	}

	quotes := quotefetcher.FetchQuotesFor(entitySecurities)

	//Calculating the total portfolio value
	totalPortfolioValue := 0.0
	for _, lot := range positions {
		for _, quote := range quotes {
			var stock security.Stock
			for _, s := range securities {
				if s.GetId() == lot.SecurityId {
					stock = security.Stock{
						Id:           s.GetId(),
						CompanyName:  s.GetCompanyName(),
						Figi:         s.GetFigi(),
						Isin:         s.GetIsin(),
						SecurityType: s.GetSecurityType(),
						Country:      s.GetCountry(),
						Ticker:       s.GetTicker(),
						IssueSize:    s.GetIssueSize(),
						Sector:       s.GetSector(),
						MIC:          s.GetMic(),
					}
					lotsWithSecurities = append(lotsWithSecurities, LotWithSecurity{lot: lot, stock: stock})
				}
			}

			if stock.GetFigi() == quote.Figi() {
				marketValue, _ := lot.MarketValue(quote)
				totalPortfolioValue += marketValue
			}
		}
	}

	// slices.SortFunc(positions, func(a lot.Lot, b lot.Lot) int {
	// 	return 1
	// })

	//Displaying the portfolio
	for _, lot := range lotsWithSecurities {
		profitOrLoss := 0.0
		var stockQuote entity.SimpleQuote
		didFindQuote := false
		for _, quote := range quotes {
			if lot.stock.GetFigi() == quote.Figi() {
				profitOrLoss = lot.lot.CurrentReturn(quote)
				stockQuote = quote
				didFindQuote = true
			}
		}
		if !didFindQuote {
			fmt.Println("Unable to fetch quotes for ", lot.stock.GetTicker(), "Quantity: ", lot.lot.Quantity, "Spent on position: ", lot.lot.CostBasis())
			continue
		}
		fmt.Printf("%-6s", lot.stock.GetTicker())
		fmt.Printf("Quantity: %.0f | ", lot.lot.Quantity)
		fmt.Printf("Opening Price: %.1f | ", lot.lot.PricePerUnit)
		fmt.Printf("Profit: %.2f | ", profitOrLoss*100)
		mv, _ := lot.lot.MarketValue(stockQuote)
		fmt.Printf("Percentage of portfolio: %.2f %% | ", mv/totalPortfolioValue*100)
		fmt.Printf("Market value: %.0f\n", mv)
	}
}
