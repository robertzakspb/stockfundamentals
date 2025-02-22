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

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

const DB_CONNECTION_STRING = "grpc://localhost:2136/local"

func InitialSeed() {
	ctx := context.TODO()
	db, err := ydb.Open(context.TODO(), DB_CONNECTION_STRING, ydb.WithDialTimeout(10))
	if err != nil {
		log.Fatal(err)
		panic("Failed to connect to the database")
	}
	createTables(ctx, db)
	populateTables(ctx, db)
}

const STOCK_DIRECTORY_PREFIX = "stockfundamentals/stocks"

func createTables(ctx context.Context, db *ydb.Driver) {
	client := db.Table()
	createCompanyTable(ctx, db, client)
}

func populateTables(ctx context.Context, db *ydb.Driver) error {
	client := db.Query()
	err := populateCompanyTable(ctx, client)
	if err != nil {
		fmt.Println(err)
		return err
	}
	
	return nil
}

func createCompanyTable(ctx context.Context, db *ydb.Driver, c table.Client) error {
	prefix := path.Join(db.Name(), STOCK_DIRECTORY_PREFIX)
	const stockTableName = "stock"

	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, stockTableName),
				options.WithColumn("id", types.TypeUTF8),
				options.WithColumn("company_name", types.Optional(types.TypeUTF8)),
				options.WithColumn("is_public", types.TypeBool),
				options.WithColumn("isin", types.TypeUTF8),
				options.WithColumn("security_type", types.TypeUTF8),
				options.WithColumn("country_iso2", types.TypeUTF8),
				options.WithColumn("ticker", types.TypeUTF8),
				options.WithColumn("share_count", types.TypeUint64),
				options.WithColumn("sector", types.Optional(types.TypeUTF8)),
				options.WithPrimaryKeyColumn("id", "isin", "ticker"),
			)
			if err != nil {
				fmt.Println(err)
				return err
			}

			return nil
		})
}



func populateCompanyTable(ctx context.Context, c query.Client) error {

	insertSecuritiesData, err := os.Open("dataseed/insertstock.yql")
	if err != nil {
		fmt.Println(err)
	}

	var buffer strings.Builder
	_, err = io.Copy(&buffer, insertSecuritiesData)
	if err != nil {
		fmt.Println(err)
	}
	insertSecuritiesYQL := buffer.String()

	c.Exec(ctx, insertSecuritiesYQL)

	return c.Do(ctx,
		func(ctx context.Context, s query.Session) (err error) {
			return nil
		})

}

// Example from the SDK; to be deleted
func createTabless(ctx context.Context, c table.Client, prefix string) error {
	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			err := s.CreateTable(ctx, path.Join(prefix, "series"),
				options.WithColumn("series_id", types.Optional(types.TypeUint64)),
				options.WithColumn("title", types.Optional(types.TypeUTF8)),
				options.WithColumn("series_info", types.Optional(types.TypeUTF8)),
				options.WithColumn("release_date", types.Optional(types.TypeUint64)),
				options.WithColumn("comment", types.Optional(types.TypeUTF8)),
				options.WithPrimaryKeyColumn("series_id"),
			)
			if err != nil {
				return err
			}

			err = s.CreateTable(ctx, path.Join(prefix, "seasons"),
				options.WithColumn("series_id", types.Optional(types.TypeUint64)),
				options.WithColumn("season_id", types.Optional(types.TypeUint64)),
				options.WithColumn("title", types.Optional(types.TypeUTF8)),
				options.WithColumn("first_aired", types.Optional(types.TypeUint64)),
				options.WithColumn("last_aired", types.Optional(types.TypeUint64)),
				options.WithPrimaryKeyColumn("series_id", "season_id"),
			)
			if err != nil {
				return err
			}

			err = s.CreateTable(ctx, path.Join(prefix, "episodes"),
				options.WithColumn("series_id", types.Optional(types.TypeUint64)),
				options.WithColumn("season_id", types.Optional(types.TypeUint64)),
				options.WithColumn("episode_id", types.Optional(types.TypeUint64)),
				options.WithColumn("title", types.Optional(types.TypeUTF8)),
				options.WithColumn("air_date", types.Optional(types.TypeUint64)),
				options.WithPrimaryKeyColumn("series_id", "season_id", "episode_id"),
			)
			if err != nil {
				return err
			}

			return nil
		},
	)
}

// Example from the YDB SDK
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
