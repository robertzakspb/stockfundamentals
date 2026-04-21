package bondsapi

import "time"

type BondDTO struct {
	Figi                    string      `json:"figi"`
	Isin                    string      `json:"isin"`
	Lot                     int         `json:"lot"`
	Currency                string      `json:"currency"`
	Name                    string      `json:"name"`
	CountryOfRisk           string      `json:"countryOfRisk"`
	RealExchange            string      `json:"realExchange"`
	CouponCountPerYear      int         `json:"annualCouponCount"`
	MaturityDate            time.Time   `json:"maturityDate"`
	NominalValue            float64     `json:"nominalValue"`
	NominalCurrency         string      `json:"nominalCurrency"`
	InitialNominalValue     float64     `json:"initialNominalValue"`
	InitialNominalCurrency  string      `json:"initialNominalCurrency"`
	RegistrationDate        time.Time   `json:"registrationDate"`
	PlacementDate           time.Time   `json:"placementDAte"`
	PlacementPrice          float64     `json:"placementPrice"`
	PlacementCurrency       string      `json:"placementCurrency"`
	AccruedInterest         float64     `json:"accruedInterest"`
	IssueSize               int         `json:"issueSize"`
	IssueSizePlan           int         `json:"issueSizePlan"`
	HasFloatingCoupon       bool        `json:"hasFloatingCoupon"`
	IsPerpetual             bool        `json:"isPerpetual"`
	HasAmortization         bool        `json:"hasAmortization"`
	IsAvailableForIis       bool        `json:"isAvaialbleForIis"`
	IsForQualifiedInvestors bool        `json:"qualifiedInvestorsOnly"`
	IsSubordinated          bool        `json:"isSubordinated"`
	RiskLevel               string      `json:"riskLevel"`
	BondType                string      `json:"bondType"`
	CallOptionExerciseDate  time.Time   `json:"callOptionExerciseDate"`
	YieldToMaturity         float64     `json:"yieldToMaturity"`
	YieldToCallOption       float64     `json:"yieldToCallOption"`
	Coupons                 []CouponDTO `json:"coupons"`
}

type CouponDTO struct {
	Figi               string    `json:"figi"`
	Date               time.Time `json:"date"`
	CouponNumber       int       `json:"couponNumber"`
	RecordDate         time.Time `json:"recordDate"`
	PerBondAmount      float64   `json:"perBondAmount"`
	CouponType         string    `json:"couponType"`
	CouponStartDate    time.Time `json:"couponStartDate"`
	CouponEndDate      time.Time `json:"couponEndDate"`
	CouponPeriodInDays int       `json:"couponPeriodInDays"`
}
