package positionsnapshot

import (
	"time"

	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	stockportfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	positionsnapshotdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/position-snapshot"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func SaveStockPositionLotSnapshots() error {
	lots, err := portfolio.GetAllStockLots()
	if err != nil {
		logger.LogError(err, logger.ERROR)
		return err
	}

	//Need to populate securities to extract ISINs
	lots, err = portfolio.PopulateLotSecurities(lots)
	if err != nil {
		logger.LogError(err, logger.ERROR)
		return err
	}

	accountPositions := stockportfolio.GroupLotsByAccount(lots)

	//Collapsing each account's lots into security-level positions
	for accountId := range accountPositions {
		accountPositions[accountId] = stockportfolio.CollapseLotsIntoPositionsUsingFigi(accountPositions[accountId])
	}

	snapshots := []StockPositionSnapshot{}
	for _, positions := range accountPositions {
		for i := range positions {
			snapshot := StockPositionSnapshot{
				Figi:      positions[i].Figi,
				Isin:      positions[i].Stock.Isin,
				AccountId: positions[i].AccountId,
				Date:      time.Now(),
				Quantity:  positions[i].Quantity,
			}
			snapshots = append(snapshots, snapshot)
		}
	}

	dbModels := mapPositionSnapshotsToDbModels(snapshots)

	err = positionsnapshotdb.SaveStockPositionSnapshots(dbModels)
	if err != nil {
		logger.LogError(err, logger.ERROR)
	}

	return err
}
