package bondservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	pb "opensource.tbank.ru/invest/invest-go/proto"
)

func mapTinkoffBondToBond(tinkoffBond *pb.Bond) bonds.Bond {
	bond := bonds.Bond{
		Id:                      uuid.New(),
		Figi:                    tinkoffBond.Figi,
		Isin:                    tinkoffBond.Isin,
		Lot:                     int(tinkoffBond.Lot),
		Currency:                tinkoffBond.Currency,
		Name:                    tinkoffBond.Name,
		CountryOfRisk:           tinkoffBond.CountryOfRisk,
		RealExchange:            tinkoffBond.RealExchange.String(),
		CouponCountPerYear:      int(tinkoffBond.CouponQuantityPerYear),
		MaturityDate:            tinkoffBond.MaturityDate.AsTime(),
		NominalValue:            tinkoffBond.Nominal.ToFloat(),
		NominalCurrency:         tinkoffBond.Nominal.GetCurrency(),
		InitialNominalValue:     tinkoffBond.InitialNominal.ToFloat(),
		InitialNominalCurrency:  tinkoffBond.InitialNominal.GetCurrency(),
		RegistrationDate:        tinkoffBond.StateRegDate.AsTime(),
		PlacementDate:           tinkoffBond.PlacementDate.AsTime(),
		PlacementPrice:          tinkoffBond.PlacementPrice.ToFloat(),
		PlacementCurrency:       tinkoffBond.PlacementPrice.GetCurrency(),
		AccruedInterest:         tinkoffBond.AciValue.ToFloat(),
		IssueSize:               int(tinkoffBond.IssueSize),
		IssueSizePlan:           int(tinkoffBond.IssueSizePlan),
		HasFloatingCoupon:       tinkoffBond.FloatingCouponFlag,
		IsPerpetual:             tinkoffBond.PerpetualFlag,
		HasAmortization:         tinkoffBond.AmortizationFlag,
		IsAvailableForIis:       tinkoffBond.ForIisFlag,
		IsForQualifiedInvestors: tinkoffBond.ForQualInvestorFlag,
		IsSubordinated:          tinkoffBond.SubordinatedFlag,
		RiskLevel:               mapTinkoffRiskLevel(tinkoffBond.RiskLevel),
		BondType:                bonds.BondType(tinkoffBond.BondType),
		CallOptionExerciseDate:  tinkoffBond.CallDate.AsTime(),
	}
	return bond
}

func mapTinkoffRiskLevel(rl pb.RiskLevel) bonds.RiskLevel {
	switch rl {
	case pb.RiskLevel_RISK_LEVEL_UNSPECIFIED:
		return bonds.UNSPECIFIED_RISK_LEVEL
	case pb.RiskLevel_RISK_LEVEL_LOW:
		return bonds.LOW_RISK_LEVEL
	case pb.RiskLevel_RISK_LEVEL_MODERATE:
		return bonds.MODERATE_RISK_LEVEL
	case pb.RiskLevel_RISK_LEVEL_HIGH:
		return bonds.HIGH_RISK_LEVEL
	default:
		logger.Log("Unknown pb.RiskLevel value: "+rl.String(), logger.ALERT)
		return bonds.UNSPECIFIED_RISK_LEVEL
	}
}

func mapBondToDbBond(bond bonds.Bond) bondsdb.BondDbModel {
	dbBond := bondsdb.BondDbModel{
		Id:                      bond.Id,
		Figi:                    bond.Figi,
		Isin:                    bond.Isin,
		Lot:                     int64(bond.Lot),
		Currency:                bond.Currency,
		Name:                    bond.Name,
		CountryOfRisk:           bond.CountryOfRisk,
		RealExchange:            bond.RealExchange,
		CouponCountPerYear:      int64(bond.CouponCountPerYear),
		MaturityDate:            bond.MaturityDate,
		NominalValue:            bond.NominalValue,
		NominalCurrency:         bond.NominalCurrency,
		InitialNominalValue:     bond.InitialNominalValue,
		InitialNominalCurrency:  bond.InitialNominalCurrency,
		RegistrationDate:        bond.RegistrationDate,
		PlacementDate:           bond.PlacementDate,
		PlacementPrice:          bond.PlacementPrice,
		PlacementCurrency:       bond.PlacementCurrency,
		AccruedInterest:         bond.AccruedInterest,
		IssueSize:               int64(bond.IssueSize),
		IssueSizePlan:           int64(bond.IssueSizePlan),
		HasFloatingCoupon:       bond.HasFloatingCoupon,
		IsPerpetual:             bond.IsPerpetual,
		HasAmortization:         bond.HasAmortization,
		IsAvailableForIis:       bond.IsAvailableForIis,
		IsForQualifiedInvestors: bond.IsForQualifiedInvestors,
		IsSubordinated:          bond.IsSubordinated,
		RiskLevel:               bonds.RiskLevel_name[int32(bond.RiskLevel)],
		BondType:                bonds.BondType_name[int32(bond.BondType)],
		CallOptionExerciseDate:  bond.CallOptionExerciseDate,
	}

	return dbBond
}

func mapDbBondToBond(dbModel bondsdb.BondDbModel) bonds.Bond {
	domain := bonds.Bond{
		Figi:                    dbModel.Figi,
		Id:                      dbModel.Id,
		Isin:                    dbModel.Isin,
		Lot:                     int(dbModel.Lot),
		Currency:                dbModel.Currency,
		Name:                    dbModel.Name,
		CountryOfRisk:           dbModel.CountryOfRisk,
		RealExchange:            dbModel.RealExchange,
		CouponCountPerYear:      int(dbModel.CouponCountPerYear),
		MaturityDate:            dbModel.MaturityDate,
		NominalValue:            dbModel.NominalValue,
		NominalCurrency:         dbModel.NominalCurrency,
		InitialNominalValue:     dbModel.InitialNominalValue,
		InitialNominalCurrency:  dbModel.InitialNominalCurrency,
		RegistrationDate:        dbModel.RegistrationDate,
		PlacementDate:           dbModel.PlacementDate,
		PlacementPrice:          dbModel.PlacementPrice,
		PlacementCurrency:       dbModel.PlacementCurrency,
		AccruedInterest:         dbModel.AccruedInterest,
		IssueSize:               int(dbModel.IssueSize),
		IssueSizePlan:           int(dbModel.IssueSizePlan),
		HasFloatingCoupon:       dbModel.HasFloatingCoupon,
		IsPerpetual:             dbModel.IsPerpetual,
		HasAmortization:         dbModel.HasAmortization,
		IsAvailableForIis:       dbModel.IsAvailableForIis,
		IsForQualifiedInvestors: dbModel.IsForQualifiedInvestors,
		IsSubordinated:          dbModel.IsSubordinated,
		RiskLevel:               bonds.RiskLevel(bonds.RiskLevel_value[dbModel.RiskLevel]),
		BondType:                bonds.BondType(bonds.BondType_value[dbModel.BondType]),
		CallOptionExerciseDate:  dbModel.CallOptionExerciseDate,
	}

	return domain
}

func mapTinkoffCouponToCoupon(figi string, tinkoffCoupon *pb.Coupon) bonds.Coupon {
	coupon := bonds.Coupon{
		Id:              uuid.New(),
		Figi:            figi,
		CouponDate:      tinkoffCoupon.CouponDate.AsTime(),
		CouponNumber:    int(tinkoffCoupon.CouponNumber),
		RecordDate:      tinkoffCoupon.FixDate.AsTime(),
		PerBondAmount:   tinkoffCoupon.GetPayOneBond().ToFloat(),
		CouponType:      bonds.CouponType(bonds.CouponType_value[pb.CouponType_name[int32(tinkoffCoupon.CouponType)]]),
		CouponStartDate: tinkoffCoupon.CouponStartDate.AsTime(),
		CouponEndDate:   tinkoffCoupon.CouponEndDate.AsTime(),
		CouponPeriod:    int(tinkoffCoupon.CouponPeriod),
	}
	return coupon
}

func mapCouponToDbModel(coupon bonds.Coupon) bondsdb.CouponDbModel {
	dbModel := bondsdb.CouponDbModel{
		Id:              coupon.Id,
		Figi:            coupon.Figi,
		CouponDate:      coupon.CouponDate,
		CouponNumber:    int64(coupon.CouponNumber),
		RecordDate:      coupon.RecordDate,
		PerBondAmount:   coupon.PerBondAmount,
		CouponType:      bonds.CouponType_name[int32(coupon.CouponType)],
		CouponStartDate: coupon.CouponStartDate,
		CouponEndDate:   coupon.CouponEndDate,
		CouponPeriod:    int64(coupon.CouponPeriod),
	}

	return dbModel
}

func mapCouponDbModelToDomain(coupon bondsdb.CouponDbModel) bonds.Coupon {
	domain := bonds.Coupon{
		Id:              coupon.Id,
		Figi:            coupon.Figi,
		CouponDate:      coupon.CouponDate,
		CouponNumber:    int(coupon.CouponNumber),
		RecordDate:      coupon.RecordDate,
		PerBondAmount:   coupon.PerBondAmount,
		CouponType:      bonds.CouponType(bonds.CouponType_value[coupon.CouponType]),
		CouponStartDate: coupon.CouponStartDate,
		CouponEndDate:   coupon.CouponEndDate,
		CouponPeriod:    int(coupon.CouponPeriod),
	}

	return domain
}
