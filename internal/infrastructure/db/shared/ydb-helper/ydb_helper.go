package ydbhelper

import (
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func ConvertStringsToYdbList(stringVals []string) types.Value {
	textVals := []types.Value{}
	for _, string := range stringVals {
		textVals = append(textVals, types.TextValue(string))
	}
	return types.ListValue(textVals...)
}

func ConvertToYdbDateTime(timestamp time.Time) types.Value {
	return types.DatetimeValue(uint32(timestamp.Unix()))
}

const secondsInADay = 60*60*24

func ConvertToYdbDate(date time.Time) types.Value {
	return types.DateValue(uint32(date.Unix() / secondsInADay))
}

func ConvertToOptionalYDBdate(date time.Time) types.Value {
	if date.Unix() == 0 || date.Unix() == -62135596800 {
		return types.NullValue(types.TypeDate)
	}

	return types.DateValue(uint32(date.Unix() / secondsInADay))
}
