package security

import (
	"context"
	// "time"

	"github.com/compoundinvest/stockfundamentals/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

func FetchSecuritiesFromDB() ([]Stock, error) {
	// ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	// defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		return []Stock{}, err
	}

	db, err := ydb.Open(context.TODO(), config.DB.ConnectionString)

	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}

	return FetchSecuritiesFromDBWithDriver(db)
}
