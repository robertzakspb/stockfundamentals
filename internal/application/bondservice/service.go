package bondservice

import (
	"time"

	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

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
		if len(bond.Coupons) > 0 {
			logger.Log("The provided bond appears to have coupons, will not populate its coupons", logger.WARNING)
			continue
		}
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
	bondList, err := security_master.GetAllBonds()
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

	dbBonds := security_master.MapBondsToDbBonds(bondList)

	err = bondsdb.SaveBonds(dbBonds)
	if err != nil {
		return err
	}

	logger.Log("Completed the accrued interest update job", logger.INFORMATION)

	return nil
}
