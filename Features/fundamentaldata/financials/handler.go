package financials

import (
	"context"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

func FetchFinancialsFromDB() ([]FinancialMetric, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return []FinancialMetric{}, err
	}

	db, err := ydb.Open(context.TODO(), config.DB.ConnectionString)

	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}

	return fetchFinancialMetricsFromDbWithDriver(db)
}
