package stockoverview

import (
	"context"
	"errors"
	"os/signal"
	"syscall"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"opensource.tbank.ru/invest/invest-go/investgo"
)

func GetFinancialOverviewFinancials(figis []string) ([]StockOverview, error) {
	if len(figis) == 0 {
		logger.Log("Unexpectedly received 0 figis when attempting to fetch stocks' fundamentals", logger.ERROR)
	}

	config, err := investgo.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return []StockOverview{}, err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	client, err := investgo.NewClient(ctx, config, nil)
	if err != nil {
		logger.Log("Failed to initialize the Tinkoff API client: ", logger.ALERT)
		return []StockOverview{}, err
	}

	finService := client.NewInstrumentsServiceClient()
	if finService == nil {
		logger.Log("The financials service is unexpectedly nil", logger.ALERT)
		return []StockOverview{}, err
	}

	response, err := finService.GetAssetFundamentals(figis)
	if err != nil {
		logger.Log("Failed to receive fundamentals for figis due to "+err.Error(), logger.ERROR)
		return []StockOverview{}, err
	}
	if response == nil || response.Fundamentals == nil {
		return []StockOverview{}, errors.New("Response or response.Fundamentals are nil while attempting to fetch the fundamentals")
	}

	overviews := []StockOverview{}
	for i := range response.Fundamentals {
		overview := mapTinkoffOverviewToInternal(response.Fundamentals[i])
		overviews = append(overviews, overview)
	}

	return overviews, nil
}
