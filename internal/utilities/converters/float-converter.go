package typeconverter

import (
	"errors"
	"math"
	"strconv"
)

func GetFloat(unk any) (float64, error) {
	switch v := unk.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return math.NaN(), err
		}
		return f, nil
	default:
		return math.NaN(), errors.New("unknown value is of incompatible type")
	}
}
