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
		Figi:                   lot.Figi,
		Isin:                   lot.Isin,
		Name:                   lot.Bond.Name,
		OpeningDate:            lot.OpeningDate,
		ModificationDate:       lot.ModificationDate,
		AccountId:              lot.AccountId,
		Quantity:               lot.Quantity,
		PricePerUnitPercentage: lot.PricePerUnitPercentage,
		MarketValueInRUB:       lot.MarketValueInRUB,
		Ytm:                    lot.Bond.YieldTomaturity,
		YieldToCallOption:      lot.Bond.YieldToCallOption,
		QuoteInPercentage:      lot.Bond.QuoteInPercentage,
		CurrentCouponYield:     lot.Bond.CurrentCouponYield(),
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
			Amount:    item.Amount,
			Currency:  item.Currency,
		}
		dtos = append(dtos, dto)
	}
	return dtos
}
