package dividend

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type Dividend struct {
	Id                string    `sql:"id"`
	StockID           string    `sql:"stock_id"`
	ActualDPS         float64   `sql:"actual_DPS"`
	ExpectedDPS       float64   `sql:"expected_DPS"`
	Currency          string    `sql:"currency"`
	AnnouncementDate  time.Time `sql:"announcement_date"`
	RecordDate        time.Time `sql:"record_date"`
	PayoutDate        time.Time `sql:"payout_date"`
	PaymentPeriod     string    `sql:"payment_period"`
	ManagementComment string    `sql:"management_comment"`
}

const STOCK_DIRECTORY_PREFIX = "stockfundamentals/stocks"
const DIVIDEND_PAYMENT_TABLE_NAME = "dividend_payment"

func saveDividendsToDB(dividends []Dividend, db *ydb.Driver) error {
	logger.Log(strconv.Itoa(len(dividends)), logger.INFORMATION)

	if db == nil {
		logger.Log("Database driver is nil while attempting to save dividends to the DB", logger.ALERT)
	}

	ydbDividends := []types.Value{}
	for _, dividend := range dividends {
		ydbDividend := types.StructValue(
			types.StructFieldValue("id", types.TextValue(dividend.Id)),
			types.StructFieldValue("stock_id", types.TextValue(dividend.StockID)),
			types.StructFieldValue("actual_DPS", types.DoubleValue(dividend.ActualDPS)),
			types.StructFieldValue("expected_DPS", types.DoubleValue(dividend.ExpectedDPS)),
			types.StructFieldValue("currency", types.TextValue(dividend.Currency)),
			types.StructFieldValue("announcement_date", types.DateValue(uint32(dividend.AnnouncementDate.Unix()))),
			types.StructFieldValue("record_date", types.DateValue(uint32(dividend.RecordDate.Unix()))),
			types.StructFieldValue("payout_date", types.DateValue(uint32(dividend.PayoutDate.Unix()))),
			types.StructFieldValue("payment_period", types.TextValue(dividend.PaymentPeriod)),
			types.StructFieldValue("management_comment", types.TextValue(dividend.ManagementComment)),
		)
		ydbDividends = append(ydbDividends, ydbDividend)
	}

	tableName := path.Join(db.Name(), STOCK_DIRECTORY_PREFIX, DIVIDEND_PAYMENT_TABLE_NAME)

	err := db.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(ydbDividends...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	return nil
}
