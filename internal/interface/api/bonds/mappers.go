package bondsapi

import "github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"

func mapBondsToDTOs(bondList []bonds.Bond, includeCoupons bool) []BondDTO {
	bondDtos := make([]BondDTO, len(bondList))

	for i := range bondList {
		bondDto := BondDTO{
			Figi:                    bondList[i].Figi,
			Isin:                    bondList[i].Isin,
			Lot:                     bondList[i].Lot,
			Currency:                bondList[i].Currency,
			Name:                    bondList[i].Name,
			CountryOfRisk:           bondList[i].CountryOfRisk,
			RealExchange:            bondList[i].CountryOfRisk,
			CouponCountPerYear:      bondList[i].CouponCountPerYear,
			MaturityDate:            bondList[i].MaturityDate,
			NominalValue:            bondList[i].NominalValue,
			NominalCurrency:         bondList[i].NominalCurrency,
			InitialNominalValue:     bondList[i].InitialNominalValue,
			InitialNominalCurrency:  bondList[i].InitialNominalCurrency,
			RegistrationDate:        bondList[i].RegistrationDate,
			PlacementDate:           bondList[i].PlacementDate,
			PlacementPrice:          bondList[i].PlacementPrice,
			PlacementCurrency:       bondList[i].PlacementCurrency,
			AccruedInterest:         bondList[i].AccruedInterest,
			IssueSize:               bondList[i].IssueSize,
			IssueSizePlan:           bondList[i].IssueSizePlan,
			HasFloatingCoupon:       bondList[i].HasFloatingCoupon,
			IsPerpetual:             bondList[i].IsPerpetual,
			HasAmortization:         bondList[i].HasAmortization,
			IsAvailableForIis:       bondList[i].IsAvailableForIis,
			IsForQualifiedInvestors: bondList[i].IsForQualifiedInvestors,
			IsSubordinated:          bondList[i].IsSubordinated,
			RiskLevel:               bonds.RiskLevel_name[int32(bondList[i].RiskLevel)],
			BondType:                bonds.BondType_name[int32(bondList[i].BondType)],
			CallOptionExerciseDate:  bondList[i].CallOptionExerciseDate,
			YieldToMaturity:         bondList[i].YieldToMaturity,
			YieldToCallOption:       bondList[i].YieldToCallOption,
		}
		if includeCoupons {
			bondDto.Coupons = mapCouponsToCouponDTOs(bondList[i].Coupons)
		} else {
			bondDto.Coupons = []CouponDTO{}
		}
		bondDtos[i] = bondDto
	}

	return bondDtos
}

func mapCouponsToCouponDTOs(coupons []bonds.Coupon) []CouponDTO {
	couponDtos := make([]CouponDTO, len(coupons))

	for i := range coupons {
		couponDTO := CouponDTO{
			Figi:               coupons[i].Figi,
			Date:               coupons[i].CouponDate,
			CouponNumber:       coupons[i].CouponNumber,
			RecordDate:         coupons[i].RecordDate,
			PerBondAmount:      coupons[i].PerBondAmount,
			CouponType:         bonds.CouponType_name[int32(coupons[i].CouponType)],
			CouponStartDate:    coupons[i].CouponStartDate,
			CouponEndDate:      coupons[i].CouponEndDate,
			CouponPeriodInDays: coupons[i].CouponPeriod,
		}
		couponDtos[i] = couponDTO
	}

	return couponDtos
}
