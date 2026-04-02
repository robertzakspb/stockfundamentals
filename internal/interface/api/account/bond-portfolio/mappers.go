package bondportfolioapi

import (
	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/google/uuid"
)

func mapBondLotDtoToDomain(dto bondPositionLotDto) bonds.BondLot {
	accountId, _ := uuid.Parse(dto.AccountId)
	domain := bonds.BondLot{
		Id:               uuid.New(),
		Figi:             dto.Figi,
		Isin:             dto.Isin,
		OpeningDate:      dto.OpeningDate,
		ModificationDate: dto.ModificationDate,
		AccountId:        accountId,
		Quantity:         dto.Quantity,
		PricePerUnit:     dto.PricePerUnit,
		PricePerUnitInRUB: dto.PricePerUnitInRUB,
	}

	return domain
}

func mapBondLotToDto(lot bonds.BondLot) bondPositionLotDto {
	dto := bondPositionLotDto{
		Figi:              lot.Figi,
		Isin:              lot.Isin,
		OpeningDate:       lot.OpeningDate,
		ModificationDate:  lot.ModificationDate,
		AccountId:         lot.AccountId.String(),
		Quantity:          lot.Quantity,
		PricePerUnit:      lot.PricePerUnit,
		CurrentYtm:        lot.Bond.YieldToMaturity,
		YieldToCallOption: lot.Bond.YieldToCallOption,
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
