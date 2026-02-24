package ydbfilter

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type YdbFilter struct {
	YqlColumnName  string
	Condition      YdbFilterCondition
	ConditionValue types.Value
}

type YdbFilterCondition string

const (
	GreaterThan          YdbFilterCondition = ">"
	GreaterThanOrEqualTo YdbFilterCondition = ">="
	LessThan             YdbFilterCondition = "<"
	LessThanOrEqualTo    YdbFilterCondition = "<="
	Contains             YdbFilterCondition = "IN"
	Equal                YdbFilterCondition = "="
)

var ydbConditions = map[string]YdbFilterCondition{
	">":  GreaterThan,
	">=": GreaterThanOrEqualTo,
	"<":  LessThan,
	"<=": LessThanOrEqualTo,
	"IN": Contains,
	"=":  Equal,
}

// This function only sets the values of query parameters themeselves but does not add them WHERE (done by AddWhereClause())
func SetQueryParams(filters []YdbFilter) *table.QueryParameters {
	params := []table.ParameterOption{}

	for _, filter := range filters {
		param := table.ValueParam(MakeColumnFilterName(filter.YqlColumnName), filter.ConditionValue)
		params = append(params, param)
	}

	return table.NewQueryParameters(params...)
}

func MakeWhereClause(filters []YdbFilter) string {
	if len(filters) == 0 {
		return ""
	}

	b := strings.Builder{}

	b.WriteString(" WHERE\n ")

	for i, filter := range filters {
		b.WriteString(filter.YqlColumnName)
		b.WriteString(" ")
		b.WriteString(string(filter.Condition))
		b.WriteString(" ")
		b.WriteString(MakeColumnFilterName(filter.YqlColumnName))
		if i < len(filters)-1 {
			b.WriteString(" ")
			b.WriteString("AND")
			b.WriteString(" ")
		}
	}

	return b.String()
}

func AddYqlVarDeclarations(filters []YdbFilter) string {
	if len(filters) == 0 {
		return ""
	}

	b := strings.Builder{}

	for _, filter := range filters {
		yqlVarName := MakeColumnFilterName(filter.YqlColumnName)
		b.WriteString(Declare(yqlVarName, filter.ConditionValue))
	}

	return b.String()
}

func MapQueryFiltersToYdb(filters map[string][]string, entity any) []YdbFilter {
	ydbFilters := []YdbFilter{}

	for parameter, values := range filters {
		v := reflect.ValueOf(entity).Elem()

		for i := 0; i < v.NumField(); i++ {
			jsonTagValue, found := v.Type().Field(i).Tag.Lookup("json")
			if !found {
				logger.Log("Failed to find the json tag in "+v.Type().Name()+" for field "+v.Type().Field(i).Name+" which is unexpected for a DTO struct", logger.ERROR)
				continue
			}
			if jsonTagValue == parameter {
				sqlTagValue, found := v.Type().Field(i).Tag.Lookup("sql")
				if !found {
					logger.Log("Failed to find the sql tag in "+v.Type().Name()+" for field "+v.Type().Field(i).Name+" which is unexpected for a DTO struct", logger.ERROR)
					continue
				}

				condition, err := mapQueryConditionToYdb(values[0])
				if err != nil {
					logger.Log("Failed to map the API query parameter – "+values[0]+"to a YDB filter. Fix the API call.", logger.ERROR)
					continue
				}

				filterValues, err := mapQueryValuesToYdbFilterValues(condition, values[1:])
				if err != nil {
					logger.Log("Failed to generate filter values", logger.ERROR)
					continue
				}

				ydbFilters = append(ydbFilters, YdbFilter{
					sqlTagValue, //FIXME
					condition,
					filterValues,
				})
			}
		}

	}

	return ydbFilters
}

func mapQueryConditionToYdb(condition string) (YdbFilterCondition, error) {
	ydbCondition, found := ydbConditions[condition]
	if !found {
		return GreaterThan, errors.New("Unknown condition")
	}

	return ydbCondition, nil
}

func mapQueryValuesToYdbFilterValues(condition YdbFilterCondition, values []string) (types.Value, error) {
	switch condition {
	case GreaterThan, GreaterThanOrEqualTo, LessThan, LessThanOrEqualTo:
		f, err := strconv.ParseFloat(values[0], 64)
		if err == nil {
			return types.DoubleValue(f), nil
		}

		i, err := strconv.Atoi(values[0])
		if err == nil {
			return types.Int64Value(int64(i)), nil
		}

		return types.UTF8Value(values[0]), nil

	case Equal:
		b, err := strconv.ParseBool(values[0])
		if err == nil {
			return types.BoolValue(b), nil
		}

		f, err := strconv.ParseFloat(values[0], 64)
		if err == nil {
			return types.DoubleValue(f), nil
		}

		i, err := strconv.Atoi(values[0])
		if err == nil {
			return types.Int64Value(int64(i)), nil
		}

		return types.UTF8Value(values[0]), nil
	case Contains:
		return parseArrayFromQueryParameters(values)
	}

	return types.NullValue(types.TypeBool), errors.New("failed to map query parameters to Ydb filter values")
}

func parseArrayFromQueryParameters(values []string) (types.Value, error) {
	if len(values) == 0 {
		return types.NullValue(types.TypeInt8), errors.New("Failed to parse query parameters, as the array is empty")
	}

	ydbValues := []types.Value{}

	//FIXME: Verify that all elements in the values []string array are of the same type; otherwise throw an error

	//Determine the array element type

	//First see if it's a float
	_, err := strconv.ParseFloat(values[0], 64)
	if err == nil {
		for _, value := range values {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return types.ListValue(), errors.New("First element in the array is a float but one subsequent value is not a float")
			}
			ydbValues = append(ydbValues, types.DoubleValue(f))
		}
	}

	//Then see if it's an int
	_, err = strconv.Atoi(values[0])
	if err == nil {
		for _, value := range values {
			i, err := strconv.Atoi(value)
			if err != nil {
				return types.ListValue(), errors.New("First element in the array is an integer but one subsequent value is not an integer")
			}
			ydbValues = append(ydbValues, types.Int64Value(int64(i)))
		}
	}

	//Then assume it's a string. Other parameter types will be implemented later
	for _, value := range values {
		ydbValues = append(ydbValues, types.UTF8Value(value))
	}

	return types.ListValue(), nil
}
