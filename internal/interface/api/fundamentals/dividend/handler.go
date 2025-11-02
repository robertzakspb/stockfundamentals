package apidividend

import (
	"context"
	"time"

	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

func FetchAndSaveAllDividends() error {
	dividends := appdividend.FetchDividendsForAllStocks()

	//TODO: Extract the DB initialization in a single method
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

	err = dbdividend.SaveDividendsToDB(dividends, db)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}

	return nil
}

func GetAllDividends() ([]dividend.Dividend, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		return []dividend.Dividend{}, err
	}

	db, err := ydb.Open(ctx, config.DB.ConnectionString)

	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}
	dividends, _ := dbdividend.GetAllDividends(db)

	return dividends, nil
}
