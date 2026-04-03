package bondsdb

import (
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type BondDbModel struct {
	Id                      uuid.UUID `sql:"id"`
	Figi                    string    `sql:"figi"`
	Isin                    string    `sql:"isin"`
	Lot                     int64     `sql:"lot"`
	Currency                string    `sql:"currency"`
	Name                    string    `sql:"name"`
	CountryOfRisk           string    `sql:"country_of_risk"`
	RealExchange            string    `sql:"real_exchange"`
	CouponCountPerYear      int64     `sql:"coupon_count_per_year"`
	MaturityDate            time.Time `sql:"maturity_date"`
	NominalValue            float64   `sql:"nominal_value"`
	NominalCurrency         string    `sql:"nominal_currency"`
	InitialNominalValue     float64   `sql:"initial_nominal_value"`
	InitialNominalCurrency  string    `sql:"initial_nominal_currency"`
	RegistrationDate        time.Time `sql:"registration_date"`
	PlacementDate           time.Time `sql:"placement_date"`
	PlacementPrice          float64   `sql:"placement_price"`
	PlacementCurrency       string    `sql:"placement_currency"`
	AccruedInterest         float64   `sql:"accumulated_coupon_income"`
	IssueSize               int64     `sql:"issue_size"`
	IssueSizePlan           int64     `sql:"issue_size_plan"`
	HasFloatingCoupon       bool      `sql:"has_floating_coupon"`
	IsPerpetual             bool      `sql:"is_perpetual"`
	HasAmortization         bool      `sql:"has_amortization"`
	IsAvailableForIis       bool      `sql:"is_available_for_iis"`
	IsForQualifiedInvestors bool      `sql:"is_for_qualified_investors"`
	IsSubordinated          bool      `sql:"is_subordinated"`
	RiskLevel               string    `sql:"risk_level"`
	BondType                string    `sql:"bond_type"`
	CallOptionExerciseDate  time.Time `sql:"call_option_exercise_date"`
}

/*
Perpetual bonds returned by the Tinkoff API tend to have the maturity date's year set to 2111.
The problem is that YDB's Date type's max value is 00:00 01.01.2106, meaning that supplying 2111 leads to a DB-level error while saving.
The workaround is to save NULL in the DB if the bond is perpetual, which should not impact the accuracy,
as technically perpetual bonds do not have any
maturity date whatsoever.
read more on YDB Date constraints: https://ydb.tech/docs/en/yql/reference/types/primitive?version=v25.2#datetime
*/
func (bond BondDbModel) CorrectMaturityDate() types.Value {
	if bond.IsPerpetual {
		return types.NullValue(types.TypeDate)
	}

	if bond.MaturityDate.Year() > 2105 {
		maximumAllowedYdbValue, _ := time.Parse(time.DateOnly, "2105-12-31")
		return db.ConvertToYdbDate(maximumAllowedYdbValue)
	}

	return db.ConvertToYdbDate(bond.MaturityDate)
}

type CouponDbModel struct {
	Id              uuid.UUID `sql:"id"`
	Figi            string    `sql:"figi"`
	CouponDate      time.Time `sql:"coupon_date"`
	CouponNumber    int64     `sql:"coupon_number"`
	RecordDate      time.Time `sql:"record_date"`
	PerBondAmount   float64   `sql:"per_bond_amount"`
	CouponType      string    `sql:"coupon_type"`
	CouponStartDate time.Time `sql:"coupon_start_date"`
	CouponEndDate   time.Time `sql:"coupon_end_date"`
	CouponPeriod    int64     `sql:"coupon_period"`
}

type BondPositionLotDb struct {
	Id                uuid.UUID `sql:"id"`
	Figi              string    `sql:"figi"`
	Isin              string    `sql:"isin"`
	OpeningDate       time.Time `sql:"opening_date"`
	ModificationDate  time.Time `sql:"modification_date"`
	AccountId         uuid.UUID `sql:"account_id"`
	Quantity          float64   `sql:"quantity"`
	PricePerUnit      float64   `sql:"price_per_unit"`
	PricePerUnitInRUB float64   `sql:"price_per_unit_rub"`
	AccruedInterest   float64   `sql:"accumulated_coupon_income"`
}
