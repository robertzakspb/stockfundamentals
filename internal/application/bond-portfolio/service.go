package bondportfolio

import (
	"errors"
	"time"

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

func validateLot(lot bonds.BondLot) error {
	if lot.Figi != "" {
		_, err := bondservice.GetBondByFigi(lot.Figi)
		if err != nil {
			return err
		}
	} else if lot.Isin != "" {
		_, err := bondservice.GetBondByIsin(lot.Isin)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Missing both figi and ISIN in the bond")
	}

	if lot.Quantity < 0 {
		return errors.New("Invalid quantity")
	}
	if lot.OpeningDate.After(time.Now()) {
		return errors.New("Invalid opening date")
	}
	if lot.ModificationDate.After(time.Now()) {
		return errors.New("Invalid modification date")
	}
	if lot.PricePerUnit < 0 {
		return errors.New("Invalid price per unit")
	}

	return nil
}
