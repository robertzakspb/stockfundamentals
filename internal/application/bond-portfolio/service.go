package bondportfolio

import (
	"fmt"
	"sort"
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
		quotes, err = bondquote.FetchQuotesForFigis(figis, config)
	})

	var bondList []bonds.Bond
	wg.Go(func() {
		bondList, err = bondservice.GetBondsByFigi(figis)
	})
	wg.Wait()
	fmt.Println("The wait group is finished!")

	bondList = bondservice.PopulateBondCoupons(bondList)

	bondList = bondservice.CalculateYtmForBonds(bondList, quotes)

	lots = matchLotsWithBonds(lots, bondList)

	sort.Slice(lots, func(i, j int) bool {
		return lots[i].Bond.YieldToMaturity > lots[j].Bond.YieldToMaturity
	})

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

func PopulateLotsWithCoupons(lots []bonds.BondLot) []bonds.BondLot {
	bonds := bondservice.PopulateBondCoupons(getLotBonds(lots))
	lotsWithCoupons := matchLotsWithBonds(lots, bonds)
	return lotsWithCoupons
}

func GetAccountTimeline() ([]TimeLineItem, error) {
	lots, err := GetAllPositionLots()
	if err != nil {
		return []TimeLineItem{}, err
	}

	lots, err = PopulateLotsWithBonds(lots)
	if err != nil {
		return []TimeLineItem{}, err
	}

	accountTimeline, err := generateTimeLineForLots(lots)
	if err != nil {
		return []TimeLineItem{}, err
	}

	return accountTimeline, nil
}

// This function assumes that the lots' bonds have already been populated with coupons
// func CalculateAciForLots(lots []bonds.BondLot) ([]bonds.BondLot, error) {
// 	for i, lot := range lots {
// 		aci, err := CalcAccumulatedCouponIncomeForLot(lot)
// 		if err != nil {
// 			logger.Log(err.Error(), logger.ERROR)
// 		}
// 		lots[i].AccumulatedCouponIncome = aci
// 	}

// 	return lots, nil
// }
