package bondsapi

import (
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
	}

	return domain
}

func mapBondLotToDto(lot bonds.BondLot) bondPositionLotDto {
	dto := bondPositionLotDto{
		Figi:             lot.Figi,
		Isin:             lot.Isin,
		OpeningDate:      lot.OpeningDate,
		ModificationDate: lot.ModificationDate,
		AccountId:        lot.AccountId.String(),
		Quantity:         lot.Quantity,
		PricePerUnit:     lot.PricePerUnit,
	}

	return dto
}
