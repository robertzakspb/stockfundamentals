package divcalapi

import (
	"time"

	dividendcalendar "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/dividend-calendar"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/google/uuid"
)

type DividendCalendar dividendcalendar.DividendCalendar
type Payout dividend.Payout

type DividendCalendarDto struct {
	AccountIds    uuid.UUIDs  `json:"accountIds"`
	FuturePayouts []PayoutDto `json:"futurePayouts"`
}

type PayoutDto struct {
	Id               uuid.UUID `json:"id"`
	DividendId       uuid.UUID `json:"dividendId"`
	AccountId        uuid.UUID `json:"accountId"`
	Amount           float64   `json:"amount"`
	Figi             string    `json:"figi"`
	Isin             string    `json:"isin"`
	Ticker           string    `json:"ticker"`
	ActualDPS        float64   `json:"actualDPS"`
	AnnouncementDate time.Time `json:"announcementDate"`
	RecordDate       time.Time `json:"recordDate"`
	PayoutDate       time.Time `json:"payoutDate"`
}

func mapDivCalToDto(divcal DividendCalendar) DividendCalendarDto {
	payoutDtos := make([]PayoutDto, len(divcal.FuturePayouts))
	for i, payout := range divcal.FuturePayouts {
		dto := mapPayoutToDto(Payout(payout))
		payoutDtos[i] = dto
	}
	divCalDto := DividendCalendarDto{
		AccountIds:    divcal.AccountIds,
		FuturePayouts: payoutDtos,
	}

	return divCalDto
}

func mapPayoutToDto(payout Payout) PayoutDto {
	dto := PayoutDto{
		Id:               payout.Id,
		DividendId:       payout.DividendId,
		AccountId:        payout.AccountId,
		Amount:           payout.Amount,
		Figi:             payout.Dividend.Figi,
		Isin:             payout.Dividend.Security.Isin,
		Ticker:           payout.Dividend.Security.Ticker,
		ActualDPS:        payout.Dividend.ActualDPS,
		AnnouncementDate: payout.Dividend.AnnouncementDate,
		RecordDate:       payout.Dividend.RecordDate,
		PayoutDate:       payout.Dividend.PayoutDate,
	}

	return dto
}
