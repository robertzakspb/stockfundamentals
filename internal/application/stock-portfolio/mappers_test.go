package portfolio

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/portfoliodb"
	"github.com/google/uuid"
)

func TestDbLotToLotMapping(t *testing.T) {
	id := uuid.New()
	createdAt := time.Now()
	updateAt := time.Now()
	figi := "abcdef"
	accountId := uuid.New()
	quantity := 10.0
	pricePerUnit := 25.2
	currency := "USD"

	dbLot := portfoliodb.LotDb{
		Id:           id,
		Figi:         figi,
		Ticker:       "test",
		CompanyName:  "test",
		AccountId:    accountId,
		CreatedAt:    createdAt,
		UpdatedAt:    updateAt,
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		Currency:     currency,
	}

	mappedLot := mapLotDbToLot(dbLot)

	if mappedLot.Id != id {
		t.Errorf("The specified lot's ID is incorrect")
	}
	if mappedLot.CreatedAt != createdAt {
		t.Errorf("The specified lot's created at is incorrect")
	}
	if mappedLot.UpdatedAt != updateAt {
		t.Errorf("The specified lot's ID update at incorrect")
	}
	if mappedLot.SecurityId != figi {
		t.Errorf("The specified lot's figi is incorrect")
	}
	if mappedLot.AccountId != accountId {
		t.Errorf("The specified lot's account ID is incorrect")
	}
	if mappedLot.Quantity != quantity {
		t.Errorf("The specified lot's quantity is incorrect")
	}
	if mappedLot.PricePerUnit != pricePerUnit {
		t.Errorf("The specified lot's price per unit is incorrect")
	}
	if mappedLot.Currency != currency {
		t.Errorf("The specified lot's currency is incorrect")
	}
}

func TestLotToDbLotMapping(t *testing.T) {
	id := uuid.New()
	createdAt := time.Now()
	updateAt := time.Now()
	figi := "abcdef"
	accountId := uuid.New()
	quantity := 10.0
	pricePerUnit := 25.2
	currency := "USD"

	sampleLot := lot.Lot{
		Id:           id,
		SecurityId:   figi,
		AccountId:    accountId,
		CreatedAt:    createdAt,
		UpdatedAt:    updateAt,
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		Currency:     currency,
	}

	mappedLot := mapLotToDbLot([]lot.Lot{sampleLot})[0]

	if mappedLot.Id != id {
		t.Errorf("The specified lot's ID is incorrect")
	}
	if mappedLot.CreatedAt != createdAt {
		t.Errorf("The specified lot's created at is incorrect")
	}
	if mappedLot.UpdatedAt != updateAt {
		t.Errorf("The specified lot's ID update at incorrect")
	}
	if mappedLot.Figi != figi {
		t.Errorf("The specified lot's figi is incorrect")
	}
	if mappedLot.AccountId != accountId {
		t.Errorf("The specified lot's account ID is incorrect")
	}
	if mappedLot.Quantity != quantity {
		t.Errorf("The specified lot's quantity is incorrect")
	}
	if mappedLot.PricePerUnit != pricePerUnit {
		t.Errorf("The specified lot's price per unit is incorrect")
	}
	if mappedLot.Currency != currency {
		t.Errorf("The specified lot's currency is incorrect")
	}
}
