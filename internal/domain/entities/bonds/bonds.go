package bonds

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
	"github.com/google/uuid"
)

type Bond struct {
	Id                      uuid.UUID
	Figi                    string
	Isin                    string
	Ticker                  string
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
	AccruedInterest         float64
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
	Coupons                 []Coupon
	SimpleYieldToMaturity   float64
	SimpleYieldToCallOption float64
	YieldTomaturity         float64
	YieldToCallOption       float64
	MarketValueInRUB        float64
	QuoteInPercentage       float64
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
	BondType_BOND_TYPE_UNSPECIFIED BondType = 0
	BondType_BOND_TYPE_REPLACED    BondType = 1 // Replaced bonds (2022)
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

func (b Bond) Validate() error {
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
	if b.Currency == "" || !forexservice.IsSupportedCurrency(b.Currency) {
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
	if b.NominalCurrency == "" || !forexservice.IsSupportedCurrency(b.NominalCurrency) {
		return errors.New("Missing or unsupported nominal currency " + b.NominalCurrency)
	}
	if b.InitialNominalValue <= 0 {
		return errors.New("Invalid initial nominal value for the bond")
	}
	if b.InitialNominalCurrency == "" || !forexservice.IsSupportedCurrency(b.InitialNominalCurrency) {
		return errors.New("Missing or unsupported initial nominal currency " + b.InitialNominalCurrency)
	}
	if b.PlacementPrice <= 0.0 {
		return errors.New("Invalid placement price ")
	}
	if b.PlacementCurrency == "" || !forexservice.IsSupportedCurrency(b.PlacementCurrency) {
		return errors.New("Missing or unsupported placement currency " + b.NominalCurrency)
	}
	if b.AccruedInterest < 0 {
		return errors.New("Invalid accumulated coupon value")
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

func (b *Bond) HasCallOption() bool {
	return !b.CallOptionExerciseDate.IsZero()
}

func (b *Bond) IsRubleBond() bool {
	return strings.ToUpper(b.Currency) == "RUB" && b.Currency == b.NominalCurrency
}

func (b *Bond) IsBondWithDifferentNominalCurrencyAndCurrency() bool {
	return b.NominalCurrency != b.Currency
}

func (b *Bond) MarketValue(quoteAsPercentage, fxRate float64) float64 {
	mv := b.MarketPriceInCurrency(quoteAsPercentage) + b.AccruedInterest
	mvInCurrency := mv * fxRate
	return mvInCurrency
}

func (b *Bond) MarketPriceInCurrency(quoteAsPercentage float64) float64 {
	marketPriceInCurrency := quoteAsPercentage * b.NominalValue / 100
	return marketPriceInCurrency
}

func (b *Bond) HasFixedCoupon() (bool, error) {
	if len(b.Coupons) < 1 {
		return false, errors.New("Unable to determine whether a bond has a fixed coupon due to missing coupons")
	}

	isFixedCoupon := b.Coupons[0].CouponType == CouponType_COUPON_TYPE_FIX || b.Coupons[0].CouponType == CouponType_COUPON_TYPE_CONSTANT
	return isFixedCoupon, nil
}

func (b *Bond) CurrentCouponYield() float64 {
	if len(b.Coupons) == 0 || b.MarketPriceInCurrency(b.QuoteInPercentage) == 0 {
		return 0
	}
	hasFixedCoupon, err := b.HasFixedCoupon()
	if err != nil || !hasFixedCoupon {
		return 0
	}
	currentCouponYield := float64(b.CouponCountPerYear) * b.Coupons[0].PerBondAmount / b.MarketPriceInCurrency(b.QuoteInPercentage)
	return currentCouponYield
}

// Returns the sum total of all future coupon payments.
func (b *Bond) TotalFutureCoupons(onlyTillCallOptionExerciseDate bool) (float64, error) {
	if len(b.Coupons) == 0 {
		return -1, errors.New("Unable to calculate cumulative cashflows due to missing coupons")
	}
	isFixed, err := b.HasFixedCoupon()
	if err != nil {
		return -1, err
	}
	if !isFixed {
		return -1, errors.New("Cashflows can be calculated only for fixed or constant coupons")
	}

	var tillDate time.Time
	if onlyTillCallOptionExerciseDate {
		tillDate = b.CallOptionExerciseDate
	} else {
		tillDate = b.MaturityDate
	}
	futureCoupons := TotalCouponIncome(b.Coupons, false, tillDate)

	return futureCoupons, nil
}

// Returns the sum total of all future cash flows of a bond if it were purchased today, including redemption, adjusted for the accrued interest for maximum accuracy
func (b *Bond) TotalFutureCashflows(onlyTillCallOptionExerciseDate bool, quote float64) (float64, error) {
	futureCashFlows := 0.0
	futureCashFlows += -(quote + b.AccruedInterest) //The full market price of the bond

	futureCoupons, err := b.TotalFutureCoupons(onlyTillCallOptionExerciseDate)
	if err != nil {
		return futureCoupons, err
	}
	futureCashFlows += futureCoupons

	futureCashFlows += b.NominalValue

	return futureCashFlows, nil
}

func (b *Bond) FutureCouponsTillDate(tillDate time.Time) ([]Coupon, error) {
	if len(b.Coupons) == 0 {
		return []Coupon{}, errors.New("Unable to get future coupons due to missing coupons")
	}
	coupons := []Coupon{}

	for i := range b.Coupons {
		cd := b.Coupons[i].CouponDate
		//Skipping the past coupons and coupons past the provided latest date
		if timehelpers.DateIsEarlierOrSameDate(cd, time.Now()) || timehelpers.DateIsLater(cd, tillDate) {

			continue
		}
		coupons = append(coupons, b.Coupons[i])
	}
	return coupons, nil
}

// Returns all cashflows of a bond, starting from its acquisition on the provided date and ending with redemption on the maturity date or call option exercise date
func (b *Bond) FutureCashflowsWithDates(onlyTillCallOptionExerciseDate bool, quote float64, acquisitionDate time.Time) (cashflows []float64, dates []time.Time, err error) {
	cashflows = append(cashflows, -(quote + b.AccruedInterest)) //Full market price of the bond
	dates = append(dates, acquisitionDate)                                 //THe date on which the bond is to be acquired

	//Adding each future coupon to the cashflow stream
	var tillDate time.Time
	if onlyTillCallOptionExerciseDate {
		tillDate = b.CallOptionExerciseDate
	} else {
		tillDate = b.MaturityDate
	}
	coupons, err := b.FutureCouponsTillDate(tillDate)
	if err != nil {
		return cashflows, dates, err
	}
	for i := range coupons {
		cashflows = append(cashflows, coupons[i].PerBondAmount)
		dates = append(dates, coupons[i].CouponDate)
	}

	//Adding the redemption
	cashflows = append(cashflows, b.NominalValue)
	dates = append(dates, tillDate)

	return cashflows, dates, nil
}
