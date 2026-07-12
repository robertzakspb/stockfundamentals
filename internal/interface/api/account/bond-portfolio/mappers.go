package bondportfolioapi

import (
	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/google/uuid"
)

func mapBondLotDtoToDomain(dto bondPositionLotDto) bonds.BondLot {
	domain := bonds.BondLot{
		Id:                     uuid.New(),
		Figi:                   dto.Figi,
		Isin:                   dto.Isin,
		OpeningDate:            dto.OpeningDate,
		ModificationDate:       dto.ModificationDate,
		AccountId:              dto.AccountId,
		Quantity:               dto.Quantity,
		PricePerUnitPercentage: dto.PricePerUnitPercentage,
		MarketValueInRUB:       dto.MarketValueInRUB,
	}

	return domain
}

func mapBondLotToDto(lot bonds.BondLot) bondPositionLotDto {
	dto := bondPositionLotDto{
		Figi:                    lot.Figi,
		Isin:                    lot.Isin,
		Name:                    lot.Bond.Name,
		OpeningDate:             lot.OpeningDate,
		ModificationDate:        lot.ModificationDate,
		AccountId:               lot.AccountId,
		Quantity:                lot.Quantity,
		PricePerUnitPercentage:  lot.PricePerUnitPercentage,
		SimpleYtm:               lot.Bond.SimpleYieldToMaturity,
		SimpleYieldToCallOption: lot.Bond.SimpleYieldToCallOption,
		MarketValueInRUB:        lot.MarketValueInRUB,
		Ytm:                     lot.Bond.YieldTomaturity,
		QuoteInPercentage:       lot.Bond.QuoteInPercentage,
	}

	return dto
}

func mapTimeLineItemsToDtos(items []bondportfolio.TimeLineItem) []timeLineItemDto {
	dtos := []timeLineItemDto{}

	for _, item := range items {
		dto := timeLineItemDto{
			Timestamp: item.Timestamp,
			EventName: item.EventName,
			BondName:  item.BondName,
		}
		dtos = append(dtos, dto)
	}
	return dtos
}
