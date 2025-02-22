package portfolio

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	tinkoffapi "github.com/russianinvestments/invest-api-go-sdk/investgo"
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
	config, err := tinkoffapi.LoadConfig("config.yaml")
	if err != nil {
		println("Failed to initialize the configuration file: ", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	client, err := tinkoffapi.NewClient(ctx, config, nil)
	if err != nil {
		fmt.Println("Failed to initialize the Tinkoff API client: ", err)
	}

	usersService := client.NewUsersServiceClient()
	status := pb.AccountStatus_ACCOUNT_STATUS_ALL
	accsResp, err := usersService.GetAccounts(&status)
	if err != nil {
		fmt.Println(err.Error())
	}

	accounts := accsResp.GetAccounts()
	if len(accounts) == 0 {
		fmt.Println("No accounts found in Tinkoff API")
	}
	defaultAccountID := accounts[0].GetId()

	securityService := client.NewInstrumentsServiceClient()

	positionService := client.NewOperationsServiceClient()
	portfolio, _ := positionService.GetPortfolio(defaultAccountID, pb.PortfolioRequest_RUB)

	lots := []Lot{}
	for _, position := range portfolio.GetPositions() {
		if position.GetInstrumentType() != "share" {
			continue
		}

		asset, _ := securityService.ShareByFigi(position.GetFigi())
		ticker := asset.Instrument.GetTicker()
		newLot := Lot{
			SecurityID:   "",
			Ticker:       ticker,
			Quantity:     float64(position.Quantity.ToFloat()),
			OpeningPrice: position.AveragePositionPrice.ToFloat(),
			Currency:     position.AveragePositionPrice.Currency,
			CompanyName:  "",
			BrokerName:   "Tinkoff",
			MIC:          "MISX",
		}
		lots = append(lots, newLot)
	}

	return lots
}
