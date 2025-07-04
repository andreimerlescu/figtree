package figtree

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// INTERFACE TYPE CONVERSIONS

// toInt returns an interface{} as an int or returns an error
func toInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case *Value:
		return toInt(v.Value)
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
		return 0, ErrConversion{MutagenesisOf(value), tInt, value}
	}
}

// toInt64 returns an interface{} as an int64 or returns an error
func toInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case *Value:
		return toInt64(v.Value)
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
		return 0, ErrConversion{MutagenesisOf(value), tInt64, value}
	}
}

// toFloat64 returns an interface{} as an float64 or returns an error
func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case *Value:
		return toFloat64(v.Value)
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
		return 0, ErrConversion{MutagenesisOf(value), tFloat64, value}
	}
}

// toString returns an interface{} as a string or returns an error
func toString(value interface{}) (string, error) {
	switch v := value.(type) {
	case MapFlag:
		return toString(v.values)
	case *MapFlag:
		return toString(v.values)
	case ListFlag:
		return toString(v.values)
	case *ListFlag:
		return toString(v.values)
	case *Value:
		return v.String(), nil
	case Value:
		return v.String(), nil
	case *figFlesh:
		return toString(v.AsIs())
	case figFlesh:
		return toString(v.AsIs())
	case *string:
		return *v, nil
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case *int:
		return strconv.Itoa(*v), nil
	case time.Duration:
		return v.String(), nil
	case *time.Duration:
		return v.String(), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case *int64:
		return strconv.FormatInt(*v, 10), nil
	case *float64:
		return strconv.FormatFloat(*v, 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case *bool:
		return strconv.FormatBool(*v), nil
	case bool:
		return strconv.FormatBool(v), nil
	case []string:
		return strings.Join(v, ListSeparator), nil
	case *[]string:
		return strings.Join(*v, ListSeparator), nil
	case *map[string]string:
		parts := make([]string, 0, len(*v))
		for k, x := range *v {
			parts = append(parts, fmt.Sprintf("%s%s%s", k, MapKeySeparator, x))
		}
		return strings.Join(parts, MapSeparator), nil
	case map[string]string:
		parts := make([]string, 0, len(v))
		for k, x := range v {
			parts = append(parts, fmt.Sprintf("%s%s%s", k, MapKeySeparator, x))
		}
		return strings.Join(parts, MapSeparator), nil
	default:
		return "", ErrConversion{MutagenesisOf(value), tString, value}
	}
}

// toBool returns an interface{} as a bool or returns an error
func toBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case *Value:
		return toBool(v.Value)
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
		return false, ErrConversion{MutagenesisOf(value), tBool, value}
	}
}

// toStringSlice returns an interface{} as a []string{} or returns an error
func toStringSlice(value interface{}) ([]string, error) {
	switch v := value.(type) {
	case *Value:
		return toStringSlice(v.Value)
	case *figFlesh:
		return toStringSlice(v.AsIs())
	case *ListFlag:
		return toStringSlice(v.values)
	case ListFlag:
		return toStringSlice(v.values)
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
		if strings.Contains(*v, MapKeySeparator) {
			return nil, ErrConversion{MutagenesisOf(value), tList, value}
		}
		return strings.Split(*v, ListSeparator), nil
	case string:
		if v == "" {
			return []string{}, nil
		}
		if strings.Contains(v, MapSeparator) && strings.Contains(v, MapKeySeparator) {
			return nil, ErrConversion{MutagenesisOf(value), tList, value}
		}
		return strings.Split(v, ListSeparator), nil
	default:
		return nil, ErrConversion{MutagenesisOf(value), tList, value}
	}
}

// toStringMap returns an interface{} as a map[string]string or returns an error
func toStringMap(value interface{}) (map[string]string, error) {
	switch v := value.(type) {
	case *Value:
		return toStringMap(v.Value)
	case *figFlesh:
		return toStringMap(v.AsIs())
	case *MapFlag:
		return toStringMap(v.values)
	case MapFlag:
		return toStringMap(v.values)
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
		pairs := strings.Split(*v, MapSeparator)
		result := make(map[string]string)
		for _, pair := range pairs {
			kv := strings.SplitN(pair, MapKeySeparator, 2)
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
		pairs := strings.Split(v, MapSeparator)
		result := make(map[string]string)
		for _, pair := range pairs {
			kv := strings.SplitN(pair, MapKeySeparator, 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid map item: %s", pair)
			}
			result[kv[0]] = kv[1]
		}
		return result, nil
	default:
		return nil, ErrConversion{MutagenesisOf(value), tMap, value}
	}
}
