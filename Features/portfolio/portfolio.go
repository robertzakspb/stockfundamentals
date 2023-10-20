package portfolio

import (
	"github.com/compoundinvest/stockfundamentals/Features/portfolio/lot"
)

type Lot = lot.Lot

type Portfolio struct {
	Lots []Lot `json:"lots"`
}

// func ConvertToPortfolio(lots []lot.Lot) Portfolio {
// 	positions := []position.Position{}
// 	for _, lot := range lots {
// 		if len(positions) == 0 {
// 			initialLots := []Lot{}
// 			newPosition := position.Position{Ticker: lot.Ticker, CompanyName: lot.CompanyName, Lots: initialLots}
// 			positions = append(positions, newPosition)
// 			continue
// 		}

// 		didFindTheTickerInPositions := false
// 		for i := 0; i < len(positions); i++ {
// 			if lot.Ticker == positions[i].Ticker {
// 				didFindTheTickerInPositions = true
// 				positions[i].Lots = append(positions[i].Lots, lot)
// 				break
// 			}
// 		}

// 		if !didFindTheTickerInPositions {
// 			newPosition := position.Position{Ticker: lot.Ticker, CompanyName: lot.CompanyName, Lots: []Lot{lot}}
// 			positions = append(positions, newPosition)
// 		}
// 	}

// 	return Portfolio{positions}
// }
