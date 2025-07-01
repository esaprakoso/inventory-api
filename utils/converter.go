package utils

import (
	"fmt"
	"strconv"
)

func GetString(val any) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%.0f", v), nil
	case int:
		return strconv.Itoa(v), nil
	default:
		return "", fmt.Errorf("cannot convert %T to string", val)
	}
}

func GetInt(val any) (int, error) {
	switch v := val.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", val)
	}
}

func GetFloat64(val any) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", val)
	}
}
