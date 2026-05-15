package dbfinancials

import (
	"context"
	"path"

	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	ydbtemplate "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-template"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

const stock_directory_prefix = "stockfundamentals/stocks"
const financial_metrics_table_name = "financial_metric"

type FinancialMetricDbModel struct {
	Id              uuid.UUID `sql:"id"`
	StockId         uuid.UUID `sql:"stock_id"`
	Name            string    `sql:"metric"`
	ReportingPeriod string    `sql:"reporting_period"`
	Year            int64     `sql:"year"`
	Value           int64     `sql:"metric_value"`
	Currency        string    `sql:"metric_currency"`
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
			types.StructFieldValue("reporting_period", types.UTF8Value(string(metric.ReportingPeriod))),
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
	filters := []ydbfilter.YdbFilter{}
	tablePath := ydbhelper.GenerateTablePath(db.STOCK_DIRECTORY_PREFIX, db.FINANCIAL_METRIC_TABLE_NAME)

	dbMetrics, err := ydbtemplate.GetEntity[FinancialMetricDbModel](filters, tablePath)

	return dbMetrics, err
}
