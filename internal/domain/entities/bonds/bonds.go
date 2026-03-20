package bonds

import (
	"errors"
	"strconv"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/compoundinterest"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/forex"
	"github.com/google/uuid"
)

type Bond struct {
	Id                      uuid.UUID
	Figi                    string
	Isin                    string
	Lot                     int
	Currency                string
	Name                    string
	CountryOfRisk           string
	RealExchange            string
	CouponCountPerYear      int
	MaturityDate            time.Time
	NominalValue            float64
	NominalCurrency         string
	InitialNominalValue     float64
	InitialNominalCurrency  string
	RegistrationDate        time.Time
	PlacementDate           time.Time
	PlacementPrice          float64
	PlacementCurrency       string
	AccumulatedCouponIncome float64
	IssueSize               int
	IssueSizePlan           int
	HasFloatingCoupon       bool
	IsPerpetual             bool
	HasAmortization         bool
	IsAvailableForIis       bool
	IsForQualifiedInvestors bool
	IsSubordinated          bool
	RiskLevel               RiskLevel
	BondType                BondType
	CallOptionExerciseDate  time.Time
}

type RiskLevel int

const (
	UNSPECIFIED_RISK_LEVEL RiskLevel = 0 //Не указан.
	LOW_RISK_LEVEL         RiskLevel = 1 //Низкий уровень риска.
	MODERATE_RISK_LEVEL    RiskLevel = 2 //Средний уровень риска.
	HIGH_RISK_LEVEL        RiskLevel = 3 //Высокий уровень риска.
)

// Enum value maps for RiskLevel.
var (
	RiskLevel_name = map[int32]string{
		0: "UNSPECIFIED_RISK_LEVEL",
		1: "LOW_RISK_LEVEL",
		2: "MODERATE_RISK_LEVEL",
		3: "HIGH_RISK_LEVEL",
	}
	RiskLevel_value = map[string]int32{
		"UNSPECIFIED_RISK_LEVEL": 0,
		"LOW_RISK_LEVEL":         1,
		"MODERATE_RISK_LEVEL":    2,
		"HIGH_RISK_LEVEL":        3,
	}
)

type BondType int

const (
	BondType_BOND_TYPE_UNSPECIFIED BondType = 0 // Тип облигации не определен.
	BondType_BOND_TYPE_REPLACED    BondType = 1 // Замещающая облигация.
)

// Enum value maps for BondType.
var (
	BondType_name = map[int32]string{
		0: "BOND_TYPE_UNSPECIFIED",
		1: "BOND_TYPE_REPLACED",
	}
	BondType_value = map[string]int32{
		"BOND_TYPE_UNSPECIFIED": 0,
		"BOND_TYPE_REPLACED":    1,
	}
)

func (b Bond) validate() error {
	if b.Id == uuid.Nil {
		return errors.New("Nil Id in the bond")
	}
	if b.Figi == "" {
		return errors.New("Missing figi in the bond")
	}
	if b.Isin == "" {
		return errors.New("Missing ISIN in the bond")
	}
	if b.Lot <= 0 {
		return errors.New("Invalid lot value for bond: " + strconv.Itoa(b.Lot))
	}
	forexDP := forex.ForexDP{}
	if b.Currency == "" || !forexDP.IsSupportedCurrency(b.Currency) {
		return errors.New("Missing or unsupported currency " + b.Currency)
	}
	if b.CouponCountPerYear <= 0 {
		return errors.New("Invalid coupon count for the bond: " + strconv.Itoa(b.CouponCountPerYear))
	}
	if b.MaturityDate.IsZero() {
		return errors.New("Invalid maturity date for the bond: " + b.MaturityDate.String())
	}
	if b.NominalValue <= 0 {
		return errors.New("Invalid nominal value for the bond")
	}
	if b.NominalCurrency == "" || !forexDP.IsSupportedCurrency(b.NominalCurrency) {
		return errors.New("Missing or unsupported nominal currency " + b.NominalCurrency)
	}
	if b.InitialNominalValue <= 0 {
		return errors.New("Invalid initial nominal value for the bond")
	}
	if b.InitialNominalCurrency == "" || !forexDP.IsSupportedCurrency(b.NominalCurrency) {
		return errors.New("Missing or unsupported initial nominal currency " + b.InitialNominalCurrency)
	}
	if b.PlacementPrice <= 0.0 {
		return errors.New("Invalid placement price ")
	}
	if b.PlacementCurrency == "" || !forexDP.IsSupportedCurrency(b.NominalCurrency) {
		return errors.New("Missing or unsupported placement currency " + b.NominalCurrency)
	}
	if b.AccumulatedCouponIncome <= 0 {
		return errors.New("Invalud accumulated coupon value")
	}
	if b.IssueSize <= 0 {
		return errors.New("Invalid issue size: " + strconv.Itoa(b.IssueSize))
	}
	if b.IssueSizePlan <= 0 {
		return errors.New("Invalid issue size plan: " + strconv.Itoa(b.IssueSizePlan))
	}
	_, found := RiskLevel_name[int32(b.RiskLevel)]
	if !found {
		return errors.New("Unsupported risk level: " + RiskLevel_name[int32(b.RiskLevel)])
	}
	_, found = BondType_name[int32(b.BondType)]
	if !found {
		return errors.New("Unsupported bond type: " + BondType_name[int32(b.BondType)])
	}

	return nil
}

func (b Bond) YieldToMaturity(coupons []Coupon, marketPrice float64) (float64, error) {
	yield, err := calculateYield(b, coupons, marketPrice, time.Now(), b.MaturityDate)
	if err != nil {
		return -1, err
	}

	return yield, nil
}

func (b Bond) YieldToCallOption(coupons []Coupon, marketPrice float64) (float64, error) {
	if b.CallOptionExerciseDate.IsZero() {
		return -1, errors.New("Attempting to calculate a yield to call option for a bond without a call exercise date")
	}

	yield, err := calculateYield(b, coupons, marketPrice, time.Now(), b.CallOptionExerciseDate)
	if err != nil {
		return -1, err
	}

	return yield, nil
}

// Calculates the return realized on the bond given a market price, including coupons and redemption
// Coupon reinvestment is not assumed
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func totalBondReturn(bond Bond, coupons []Coupon, marketPrice float64, acquisitionDate, redemptionDate time.Time) (float64, error) {
	if marketPrice == 0 {
		return -1, errors.New("Invalid market price")
	}

	futureCashflows := totalCouponIncome(coupons, false) + bond.NominalValue

	cumulativeReturn := futureCashflows/marketPrice - 1
	totalReturn := compoundinterest.CalcAnnualizedReturn(cumulativeReturn, acquisitionDate, redemptionDate)

	return totalReturn, nil
}

func calculateYield(b Bond, coupons []Coupon, marketPrice float64, acquisitionDate, redemptionDate time.Time) (float64, error) {
	if len(coupons) == 0 {
		return -1, errors.New("Failed to calculate the yield due to missing couponsa~Z``````````````````")
	}

	if !(coupons[0].CouponType == CouponType_COUPON_TYPE_FIX || coupons[0].CouponType == CouponType_COUPON_TYPE_CONSTANT) {
		return -1, errors.New("Unable to calculate the YTM for non-fixed and non-constant coupons")
	}
	holdingPeriod := redemptionDate.Sub(acquisitionDate).Hours() / 24

	yield := (b.NominalValue - marketPrice + totalCouponIncome(coupons, false)) / marketPrice * 365 / holdingPeriod * 100
	return yield, nil
}
