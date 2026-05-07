package account

import (
	"errors"
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

func (a Account) validate() error {
	if a.Id == uuid.Nil {
		return errors.New("The account is missing an ID")
	}
	if a.OpeningDate.After(time.Now()) {
		return errors.New("The account's opening date cannot be in the future")
	}
	if a.Type == "" {
		return errors.New("The account is missing its type")
	}
	if a.Broker == "" {
		return errors.New("The account is missing the associated brokerage")
	}
	if a.Holder == "" {
		return errors.New("The account is lacking its holder information")
	}
	if a.PrimaryCurrency == "" {
		return errors.New("The account is missing the primary currency")
	}

	return nil
}
