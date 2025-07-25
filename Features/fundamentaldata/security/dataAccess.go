package security

import (
	"context"

	"errors"
	"fmt"
	"io"
	"path"

	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"

	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"

	"github.com/ydb-platform/ydb-go-sdk/v3/query"
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

func FetchSecuritiesFromDBWithDriver(db *ydb.Driver) ([]Stock, error) {
	dbStocks := []StockDbModel{}
	parsedStocks := []Stock{}
	err := db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, fmt.Sprintf(`
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
							sector
						FROM
							%s
					`, "`"+path.Join(stock_directory_prefix, stock_table_name)+"`"),
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
	}

	return stock
}
