package stockportfolio

import (
	"errors"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	stringhelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/string-helpers"
)

type Portfolio struct {
	Lots []lot.Lot
	Cash []CashPosition //TODO: Populate this
}

func (portfolio Portfolio) UniquePositions() []lot.Lot {
	uniquePositions := []lot.Lot{}
	for _, lot := range portfolio.Lots {
		foundLotWithSameTicker := false
		lotWithSameTickerIndex := 0
		for i, uniquePosition := range uniquePositions {
			if lot.Figi == uniquePosition.Figi {
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

func (portfolio Portfolio) GetEtfLotFigis() []string {
	etfLotFigis := []string{}

	for _, lot := range portfolio.Lots {
		if lot.Stock.SecurityType == security.ETF {
			etfLotFigis = append(etfLotFigis, lot.Figi)
		}
	}

	etfLotFigis = stringhelpers.RemoveDuplicatesFrom(etfLotFigis)

	return etfLotFigis
}

func LotFigis(lots []lot.Lot) []string {
	figis := []string{}
	for _, lot := range lots {
		figis = append(figis, lot.Figi)
	}

	figis = stringhelpers.RemoveDuplicatesFrom(figis)

	return figis
}

func (p Portfolio) PositionCurrencies() []string {
	currenciesInPortfolio := []string{}
	for _, position := range p.UniquePositions() {
		alreadyInArray := false
		for _, currency := range currenciesInPortfolio {
			if position.Currency == currency {
				alreadyInArray = true
			}
		}
		if !alreadyInArray {
			currenciesInPortfolio = append(currenciesInPortfolio, position.Currency)
		}
	}
	return currenciesInPortfolio
}

func MatchLotsWithStocks(lots []lot.Lot, securities []security.Stock) ([]lot.Lot, []error) {
	errorList := []error{}
	for i := range lots {
		foundSecurity := false
		for _, s := range securities {
			if lots[i].Figi == s.GetFigi() {
				foundSecurity = true
				lots[i].Stock = s
			}
		}
		if !foundSecurity {
			errorList = append(errorList, errors.New("Failed to find a security for lot "+lots[i].Figi))
		}
	}

	return lots, errorList
}

func MatchLotsWithQuotes(lots []lot.Lot, quotes []entity.SimpleQuote) ([]lot.Lot, error) {
	var err error
	for i := range lots {
		foundQuote := false
		for _, q := range quotes {
			if lots[i].Figi == q.Figi() {
				foundQuote = true
				lots[i].Quote = q.Quote()
			}
		}
		if !foundQuote {
			err = errors.New("Failed to find the quote for " + lots[i].Figi)
			continue
		}
	}
	return lots, err
}

func LotStocks(lots []lot.Lot) []security.Stock {
	stocks := []security.Stock{}
	for i := range lots {
		stocks = append(stocks, lots[i].Stock)
	}
	return stocks
}
