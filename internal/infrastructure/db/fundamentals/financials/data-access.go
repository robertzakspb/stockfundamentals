package dbfinancials

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

const stock_directory_prefix = "stockfundamentals/stocks"
const financial_metrics_table_name = "financial_metric"

type FinancialMetricDbModel struct {
	Id       uuid.UUID `sql:"id"`
	StockId  uuid.UUID `sql:"stock_id"`
	Name     string    `sql:"metric"`
	Period   string    `sql:"reporting_period"`
	Year     int64     `sql:"year"`
	Value    int64     `sql:"metric_value"`
	Currency string    `sql:"metric_currency"`
}

func SaveFinancialMetricsToDb(dbModels []FinancialMetricDbModel) error {
	ydbFinancials := []types.Value{}

	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		return err
	}
	defer db.ReleaseDriver(dbConnection)

	for _, metric := range dbModels {
		ydbMetric := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(metric.Id)),
			types.StructFieldValue("stock_id", types.UuidValue(metric.StockId)),
			types.StructFieldValue("metric", types.UTF8Value(metric.Name)),
			types.StructFieldValue("reporting_period", types.UTF8Value(string(metric.Period))),
			types.StructFieldValue("year", types.Int64Value(int64(metric.Year))),
			types.StructFieldValue("metric_value", types.Int64Value(int64(metric.Value))),
			types.StructFieldValue("metric_currency", types.UTF8Value(metric.Currency)),
		)
		ydbFinancials = append(ydbFinancials, ydbMetric)
	}

	financialsTableName := path.Join(dbConnection.Name(), stock_directory_prefix, financial_metrics_table_name)
	err = dbConnection.Table().BulkUpsert(
		context.TODO(),
		financialsTableName,
		table.BulkUpsertDataRows(types.ListValue(ydbFinancials...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	return nil
}

func FetchFinancialMetrics() ([]FinancialMetricDbModel, error) {
	//TODO: Refactor to use the new generic method
	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		return []FinancialMetricDbModel{}, err
	}
	defer db.ReleaseDriver(dbConnection)

	dbMetrics := []FinancialMetricDbModel{}

	err = dbConnection.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, fmt.Sprintf(`
						SELECT
							year, 
                            id, 
							stock_id, 
							metric, 
							reporting_period,
							metric_value,
							metric_currency
						FROM
							%s
					`, "`"+path.Join(stock_directory_prefix, financial_metrics_table_name)+"`"),
				query.WithTxControl(query.TxControl(query.BeginTx(query.WithSnapshotReadOnly()))),
			)
			if err != nil {
				return err
			}

			defer func() {
				_ = result.Close(ctx)
			}()

			for {
				resultSet, err := result.NextResultSet(ctx)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}

					return err
				}

				for row, err := range sugar.UnmarshalRows[FinancialMetricDbModel](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}

					dbMetrics = append(dbMetrics, row)
				}
			}

			for _, metric := range dbMetrics {
				dbMetrics = append(dbMetrics, metric)
			}

			return nil
		},
	)
	if err != nil {
		return []FinancialMetricDbModel{}, err
	}

	return dbMetrics, nil
}
