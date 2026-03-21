package bondportfolio

import (
	"errors"
	"fmt"
	"sort"
	"time"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

type TimeLineItem struct {
	Timestamp time.Time
	EventName string
}

func generateTimeLineForLot(lots []bonds.BondLot) ([]TimeLineItem, error) {
	figis := []string{}
	for _, lot := range lots {
		figis = append(figis, lot.Figi)
	}

	bondList, err := bondservice.GetBondsByFigi(figis)
	if err != nil {
		return []TimeLineItem{}, err
	}

	timeline := []TimeLineItem{}
	for _, lot := range lots {
		bond, err := findBondByFigi(lot.Figi, bondList)
		if err != nil {
			logger.Log("Failed to find a bond for figi "+lot.Figi, logger.ERROR)
			continue
		}
		if bond.RegistrationDate.IsZero() == false {
			timeline = append(timeline, TimeLineItem{
				Timestamp: bond.RegistrationDate,
				EventName: "Дата Регистрации Облигации",
			})
		}
		if bond.PlacementDate.IsZero() == false {
			timeline = append(timeline, TimeLineItem{
				Timestamp: bond.PlacementDate,
				EventName: "Дата Размещения Облигации",
			})
		}
		if bond.CallOptionExerciseDate.IsZero() == false {
			timeline = append(timeline, TimeLineItem{
				Timestamp: bond.CallOptionExerciseDate,
				EventName: "Дата Колл-опциона",
			})
		}
		timeline = append(timeline, TimeLineItem{
			Timestamp: bond.RegistrationDate,
			EventName: "Дата Погашения Облигации. Возврат денежных средств: " + bond.Currency + fmt.Sprint(lot.TotalPrincipalRedemption(bond)),
		})

		coupons, _ := bondservice.GetCouponsByFigi(bond.Figi)
		for _, coupon := range coupons {
			timeline = append(timeline, TimeLineItem{
				Timestamp: coupon.CouponDate,
				EventName: "Выплата купона: " + bond.Currency + fmt.Sprint(coupon.PerBondAmount) + " * " + fmt.Sprint(lot.Quantity) + " = " + bond.Currency + fmt.Sprint(lot.CouponPayoutForPosition(coupon)),
			})
		}
	}

	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Timestamp.After(timeline[j].Timestamp)
	})

	return timeline, nil
}

func findBondByFigi(figi string, bondList []bonds.Bond) (bonds.Bond, error) {
	for _, bond := range bondList {
		if bond.Figi == figi {
			return bond, nil
		}
	}

	return bonds.Bond{}, errors.New("Failed to find the target bonds")
}
