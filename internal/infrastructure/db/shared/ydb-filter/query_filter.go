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

func MapQueryFiltersToYdb[API any](filters map[string][]string) ([]YdbFilter, error) {
	ydbFilters := []YdbFilter{}

	for parameter, queryValues := range filters {
		for _, queryValue := range queryValues {
			values := strings.Split(queryValue, ",") //We assume that any given query parameter has only 1 string
			//The parameter must contain at least two parameters: 1) the condition, 2) the filter value
			if len(values) < 2 || values[1] == "" {
				return ydbFilters, errors.New("Expected at least two comma-delineated parameters in the query value: " + parameter + ": " + queryValue)
			}
			filter, err := convertQueryParamToYdbFilter[API](parameter, values)
			if err != nil {
				logger.Log(err.Error(), logger.ERROR)
				continue
			}
			ydbFilters = append(ydbFilters, filter)
		}
	}

	return ydbFilters, nil
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

		filterValues, err := mapQueryValuesToYdbFilterValues(condition, values[1:], jsonReflection.Type().Field(i).Type.String())
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

func mapQueryValuesToYdbFilterValues(condition YdbFilterCondition, values []string, typeName string) (types.Value, error) {
	if condition == Contains {
		return parseArrayFromQueryParameters(values, typeName)
	}

	ydbTypeValue, err := convertParameterToYdbTypeValue(typeName, values[0])
	return ydbTypeValue, err
}

func convertParameterToYdbTypeValue(typeName string, value string) (types.Value, error) {
	switch typeName {
	case "bool":
		b, err := strconv.ParseBool(value)
		if err == nil {
			return types.BoolValue(b), nil
		}
	case "string":
		return types.TextValue(value), nil
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		i, err := strconv.Atoi(value)
		if err != nil {
			return types.TextValue(value), err
		}
		return types.Int64Value(int64(i)), nil
	case "float32", "float64":
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return types.TextValue(value), err
		}
		return types.DoubleValue(f), nil

	case "time.Time":
		date, err := time.Parse("2006-01-02", value)
		if err == nil {
			return ydbhelper.ConvertToYdbDate(date), nil
		}
		timestamp, err := time.Parse(time.RFC3339, value)
		if err == nil {
			return ydbhelper.ConvertToYdbDateTime(timestamp), nil
		}
		return types.TextValue(value), err

	case "uuid.UUID":
		id, err := uuid.Parse(value)
		if err != nil {
			return types.TextValue(value), err
		}
		return types.UuidValue(id), nil
	}

	return types.TextValue(value), nil
}

func parseArrayFromQueryParameters(values []string, typeName string) (types.Value, error) {
	if len(values) == 0 {
		return types.NullValue(types.TypeInt8), errors.New("Failed to parse query parameters, as the array is empty")
	}
	ydbValues := make([]types.Value, len(values))

	for i := range values {
		ydbTypeValue, err := convertParameterToYdbTypeValue(typeName, values[i])
		if err != nil {
			return types.NullValue(types.TypeBool), err
		}
		ydbValues[i] = ydbTypeValue
	}

	return types.ListValue(ydbValues...), nil
}
