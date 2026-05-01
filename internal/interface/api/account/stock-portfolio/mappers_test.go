package portfolio

import (
	"testing"
	"time"

	stockportfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_mapLotToDto(t *testing.T) {
	id := uuid.New()
	accountId := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	quantity := 10.0
	pricePerUnit := 7.5
	currency := "USD"
	figi := "testFigi"
	isin := "testIsin"
	ticker := "testTicker"
	quote := 20.0

	lot := lot.Lot{
		Id:           id,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		Currency:     currency,
		AccountId:    accountId,
		Figi:         figi,
		Quote:        quote,
		Stock: security.Stock{
			Isin:   isin,
			Ticker: ticker,
		},
	}
	mappedLot := mapLotToDto(lot)

	test.AssertEqual(t, id, mappedLot.Id)
	test.AssertEqual(t, accountId, mappedLot.AccountId)
	test.AssertEqual(t, createdAt, mappedLot.CreatedAt)
	test.AssertEqual(t, updatedAt, mappedLot.UpdatedAt)
	test.AssertEqual(t, quantity, mappedLot.Quantity)
	test.AssertEqual(t, pricePerUnit, mappedLot.PricePerUnit)
	test.AssertEqual(t, currency, mappedLot.Currency)
	test.AssertEqual(t, figi, mappedLot.Figi)
	test.AssertEqual(t, isin, mappedLot.Isin)
	test.AssertEqual(t, ticker, mappedLot.Ticker)
	test.AssertEqual(t, quote, mappedLot.Quote)
	test.AssertEqual(t, 125, mappedLot.CurrentPL)
	test.AssertEqual(t, 200, mappedLot.MarketValue)
}

func Test_mapPortfolioToDto(t *testing.T) {
	id := uuid.New()
	accountId := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	quantity := 10.0
	pricePerUnit := 7.5
	currency := "USD"
	figi := "testFigi"
	isin := "testIsin"
	ticker := "testTicker"
	quote := 20.0

	sampleLot := lot.Lot{
		Id:           id,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		Currency:     currency,
		AccountId:    accountId,
		Figi:         figi,
		Quote:        quote,
		Stock: security.Stock{
			Isin:   isin,
			Ticker: ticker,
		},
	}
	portfolio := stockportfolio.Portfolio{Lots: []lot.Lot{sampleLot}}
	dto := mapPortfolioToDto(portfolio)

	test.AssertEqual(t, 1, len(dto.Lots))
	test.AssertEqual(t, "testFigi", dto.Lots[0].Figi)
}
