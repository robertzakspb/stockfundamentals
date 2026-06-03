package accountdb

import (
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func mapAccountDbModelToYdbEntity(dbModels []AccountDbModel) types.Value {
	dbAccounts := make([]types.Value, len(dbModels))
	for i, account := range dbModels {
		dbAccount := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(account.Id)),
			types.StructFieldValue("opening_date", ydbhelper.ConvertToYdbDate(account.OpeningDate)),
			types.StructFieldValue("type", types.TextValue(account.Type)),
			types.StructFieldValue("broker", types.TextValue(account.Broker)),
			types.StructFieldValue("holder", types.TextValue(account.Holder)),
			types.StructFieldValue("primary_currency", types.TextValue(account.PrimaryCurrency)),
			types.StructFieldValue("cash_balance", types.DoubleValue(account.CashBalance)),
		)
		dbAccounts[i] = dbAccount
	}

	return types.ListValue(dbAccounts...)

}
