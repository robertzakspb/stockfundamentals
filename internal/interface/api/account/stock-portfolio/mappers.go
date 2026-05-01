package portfolio

import (
	stockportfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
)

func mapPortfolioToDto(portfolio stockportfolio.Portfolio) PortfolioDto {
	dtos := make([]LotDto, len(portfolio.Lots))
	for i := range portfolio.Lots {
		dtos[i] = mapLotToDto(portfolio.Lots[i])
	}

	return PortfolioDto{Lots: dtos}
}

func mapLotToDto(lot lot.Lot) LotDto {
	mv, _ := lot.MarketValue()
	dto := LotDto{
		Id:            lot.Id,
		CreatedAt:     lot.CreatedAt,
		UpdatedAt:     lot.UpdatedAt,
		Quantity:      lot.Quantity,
		PricePerUnit:  lot.PricePerUnit,
		Currency:      lot.Currency,
		AccountId:     lot.AccountId,
		Figi:          lot.Figi,
		Quote:         lot.Quote,
		CurrentPL:     lot.CurrentPL(),
		CurrentReturn: lot.CurrentReturn(),
		MarketValue:   mv,
		Isin:          lot.Stock.Isin,
		Ticker:        lot.Stock.Ticker,
	}

	return dto
}
