package figtree

import (
	"fmt"
	"strconv"
	"strings"
)

// INTERFACE TYPE CONVERSIONS

// toInt returns an interface{} as an int or returns an error
func toInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case *figFlesh:
		return toInt(v.AsIs())
	case int:
		return v, nil
	case *int:
		return *v, nil
	case int64:
		return int(v), nil
	case *int64:
		return int(*v), nil
	case *float64:
		return int(*v), nil
	case float64:
		return int(v), nil
	case *string:
		if f, err := strconv.ParseFloat(*v, 64); err == nil {
			return int(f), nil
		}
		return strconv.Atoi(*v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return int(f), nil
		}
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %v to int", value)
	}
}

// toInt64 returns an interface{} as an int64 or returns an error
func toInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case *figFlesh:
		return toInt64(v.AsIs())
	case int:
		return int64(v), nil
	case *int:
		return int64(*v), nil
	case int64:
		return v, nil
	case *int64:
		return int64(*v), nil
	case *float64:
		return int64(*v), nil
	case float64:
		return int64(v), nil
	case *string:
		if f, err := strconv.ParseFloat(*v, 64); err == nil {
			return int64(f), nil
		}
		return strconv.ParseInt(*v, 10, 64)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return int64(f), nil
		}
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %v to int64", value)
	}
}

// toFloat64 returns an interface{} as an float64 or returns an error
func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case *figFlesh:
		return toFloat64(v.AsIs())
	case *float64:
		return *v, nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case *int:
		return float64(*v), nil
	case int64:
		return float64(v), nil
	case *int64:
		return float64(*v), nil
	case *string:
		return strconv.ParseFloat(*v, 64)
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %v to float64", value)
	}
}

// toString returns an interface{} as a string or returns an error
func toString(value interface{}) (string, error) {
	switch v := value.(type) {
	case *figFlesh:
		return toString(v.AsIs())
	case *string:
		return *v, nil
	case string:
		return v, nil
	case *float64:
		return strconv.FormatFloat(*v, 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case *bool:
		return strconv.FormatBool(*v), nil
	case bool:
		return strconv.FormatBool(v), nil
	case []string:
		return strings.Join(v, ","), nil
	case *[]string:
		return strings.Join(*v, ","), nil
	case *map[string]string:
		parts := make([]string, 0, len(*v))
		for k, x := range *v {
			parts = append(parts, fmt.Sprintf("%s=%s", k, x))
		}
		return strings.Join(parts, ","), nil
	case map[string]string:
		parts := make([]string, 0, len(v))
		for k, x := range v {
			parts = append(parts, fmt.Sprintf("%s=%s", k, x))
		}
		return strings.Join(parts, ","), nil
	default:
		return "", fmt.Errorf("cannot convert %v to string", value)
	}
}

// toBool returns an interface{} as a bool or returns an error
func toBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case *figFlesh:
		return toBool(v.AsIs())
	case *string:
		return strconv.ParseBool(*v)
	case string:
		return strconv.ParseBool(v)
	case *bool:
		return *v, nil
	case bool:
		return v, nil
	default:
		return false, fmt.Errorf("cannot convert %v to bool", value)
	}
}

// toStringSlice returns an interface{} as a []string{} or returns an error
func toStringSlice(value interface{}) ([]string, error) {
	switch v := value.(type) {
	case *figFlesh:
		return toStringSlice(v.AsIs())
	case []string:
		return v, nil
	case *[]string:
		return *v, nil
	case []interface{}:
		var result []string
		for _, item := range v {
			str, err := toString(item)
			if err != nil {
				return nil, err
			}
			result = append(result, str)
		}
		return result, nil
	case *string:
		if *v == "" {
			return []string{}, nil
		}
		return strings.Split(*v, ","), nil
	case string:
		if v == "" {
			return []string{}, nil
		}
		return strings.Split(v, ","), nil
	default:
		return nil, fmt.Errorf("cannot convert %v to []string", value)
	}
}

// toStringMap returns an interface{} as a map[string]string or returns an error
func toStringMap(value interface{}) (map[string]string, error) {
	switch v := value.(type) {
	case *figFlesh:
		return toStringMap(v.AsIs())
	case map[string]string:
		return v, nil
	case *map[string]string:
		return *v, nil
	case map[string]interface{}:
		result := make(map[string]string)
		for key, val := range v {
			strVal, err := toString(val)
			if err != nil {
				return nil, err
			}
			result[key] = strVal
		}
		return result, nil
	case *string:
		if *v == "" {
			return map[string]string{}, nil
		}
		pairs := strings.Split(*v, ",")
		result := make(map[string]string)
		for _, pair := range pairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid map item: %s", pair)
			}
			result[kv[0]] = kv[1]
		}
		return result, nil
	case string:
		if v == "" {
			return map[string]string{}, nil
		}
		pairs := strings.Split(v, ",")
		result := make(map[string]string)
		for _, pair := range pairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid map item: %s", pair)
			}
			result[kv[0]] = kv[1]
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %v to map[string]string", value)
	}
}
