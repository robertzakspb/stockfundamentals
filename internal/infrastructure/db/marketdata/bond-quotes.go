package timeseriesdb

import (
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	ydbtemplate "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-template"
)

func SaveBondQuotes(quotes []BondQuoteDB) error {
	ydbBondModels := mapBondQuoteDbModelToYdbEntity(quotes)
	tablePath := ydbhelper.GenerateTablePath(db.MARKET_DATA_DIRECTORY_PREFIX, db.BOND_TIME_SERIES_TABLE_NAME)

	err := ydbtemplate.SaveEntity(ydbBondModels, tablePath)

	return err
}
