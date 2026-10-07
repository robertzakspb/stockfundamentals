package security_master

import (
	"strings"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
)

func mapDbBondToBond(dbModel bondsdb.BondDbModel) bonds.Bond {
	domain := bonds.Bond{
		Figi:                    dbModel.Figi,
		Id:                      dbModel.Id,
		Isin:                    dbModel.Isin,
		Ticker:                  dbModel.Ticker,
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

func MapBondsToDbBonds(bondList []bonds.Bond) []bondsdb.BondDbModel {
	dbBonds := make([]bondsdb.BondDbModel, len(bondList))

	for i, bond := range bondList {
		dbBond := bondsdb.BondDbModel{
			Id:                      bond.Id,
			Figi:                    bond.Figi,
			Isin:                    bond.Isin,
			Ticker:                  bond.Ticker,
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
