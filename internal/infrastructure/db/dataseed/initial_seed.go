package dataseed

import (
	"context"
	"encoding/csv"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	entity "github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	dbfinancials "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/financials"
	dbsecurity "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ydb-platform/ydb-go-sdk/v3"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func InitialSeed(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to proceed due to internal configuration issues")
	}

	db, err := ydb.Open(ctx, config.DB.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to proceed due to database issues")
		panic("Failed to connect to the database")
	}

	err = createTables(ctx, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to create tables")
	}

	err = populateTables(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to populate tables")
	}
}

func createTables(ctx context.Context, db *ydb.Driver) error {
	client := db.Table()

	err := createStockTables(ctx, db, client)
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

	err = CreateDividendForecastTable(ctx, db, client) //FIXME: make it private

	return nil
}

func populateTables(db *ydb.Driver) error {
	err := populateAllTables(db)
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

func createPortfolioTable(ctx context.Context, db *ydb.Driver, c table.Client) error {
	prefix := path.Join(db.Name(), shared.USER_DIRECTORY_PREFIX)
	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, shared.POSITION_LOT_TABLE_NAME),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("figi", types.TypeUTF8),
				options.WithColumn("account_id", types.TypeUUID),
				options.WithColumn("created_at", types.TypeDatetime),
				options.WithColumn("updated_at", types.TypeDatetime),
				options.WithColumn("quantity", types.TypeDouble),
				options.WithColumn("price_per_unit", types.TypeDouble),
				options.WithColumn("currency", types.TypeUTF8),
				options.WithPrimaryKeyColumn("id"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

func CreateDividendForecastTable(ctx context.Context, db *ydb.Driver, c table.Client) error {
	prefix := path.Join(db.Name(), shared.STOCK_DIRECTORY_PREFIX)
	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, shared.DIVIDEND_FORECAST_TABLE_NAME),
				options.WithColumn("id", types.TypeUUID),
				options.WithColumn("figi", types.TypeText),
				options.WithColumn("expected_DPS", types.Optional(types.TypeInt64)),
				options.WithColumn("currency", types.TypeText),
				options.WithColumn("payment_period", types.TypeText),
				options.WithColumn("forecast_author", types.TypeText),
				options.WithColumn("comment", types.TypeText),
				options.WithPrimaryKeyColumn("figi", "payment_period"),
			)
			if err != nil {
				logger.Log(err.Error(), logger.ALERT)
				return err
			}

			return nil
		})
}

const seedDataFolder = "internal/infrastructure/db/dataseed/seed-data/"

func populateAllTables(db *ydb.Driver) error {
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
			seedError = populateStockTable(csvReader, db)
		case seedDataFolder + "dividend-seed.csv":
			seedError = populateDividendTable(csvReader, db)
		case seedDataFolder + "revenue-income-seed.csv":
			seedError = populateFinancialMetricsTable(csvReader, db)
		default:
			logger.Log("Attempting to seed data from an unknow file: "+fileName, logger.ALERT)
		}
		if seedError != nil {
			return seedError
		}
	}

	return nil
}

func populateStockTable(reader *csv.Reader, db *ydb.Driver) error {
	seedRecords, err := reader.ReadAll()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	serbianStocks := []security.Security{}
	recordsLessHeader := seedRecords[1:]
	for _, record := range recordsLessHeader {
		isPublic, err := strconv.ParseBool(record[2])
		if err != nil {
			logger.Log("Failed to parse the is public flag "+record[2]+" in the stock seed file", logger.ALERT)
			continue
		}

		securityType, found := security.SecurityTypeMap[record[4]]
		if !found {
			logger.Log("Failed to parse the security type"+record[4]+" in the stock seed file", logger.ALERT)
			continue
		}

		issueSize, err := strconv.Atoi(record[7])
		if err != nil {
			logger.Log("Failed to parse the issue size "+record[7]+" in the stock seed file", logger.ALERT)
			continue
		}

		stock := security.Stock{
			Isin:         record[3],
			Figi:         record[9],
			CompanyName:  record[1],
			IsPublic:     isPublic,
			SecurityType: securityType,
			Country:      record[5],
			Ticker:       record[6],
			IssueSize:    issueSize,
			Sector:       record[8],
			MIC:          record[10],
		}
		serbianStocks = append(serbianStocks, stock)
	}

	err = dbsecurity.SaveSecuritiesToDB(serbianStocks, db)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}

	return nil
}

func populateDividendTable(reader *csv.Reader, db *ydb.Driver) error {
	seedRecords, err := reader.ReadAll()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	dividends := []dividend.Dividend{}
	csvDividends := seedRecords[1:]

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
	}

	err = dbdividend.SaveDividendsToDB(dividends, db)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}

	return nil
}

func populateFinancialMetricsTable(reader *csv.Reader, db *ydb.Driver) error {
	seedRecords, err := reader.ReadAll()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	csvMetrics := seedRecords[1:]
	metrics := []entity.FinancialMetric{}

	for _, csvMetric := range csvMetrics {
		parsedId, err := uuid.Parse(csvMetric[0])
		if err != nil {
			logger.Log("Failed to parse the metric ID from value "+csvMetric[0]+" in the revenue-income seed file", logger.ALERT)
			continue
		}

		parsedStockId, err := uuid.Parse(csvMetric[1])
		if err != nil {
			logger.Log("Failed to parse the stock ID from value "+csvMetric[1]+" in the revenue-income seed file", logger.ALERT)
			continue
		}

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

		metrics = append(metrics, entity.FinancialMetric{
			Id:       parsedId,
			StockId:  parsedStockId,
			Name:     csvMetric[2],
			Period:   entity.ReportingPeriodMap[csvMetric[3]],
			Year:     int(parsedYear),
			Value:    int(parsedValue),
			Currency: csvMetric[6],
		})
	}

	err = dbfinancials.SaveFinancialMetricsToDb(metrics, db)
	if err != nil {
		return err
	}

	return nil
}
