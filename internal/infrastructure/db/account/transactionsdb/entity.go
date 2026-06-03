package transactionsdb

import (
	"time"

	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type TransactionDbModel struct {
	Id           uuid.UUID `sql:"id"`
	AccountId    uuid.UUID `sql:"account_id"`
	Figi         string    `sql:"figi"`
	Type         string    `sql:"type"`
	Timestamp    time.Time `sql:"timestamp"`
	Side         string    `sql:"side"`
	Quantity     float64   `sql:"quantity"`
	PricePerUnit float64   `sql:"price_per_unit"`
	Currency     string    `sql:"currency"`
	Description  string    `sql:"description"`
}

func mapToYdbList(dbModels []TransactionDbModel) types.Value {
	dbTransactions := make([]types.Value, len(dbModels))
	for i, transaction := range dbModels {
		dbTransaction := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(transaction.Id)),
			types.StructFieldValue("account_id", types.UuidValue(transaction.AccountId)),
			types.StructFieldValue("figi", types.TextValue(transaction.Figi)),
			types.StructFieldValue("type", types.TextValue(string(transaction.Type))),
			types.StructFieldValue("timestamp", ydbhelper.ConvertToYdbDate(transaction.Timestamp)),
			types.StructFieldValue("side", types.TextValue(transaction.Side)),
			types.StructFieldValue("quantity", types.DoubleValue(transaction.Quantity)),
			types.StructFieldValue("price_per_unit", types.DoubleValue(transaction.PricePerUnit)),
			types.StructFieldValue("currency", types.TextValue(transaction.Currency)),
			types.StructFieldValue("description", types.TextValue(transaction.Description)),
		)
		dbTransactions[i] = dbTransaction
	}

	return types.ListValue(dbTransactions...)
}
