package bondservice

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_mapBondToDbBond(t *testing.T) {
	id, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Unexpected error when initializing a new UUID")
	}
	figi := ""
	isin := ""
	lot := 5
	currency := "USD"
	name := "cool new issue"
	countryOfRisk := "RU"
	realExchange := "moex"
	couponCountPerYear := 4
	maturityDate := time.Now()
	nominalValue := 1000.0
	nominalCurrency := "USD"
	initialNominalValue := 2000.0
	initialNominalCurrency := "USD"
	registrationDate := time.Now()
	placementDate := time.Now()
	placementPrice := 10.43
	placementCurrency := "USD"
	accumulatedCouponIncome := 45.4
	issueSize := 5000
	issueSizePlan := 6000
	hasFloatingCoupon := false
	isPerpetual := false
	hasAmortization := false
	isAvailableForIis := true
	isForQualifiedInvestors := false
	isSubordinated := true
	riskLevel := "HIGH_RISK_LEVEL"
	bondType := "BOND_TYPE_REPLACED"
	callOptionExerciseDate := time.Now()

	bond := bonds.Bond{
		Id:                      id,
		Figi:                    figi,
		Isin:                    isin,
		Lot:                     lot,
		Currency:                currency,
		Name:                    name,
		CountryOfRisk:           countryOfRisk,
		RealExchange:            realExchange,
		CouponCountPerYear:      couponCountPerYear,
		MaturityDate:            maturityDate,
		NominalValue:            nominalValue,
		NominalCurrency:         nominalCurrency,
		InitialNominalValue:     initialNominalValue,
		InitialNominalCurrency:  initialNominalCurrency,
		RegistrationDate:        registrationDate,
		PlacementDate:           placementDate,
		PlacementPrice:          placementPrice,
		PlacementCurrency:       placementCurrency,
		AccumulatedCouponIncome: accumulatedCouponIncome,
		IssueSize:               issueSize,
		IssueSizePlan:           issueSizePlan,
		HasFloatingCoupon:       hasFloatingCoupon,
		IsPerpetual:             isPerpetual,
		HasAmortization:         hasAmortization,
		IsAvailableForIis:       isAvailableForIis,
		IsForQualifiedInvestors: isForQualifiedInvestors,
		IsSubordinated:          isSubordinated,
		RiskLevel:               bonds.RiskLevel(bonds.RiskLevel_value[riskLevel]),
		BondType:                bonds.BondType(bonds.BondType_value[bondType]),
		CallOptionExerciseDate:  callOptionExerciseDate,
	}

	mappedDomain := mapBondToDbBond(bond)

	test.AssertEqual(t, mappedDomain.Id, id)
	test.AssertEqual(t, mappedDomain.Figi, figi)
	test.AssertEqual(t, mappedDomain.Isin, isin)
	test.AssertEqual(t, mappedDomain.Lot, int64(lot))
	test.AssertEqual(t, mappedDomain.Currency, currency)
	test.AssertEqual(t, mappedDomain.Name, name)
	test.AssertEqual(t, mappedDomain.CountryOfRisk, countryOfRisk)
	test.AssertEqual(t, mappedDomain.RealExchange, realExchange)
	test.AssertEqual(t, mappedDomain.CouponCountPerYear, int64(couponCountPerYear))
	test.AssertEqual(t, mappedDomain.MaturityDate, maturityDate)
	test.AssertEqual(t, mappedDomain.NominalValue, nominalValue)
	test.AssertEqual(t, mappedDomain.NominalCurrency, nominalCurrency)
	test.AssertEqual(t, mappedDomain.InitialNominalValue, initialNominalValue)
	test.AssertEqual(t, mappedDomain.InitialNominalCurrency, initialNominalCurrency)
	test.AssertEqual(t, mappedDomain.RegistrationDate, registrationDate)
	test.AssertEqual(t, mappedDomain.PlacementDate, placementDate)
	test.AssertEqual(t, mappedDomain.PlacementPrice, placementPrice)
	test.AssertEqual(t, mappedDomain.AccumulatedCouponIncome, accumulatedCouponIncome)
	test.AssertEqual(t, mappedDomain.IssueSize, int64(issueSize))
	test.AssertEqual(t, mappedDomain.IssueSizePlan, int64(issueSizePlan))
	test.AssertEqual(t, mappedDomain.HasFloatingCoupon, hasFloatingCoupon)
	test.AssertEqual(t, mappedDomain.IsPerpetual, isPerpetual)
	test.AssertEqual(t, mappedDomain.HasAmortization, hasAmortization)
	test.AssertEqual(t, mappedDomain.IsAvailableForIis, isAvailableForIis)
	test.AssertEqual(t, mappedDomain.IsForQualifiedInvestors, isForQualifiedInvestors)
	test.AssertEqual(t, mappedDomain.IsSubordinated, isSubordinated)
	test.AssertEqual(t, mappedDomain.RiskLevel, riskLevel)
	test.AssertEqual(t, mappedDomain.BondType, bondType)
	test.AssertEqual(t, mappedDomain.CallOptionExerciseDate, callOptionExerciseDate)
}

func Test_MapBondToDbBond(t *testing.T) {
	id, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Unexpected error when initializing a new UUID")
	}
	figi := ""
	isin := ""
	lot := 5
	currency := "USD"
	name := "cool new issue"
	countryOfRisk := "RU"
	realExchange := "moex"
	couponCountPerYear := 4
	maturityDate := time.Now()
	nominalValue := 1000.0
	nominalCurrency := "USD"
	initialNominalValue := 2000.0
	initialNominalCurrency := "USD"
	registrationDate := time.Now()
	placementDate := time.Now()
	placementPrice := 10.43
	placementCurrency := "USD"
	accumulatedCouponIncome := 45.4
	issueSize := 5000
	issueSizePlan := 6000
	hasFloatingCoupon := false
	isPerpetual := false
	hasAmortization := false
	isAvailableForIis := true
	isForQualifiedInvestors := false
	isSubordinated := true
	riskLevel := "HIGH_RISK_LEVEL"
	bondType := "BOND_TYPE_REPLACED"
	callOptionExerciseDate := time.Now()

	dbModel := bondsdb.BondDbModel{
		Id:                      id,
		Figi:                    figi,
		Isin:                    isin,
		Lot:                     int64(lot),
		Currency:                currency,
		Name:                    name,
		CountryOfRisk:           countryOfRisk,
		RealExchange:            realExchange,
		CouponCountPerYear:      int64(couponCountPerYear),
		MaturityDate:            maturityDate,
		NominalValue:            nominalValue,
		NominalCurrency:         nominalCurrency,
		InitialNominalValue:     initialNominalValue,
		InitialNominalCurrency:  initialNominalCurrency,
		RegistrationDate:        registrationDate,
		PlacementDate:           placementDate,
		PlacementPrice:          placementPrice,
		PlacementCurrency:       placementCurrency,
		AccumulatedCouponIncome: accumulatedCouponIncome,
		IssueSize:               int64(issueSize),
		IssueSizePlan:           int64(issueSizePlan),
		HasFloatingCoupon:       hasFloatingCoupon,
		IsPerpetual:             isPerpetual,
		HasAmortization:         hasAmortization,
		IsAvailableForIis:       isAvailableForIis,
		IsForQualifiedInvestors: isForQualifiedInvestors,
		IsSubordinated:          isSubordinated,
		RiskLevel:               riskLevel,
		BondType:                bondType,
		CallOptionExerciseDate:  callOptionExerciseDate,
	}

	mappedDomain := mapDbBondToBond(dbModel)

	test.AssertEqual(t, mappedDomain.Id, id)
	test.AssertEqual(t, mappedDomain.Figi, figi)
	test.AssertEqual(t, mappedDomain.Isin, isin)
	test.AssertEqual(t, mappedDomain.Lot, lot)
	test.AssertEqual(t, mappedDomain.Currency, currency)
	test.AssertEqual(t, mappedDomain.Name, name)
	test.AssertEqual(t, mappedDomain.CountryOfRisk, countryOfRisk)
	test.AssertEqual(t, mappedDomain.RealExchange, realExchange)
	test.AssertEqual(t, mappedDomain.CouponCountPerYear, couponCountPerYear)
	test.AssertEqual(t, mappedDomain.MaturityDate, maturityDate)
	test.AssertEqual(t, mappedDomain.NominalValue, nominalValue)
	test.AssertEqual(t, mappedDomain.NominalCurrency, nominalCurrency)
	test.AssertEqual(t, mappedDomain.InitialNominalValue, initialNominalValue)
	test.AssertEqual(t, mappedDomain.InitialNominalCurrency, initialNominalCurrency)
	test.AssertEqual(t, mappedDomain.RegistrationDate, registrationDate)
	test.AssertEqual(t, mappedDomain.PlacementDate, placementDate)
	test.AssertEqual(t, mappedDomain.PlacementPrice, placementPrice)
	test.AssertEqual(t, mappedDomain.AccumulatedCouponIncome, accumulatedCouponIncome)
	test.AssertEqual(t, mappedDomain.IssueSize, issueSize)
	test.AssertEqual(t, mappedDomain.IssueSizePlan, issueSizePlan)
	test.AssertEqual(t, mappedDomain.HasFloatingCoupon, hasFloatingCoupon)
	test.AssertEqual(t, mappedDomain.IsPerpetual, isPerpetual)
	test.AssertEqual(t, mappedDomain.HasAmortization, hasAmortization)
	test.AssertEqual(t, mappedDomain.IsAvailableForIis, isAvailableForIis)
	test.AssertEqual(t, mappedDomain.IsForQualifiedInvestors, isForQualifiedInvestors)
	test.AssertEqual(t, mappedDomain.IsSubordinated, isSubordinated)
	test.AssertEqual(t, bonds.RiskLevel_name[int32(mappedDomain.RiskLevel)], riskLevel)
	test.AssertEqual(t, bonds.BondType_name[int32(mappedDomain.BondType)], bondType)
	test.AssertEqual(t, mappedDomain.CallOptionExerciseDate, callOptionExerciseDate)
}
