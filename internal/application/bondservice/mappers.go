package bondservice

import (
	"strings"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	pb "opensource.tbank.ru/invest/invest-go/proto"
)

func mapTinkoffBondsToBonds(tinkoffBonds []*pb.Bond) []bonds.Bond {
	bondList := []bonds.Bond{}

	for _, tinkoffBond := range tinkoffBonds {
		if tinkoffBond == nil ||
			tinkoffBond.MaturityDate.AsTime().Before(time.Now()) { //No need to import historical bonds that have matured
			continue
		}

		bond := bonds.Bond{
			Id:                      uuid.New(),
			Figi:                    tinkoffBond.Figi,
			Isin:                    tinkoffBond.Isin,
			Lot:                     int(tinkoffBond.Lot),
			Currency:                strings.ToUpper(tinkoffBond.Currency),
			Name:                    tinkoffBond.Name,
			CountryOfRisk:           tinkoffBond.CountryOfRisk,
			RealExchange:            tinkoffBond.RealExchange.String(),
			CouponCountPerYear:      int(tinkoffBond.CouponQuantityPerYear),
			MaturityDate:            tinkoffBond.MaturityDate.AsTime(),
			NominalValue:            tinkoffBond.Nominal.ToFloat(),
			NominalCurrency:         tinkoffBond.Nominal.GetCurrency(),
			InitialNominalValue:     tinkoffBond.InitialNominal.ToFloat(),
			InitialNominalCurrency:  strings.ToUpper(tinkoffBond.InitialNominal.GetCurrency()),
			RegistrationDate:        tinkoffBond.StateRegDate.AsTime(),
			PlacementDate:           tinkoffBond.PlacementDate.AsTime(),
			PlacementPrice:          tinkoffBond.PlacementPrice.ToFloat(),
			PlacementCurrency:       strings.ToUpper(tinkoffBond.PlacementPrice.GetCurrency()),
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
		validationErr := bond.Validate()
		if validationErr != nil {
			logger.Log(validationErr.Error(), logger.WARNING)
		}
		bondList = append(bondList, bond)
	}
	return bondList
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

func mapBondsToDbBonds(bondList []bonds.Bond) []bondsdb.BondDbModel {
	dbBonds := make([]bondsdb.BondDbModel, len(bondList))

	for i, bond := range bondList {
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
			NominalCurrency:         strings.ToUpper(bond.NominalCurrency),
			InitialNominalValue:     bond.InitialNominalValue,
			InitialNominalCurrency:  strings.ToUpper(bond.InitialNominalCurrency),
			RegistrationDate:        bond.RegistrationDate,
			PlacementDate:           bond.PlacementDate,
			PlacementPrice:          bond.PlacementPrice,
			PlacementCurrency:       strings.ToUpper(bond.PlacementCurrency),
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
		dbBonds[i] = dbBond
	}

	return dbBonds
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

func mapTinkoffCouponsToCoupons(tinkoffCoupons []*pb.Coupon) []bonds.Coupon {
	coupons := []bonds.Coupon{}

	for _, tinkoffCoupon := range tinkoffCoupons {
		if tinkoffCoupon == nil {
			continue
		}
		coupon := bonds.Coupon{
			Id:              uuid.New(),
			Figi:            tinkoffCoupon.Figi,
			CouponDate:      tinkoffCoupon.CouponDate.AsTime(),
			CouponNumber:    int(tinkoffCoupon.CouponNumber),
			RecordDate:      tinkoffCoupon.FixDate.AsTime(),
			PerBondAmount:   tinkoffCoupon.GetPayOneBond().ToFloat(),
			CouponType:      bonds.CouponType(bonds.CouponType_value[pb.CouponType_name[int32(tinkoffCoupon.CouponType)]]),
			CouponStartDate: tinkoffCoupon.CouponStartDate.AsTime(),
			CouponEndDate:   tinkoffCoupon.CouponEndDate.AsTime(),
			CouponPeriod:    int(tinkoffCoupon.CouponPeriod),
		}
		coupons = append(coupons, coupon)
	}
	return coupons
}

func mapCouponsToDbModels(coupons []bonds.Coupon) []bondsdb.CouponDbModel {
	dbModels := make([]bondsdb.CouponDbModel, len(coupons))
	
	for i, coupon := range coupons {
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
		dbModels[i] = dbModel
	}

	return dbModels
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
