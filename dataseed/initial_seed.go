package dataseed

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/compoundinvest/stockfundamentals/infrastructure/config"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)
const INSERT_SCRIPTS_FOLDER = "dataseed/yql_scripts/insert_scripts/"

func InitialSeed() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := ydb.Open(ctx, config.DB.ConnectionString)
	if err != nil {
		panic("Failed to connect to the database")
	}

	err = createTables(ctx, db)
	if err != nil {
		return err
	}
	
	err = populateTables(ctx, db)
	if err != nil {
		return err
	}

	return nil
}

const STOCK_DIRECTORY_PREFIX = "stockfundamentals/stocks"

func createTables(ctx context.Context, db *ydb.Driver) error {
	client := db.Table()
	err := createAllTables(ctx, db, client)
	return err
}

func populateTables(ctx context.Context, db *ydb.Driver) error {
	client := db.Query()

	err := populateAllTables(ctx, client)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func createAllTables(ctx context.Context, db *ydb.Driver, c table.Client) error {
	prefix := path.Join(db.Name(), STOCK_DIRECTORY_PREFIX)

	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, "stock"),
				options.WithColumn("id", types.TypeUTF8),
				options.WithColumn("figi", types.Optional(types.TypeUTF8)),
				options.WithColumn("company_name", types.Optional(types.TypeUTF8)),
				options.WithColumn("is_public", types.TypeBool),
				options.WithColumn("isin", types.TypeUTF8),
				options.WithColumn("security_type", types.TypeUTF8),
				options.WithColumn("country_iso2", types.TypeUTF8),
				options.WithColumn("ticker", types.TypeUTF8),
				options.WithColumn("issue_size", types.TypeInt64),
				options.WithColumn("sector", types.Optional(types.TypeUTF8)),
				options.WithPrimaryKeyColumn( "isin"),
			)
			if err != nil {
				fmt.Println(err)
				return err
			}

			err = s.CreateTable(ctx, path.Join(prefix, "dividend_payment"),
				options.WithColumn("id", types.TypeUTF8),
				options.WithColumn("stock_id", types.TypeUTF8),
				options.WithColumn("actual_DPS", types.TypeDouble),
				options.WithColumn("expected_DPS", types.Optional(types.TypeDouble)),
				options.WithColumn("currency", types.TypeUTF8),
				options.WithColumn("announcement_date", types.Optional(types.TypeDate)),
				options.WithColumn("record_date", types.TypeDate),
				options.WithColumn("payout_date", types.Optional(types.TypeDate)),
				options.WithColumn("payment_period", types.TypeUTF8),
				options.WithColumn("management_comment", types.Optional(types.TypeUTF8)),
				options.WithPrimaryKeyColumn("id"),
			)
			if err != nil {
				fmt.Println(err)
				return err
			}

			err = s.CreateTable(ctx, path.Join(prefix, "corporate_financials"),
				options.WithColumn("id", types.TypeUTF8),
				options.WithColumn("stock_id", types.TypeUTF8),
				options.WithColumn("financial_metric", types.TypeUTF8),
				options.WithColumn("reporting_period", types.TypeUTF8),
				options.WithColumn("metric_value", types.TypeDouble),
				options.WithColumn("metric_currency", types.TypeUTF8),
				options.WithPrimaryKeyColumn("id"),
			)
			if err != nil {
				fmt.Println(err)
				return err
			}

			return nil
		})
}

func populateAllTables(ctx context.Context, c query.Client) error {
	files, _ := os.ReadDir(INSERT_SCRIPTS_FOLDER)
	for _, file := range files {
		insertScriptData, err := os.Open(INSERT_SCRIPTS_FOLDER + file.Name())
		if err != nil {
			fmt.Println(err)
		}

		var buffer strings.Builder
		_, err = io.Copy(&buffer, insertScriptData)
		if err != nil {
			fmt.Println(err)
		}
		insertScriptYQL := buffer.String()

		c.Exec(ctx, insertScriptYQL)
	}

	return c.Do(ctx,
		func(ctx context.Context, s query.Session) (err error) {
			return nil
		})
}


// TODO: Delete. Example from the YDB SDK
func read(ctx context.Context, c query.Client) error {
	return c.Do(ctx,
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, fmt.Sprintf(`
					SELECT
						*
					FROM
						%s`, "`"+path.Join(STOCK_DIRECTORY_PREFIX, "stock")+"`"),
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
				for row, err := range sugar.UnmarshalRows[string](resultSet.Rows(ctx)) {
					if err != nil {
						fmt.Println(err)
						return err
					}

					log.Printf("%v:", row)
				}
			}
			return nil
		},
	)
}
