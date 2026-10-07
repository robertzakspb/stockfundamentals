package security_master

import (
	"errors"
	"strings"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

func GetAllBonds() ([]bonds.Bond, error) {
	return GetFilteredBonds([]ydbfilter.YdbFilter{})
}

func GetFilteredBonds(filters []ydbfilter.YdbFilter) ([]bonds.Bond, error) {
	//Default filter to remove historical matured bonds
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "maturity_date",
		Condition:      ydbfilter.GreaterThan,
		ConditionValue: ydbhelper.ConvertToYdbDate(time.Now()),
	}
	filters = append(filters, filter)

	bondList, err := bondsdb.GetAllBonds(filters)
	if err != nil {
		return []bonds.Bond{}, err
	}

	if len(bondList) == 0 {
		return []bonds.Bond{}, errors.New("Found zero bonds in the DB")
	}

	mappedBonds := []bonds.Bond{}
	for _, dbBond := range bondList {
		mappedBond := mapDbBondToBond(dbBond)
		mappedBonds = append(mappedBonds, mappedBond)
	}

	return mappedBonds, nil
}

func GetBondByFigi(figi string) (bonds.Bond, error) {
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "figi",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.TextValue(figi),
	}
	bondList, err := bondsdb.GetAllBonds([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return bonds.Bond{}, err
	}
	if len(bondList) == 0 {
		return bonds.Bond{}, errors.New("Found zero bonds with the specificed figi")
	}

	mappedBond := mapDbBondToBond(bondList[0])

	return mappedBond, nil
}

func GetBondsByFigi(figis []string) ([]bonds.Bond, error) {
	ydbFigis := []types.Value{}
	for _, figi := range figis {
		ydbFigis = append(ydbFigis, types.TextValue(figi))
	}

	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "figi",
		Condition:      ydbfilter.Contains,
		ConditionValue: types.ListValue(ydbFigis...),
	}
	bondList, err := bondsdb.GetAllBonds([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return []bonds.Bond{}, err
	}
	if len(bondList) == 0 {
		return []bonds.Bond{}, errors.New("Found zero bonds with the specificed figis")
	}

	mappedBonds := []bonds.Bond{}
	for _, dbBond := range bondList {
		mappedBond := mapDbBondToBond(dbBond)
		mappedBonds = append(mappedBonds, mappedBond)
	}

	return mappedBonds, nil
}

func GetBondByIsin(isin string) (bonds.Bond, error) {
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "isin",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.TextValue(strings.ToUpper(isin)),
	}
	bondList, err := bondsdb.GetAllBonds([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return bonds.Bond{}, err
	}
	if len(bondList) == 0 {
		return bonds.Bond{}, errors.New("Found zero bonds with the specificed ISIN")
	}

	mappedBond := mapDbBondToBond(bondList[0])

	return mappedBond, nil
}

func GetBondsByIsin(isins []string) ([]bonds.Bond, error) {
	ydbIsins := []types.Value{}
	for _, ydbIsin := range isins {
		ydbIsins = append(ydbIsins, types.TextValue(ydbIsin))
	}
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "isin",
		Condition:      ydbfilter.Contains,
		ConditionValue: types.ListValue(ydbIsins...),
	}
	bondList, err := bondsdb.GetAllBonds([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return []bonds.Bond{}, err
	}
	if len(bondList) == 0 {
		return []bonds.Bond{}, errors.New("Found zero bonds with the specificed ISINs")
	}

	mappedBonds := []bonds.Bond{}
	for _, bond := range bondList {
		mappedBond := mapDbBondToBond(bond)
		mappedBonds = append(mappedBonds, mappedBond)
	}

	return mappedBonds, nil
}
