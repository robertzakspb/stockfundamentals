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
	BondName  string
}

func generateTimeLineForLots(lots []bonds.BondLot) ([]TimeLineItem, error) {
	if len(lots) == 0 {
		return []TimeLineItem{}, errors.New("Attempting to generate a timeline for 0 bonds")
	}

	bondList := []bonds.Bond{}

	figis := []string{}
	isins := []string{}
	for _, lot := range lots {
		if lot.Figi != "" {
			figis = append(figis, lot.Figi)
		}
		if lot.Isin != "" {
			isins = append(isins, lot.Isin)
		}
	}

	if len(figis) > 0 {
		bondsByFigi, err := bondservice.GetBondsByFigi(figis)
		if err != nil {
			return []TimeLineItem{}, err
		}
		bondList = append(bondList, bondsByFigi...)
	}

	if len(isins) > 0 {
		bondsByIsin, err := bondservice.GetBondsByIsin(isins)
		if err != nil {
			return []TimeLineItem{}, err
		}
		bondList = append(bondList, bondsByIsin...)
	}

	if len(bondList) == 0 {
		logger.Log("Found zero bonds for the provided lots", logger.ERROR)
	}

	timeline := []TimeLineItem{}
	for _, lot := range lots {
		bond, err := findBondByFigiOrIsin(lot, bondList)
		if err != nil {
			logger.Log("Failed to find a bond for figi "+lot.Figi, logger.ERROR)
			continue
		}
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
		})

		coupons, _ := bondservice.GetCouponsByFigi(bond.Figi)
		for _, coupon := range coupons {
			timeline = append(timeline, TimeLineItem{
				Timestamp: coupon.CouponDate,
				EventName: "Выплата купона: " + bond.Currency + fmt.Sprint(coupon.PerBondAmount) + " * " + fmt.Sprint(lot.Quantity) + " = " + bond.Currency + fmt.Sprint(lot.CouponPayoutForPosition(coupon)),
				BondName:  bond.Name,
			})
		}
	}

	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Timestamp.Before(timeline[j].Timestamp)
	})

	return timeline, nil
}

func findBondByFigiOrIsin(lot bonds.BondLot, bondList []bonds.Bond) (bonds.Bond, error) {
	for _, bond := range bondList {
		if bond.Figi == lot.Figi {
			return bond, nil
		}
		if bond.Isin == lot.Isin {
			return bond, nil
		}
	}

	return bonds.Bond{}, errors.New("Failed to find the target bonds")
}
