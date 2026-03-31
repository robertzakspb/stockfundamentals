package ydbhelper

import "github.com/ydb-platform/ydb-go-sdk/v3/table/types"

func ConvertStringsToYdbList(stringVals []string) types.Value {
	textVals := []types.Value{}
	for _, string := range stringVals {
		textVals = append(textVals, types.TextValue(string))
	}
	return types.ListValue(textVals...)
}
