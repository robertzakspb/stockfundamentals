package bondservice

import (
	"context"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	tthrottler "github.com/compoundinvest/stockfundamentals/internal/application/tinkoff-throttler"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
)

func importAllCoupons() error {
	startTime := time.Now()
	bonds, err := bondsdb.GetAllBonds([]ydbfilter.YdbFilter{})
	if err != nil {
		return err
	}

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
	for i, bond := range bonds {
		<-tthrottler.InstrumentServiceThrottle
		
		response, err := bondService.GetBondCoupons(bond.Figi, coupondPeriodStartDate, coupondPeriodEndDate)
		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
			continue
		}

		if len(response.GetEvents()) == 0 {
			continue
		}

		coupons := mapTinkoffCouponsToCoupons(response.GetEvents())
		dbCoupons := mapCouponsToDbModels(coupons)

		go bondsdb.SaveCoupons(&dbCoupons)

		logger.Log(strconv.Itoa(i+1)+" out of "+strconv.Itoa(len(bonds))+". Fetched coupons for figi "+bond.Figi, logger.INFORMATION)

	}

	jobCompletionTime := time.Now()
	logger.Log("The coupon import job has completed. Duration: "+jobCompletionTime.Sub(startTime).String(), logger.INFORMATION)

	go UpdateAllBondsAci()

	return nil
}
