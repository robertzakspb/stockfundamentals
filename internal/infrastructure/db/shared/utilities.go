package shared

import (
	"context"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func ConvertToYdbDate(date time.Time) types.Value {
	const secondsInADay = 86400
	return types.DateValue(uint32(date.Unix() / secondsInADay))
}

func MakeYdbDriver() (*ydb.Driver, error ){
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		logger.Log("Failed to fetch the configuration", logger.ALERT)
		return nil, err
	}

	db, err := ydb.Open(ctx, config.DB.ConnectionString)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}

	return db, nil
}
