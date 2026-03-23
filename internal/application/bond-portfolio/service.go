package bondportfolio

import (
	"fmt"
	"sync"
	"time"

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

func CalculateYtmForLots(lots []bonds.BondLot) ([]bonds.BondLot, error) {
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
	wg.Go(func() {
		defer wg.Done()
		quotes, err = bondquote.FetchQuotesForFigis(figis, config)
	})

	var bondList []bonds.Bond
	wg.Go(func() {
		defer wg.Done()
		bondList, err = bondservice.GetBondsByFigi(figis)
	})
	wg.Wait()
	fmt.Println("The wait group is finished!")

	bondList = bondservice.PopulateBondCoupons(bondList)

	bondList = bondservice.CalculateYtmForBonds(bondList, quotes)

	lots = matchLotsWithBonds(lots, bondList)

	return lots, nil
}

func PopulateLotsWithBonds(lots []bonds.BondLot) ([]bonds.BondLot, error) {
	figis := GetLotFigis(lots)
	bondList, err := bondservice.GetBondsByFigi(figis)
	if err != nil {
		return []bonds.BondLot{}, err
	}
	lots = matchLotsWithBonds(lots, bondList)
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

// TODO: Complete this
func GetLotsWithAci() ([]bonds.BondLot, error) {
	lots, err := GetAllPositionLots()
	if err != nil {
		return []bonds.BondLot{}, err
	}

	lots, err = PopulateLotsWithBonds(lots)
	if err != nil {
		return []bonds.BondLot{}, err
	}

	bondsWithCoupons := bondservice.PopulateBondCoupons(GetLotBonds(lots))
	lots = matchLotsWithBonds(lots, bondsWithCoupons)

	for _, lot := range lots {
		aci, err := CalcAccumulatedCouponIncomeForLot(lot)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("ACI for", lot.Figi, aci)
	}

	return lots, nil
}

func CalcAccumulatedCouponIncomeForLot(lot bonds.BondLot) (float64, error) {
	aciOnOpeningDate, err := bonds.AccumulatedCouponIncome(lot.Bond, lot.OpeningDate)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Lot's ACI: ", aciOnOpeningDate)

	currentAci, err := bonds.AccumulatedCouponIncome(lot.Bond, time.Now())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Bond's current ACI: ", currentAci)

	return aciOnOpeningDate, nil
}
