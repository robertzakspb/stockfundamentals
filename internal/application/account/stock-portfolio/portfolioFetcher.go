package portfolio

import (
	"context"
	"os/signal"
	"syscall"

	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/application/shared"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"

	"github.com/google/uuid"

	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
	pb "opensource.tbank.ru/invest/invest-go/proto"
)


func GeStockPortfolio() (portfolio.Portfolio, error) {
	hardCodedPositions := getHardCodedStockPositions()
	externalPositions, err := getExternalStockPositions()

	allPositions := append(hardCodedPositions, externalPositions...)
	return portfolio.Portfolio{Lots: allPositions}, err
}

func getExternalStockPositions() ([]lot.Lot, error) {
	return getTinkoffStockPositions()
}

func getTinkoffStockPositions() ([]lot.Lot, error) {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return []lot.Lot{}, err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	client, err := tinkoff.NewClient(ctx, config, nil)
	if err != nil {
		logger.Log("Failed to initialize the Tinkoff API client: ", logger.ALERT)
		return []lot.Lot{}, err
	}

	usersService := client.NewUsersServiceClient()
	status := pb.AccountStatus_ACCOUNT_STATUS_ALL
	accsResp, err := usersService.GetAccounts(&status)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return []lot.Lot{}, err
	}

	accounts := accsResp.GetAccounts()
	if len(accounts) == 0 {
		logger.Log("No accounts found in Tinkoff API", logger.ALERT)
		return []lot.Lot{}, err
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
	securities := FetchTinkoffPositionSecurities(allPositions)
	for _, position := range allPositions {
		if !(position.GetInstrumentType() == "share" || position.GetInstrumentType() == "etf" || position.GetInstrumentType() == "bond") {
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

		var tinkoffIisId, _ = uuid.Parse(shared.TINKOFF_IIS_ACCOUNT_ID)
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

	return lots, nil
}

func FetchTinkoffPositionSecurities(positions []*pb.PortfolioPosition) []security.Stock {
	figis := []string{}
	for _, position := range positions {
		figis = append(figis, position.Figi)
	}

	securities, err := security_master.GetSecuritiesFilteredByFigi(figis)
	if err != nil || len(securities) == 0 {
		logger.Log("Failed to find positions with the required figis: ", logger.ERROR)
		return []security.Stock{}
	}

	return securities
}
