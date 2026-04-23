package bondportfolio

import (
	"sort"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bondservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

// TODO: Optimize to remove the loop
func SaveBondPositionLots(lots []bonds.BondLot) error {
	for _, lot := range lots {
		err := SaveBondPositionLot(lot)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveBondPositionLot(lot bonds.BondLot) error {
	lot, err := validateLot(lot)
	if err != nil {
		return err
	}

	lot, err = addMissingInformationToLot(lot)
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
	return GetFilteredPositionLots([]ydbfilter.YdbFilter{})
}

func GetAccountPositions(accountId uuid.UUID) ([]bonds.BondLot, error) {
	accountFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "account_id",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.UuidValue(accountId),
	}
	return GetFilteredPositionLots([]ydbfilter.YdbFilter{accountFilter})
}

func GetFilteredPositionLots(filters []ydbfilter.YdbFilter) ([]bonds.BondLot, error) {
	lots, err := bondsdb.GetAccountBondPortfolio(filters)
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
	figis := []string{}
	for _, bond := range lots {
		figis = append(figis, bond.Figi)
	}

	bondList, err := bondservice.GetBondsByFigi(figis)
	if err != nil {
		return []bonds.BondLot{}, err
	}

	bondList = bondservice.PopulateBondsWithCouponsAndCalculateYtm(bondList)

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
	bonds := bondservice.PopulateBondCoupons(GetLotBonds(lots))
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

	lots = PopulateLotsWithCoupons(lots)

	accountTimeline, err := generateTimeLineForLots(lots, false)
	if err != nil {
		return []TimeLineItem{}, err
	}

	return accountTimeline, nil
}
