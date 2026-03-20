package bondportfolio

import (
	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
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

func GetAccountTimeline() ([]TimeLineItem, error) {
	lots, err := GetAllPositionLots()
	if err != nil {
		return []TimeLineItem{}, err
	}

	accountTimeline, err := generateTimeLineForLot(lots)
	if err != nil {
		return []TimeLineItem{}, err
	}

	return accountTimeline, nil
}

func GetPositionLotsWithYtm() ([]bonds.BondLot, error) {
	//TODO: introduce a wait group to use routines for collection of bonds, lots, and quotes
	lots, err := GetAllPositionLots()
	if err != nil {
		return []bonds.BondLot{}, err
	}

	figis := []string{}
	for _, bond := range lots {
		figis = append(figis, bond.Figi)
	}

	quotes := bondquote.FetchQuotesForFigis(figis)
	bonds := bondservice.GetBondsByFigi(figis)
	coupons := bondservice.GetCouponsByFigi()

	for i, lot := range lots {
		for _, bond := range bonds {
			if bond.Figi == lot.Figi {
				for _, quote := range quotes {
					if quote.Figi == lot.Figi {
						lots[i].YieldToMaturity = bond.YieldToMaturity()
					}
				}
			}
		}
	}
}
