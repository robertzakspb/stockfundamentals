package portfolio

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/portfoliodb"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"

	"github.com/google/uuid"

	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
	pb "opensource.tbank.ru/invest/invest-go/proto"
)

func GeMyPortfolio() portfolio.Portfolio {
	hardCodedPositions := getHardCodedStockPositions()
	externalPositions := getExternalStockPositions()

	allPositions := append(hardCodedPositions, externalPositions...)
	return portfolio.Portfolio{Lots: allPositions}
}

func GetAccountPortfolio(accountIDs uuid.UUIDs) (portfolio.Portfolio, error) {
	dbLots, err := portfoliodb.GetAccountPortfolio(accountIDs)
	if err != nil {
		return portfolio.Portfolio{}, err
	}
	lots := []lot.Lot{}
	for _, lot := range dbLots {
		lots = append(lots, mapLotDbToLot(lot))
	}

	return portfolio.Portfolio{Lots: lots}, nil
}

func getExternalStockPositions() []lot.Lot {
	externalPositions := getTinkoffStockPositions()
	return externalPositions
}

func getTinkoffStockPositions() []lot.Lot {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	client, err := tinkoff.NewClient(ctx, config, nil)
	if err != nil {
		logger.Log("Failed to initialize the Tinkoff API client: ", logger.ALERT)
	}

	usersService := client.NewUsersServiceClient()
	status := pb.AccountStatus_ACCOUNT_STATUS_ALL
	accsResp, err := usersService.GetAccounts(&status)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	accounts := accsResp.GetAccounts()
	if len(accounts) == 0 {
		logger.Log("No accounts found in Tinkoff API", logger.ALERT)
	}

	positionService := client.NewOperationsServiceClient()
	allPositions := []*pb.PortfolioPosition{}
	for _, account := range accounts {
		if account == nil {
			continue
		}

		portfolio, err := positionService.GetPortfolio(account.GetId(), pb.PortfolioRequest_RUB)
		if err != nil {
			logger.Log(err.Error(), logger.ALERT)
		}
		allPositions = append(allPositions, portfolio.GetPositions()...)
	}

	lots := []lot.Lot{}
	securities := FetchPositionSecurities(allPositions)
	for _, position := range allPositions {
		if position.GetInstrumentType() != "share" {
			continue //Skipping the cash position until it is handled separately
		}

		var stockId string
		for _, s := range securities {
			if s.GetFigi() == position.Figi {
				stockId = s.Figi
			}
		}
		if stockId == "" {
			logger.Log("Failed to find the stockId for "+position.Figi, logger.ERROR)
		}

		var tinkoffIisId, _ = uuid.Parse("3315bd1c-12a4-444e-a294-84ef339e26e1")
		newLot, err := lot.NewLot(stockId,
			float64(position.Quantity.ToFloat()),
			position.AveragePositionPrice.ToFloat(),
			position.AveragePositionPrice.Currency,
			tinkoffIisId,
		)
		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
		}

		lots = append(lots, newLot)
	}

	return lots
}
