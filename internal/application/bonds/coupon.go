package bondservice

import (
	"context"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
)

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
