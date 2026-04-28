package divcalapi

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_mapPayoutToDto(t *testing.T) {
	id := uuid.New()
	divId := uuid.New()
	accountId := uuid.New()
	amount := 10.5
	figi := "testFigi"
	dps := 2.1
	announcementDate := time.Now().AddDate(0, 0, -7)
	recordDate := time.Now().AddDate(0, 0, 10)
	payoutDate := time.Now().AddDate(0, 0, 20)

	payout := Payout{
		Id:         id,
		DividendId: divId,
		AccountId:  accountId,
		Amount:     amount,
		Dividend: dividend.Dividend{
			Figi:             figi,
			ActualDPS:        dps,
			AnnouncementDate: announcementDate,
			RecordDate:       recordDate,
			PayoutDate:       payoutDate,
		},
	}

	mappedDto := mapPayoutToDto(payout)

	test.AssertEqual(t, mappedDto.Id, id)
	test.AssertEqual(t, mappedDto.DividendId, divId)
	test.AssertEqual(t, mappedDto.AccountId, accountId)
	test.AssertEqual(t, mappedDto.Amount, amount)
	test.AssertEqual(t, mappedDto.Figi, figi)
	test.AssertEqual(t, mappedDto.ActualDPS, dps)
	test.AssertEqual(t, mappedDto.AnnouncementDate, announcementDate)
	test.AssertEqual(t, mappedDto.RecordDate, recordDate)
	test.AssertEqual(t, mappedDto.PayoutDate, payoutDate)
}

func Test_mapDivCalToDto(t *testing.T) {
	id := uuid.New()
	divId := uuid.New()
	accountId := uuid.New()
	amount := 10.5
	figi := "testFigi"
	dps := 2.1
	announcementDate := time.Now().AddDate(0, 0, -7)
	recordDate := time.Now().AddDate(0, 0, 10)
	payoutDate := time.Now().AddDate(0, 0, 20)

	payout := Payout{
		Id:         id,
		DividendId: divId,
		AccountId:  accountId,
		Amount:     amount,
		Dividend: dividend.Dividend{
			Figi:             figi,
			ActualDPS:        dps,
			AnnouncementDate: announcementDate,
			RecordDate:       recordDate,
			PayoutDate:       payoutDate,
		},
	}

	divCalendar := DividendCalendar{
		AccountIds:    []uuid.UUID{accountId},
		FuturePayouts: []dividend.Payout{dividend.Payout(payout)},
	}

	mappedDto := mapDivCalToDto(divCalendar)

	test.AssertEqual(t, 1, len(mappedDto.AccountIds))
	test.AssertEqual(t, 1, len(mappedDto.FuturePayouts))
	test.AssertEqual(t, accountId, mappedDto.AccountIds[0])
	test.AssertEqual(t, mappedDto.FuturePayouts[0].Id, id)
	test.AssertEqual(t, mappedDto.FuturePayouts[0].DividendId, divId)
	test.AssertEqual(t, mappedDto.FuturePayouts[0].AccountId, accountId)
	test.AssertEqual(t, mappedDto.FuturePayouts[0].Amount, amount)
	test.AssertEqual(t, mappedDto.FuturePayouts[0].Figi, figi)
	test.AssertEqual(t, mappedDto.FuturePayouts[0].ActualDPS, dps)
	test.AssertEqual(t, mappedDto.FuturePayouts[0].AnnouncementDate, announcementDate)
	test.AssertEqual(t, mappedDto.FuturePayouts[0].RecordDate, recordDate)
	test.AssertEqual(t, mappedDto.FuturePayouts[0].PayoutDate, payoutDate)

}
