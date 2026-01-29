package portfolio

import (
	"fmt"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/invest-core/quote/quotefetcher"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

type Portfolio struct {
	Lots []lot.Lot `json:"lots"`
	Cash []CashPosition `json:"cashPositions"` //TODO: Populate this
}

func (portfolio Portfolio) UniquePositions() []lot.Lot {
	uniquePositions := []lot.Lot{}
	for _, lot := range portfolio.Lots {
		foundLotWithSameTicker := false
		lotWithSameTickerIndex := 0
		for i, uniquePosition := range uniquePositions {
			if lot.SecurityId == uniquePosition.SecurityId {
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
	Lot   lot.Lot `json:"lot"`
	Stock security.Stock `json:"stock"`
}

func (portfolio Portfolio) WithPLs() []LotWithSecurity {
	positions := portfolio.UniquePositions()
	lotsWithSecurities := []LotWithSecurity{}

	ids := []string{}
	for _, p := range positions {
		ids = append(ids, p.SecurityId)
	}
	securities, err := security_master.GetSecuritiesFilteredByFigi(ids)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}

	entitySecurities := []entity.Security{}
	for _, s := range securities {
		if s.GetId() == "" {
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

	for _, lot := range positions {
		
			var stock security.Stock
			for _, s := range securities {
				if s.GetId() == lot.SecurityId {
					stock = security.Stock{
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
					lotsWithSecurities = append(lotsWithSecurities, LotWithSecurity{Lot: lot, Stock: stock})
				}
			}
		
	}

	for i, lot := range lotsWithSecurities {
		profitOrLoss := 0.0
		// var stockQuote entity.SimpleQuote
		didFindQuote := false
		for _, quote := range quotes {
			if lot.Stock.GetFigi() == quote.Figi() {
				profitOrLoss = lot.Lot.CurrentReturn(quote)
				lotsWithSecurities[i].Lot.CurrentPL = profitOrLoss
				// stockQuote = quote
				didFindQuote = true
			}
		}
		if !didFindQuote {
			fmt.Println("Unable to fetch quotes for ", lot.Stock.GetTicker(), "Quantity: ", lot.Lot.Quantity, "Spent on position: ", lot.Lot.CostBasis())
			continue
		}
		// fmt.Printf("%-6s", lot.stock.GetTicker())
		// fmt.Printf("Quantity: %.0f | ", lot.lot.Quantity)
		// fmt.Printf("Opening Price: %.1f | ", lot.lot.PricePerUnit)
		// fmt.Printf("Profit: %.2f | ", profitOrLoss*100)
		// mv, _ := lot.lot.MarketValue(stockQuote)
		// fmt.Printf("Market value: %.0f\n", mv)
	}

	return lotsWithSecurities
}

