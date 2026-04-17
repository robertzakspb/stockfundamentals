package bondportfolio

import (
	"errors"
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

	timeline := makeTimeLine(lots, includePastEvents)

	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Timestamp.Before(timeline[j].Timestamp)
	})

	return timeline, nil
}

func makeTimeLine(lots []bonds.BondLot, includePastEvents bool) []TimeLineItem {
	timeline := []TimeLineItem{}
	for _, lot := range lots {

		bond := lot.Bond

		if bond.RegistrationDate.IsZero() == false {
			event := TimeLineItem{
				Timestamp: bond.RegistrationDate,
				EventName: "Дата Регистрации Облигации",
				BondName:  bond.Name,
			}
			if !includePastEvents && event.Timestamp.After(time.Now()) {
				timeline = append(timeline, event)
			}

		}
		if bond.PlacementDate.IsZero() == false {
			event := TimeLineItem{
				Timestamp: bond.PlacementDate,
				EventName: "Дата Размещения Облигации",
				BondName:  bond.Name,
			}
			if !includePastEvents && event.Timestamp.After(time.Now()) {
				timeline = append(timeline, event)
			}
		}

		if bond.CallOptionExerciseDate.IsZero() == false {
			event := TimeLineItem{
				Timestamp: bond.CallOptionExerciseDate,
				EventName: "Дата Колл-опциона",
				BondName:  bond.Name,
			}
			if !includePastEvents && event.Timestamp.After(time.Now()) {
				timeline = append(timeline, event)
			}
		}

		maturityDate := TimeLineItem{
			Timestamp: bond.MaturityDate,
			EventName: "Погашениe Облигации",
			BondName:  bond.Name,
			Amount:    bond.NominalValue * lot.Quantity,
			Currency:  bond.NominalCurrency,
		}
		if !includePastEvents && maturityDate.Timestamp.After(time.Now()) {
			timeline = append(timeline, maturityDate)
		}

		for _, coupon := range lot.Bond.Coupons {
			payout := TimeLineItem{
				Timestamp: coupon.CouponDate,
				EventName: "Выплата купона",
				Amount:    coupon.PerBondAmount * lot.Quantity,
				Currency:  bond.NominalCurrency,
				BondName:  bond.Name,
			}
			if !includePastEvents && payout.Timestamp.After(time.Now()) {
				timeline = append(timeline, payout)
			}
		}
	}
	return timeline
}
