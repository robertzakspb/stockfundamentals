package bondportfolio

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/application/bondservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/shared"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
	pb "opensource.tbank.ru/invest/invest-go/proto"
)

func ImportTinkoffBondLots() error {
	bondLots, err := FetchTinkoffBondLots()
	if err != nil {
		return err
	}

	mappedDbBondLots := make([]bondsdb.BondPositionLotDb, len(bondLots))
	for i := range bondLots {
		mappedDbBondLot := mapBondLotToDbModel(bondLots[i])
		mappedDbBondLots[i] = mappedDbBondLot
	}

	err = bondsdb.SaveBondPositionLots(mappedDbBondLots)
	if err != nil {
		return err
	}

	return nil
}

func FetchTinkoffBondLots() ([]bonds.BondLot, error) {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return []bonds.BondLot{}, err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	client, err := tinkoff.NewClient(ctx, config, nil)
	if err != nil {
		logger.Log("Failed to initialize the Tinkoff API client: ", logger.ALERT)
		return []bonds.BondLot{}, err
	}

	usersService := client.NewUsersServiceClient()
	status := pb.AccountStatus_ACCOUNT_STATUS_ALL
	accsResp, err := usersService.GetAccounts(&status)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return []bonds.BondLot{}, err
	}

	accounts := accsResp.GetAccounts()
	if len(accounts) == 0 {
		logger.Log("No accounts found in Tinkoff API", logger.ALERT)
		return []bonds.BondLot{}, err
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

	lots := []bonds.BondLot{}
	// securities := FetchTinkoffPositionSecurities(allPositions)
	for _, position := range allPositions {
		if !(position.GetInstrumentType() == "bond") {
			continue //Skipping non-bond positions in the bond import job
		}

		if position.GetFigi() == "" {
			logger.Log("Missing the bond position's figi "+position.Figi, logger.ERROR)
		}

		bond, err := bondservice.GetBondByFigi(position.GetFigi())

		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
			continue
		}

		var tinkoffIisId, _ = uuid.Parse(shared.TINKOFF_IIS_ACCOUNT_ID)
		newLot := bonds.BondLot{
			Id:                     uuid.New(),
			Figi:                   position.GetFigi(),
			Isin:                   "",
			OpeningDate:            time.Now(),
			ModificationDate:       time.Now(),
			AccountId:              tinkoffIisId,
			Quantity:               position.Quantity.ToFloat(),
			PricePerUnitPercentage: position.AveragePositionPrice.ToFloat() / bond.NominalValue * 100,
			Bond:                   bond,
			AccruedInterest:        position.CurrentNkd.ToFloat(),
		}

		lots = append(lots, newLot)
	}

	return lots, nil
}
