package dataseed

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	financialsservice "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/financials"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
	entity "github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	dbsecurity "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/security"
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ydb-platform/ydb-go-sdk/v3"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func InitialSeed(c *gin.Context) {
	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ResponseError{Errors: []string{"Unable to proceed due to database issues"}})
		panic("Failed to connect to the database")
	}
	defer db.ReleaseDriver(dbConnection)

	err = createTables(context.TODO(), dbConnection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ResponseError{Errors: []string{"Failed to create tables"}})
	}

	err = PopulateTables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ResponseError{Errors: []string{"Failed to populate tables"}})
	}
}

func createTables(ctx context.Context, db *ydb.Driver) error {
	client := db.Table()

	err := createDividendForecastTable(ctx, db, client)
	if err != nil {
		return err
	}

	err = createTransactionLotRelationshipTable(ctx, db, client)
	if err != nil {
		return err
	}

	err = createTransactionsTable(ctx, db, client)
	if err != nil {
		return err
	}

	err = createBondLotTable(ctx, db, client)
	if err != nil {
		return err
	}

	err = createAccountMarketValueTable(ctx, db, client)
	if err != nil {
		return err
	}

	err = createAccountTable(ctx, db, client)
	if err != nil {
		return err
	}

	err = createFxRateTable(ctx, db, client)
	if err != nil {
		return err
	}

	err = createCouponTable(ctx, db, client)
	if err != nil {
		return err
	}

	err = createStockTables(ctx, db, client)
	if err != nil {
		return err
	}

	err = createMarketDataTables(ctx, db, client)
	if err != nil {
		return err
	}

	err = createPortfolioTable(ctx, db, client)
	if err != nil {
		return err
	}

	err = createBondTable(ctx, db, client)
	if err != nil {
		return err
	}

	return nil
}

func PopulateTables() error {
	err := PopulateAllTables()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}

	return nil
}

func createStockTables(ctx context.Context, db *ydb.Driver, c table.Client) error {
	prefix := path.Join(db.Name(), "stockfundamentals/stocks")

	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, "stock"),
				options.WithColumn("figi", types.TypeUTF8),
				options.WithColumn("company_name", types.Optional(types.TypeUTF8)),
				options.WithColumn("is_public", types.TypeBool),
				options.WithColumn("isin", types.TypeUTF8),
				options.WithColumn("security_type", types.TypeUTF8),
				options.WithColumn("country_iso2", types.TypeUTF8),
				options.WithColumn("MIC", types.TypeUTF8),
				options.WithColumn("ticker", types.TypeUTF8),
				options.WithColumn("issue_size", types.TypeInt64),
				options.WithColumn("sector", types.Optional(types.TypeUTF8)),
				options.WithPrimaryKeyColumn("figi"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			err = s.CreateTable(ctx, path.Join(prefix, "dividend_payment"),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("stock_id", types.TypeText),
				options.WithColumn("actual_DPS", types.TypeInt64),
				options.WithColumn("expected_DPS", types.Optional(types.TypeInt64)),
				options.WithColumn("currency", types.TypeUTF8),
				options.WithColumn("announcement_date", types.Optional(types.TypeDate)),
				options.WithColumn("record_date", types.TypeDate),
				options.WithColumn("payout_date", types.Optional(types.TypeDate)),
				options.WithColumn("payment_period", types.TypeUTF8),
				options.WithColumn("type", types.TypeUTF8),
				options.WithColumn("regularity", types.TypeUTF8),
				options.WithColumn("management_comment", types.Optional(types.TypeUTF8)),
				options.WithPrimaryKeyColumn("stock_id", "record_date", "actual_DPS"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			err = s.CreateTable(ctx, path.Join(prefix, "financial_metric"),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("stock_id", types.TypeUUID),
				options.WithColumn("metric", types.TypeUTF8),
				options.WithColumn("reporting_period", types.TypeUTF8),
				options.WithColumn("year", types.TypeInt64),
				options.WithColumn("metric_value", types.TypeInt64),
				options.WithColumn("metric_currency", types.TypeUTF8),
				options.WithPrimaryKeyColumn("id"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}
			return nil
		})
}

func createAccountTable(ctx context.Context, dbConnection *ydb.Driver, c table.Client) error {
	prefix := path.Join(dbConnection.Name(), db.USER_DIRECTORY_PREFIX)

	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, db.ACCOUNT_TABLE_NAME),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("opening_date", types.TypeDate),
				options.WithColumn("type", types.TypeText),
				options.WithColumn("broker", types.TypeText),
				options.WithColumn("holder", types.TypeText),
				options.WithColumn("primary_currency", types.TypeText),
				options.WithColumn("cash_balance", types.TypeDouble),
				options.WithPrimaryKeyColumn("id"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

func createAccountMarketValueTable(ctx context.Context, dbConnection *ydb.Driver, c table.Client) error {
	prefix := path.Join(dbConnection.Name(), db.USER_DIRECTORY_PREFIX)

	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, db.ACCOUNT_MARKET_VALUE_HISTORY_TABLE_NAME),
				options.WithColumn("account_id", types.TypeUUID),
				options.WithColumn("date", types.TypeDate),
				options.WithColumn("currency", types.TypeText),
				options.WithColumn("eod_value", types.TypeDouble),
				options.WithPrimaryKeyColumn("account_id", "date", "currency"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

func createTransactionLotRelationshipTable(ctx context.Context, dbConnection *ydb.Driver, c table.Client) error {
	prefix := path.Join(dbConnection.Name(), db.USER_DIRECTORY_PREFIX)

	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, db.TRANSACTION_LOT_RELATIONSHIP_TABLE_NAME),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("transaction_id", types.TypeUUID),
				options.WithColumn("stock_lot_id", types.TypeUUID),
				options.WithColumn("bond_lot_id", types.TypeUUID),
				options.WithColumn("date", types.TypeDate),
				options.WithColumn("quantity", types.TypeDouble),
				options.WithPrimaryKeyColumn("id"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

func createTransactionsTable(ctx context.Context, dbConnection *ydb.Driver, c table.Client) error {
	prefix := path.Join(dbConnection.Name(), db.USER_DIRECTORY_PREFIX)

	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, db.TRANSACTION_TABLE_NAME),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("account_id", types.TypeUUID),
				options.WithColumn("figi", types.TypeText),
				options.WithColumn("type", types.TypeText),
				options.WithColumn("timestamp", types.TypeDate),
				options.WithColumn("side", types.TypeText),
				options.WithColumn("quantity", types.TypeDouble),
				options.WithColumn("price_per_unit", types.TypeDouble),
				options.WithColumn("currency", types.TypeText),
				options.WithColumn("description", types.TypeText),
				options.WithPrimaryKeyColumn("id"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

func createMarketDataTables(ctx context.Context, db *ydb.Driver, c table.Client) error {
	prefix := path.Join(db.Name(), "marketdata/")

	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, "time_series"),
				options.WithColumn("figi", types.TypeUTF8),
				options.WithColumn("close_price", types.TypeDouble),
				options.WithColumn("date", types.TypeDate),
				options.WithPrimaryKeyColumn("figi", "date"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

func createBondLotTable(ctx context.Context, dbConnection *ydb.Driver, c table.Client) error {
	prefix := path.Join(dbConnection.Name(), db.BOND_DIRECTORY_PREFIX)
	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, db.BOND_POSITION_LOT_TABLE_NAME),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("figi", types.TypeText),
				options.WithColumn("isin", types.TypeText),
				options.WithColumn("opening_date", types.TypeDatetime),
				options.WithColumn("modification_date", types.TypeDatetime),
				options.WithColumn("account_id", types.TypeUUID),
				options.WithColumn("quantity", types.TypeDouble),
				options.WithColumn("price_per_unit_percentage", types.TypeDouble),
				options.WithPrimaryKeyColumn("figi", "account_id"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

func createFxRateTable(ctx context.Context, dbConnection *ydb.Driver, c table.Client) error {
	prefix := path.Join(dbConnection.Name(), db.FOREX_DIRECTORY_PREFIX)
	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, db.FX_RATE_TABLE_NAME),
				options.WithColumn("currency_1", types.TypeText),
				options.WithColumn("currency_2", types.TypeText),
				options.WithColumn("date", types.TypeDate),
				options.WithColumn("rate", types.TypeDouble),
				options.WithPrimaryKeyColumn("currency_1", "currency_2", "date"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

func createPortfolioTable(ctx context.Context, dbConnection *ydb.Driver, c table.Client) error {
	prefix := path.Join(dbConnection.Name(), db.USER_DIRECTORY_PREFIX)
	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, db.POSITION_LOT_TABLE_NAME),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("figi", types.TypeUTF8),
				options.WithColumn("account_id", types.TypeUUID),
				options.WithColumn("created_at", types.TypeDatetime),
				options.WithColumn("updated_at", types.TypeDatetime),
				options.WithColumn("quantity", types.TypeDouble),
				options.WithColumn("price_per_unit", types.TypeDouble),
				options.WithColumn("currency", types.TypeUTF8),
				options.WithColumn("is_closed", types.TypeBool),
				options.WithPrimaryKeyColumn("id"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

func createDividendForecastTable(ctx context.Context, dbConnection *ydb.Driver, c table.Client) error {
	prefix := path.Join(dbConnection.Name(), db.STOCK_DIRECTORY_PREFIX)
	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, db.DIVIDEND_FORECAST_TABLE_NAME),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("figi", types.TypeText),
				options.WithColumn("expected_DPS", types.TypeDouble),
				options.WithColumn("currency", types.TypeText),
				options.WithColumn("payment_period", types.TypeText),
				options.WithColumn("forecast_author", types.TypeText),
				options.WithColumn("comment", types.TypeText),
				options.WithColumn("payout_date", types.TypeDate),
				options.WithPrimaryKeyColumn("figi", "payout_date"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

func createCouponTable(ctx context.Context, dbConnection *ydb.Driver, c table.Client) error {
	prefix := path.Join(dbConnection.Name(), db.BOND_DIRECTORY_PREFIX)
	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, db.COUPON_TABLE_NAME),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("figi", types.TypeText),
				options.WithColumn("coupon_date", types.TypeDate),
				options.WithColumn("record_date", types.TypeDate),
				options.WithColumn("coupon_number", types.TypeInt64),
				options.WithColumn("per_bond_amount", types.TypeDouble),
				options.WithColumn("coupon_type", types.TypeText),
				options.WithColumn("coupon_start_date", types.Optional(types.TypeDate)),
				options.WithColumn("coupon_end_date", types.Optional(types.TypeDate)),
				options.WithColumn("coupon_period", types.TypeInt64),
				options.WithPrimaryKeyColumn("figi", "coupon_date"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

func createBondTable(ctx context.Context, dbConnection *ydb.Driver, c table.Client) error {
	prefix := path.Join(dbConnection.Name(), db.BOND_DIRECTORY_PREFIX)
	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, db.BOND_TABLE_NAME),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("figi", types.TypeText),
				options.WithColumn("isin", types.TypeText),
				options.WithColumn("lot", types.TypeInt64),
				options.WithColumn("currency", types.TypeText),
				options.WithColumn("name", types.TypeText),
				options.WithColumn("country_of_risk", types.TypeText),
				options.WithColumn("real_exchange", types.TypeText),
				options.WithColumn("coupon_count_per_year", types.TypeInt64),
				options.WithColumn("maturity_date", types.Optional(types.TypeDate)),
				options.WithColumn("nominal_value", types.TypeDouble),
				options.WithColumn("nominal_currency", types.TypeText),
				options.WithColumn("initial_nominal_value", types.TypeDouble),
				options.WithColumn("initial_nominal_currency", types.TypeText),
				options.WithColumn("registration_date", types.Optional(types.TypeDate)),
				options.WithColumn("placement_date", types.Optional(types.TypeDate)),
				options.WithColumn("placement_price", types.TypeDouble),
				options.WithColumn("placement_currency", types.TypeText),
				options.WithColumn("accumulated_coupon_income", types.TypeDouble),
				options.WithColumn("issue_size", types.TypeInt64),
				options.WithColumn("issue_size_plan", types.TypeInt64),
				options.WithColumn("has_floating_coupon", types.TypeBool),
				options.WithColumn("is_perpetual", types.TypeBool),
				options.WithColumn("has_amortization", types.TypeBool),
				options.WithColumn("is_available_for_iis", types.TypeBool),
				options.WithColumn("is_for_qualified_investors", types.TypeBool),
				options.WithColumn("is_subordinated", types.TypeBool),
				options.WithColumn("risk_level", types.TypeText),
				options.WithColumn("bond_type", types.TypeText),
				options.WithColumn("call_option_exercise_date", types.Optional(types.TypeDate)),
				options.WithPrimaryKeyColumn("figi", "isin"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

const seedDataFolder = "internal/infrastructure/db/dataseed/seed-data/"

func PopulateAllTables() error {
	files, err := os.ReadDir(seedDataFolder)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	for _, file := range files {
		fileName := seedDataFolder + file.Name()
		file, err := os.Open(fileName)
		if err != nil {
			logger.Log(err.Error(), logger.ALERT)
		}
		defer file.Close()

		csvReader := csv.NewReader(file)
		csvReader.Comma = '|'

		var seedError error = nil
		switch fileName {
		case seedDataFolder + "security-seed.csv":
			seedError = populateStockTable(csvReader)
		case seedDataFolder + "dividend-seed.csv":
			seedError = populateDividendTable(csvReader)
		case seedDataFolder + "revenue-income-seed.csv":
			seedError = populateFinancialMetricsTable(csvReader)
		default:
			logger.Log("Attempting to seed data from an unknow file: "+fileName, logger.ALERT)
		}
		if seedError != nil {
			fmt.Println(seedError)
		}
	}

	return nil
}

func populateStockTable(reader *csv.Reader) error {
	seedRecords, err := reader.ReadAll()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	serbianStocks := []security.Security{}
	recordsLessHeader := seedRecords[1:]
	str := ""
	str += "UPSERT INTO `stockfundamentals/stocks/stock` (figi, company_name, is_public, isin, security_type, country_iso2, MIC, ticker, issue_size, sector) VALUES \n"

	for _, record := range recordsLessHeader {
		isPublic, err := strconv.ParseBool(record[1])
		if err != nil {
			logger.Log("Failed to parse the is public flag "+record[1]+" in the stock seed file", logger.ALERT)
			continue
		}

		securityType, found := security.SecurityTypeMap[record[3]]
		if !found {
			logger.Log("Failed to parse the security type"+record[3]+" in the stock seed file", logger.ALERT)
			continue
		}

		issueSize, err := strconv.Atoi(record[6])
		if err != nil {
			logger.Log("Failed to parse the issue size "+record[6]+" in the stock seed file", logger.ALERT)
			continue
		}

		stock := security.Stock{
			Isin:         record[2],
			Figi:         record[8],
			CompanyName:  record[0],
			IsPublic:     isPublic,
			SecurityType: securityType,
			Country:      record[4],
			Ticker:       record[5],
			IssueSize:    issueSize,
			Sector:       record[7],
			MIC:          record[9],
		}
		serbianStocks = append(serbianStocks, stock)
		str += "("
		str += "'" + stock.Figi + "'" + ", "
		str += "'" + stock.CompanyName + "'" + ", "
		str += strconv.FormatBool(stock.IsPublic) + ", "
		str += "'" + stock.Isin + "'" + ", "
		str += "'" + string(stock.SecurityType) + "'" + ", "
		str += "'" + stock.Country + "'" + ", "
		str += "'" + stock.MIC + "'" + ", "
		str += "'" + stock.Ticker + "'" + ", "
		str += strconv.Itoa(stock.IssueSize) + ", "
		str += "'" + stock.Sector + "'"
		// if i < len(recordsLessHeader)-1 {
		// 	str += ", "
		// }
		str += "),\n"

	}
	str += ";"
	fmt.Println(str)

	err = dbsecurity.SaveSecuritiesToDB(serbianStocks)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}

	return nil
}

func populateDividendTable(reader *csv.Reader) error {
	seedRecords, err := reader.ReadAll()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	dividends := []dividend.Dividend{}
	csvDividends := seedRecords[1:]

	yql := "UPSERT INTO `stockfundamentals/stocks/dividend_payment` (id, stock_id, actual_DPS, expected_DPS, currency, record_date, payout_date, payment_period, regularity, type) VALUES\n"
	for _, csvDividend := range csvDividends {
		parsedId, err := uuid.Parse(csvDividend[0])
		if err != nil {
			logger.Log("Failed to parse the dividend ID from value "+csvDividend[0]+" in the dividend seed file", logger.ALERT)
			continue
		}

		actualDPS, err := strconv.ParseFloat(csvDividend[2], 64)
		if err != nil {
			logger.Log("Failed to parse the actual DPS from value "+csvDividend[2]+" in the dividend seed file", logger.ALERT)
			continue
		}

		expectedDPS, err := strconv.ParseFloat(csvDividend[3], 64)
		if err != nil {
			logger.Log("Failed to parse the expected DPS from value "+csvDividend[3]+" in the dividend seed file", logger.ALERT)
			continue
		}

		recordDate, err := time.Parse("2006-01-02", csvDividend[5])
		if err != nil {
			logger.Log("Failed to parse the record date from value "+csvDividend[5]+" in the dividend seed file", logger.ALERT)
			continue
		}

		payoutDate, err := time.Parse("2006-01-02", csvDividend[6])
		if err != nil {
			logger.Log("Failed to parse the payout date from value "+csvDividend[6]+" in the dividend seed file", logger.WARNING)
			payoutDate = time.Unix(0, 0)
		}

		div := dividend.Dividend{
			Id:            parsedId,
			Figi:          csvDividend[1],
			ActualDPS:     actualDPS,
			ExpectedDPS:   expectedDPS,
			Currency:      csvDividend[4],
			RecordDate:    recordDate,
			PayoutDate:    payoutDate,
			PaymentPeriod: csvDividend[7],
		}
		dividends = append(dividends, div)

		yql += "("
		yql += "Uuid('" + div.Id.String() + "'), "
		yql += "'" + div.Figi + "'" + ", "
		yql += strconv.Itoa(int(div.ActualDPS*1_000_000)) + ", "
		yql += strconv.Itoa(int(div.ExpectedDPS*1_000_000)) + ", "
		yql += "'" + div.Currency + "'" + ", "
		yql += "Date('" + dateToISOString(div.RecordDate) + "'), "
		yql += "Date('" + dateToISOString(div.PayoutDate) + "'), "
		yql += "'" + div.PaymentPeriod + "', "
		yql += "'" + div.Regularity + "', "
		yql += "'" + div.Type + "'"
		// if i < len(csvDividends)-1 {
		// 	yql += ", "
		// }
		yql += "),\n"
	}
	yql += ";"
	fmt.Println(yql)

	err = dbdividend.SaveDividendsToDB(&dividends)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}

	return nil
}

func dateToISOString(date time.Time) string {
	year, month, day := date.Date()
	monthStr := strconv.Itoa(int(month))
	if len(monthStr) == 1 {
		monthStr = "0" + monthStr
	}
	dayStr := strconv.Itoa(day)
	if len(monthStr) == 1 {
		dayStr = "0" + dayStr
	}
	return strconv.Itoa(year) + "-" + monthStr + "-" + dayStr
}

func populateFinancialMetricsTable(reader *csv.Reader) error {
	seedRecords, err := reader.ReadAll()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	csvMetrics := seedRecords[1:]
	metrics := []entity.FinancialMetric{}

	yql := "UPSERT INTO `stockfundamentals/stocks/financial_metric` (id, figi, metric, reporting_period, year, metric_value, metric_currency) VALUES\n"
	for _, csvMetric := range csvMetrics {
		parsedId, err := uuid.Parse(csvMetric[0])
		if err != nil {
			logger.Log("Failed to parse the metric ID from value "+csvMetric[0]+" in the revenue-income seed file", logger.ALERT)
			continue
		}

		parsedStockId := csvMetric[1]

		parsedYear, err := strconv.ParseInt(csvMetric[4], 0, 64)
		if err != nil {
			logger.Log("Failed to parse the year from value "+csvMetric[4]+" in the revenue-income seed file", logger.ALERT)
			continue
		}

		parsedValue, err := strconv.ParseInt(csvMetric[5], 0, 64)
		if err != nil {
			logger.Log("Failed to parse the metric value from value "+csvMetric[5]+" in the revenue-income seed file", logger.ALERT)
			continue
		}

		reportingPeriod, found := entity.ReportingPeriodMap[csvMetric[3]]
		if !found {
			logger.Log("Attempting to save a financial metric with an unparsable reporting period: "+csvMetric[3], logger.ERROR)
			continue
		}
		metric := entity.FinancialMetric{
			Id:              parsedId,
			StockId:         parsedStockId,
			Name:            financials.MetricMap[csvMetric[2]],
			ReportingPeriod: reportingPeriod,
			Year:            int(parsedYear),
			Value:           int(parsedValue),
			Currency:        csvMetric[6],
		}
		metrics = append(metrics, metric)

		yql += "("
		yql += "Uuid('" + metric.Id.String() + "'), "
		yql += "'" + metric.StockId + "'" + ", "
		yql += "'" + string(metric.Name) + "'" + ", "
		yql += "'" + string(metric.ReportingPeriod) + "'" + ", "
		yql += strconv.Itoa(metric.Year) + ", "
		yql += strconv.Itoa(metric.Value) + ", "
		yql += "'" + metric.Currency + "'"
		// if i < len(csvMetrics)-1 {
		// 	yql += ", "
		// }
		yql += "),\n"
	}
	yql += ";"
	fmt.Println(yql)

	err = financialsservice.SaveFinancialMetrics(metrics)
	if err != nil {
		return err
	}

	return nil
}
