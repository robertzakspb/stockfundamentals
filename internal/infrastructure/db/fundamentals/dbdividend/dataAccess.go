package dbdividend

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	utilities "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

const STOCK_DIRECTORY_PREFIX = "stockfundamentals/stocks" //FIXME: Extract into shared
const DIVIDEND_PAYMENT_TABLE_NAME = "dividend_payment"

type dividendDbModel struct {
	Id   uuid.UUID `sql:"id"`
	Figi string    `sql:"stock_id"`
	//For DB-related reasons, expected and actual DPS are converted to integers to remove the fractional part. Multiplied by a million for maximum accuracy. When reading the value, it must consequently be divided by a million
	ExpectedDpsTimesMillion int64     `sql:"expected_DPS"`
	ActualDPSTimesMillion   int64     `sql:"actual_DPS"`
	Currency                string    `sql:"currency"`
	AnnouncementDate        time.Time `sql:"announcement_date"`
	RecordDate              time.Time `sql:"record_date"`
	PayoutDate              time.Time `sql:"payout_date"`
	PaymentPeriod           string    `sql:"payment_period"`
	ManagementComment       string    `sql:"management_comment"`
}

func SaveDividendsToDB(dividends []dividend.Dividend, db *ydb.Driver) error {
	if len(dividends) == 0 {
		logger.Log("Attempting to save 0 dividends", logger.WARNING)
	}
	if db == nil {
		logger.Log("Database driver is nil while attempting to save dividends to the DB", logger.ALERT)
		return errors.New("Database issues")
	}

	dbModels := mapDividendToDbModel(dividends)

	ydbDividends := []types.Value{}
	for _, dividend := range dbModels {
		ydbDividend := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(dividend.Id)),
			types.StructFieldValue("stock_id", types.TextValue(dividend.Figi)),
			types.StructFieldValue("actual_DPS", types.Int64Value(int64(dividend.ActualDPSTimesMillion))),
			types.StructFieldValue("expected_DPS", types.Int64Value(int64(dividend.ExpectedDpsTimesMillion))),
			types.StructFieldValue("currency", types.TextValue(dividend.Currency)),
			types.StructFieldValue("announcement_date", convertToOptionalYDBdate(dividend.AnnouncementDate)),
			types.StructFieldValue("record_date", convertToOptionalYDBdate(dividend.RecordDate)),
			types.StructFieldValue("payout_date", convertToOptionalYDBdate(dividend.PayoutDate)),
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
		return errors.New("Failed to save dividends to the database")
	}

	return nil
}

func convertToOptionalYDBdate(date time.Time) types.Value {
	if date.Unix() == 0 || date.Unix() == -62135596800 {
		return types.NullValue(types.TypeDate)
	}

	const secondsInADay = 86400
	return types.DateValue(uint32(date.Unix() / secondsInADay))
}

func mapDividendToDbModel(dividends []dividend.Dividend) []dividendDbModel {
	dbModels := []dividendDbModel{}
	for _, dividend := range dividends {
		dbModel := dividendDbModel{
			Id:                      dividend.Id,
			Figi:                    dividend.Figi,
			ExpectedDpsTimesMillion: int64(dividend.ExpectedDPS * 1_000_000),
			ActualDPSTimesMillion:   int64(dividend.ActualDPS * 1_000_000),
			Currency:                dividend.Currency,
			AnnouncementDate:        dividend.AnnouncementDate,
			RecordDate:              dividend.RecordDate,
			PayoutDate:              dividend.PayoutDate,
			PaymentPeriod:           dividend.PaymentPeriod,
			ManagementComment:       dividend.ManagementComment,
		}
		dbModels = append(dbModels, dbModel)
	}

	return dbModels
}

func mapDbModelToDividend(dbModelds []dividendDbModel) []dividend.Dividend {
	dividends := []dividend.Dividend{}
	for _, dbModel := range dbModelds {
		newDiv := dividend.Dividend{
			Id:                dbModel.Id,
			Figi:              dbModel.Figi,
			ExpectedDPS:       float64(dbModel.ExpectedDpsTimesMillion) / 1_000_000,
			ActualDPS:         float64(dbModel.ActualDPSTimesMillion) / 1_000_000,
			Currency:          dbModel.Currency,
			AnnouncementDate:  dbModel.AnnouncementDate,
			RecordDate:        dbModel.RecordDate,
			PayoutDate:        dbModel.PayoutDate,
			PaymentPeriod:     dbModel.PaymentPeriod,
			ManagementComment: dbModel.ManagementComment,
		}
		dividends = append(dividends, newDiv)
	}

	return dividends
}

func GetAllDividends() ([]dividend.Dividend, error) {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return []dividend.Dividend{}, err
	}

	userDividendsDbModels := []dividendDbModel{}

	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, fmt.Sprintf(`
						SELECT
							id,
							stock_id,
							actual_DPS,
							expected_DPS,
							currency,
							announcement_date,
							record_date,
							payout_date,
							payment_period,
							management_comment
						FROM
							%s
					`, "`"+path.Join(STOCK_DIRECTORY_PREFIX, DIVIDEND_PAYMENT_TABLE_NAME)+"`"),
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

				for row, err := range sugar.UnmarshalRows[dividendDbModel](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}

					userDividendsDbModels = append(userDividendsDbModels, row)
				}
			}

			return nil
		},
	)
	if err != nil {
		fmt.Println(err)
		return []dividend.Dividend{}, err
	}

	return mapDbModelToDividend(userDividendsDbModels), nil
}

func GetUpcomingDividends() ([]dividend.Dividend, error) {
	allDividends, err := GetAllDividends()
	if err != nil {
		return []dividend.Dividend{}, err
	}

	upcomingDivs := []dividend.Dividend{}
	for _, div := range allDividends {
		if div.PayoutDate.After(time.Now()) {
			upcomingDivs = append(upcomingDivs, div)
		}
	}

	return upcomingDivs, nil
}
