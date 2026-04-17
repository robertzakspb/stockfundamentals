package bondportfolio

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
)

type TimeLineItem struct {
	Timestamp time.Time
	EventName string
	BondName  string
	Amount    float64
	Currency  string
}

func generateTimeLineForLots(lots []bonds.BondLot, includePastEvents bool) ([]TimeLineItem, error) {
	if len(lots) == 0 {
		return []TimeLineItem{}, errors.New("Attempting to generate a timeline for 0 lots")
	}

	bondList := GetLotBonds(lots)
	if len(bondList) == 0 {
		return []TimeLineItem{}, errors.New("Provided position lots have no corresponding bonds")
	}

	timeline := makeTimeLine(lots)

	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Timestamp.Before(timeline[j].Timestamp)
	})

	return timeline, nil
}

func makeTimeLine(lots []bonds.BondLot) []TimeLineItem {
	timeline := []TimeLineItem{}
	for _, lot := range lots {
		bond := lot.Bond

		if bond.RegistrationDate.IsZero() == false {
			timeline = append(timeline, TimeLineItem{
				Timestamp: bond.RegistrationDate,
				EventName: "Дата Регистрации Облигации",
				BondName:  bond.Name,
			})
		}
		if bond.PlacementDate.IsZero() == false {
			timeline = append(timeline, TimeLineItem{
				Timestamp: bond.PlacementDate,
				EventName: "Дата Размещения Облигации",
				BondName:  bond.Name,
			})
		}
		if bond.CallOptionExerciseDate.IsZero() == false {
			timeline = append(timeline, TimeLineItem{
				Timestamp: bond.CallOptionExerciseDate,
				EventName: "Дата Колл-опциона",
				BondName:  bond.Name,
			})
		}
		timeline = append(timeline, TimeLineItem{
			Timestamp: bond.MaturityDate,
			EventName: "Дата Погашения Облигации. Возврат денежных средств: " + bond.NominalCurrency + fmt.Sprint(lot.TotalPrincipalRedemption(bond)),
			BondName:  bond.Name,
			Amount:    bond.NominalValue * lot.Quantity,
			Currency:  bond.NominalCurrency,
		})

		for _, coupon := range lot.Bond.Coupons {
			timeline = append(timeline, TimeLineItem{
				Timestamp: coupon.CouponDate,
				EventName: "Выплата купона: " + bond.Currency + fmt.Sprint(coupon.PerBondAmount) + " * " + fmt.Sprint(lot.Quantity) + " = " + bond.Currency + fmt.Sprint(lot.CouponPayoutForPosition(coupon)),
				BondName:  bond.Name,
				Amount:    coupon.PerBondAmount,
				Currency:  bond.NominalCurrency,
			})
		}
	}
	return timeline
}

func totalPayoutForEventsInCurrency(events []TimeLineItem, currency string) float64 {
	totalPayout := 0.0
	for _, event := range events {
		if event.Currency == currency {
			totalPayout += event.Amount
		}
	}
	return totalPayout
}
