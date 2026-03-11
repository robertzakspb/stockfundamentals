package bonds

import (
	"fmt"
	"sort"
	"time"
)

type BondTimeLine struct {
	Bond   Bond
	events []TimeLineItem
}

type TimeLineItem struct {
	timestamp time.Time
	eventName string
}

func (b Bond) generateTimeLine(coupons []Coupon) BondTimeLine {
	timeline := BondTimeLine{
		Bond: b,
	}

	if b.RegistrationDate.IsZero() == false {
		timeline.events = append(timeline.events, TimeLineItem{
			timestamp: b.RegistrationDate,
			eventName: "Дата Регистрации Облигации",
		})
	}
	if b.PlacementDate.IsZero() == false {
		timeline.events = append(timeline.events, TimeLineItem{
			timestamp: b.PlacementDate,
			eventName: "Дата Размещения Облигации",
		})
	}
	if b.CallOptionExerciseDate.IsZero() == false {
		timeline.events = append(timeline.events, TimeLineItem{
			timestamp: b.CallOptionExerciseDate,
			eventName: "Дата Колл-опциона",
		})
	}
	timeline.events = append(timeline.events, TimeLineItem{
		timestamp: b.RegistrationDate,
		eventName: "Дата Погашения Облигации",
	})

	for _, coupon := range coupons {
		timeline.events = append(timeline.events, TimeLineItem{
			timestamp: coupon.CouponDate,
			eventName: "Выплата купона: " + fmt.Sprint(coupon.PerBondAmount),
		})
	}

	sort.Slice(timeline.events, func(i, j int) bool {
		return timeline.events[i].timestamp.After(timeline.events[j].timestamp)
	})

	return timeline
}
