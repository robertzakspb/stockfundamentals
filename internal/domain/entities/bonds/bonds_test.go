package bonds

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_validate_ProperBond(t *testing.T) {
	validBond := generateMockValidBond()

	err := validBond.Validate()

	test.AssertNoError(t, err)
}

func Test_Validate_NilId(t *testing.T) {
	bond := generateMockValidBond()
	bond.Id = uuid.Nil

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_MissingFigi(t *testing.T) {
	bond := generateMockValidBond()
	bond.Figi = ""

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_MissingIsin(t *testing.T) {
	bond := generateMockValidBond()
	bond.Isin = ""

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidLot(t *testing.T) {
	bond := generateMockValidBond()
	bond.Lot = 0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidCurrency(t *testing.T) {
	bond := generateMockValidBond()
	bond.Currency = "ABC"

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidCouponCount(t *testing.T) {
	bond := generateMockValidBond()
	bond.CouponCountPerYear = -1

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidMaturityDate(t *testing.T) {
	bond := generateMockValidBond()
	bond.MaturityDate = time.Time{}

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidNominalValue(t *testing.T) {
	bond := generateMockValidBond()
	bond.NominalValue = 0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidNominalCurrency(t *testing.T) {
	bond := generateMockValidBond()
	bond.NominalCurrency = "WER"

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidInitialNominalValue(t *testing.T) {
	bond := generateMockValidBond()
	bond.InitialNominalValue = 0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidInitialNominalCurrency(t *testing.T) {
	bond := generateMockValidBond()
	bond.InitialNominalCurrency = "TEST"

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidPlacementPrice(t *testing.T) {
	bond := generateMockValidBond()
	bond.PlacementPrice = 0.0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidPlacementCurrency(t *testing.T) {
	bond := generateMockValidBond()
	bond.PlacementCurrency = "PPP"

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidAccruedInterest(t *testing.T) {
	bond := generateMockValidBond()
	bond.AccruedInterest = -1

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidIssueSize(t *testing.T) {
	bond := generateMockValidBond()
	bond.IssueSize = 0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidIssueSizePlan(t *testing.T) {
	bond := generateMockValidBond()
	bond.IssueSizePlan = 0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_HasCallOption_Negative(t *testing.T) {
	bond := generateMockValidBond()
	bond.CallOptionExerciseDate = time.Time{}

	test.AssertFalse(t, bond.HasCallOption())
}

func Test_HasCallOption_Positive(t *testing.T) {
	bond := generateMockValidBond()
	bond.CallOptionExerciseDate = time.Now()

	test.AssertTrue(t, bond.HasCallOption())
}

func Test_IsRubleBond_Negative(t *testing.T) {
	bond := generateMockValidBond()
	bond.Currency = "RUB"
	bond.NominalCurrency = "USD"

	test.AssertFalse(t, bond.IsRubleBond())
}

func Test_IsRubleBond_Positive(t *testing.T) {
	bond := generateMockValidBond()
	bond.Currency = "RUB"
	bond.NominalCurrency = "RUB"

	test.AssertTrue(t, bond.IsRubleBond())
	
}

func generateMockValidBond() Bond {
	bond := Bond{
		Id:                     uuid.New(),
		Figi:                   "testFigi",
		Isin:                   "testIsin",
		Lot:                    10,
		Currency:               "USD",
		CouponCountPerYear:     10,
		MaturityDate:           time.Now(),
		NominalValue:           1000,
		NominalCurrency:        "EUR",
		InitialNominalValue:    1000,
		InitialNominalCurrency: "EUR",
		PlacementPrice:         1005,
		PlacementCurrency:      "EUR",
		AccruedInterest:        10,
		IssueSize:              1_000_000,
		IssueSizePlan:          5_000_000,
		RiskLevel:              HIGH_RISK_LEVEL,
		BondType:               BondType_BOND_TYPE_UNSPECIFIED,
	}

	return bond
}
