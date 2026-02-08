package utilities

import (
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func ConvertToYdbDate(date time.Time) types.Value {
	const secondsInADay = 86400
	return types.DateValue(uint32(date.Unix() / secondsInADay))
}
