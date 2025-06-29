package financials

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

const stock_directory_prefix = "stockfundamentals/stocks"
const financial_metrics_table_name = "financial_metric"

func SaveFinancialMetricsToDb(metrics []FinancialMetric, db *ydb.Driver) error {
	ydbFinancials := []types.Value{}
	dbModels := mapFinancialMetricModelToDbModel(metrics)

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

	financialsTableName := path.Join(db.Name(), stock_directory_prefix, financial_metrics_table_name)
	err := db.Table().BulkUpsert(
		context.TODO(),
		financialsTableName,
		table.BulkUpsertDataRows(types.ListValue(ydbFinancials...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	return nil
}

func fetchFinancialMetricsFromDbWithDriver(db *ydb.Driver) ([]FinancialMetric, error) {
	dbMetrics := []FinancialMetricDbModel{}
	parsedMetrics := []FinancialMetric{}

	err := db.Query().Do(context.TODO(),
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
				parsedMetrics = append(parsedMetrics, mapYdbMetricToMetric(metric))
			}

			return nil
		},
	)
	if err != nil {
		return []FinancialMetric{}, err
	}

	return parsedMetrics, nil
}

func mapYdbMetricToMetric(dbMetric FinancialMetricDbModel) FinancialMetric {
	return FinancialMetric{
		Id: dbMetric.Id,
		StockId: dbMetric.StockId,
		Name: dbMetric.Name,
		Period: ReportingPeriodMap[dbMetric.Period],
		Year: int(dbMetric.Year),
		Value: int(dbMetric.Value),
		Currency: dbMetric.Currency,
	}
}