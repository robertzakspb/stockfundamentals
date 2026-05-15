package financialsservice

import (
	"context"
	"os/signal"
	"strconv"
	"syscall"

	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	tthrottler "github.com/compoundinvest/stockfundamentals/internal/application/tinkoff-throttler"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	stringhelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/string-helpers"
	"opensource.tbank.ru/invest/invest-go/investgo"
)

func StartTinkoffFinancialsImportJob() {
	go ExecuteImportFinancialsJob()
}

func ExecuteImportFinancialsJob() {
	config, err := investgo.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	client, err := investgo.NewClient(ctx, config, nil)
	if err != nil {
		logger.Log("Failed to initialize the Tinkoff API client: ", logger.ALERT)
		return
	}

	finService := client.NewInstrumentsServiceClient()
	if finService == nil {
		logger.Log("The financials service is unexpectedly nil", logger.ALERT)
		return
	}

	securities, err := security_master.GetAllSecuritiesFromDB()
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}
	figis := security_master.ExtractFigisFromSecurities(securities)

	if len(figis) == 0 {
		logger.Log("Unexpectedly received 0 figis when attempting to fetch stocks' fundamentals", logger.ERROR)
	}

	figiBatches := stringhelpers.SplitInBatchesOf(100, figis)

	for i := range figiBatches {
		<-tthrottler.InstrumentServiceThrottle
		response, err := finService.GetAssetFundamentals(figiBatches[i])
		if err != nil || response == nil {
			logger.Log("Failed to receive fundamentals for figi batch "+strconv.Itoa(i), logger.ERROR)
		}
	}
}
