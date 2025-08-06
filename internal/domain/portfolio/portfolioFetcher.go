package portfolio

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/compoundinvest/stockfundamentals/features/fundamentaldata/security"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"

	tinkoff "github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func GeMyPortfolio() Portfolio {
	hardCodedPositions := getHardCodedStockPositions()
	externalPositions := getExternalStockPositions()

	allPositions := append(hardCodedPositions, externalPositions...)
	return Portfolio{Lots: allPositions}
}

func getExternalStockPositions() []Lot {
	externalPositions := getTinkoffStockPositions()
	return externalPositions
}

func getTinkoffStockPositions() []Lot {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		println("Failed to initialize the configuration file: ", err)
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

	// securityService := client.NewInstrumentsServiceClient()
	lots := []Lot{}
	securities := fetchPositionSecurities(allPositions)
	for _, position := range allPositions {
		if position.GetInstrumentType() != "share" {
			continue //Skipping the cash position until it is handled separately
		}

		security, err := security.FindStockByFigi(securities, position.GetFigi())
		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
		}

		newLot := Lot{
			//TODO: Use the NewLot method
			Quantity:     float64(position.Quantity.ToFloat()),
			PricePerUnit: position.AveragePositionPrice.ToFloat(),
			Currency:     position.AveragePositionPrice.Currency,
			AccountId:    tinkoffIisId,
			Security: security,
		}
		lots = append(lots, newLot)
	}

	return lots
}

func fetchPositionSecurities(positions []*pb.PortfolioPosition) []security.Stock {
	figis := []string{}
	for _, position := range positions {
		figis = append(figis, position.Figi)
	}

	securities, err := security.GetSecuritiesFilteredByFigi(figis)
	if err != nil || len(securities) == 0 {
		logger.Log("Failed to find positions with the required figis: ", logger.ERROR)
		return []security.Stock{}
	}

	return  securities
}
