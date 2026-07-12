package ydbfilter

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

func MapQueryFiltersToYdb[API any](filters map[string][]string) []YdbFilter {
	ydbFilters := []YdbFilter{}

	for parameter, queryValues := range filters {
		for _, queryValue := range queryValues {
			values := strings.Split(queryValue, ",") //We assume that any given query parameter has only 1 string
			filter, err := convertQueryParamToYdbFilter[API](parameter, values)
			if err != nil {
				logger.Log(err.Error(), logger.ERROR)
				continue
			}
			ydbFilters = append(ydbFilters, filter)
		}
	}

	return ydbFilters
}

func convertQueryParamToYdbFilter[API any](parameter string, values []string) (YdbFilter, error) {
	var jsonEntity API
	jsonReflection := reflect.ValueOf(jsonEntity)
	for i := 0; i < jsonReflection.NumField(); i++ {
		jsonTagValue, found := jsonReflection.Type().Field(i).Tag.Lookup("json")
		if !found {
			logger.Log("Failed to find the json tag in "+jsonReflection.Type().Name()+" for field "+jsonReflection.Type().Field(i).Name+" which is unexpected for a DTO struct", logger.WARNING)
			continue
		}
		if jsonTagValue != parameter {
			continue
		}

		sqlTagValue, found := jsonReflection.Type().Field(i).Tag.Lookup("sql")
		if !found {
			logger.Log("Failed to find the sql tag in "+jsonReflection.Type().Name()+" for field "+jsonReflection.Type().Field(i).Name+" which is unexpected for a DTO struct", logger.ERROR)
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

		filter := YdbFilter{
			sqlTagValue,
			condition,
			filterValues,
		}
		return filter, nil

	}

	return YdbFilter{}, errors.New("Failed to generate a YDB filter with provided values")
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
		date, err := time.Parse("2006-01-02", values[0])
		if err == nil {
			return ydbhelper.ConvertToYdbDate(date), nil
		}

		f, err := strconv.ParseFloat(values[0], 64)
		if err == nil {
			return types.DoubleValue(f), nil
		}

		i, err := strconv.Atoi(values[0])
		if err == nil {
			return types.Int64Value(int64(i)), nil
		}

		id, err := uuid.Parse(values[0])
		if err == nil {
			return types.UuidValue(id), nil
		}

		return types.TextValue(values[0]), nil

	case Equal:
		date, err := time.Parse("2006-01-02", values[0])
		if err == nil {
			return ydbhelper.ConvertToYdbDate(date), nil
		}

		b, err := strconv.ParseBool(values[0])
		if err == nil {
			return types.BoolValue(b), nil
		}

		f, err := strconv.ParseFloat(values[0], 64)
		if err == nil {
			return types.DoubleValue(f), nil
		}

		id, err := uuid.Parse(values[0])
		if err == nil {
			return types.UuidValue(id), nil
		}

		i, err := strconv.Atoi(values[0])
		if err == nil {
			return types.Int64Value(int64(i)), nil
		}

		return types.TextValue(values[0]), nil
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
		return types.ListValue(ydbValues...), nil
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
		return types.ListValue(ydbValues...), nil

	}

	_, err = uuid.Parse(values[0])
	if err == nil {
		for _, value := range values {
			id, err := uuid.Parse(value)
			if err != nil {
				return types.ListValue(), errors.New("First element in the array is a UUID but one subsequent value is not an UUID")
			}
			ydbValues = append(ydbValues, types.UuidValue(id))
		}
		return types.ListValue(ydbValues...), nil

	}

	_, err = time.Parse(time.RFC3339, values[0])
	if err == nil {
		for _, value := range values {
			date, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return types.ListValue(), errors.New("First element in the array is a date but one subsequent value is not an date")
			}
			ydbValues = append(ydbValues, ydbhelper.ConvertToYdbDate(date))
		}
		return types.ListValue(ydbValues...), nil

	}

	//Then assume it's a string. Other parameter types will be implemented later
	for _, value := range values {
		ydbValues = append(ydbValues, types.TextValue(value))
	}

	return types.ListValue(ydbValues...), nil
}
