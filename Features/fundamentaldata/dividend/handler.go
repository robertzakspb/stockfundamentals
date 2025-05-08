package dividend

import (
	"context"
	"time"

	"github.com/compoundinvest/stockfundamentals/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

func FetchAndSaveAllDividends() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := ydb.Open(ctx, config.DB.ConnectionString)

	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}

	dividends := fetchDividendsForAllStocks(db)

	err = saveDividendsToDB(dividends, db)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}

	return nil
}
