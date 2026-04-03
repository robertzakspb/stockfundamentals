package dbdividend

import (
	"fmt"
	"path"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
)

func getDividendQuery(filters []ydbfilter.YdbFilter) string {
	return fmt.Sprintf(`
						%s

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
						%s
					`,
				ydbfilter.AddYqlVarDeclarations(filters),
				"`"+path.Join(db.STOCK_DIRECTORY_PREFIX, db.DIVIDEND_PAYMENT_TABLE_NAME)+"`",
				ydbfilter.MakeWhereClause(filters))
}