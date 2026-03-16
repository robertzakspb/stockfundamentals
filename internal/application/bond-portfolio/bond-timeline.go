package bondportfolio

import (
	"fmt"
	"sort"
	"time"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
)

type TimeLineItem struct {
	timestamp time.Time
	eventName string
}

func generateTimeLineForLot(lot bonds.BondLot) ([]TimeLineItem, error) {
	bond, err := bondservice.GetBondByFigi(lot.Figi)
	if err != nil {
		return []TimeLineItem{}, err
	}

	timeline := []TimeLineItem{}

	if bond.RegistrationDate.IsZero() == false {
		timeline = append(timeline, TimeLineItem{
			timestamp: bond.RegistrationDate,
			eventName: "Дата Регистрации Облигации",
		})
	}
	if bond.PlacementDate.IsZero() == false {
		timeline = append(timeline, TimeLineItem{
			timestamp: bond.PlacementDate,
			eventName: "Дата Размещения Облигации",
		})
	}
	if bond.CallOptionExerciseDate.IsZero() == false {
		timeline = append(timeline, TimeLineItem{
			timestamp: bond.CallOptionExerciseDate,
			eventName: "Дата Колл-опциона",
		})
	}
	timeline = append(timeline, TimeLineItem{
		timestamp: bond.RegistrationDate,
		eventName: "Дата Погашения Облигации. Возврат денежных средств: " + bond.Currency + fmt.Sprint(lot.TotalPrincipalRedemption(bond)),
	})

	coupons, _ := bondservice.GetCouponsByFigi(bond.Figi)
	for _, coupon := range coupons {
		timeline = append(timeline, TimeLineItem{
			timestamp: coupon.CouponDate,
			eventName: "Выплата купона: " + bond.Currency + fmt.Sprint(coupon.PerBondAmount) + " * " + fmt.Sprint(lot.Quantity) + " = " + bond.Currency + fmt.Sprint(lot.CouponPayoutForPosition(coupon)),
		})
	}

	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].timestamp.After(timeline[j].timestamp)
	})

	return timeline, nil
}
