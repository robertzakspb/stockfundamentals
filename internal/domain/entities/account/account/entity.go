package account

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id              uuid.UUID
	OpeningDate     time.Time
	Type            string
	Broker          string
	Holder          string
	PrimaryCurrency string
}

type AccountType string

const (
	STANDARD AccountType = "STANDARD"
	IIS_1    AccountType = "IIS_1"
	IIS_2    AccountType = "IIS_2"
	IIS_3    AccountType = "IIS_3"
)

var (
	AccountType_Map = map[string]AccountType{
		"STANDARD": STANDARD,
		"IIS_1":    IIS_1,
		"IIS_2":    IIS_2,
		"IIS_3":    IIS_3,
	}
)
