package bondportfolio

import (
	"fmt"
	"sort"
	"time"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
)

type BondTimeLine struct {
	Bond   bonds.Bond
	events []TimeLineItem
}

type TimeLineItem struct {
	timestamp time.Time
	eventName string
}

func generateTimeLineForLot(lot bonds.BondLot) (BondTimeLine, error) {
	bond, err := bondservice.GetBondByFigi(lot.Figi)
	if err != nil {
		return BondTimeLine{}, err
	}

	timeline := BondTimeLine{
		Bond: bond,
	}

	if bond.RegistrationDate.IsZero() == false {
		timeline.events = append(timeline.events, TimeLineItem{
			timestamp: bond.RegistrationDate,
			eventName: "Дата Регистрации Облигации",
		})
	}
	if bond.PlacementDate.IsZero() == false {
		timeline.events = append(timeline.events, TimeLineItem{
			timestamp: bond.PlacementDate,
			eventName: "Дата Размещения Облигации",
		})
	}
	if bond.CallOptionExerciseDate.IsZero() == false {
		timeline.events = append(timeline.events, TimeLineItem{
			timestamp: bond.CallOptionExerciseDate,
			eventName: "Дата Колл-опциона",
		})
	}
	timeline.events = append(timeline.events, TimeLineItem{
		timestamp: bond.RegistrationDate,
		eventName: "Дата Погашения Облигации. Возврат денежных средств: " + bond.Currency + fmt.Sprint(lot.TotalPrincipalRedemption(bond)),
	})

	coupons, _ := bondservice.GetCouponsByFigi(bond.Figi)
	for _, coupon := range coupons {
		timeline.events = append(timeline.events, TimeLineItem{
			timestamp: coupon.CouponDate,
			eventName: "Выплата купона: " + bond.Currency + fmt.Sprint(coupon.PerBondAmount) + " * " + fmt.Sprint(lot.Quantity) + " = " + bond.Currency + fmt.Sprint(lot.CouponPayoutForPosition(coupon)),
		})
	}

	sort.Slice(timeline.events, func(i, j int) bool {
		return timeline.events[i].timestamp.After(timeline.events[j].timestamp)
	})

	return timeline, nil
}
