package bonds

import (
	"time"

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
