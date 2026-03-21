package bondservice

import (
	"context"
	"errors"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
	pb "opensource.tbank.ru/invest/invest-go/proto"
)

func ImportAllBondsAndCoupons() error {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	client, err := tinkoff.NewClient(ctx, config, nil)
	if err != nil {
		logger.Log("Failed to initialize the Tinkoff API client: ", logger.ALERT)
		return err
	}

	bondService := client.NewInstrumentsServiceClient()
	response, err := bondService.Bonds(pb.InstrumentStatus_INSTRUMENT_STATUS_ALL)
	if response == nil {
		logger.Log("Unexpectedly received a nil response from Tinkoff API", logger.ALERT)
	}

	dbBonds := []bondsdb.BondDbModel{}
	for _, tinkoffBond := range response.Instruments {
		if tinkoffBond.MaturityDate.AsTime().Before(time.Now()) {
			//No need to import historical bonds that have matured
			continue
		}
		bond := mapTinkoffBondToBond(tinkoffBond)
		validationErr := bond.Validate()
		if validationErr != nil {
			logger.Log(validationErr.Error(), logger.WARNING)
		}
		dbBond := mapBondToDbBond(bond)
		dbBonds = append(dbBonds, dbBond)
	}

	err = bondsdb.SaveBonds(dbBonds)
	if err != nil {
		return err
	}

	go importAllCoupons()

	return nil
}

func importAllCoupons() error {
	bonds, err := bondsdb.GetAllBonds([]ydbfilter.YdbFilter{})
	if err != nil {
		return err
	}

	rateLimit := time.Second //To comply with the Tinkoff API rate limits
	throttle := time.Tick(rateLimit)
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return nil
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	client, err := tinkoff.NewClient(ctx, config, nil)
	if err != nil {
		logger.Log("Failed to initialize the Tinkoff API client: ", logger.ALERT)
		return nil
	}

	bondService := client.NewInstrumentsServiceClient()
	if bondService == nil {
		logger.Log("The bond service is unexpectedly nil", logger.ALERT)
		return nil
	}

	dbCoupons := []bondsdb.CouponDbModel{}
	coupondPeriodEndDate, _ := time.Parse(time.DateOnly, "2100-01-01")
	coupondPeriodStartDate, _ := time.Parse(time.DateOnly, "1970-01-01")
	for i, bond := range bonds {
		response, err := bondService.GetBondCoupons(bond.Figi, coupondPeriodStartDate, coupondPeriodEndDate)
		if err != nil {
			return err
		}

		for _, tinkoffCoupon := range response.GetEvents() {
			coupon := mapTinkoffCouponToCoupon(bond.Figi, tinkoffCoupon)
			dbCoupon := mapCouponToDbModel(coupon)
			dbCoupons = append(dbCoupons, dbCoupon)
		}
		logger.Log(strconv.Itoa(i+1)+" out of "+strconv.Itoa(len(bonds))+". Fetched coupons for figi "+bond.Figi, logger.INFORMATION)

		<-throttle
	}

	err = bondsdb.SaveCoupons(dbCoupons)
	if err != nil {
		return err
	}

	return nil
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
		Condition:      ydbfilter.Equal,
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
		ConditionValue: types.TextValue(isin),
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
		Condition:      ydbfilter.Equal,
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
