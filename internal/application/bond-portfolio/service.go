package bondportfolio

import (
	"fmt"
	"sync"

	"github.com/compoundinvest/invest-core/quote/bondquote"
	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"opensource.tbank.ru/invest/invest-go/investgo"
)

func SaveBondPositionLot(lot bonds.BondLot) error {
	err := validateLot(lot)
	if err != nil {
		return err
	}

	mappedLot := mapBondLotToDbModel(lot)

	err = bondsdb.SaveBondPositionLots([]bondsdb.BondPositionLotDb{mappedLot})
	if err != nil {
		return err
	}

	return nil
}

func GetAllPositionLots() ([]bonds.BondLot, error) {
	//FIXME: This UUID should eventually be moved to a separate account table
	hardCodedAccountId, _ := uuid.Parse("129274f9-ee80-4e74-aa1c-fea578bac6e6")
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "account_id",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.UuidValue(hardCodedAccountId),
	}
	lots, err := bondsdb.GetAccountBondPortfolio([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return []bonds.BondLot{}, err
	}

	mappedLots := []bonds.BondLot{}
	for _, lot := range lots {
		mappedLot := mapDbBondToDomain(lot)
		mappedLots = append(mappedLots, mappedLot)
	}

	return mappedLots, nil
}

func GetPositionLotsWithYtm() ([]bonds.BondLot, error) {
	lots, err := GetAllPositionLots()
	if err != nil {
		return []bonds.BondLot{}, err
	}

	config, err := investgo.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return []bonds.BondLot{}, err
	}

	figis := []string{}
	for _, bond := range lots {
		figis = append(figis, bond.Figi)
	}

	wg := sync.WaitGroup{}

	var quotes []bondquote.TinkoffBondQuote
	wg.Add(1)
	go func() {
		defer wg.Done()
		quotes, err = bondquote.FetchQuotesForFigis(figis, config)
	}()

	var bondList []bonds.Bond
	wg.Add(1)
	go func() {
		defer wg.Done()
		bondList, err = bondservice.GetBondsByFigi(figis)
	}()

	var coupons []bonds.Coupon
	wg.Add(1)
	go func() {
		defer wg.Done()
		coupons, err = bondservice.GetCouponsByFigis(figis)
	}()

	wg.Wait()

	fmt.Println("The wait group is finished!")

	//Populating bond coupons
	for _, coupon := range coupons {
		for _, b := range bondList {
			if coupon.Figi == b.Figi {
				b.Coupons = append(b.Coupons, coupon)
			}
		}
	}

	//Calculating each bond's yield to maturity
	for _, quote := range quotes {
		for i, b := range bondList {
			if quote.Figi() == b.Figi {
				ytm, err := b.CalcYieldToMaturity(b.Coupons, quote.QuoteAsPercentage())
				if err != nil {
					logger.Log(err.Error(), logger.ERROR)
					continue
				}
				bondList[i].YieldToMaturity = ytm
			}
		}
	}

	//Populating the lots' bonds
	for i, lot := range lots {
		for _, b := range bondList {
			if lot.Figi == b.Figi {
				lots[i].Bond = b
			}
		}
	}

	return lots, nil
}

func GetAccountTimeline() ([]TimeLineItem, error) {
	lots, err := GetAllPositionLots()
	if err != nil {
		return []TimeLineItem{}, err
	}

	accountTimeline, err := generateTimeLineForLots(lots)
	if err != nil {
		return []TimeLineItem{}, err
	}

	return accountTimeline, nil
}
