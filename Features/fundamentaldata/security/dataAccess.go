package security

import (
	"context"
	"strings"
	"time"

	"errors"
	"fmt"
	"io"
	"path"

	"github.com/compoundinvest/stockfundamentals/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"

	"github.com/google/uuid"

	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

const stock_directory_prefix = "stockfundamentals/stocks"
const stock_table_name = "stock"

func SaveSecuritiesToDB(securities []Security, db *ydb.Driver) error {
	ydbStocks := []types.Value{}
	for _, stock := range securities {
		var id = stock.GetId()
		if stock.GetId() == uuid.Nil {
			id = uuid.New()
		}

		ydbStock := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(id)),
			types.StructFieldValue("company_name", types.TextValue(stock.GetCompanyName())),
			types.StructFieldValue("is_public", types.BoolValue(true)),
			types.StructFieldValue("isin", types.TextValue(stock.GetIsin())),
			types.StructFieldValue("figi", types.TextValue(stock.GetFigi())),
			types.StructFieldValue("security_type", types.TextValue(string(stock.GetSecurityType()))),
			types.StructFieldValue("country_iso2", types.TextValue(stock.GetCountry())),
			types.StructFieldValue("ticker", types.TextValue(stock.GetTicker())),
			types.StructFieldValue("issue_size", types.Int64Value(int64(stock.GetIssueSize()))),
			types.StructFieldValue("sector", types.TextValue(stock.GetSector())),
			types.StructFieldValue("MIC", types.TextValue(stock.GetMic())),
		)

		ydbStocks = append(ydbStocks, ydbStock)
	}

	securityTableName := path.Join(db.Name(), stock_directory_prefix, stock_table_name)
	err := db.Table().BulkUpsert(
		context.TODO(),
		securityTableName,
		table.BulkUpsertDataRows(types.ListValue(ydbStocks...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	return nil
}

func GetAllSecuritiesFromDB() ([]Stock, error) {
	return FetchSecuritiesFromDBWithDriver(getSecuritiesBaseQuery())
}

func GetSecuritiesFilteredByFigi(figis []string) ([]Stock, error) {
	return FetchSecuritiesFromDBWithDriver(getSecuritiesFilteredByFigiQuery(figis))
}

func FetchSecuritiesFromDBWithDriver(yqlQuery string) ([]Stock, error) {

	config, err := config.LoadConfig()
	if err != nil {
		return []Stock{}, err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*15)
	defer cancel()
	db, err := ydb.Open(ctx, config.DB.ConnectionString)

	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}

	dbStocks := []StockDbModel{}
	parsedStocks := []Stock{}
	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx,
				yqlQuery,
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

				for row, err := range sugar.UnmarshalRows[StockDbModel](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}

					dbStocks = append(dbStocks, row)
				}
			}

			for _, stock := range dbStocks {
				parsedStocks = append(parsedStocks, mapYdbStockToStock(stock))
			}

			return nil
		},
	)
	if err != nil {
		return []Stock{}, err
	}

	return parsedStocks, nil
}

func getSecuritiesBaseQuery() string {
	yqlQuery := fmt.Sprintf(`
						SELECT
							id,
							isin,
							figi,
							company_name,
							is_public,
							security_type,
							country_iso2,
							ticker,
							issue_size,
							sector,
							MIC
						FROM
							%s
					`, securityPath())
	return yqlQuery
}

func getSecuritiesFilteredByFigiQuery(figis []string) string {
	yqlQuery := getSecuritiesBaseQuery()

	if len(figis) > 0 {
		yqlQuery += "WHERE figi IN " + convertSliceToYqlInExpression(figis)
	}

	return yqlQuery
}

func securityPath() string {
	return "`" + path.Join(stock_directory_prefix, stock_table_name) + "`"
}

func mapYdbStockToStock(dbStock StockDbModel) Stock {
	securityType, found := SecurityTypeMap[dbStock.SecurityType]
	if !found {
		logger.Log("Unable to parse the security type from the value: "+dbStock.SecurityType, logger.ERROR)
	}

	stock := Stock{
		Id:           dbStock.Id,
		CompanyName:  dbStock.CompanyName,
		IsPublic:     true,
		Isin:         dbStock.Isin,
		Figi:         dbStock.Figi,
		SecurityType: securityType,
		Country:      dbStock.Country,
		Ticker:       dbStock.Ticker,
		IssueSize:    int(dbStock.IssueSize),
		Sector:       dbStock.Sector,
		MIC:          dbStock.MIC,
	}

	return stock
}

// Converts slices like ["apple", "banana"] to an IN expression like IN ('apple', 'banana') for YQL filtering
func convertSliceToYqlInExpression(filterSlice []string) string {
	var sb strings.Builder
	sb.WriteString("(")

	for _, value := range filterSlice {
		sb.WriteString("'" + value + "', ")
	}

	str := sb.String()[:sb.Len()-2] + ")"
	return str
}
