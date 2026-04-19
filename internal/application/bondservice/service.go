package bondservice

import (
	"errors"
	"strings"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
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

func GetCouponsByFigi(figi string) ([]bonds.Coupon, error) {
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "figi",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.TextValue(figi),
	}

	coupons, err := bondsdb.GetBondCoupons([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return []bonds.Coupon{}, err
	}

	mappedCoupons := []bonds.Coupon{}
	for _, coupon := range coupons {
		mappedCoupon := mapCouponDbModelToDomain(coupon)
		mappedCoupons = append(mappedCoupons, mappedCoupon)
	}
	return mappedCoupons, nil
}

func GetCouponsByFigis(figis []string) ([]bonds.Coupon, error) {
	ydbFigis := []types.Value{}
	for _, figi := range figis {
		ydbFigis = append(ydbFigis, types.TextValue(figi))
	}
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "figi",
		Condition:      ydbfilter.Contains,
		ConditionValue: types.ListValue(ydbFigis...),
	}

	coupons, err := bondsdb.GetBondCoupons([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return []bonds.Coupon{}, err
	}

	mappedCoupons := []bonds.Coupon{}
	for _, coupon := range coupons {
		mappedCoupon := mapCouponDbModelToDomain(coupon)
		mappedCoupons = append(mappedCoupons, mappedCoupon)
	}
	return mappedCoupons, nil
}

func PopulateBondCoupons(bondList []bonds.Bond) []bonds.Bond {
	figis := []string{}
	for _, bond := range bondList {
		figis = append(figis, bond.Figi)
	}
	coupons, err := GetCouponsByFigis(figis)
	if err != nil {
		logger.Log("Failed to fetch coupons for the provided bonds", logger.ERROR)
		return bondList
	}
	bondsWithCoupons := MatchCouponsWithBonds(coupons, bondList)
	return bondsWithCoupons
}

func UpdateAllBondsAci() error {
	bondList, err := GetAllBonds()
	if err != nil {
		return err
	}

	bondList = PopulateBondCoupons(bondList)

	for i, bond := range bondList {
		aci, err := bonds.AccruedInterest(bond, time.Now())
		if err != nil {
			logger.Log(err.Error(), logger.WARNING)
			continue
		}
		bondList[i].AccruedInterest = aci
	}

	dbBonds := mapBondsToDbBonds(bondList)

	err = bondsdb.SaveBonds(dbBonds)
	if err != nil {
		return err
	}

	logger.Log("Completed the accrued interest update job", logger.INFORMATION)

	return nil
}
