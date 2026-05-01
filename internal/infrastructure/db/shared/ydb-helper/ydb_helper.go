package ydbhelper

import (
	"time"

	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

const secondsInADay = 60 * 60 * 24

func ConvertStringsToYdbList(stringVals []string) types.Value {
	textVals := make([]types.Value, len(stringVals))
	for i := range stringVals {
		textVals[i] = types.TextValue(stringVals[i])
	}

	return types.ListValue(textVals...)
}

func ConvertUUIDsToYdbList(uuids []uuid.UUID) types.Value {
	uuidValues := []types.Value{}
	for i := range uuids {
		uuidValues[i] = types.UuidValue(uuids[i])
	}

	return types.ListValue(uuidValues...)
}

func ConverTimestampsToYdbDates(timestamps []time.Time) types.Value {
	dateValues := make([]types.Value, len(timestamps))
	for i := range timestamps {
		dateValues[i] = ConvertToYdbDate(timestamps[i])
	}
	return types.ListValue(dateValues...)
}

func ConvertToYdbDateTime(timestamp time.Time) types.Value {
	return types.DatetimeValue(uint32(timestamp.Unix()))
}

func ConvertToYdbDate(date time.Time) types.Value {
	return types.DateValue(uint32(date.Unix() / secondsInADay))
}

func ConvertToOptionalYDBdate(date time.Time) types.Value {
	if date.Unix() == 0 || date.Unix() == -62135596800 {
		return types.NullValue(types.TypeDate)
	}

	return types.DateValue(uint32(date.Unix() / secondsInADay))
}
