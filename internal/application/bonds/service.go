package bondservice

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/bondsdb"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
	pb "opensource.tbank.ru/invest/invest-go/proto"
)

func ImportAllBonds() error {
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

	return nil
}
