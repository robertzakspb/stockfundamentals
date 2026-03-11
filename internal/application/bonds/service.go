package bondservice

import (
	"context"
	"os/signal"
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

	dbCoupons := []bondsdb.CouponDbModel{}
	for _, bond := range bonds {
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

		coupondPeriodEndDate, _ := time.Parse(time.DateOnly, "2100-01-01")
		coupondPeriodStartDate, _ := time.Parse(time.DateOnly, "1970-01-01")
		response, err := bondService.GetBondCoupons(bond.Figi, coupondPeriodStartDate, coupondPeriodEndDate)

		for _, tinkoffCoupon := range response.GetEvents() {
			coupon := mapTinkoffCouponToCoupon(bond.Figi, tinkoffCoupon)
			dbCoupon := mapCouponToDbModel(coupon)
			dbCoupons = append(dbCoupons, dbCoupon)
		}
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
	bond, err := bondsdb.GetAllBonds([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return bonds.Bond{}, err
	}

	mappedBond := mapDbBondToBond(bond[0])

	return mappedBond, nil
}
