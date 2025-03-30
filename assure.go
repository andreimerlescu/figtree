package figtree

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// AssureStringHasSuffix ensures a string ends with a suffix
// Returns a ValidatorFunc that checks for the substring (case-sensitive).
var AssureStringHasSuffix = func(suffix string) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(string); ok {
			if !strings.HasSuffix(v, suffix) {
				return fmt.Errorf("string must have suffix substring %q, got %q", suffix, v)
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringHasPrefix ensures a string begins with a prefix
// Returns a ValidatorFunc that checks for the substring (case-sensitive).
var AssureStringHasPrefix = func(prefix string) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(string); ok {
			if !strings.HasPrefix(v, prefix) {
				return fmt.Errorf("string must have prefix substring %q, got %q", prefix, v)
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringLengthLessThan ensures the string is less than an int
// Returns a ValidatorFunc that checks for the substring (case-sensitive).
var AssureStringLengthLessThan = func(length int) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(string); ok {
			if len(v) > length {
				return fmt.Errorf("string must be less than %d chars, got %d", length, len(v))
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringLengthGreaterThan ensures the string is greater than an int
// Returns a ValidatorFunc that checks for the substring (case-sensitive).
var AssureStringLengthGreaterThan = func(length int) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(string); ok {
			if len(v) < length {
				return fmt.Errorf("string must be greater than %d chars, got %d", length, len(v))
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringSubstring ensures a string contains a specific substring.
// Returns a ValidatorFunc that checks for the substring (case-sensitive).
var AssureStringSubstring = func(sub string) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(string); ok {
			if !strings.Contains(v, sub) {
				return fmt.Errorf("string must contain substring %q, got %q", sub, v)
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringLength ensures a string contains a specific substring.
// Returns a ValidatorFunc that checks for the substring (case-sensitive).
var AssureStringLength = func(length int) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(string); ok {
			if len(v) < length {
				return fmt.Errorf("string must be at least %d chars, got %q", length, len(v))
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringNotLength ensures a string contains a specific substring.
// Returns a ValidatorFunc that checks for the substring (case-sensitive).
var AssureStringNotLength = func(length int) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(string); ok {
			if len(v) == length {
				return fmt.Errorf("string must not be %d chars, got %q", length, len(v))
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringNotEmpty ensures a string is not empty.
// Returns an error if the value is an empty string or not a string.
var AssureStringNotEmpty = func(value interface{}) error {
	if v, ok := value.(string); ok {
		if len(v) == 0 {
			return fmt.Errorf("empty string")
		}
		return nil
	}
	return fmt.Errorf("invalid type, got %T", value)
}

// AssureStringContains ensures a string contains a specific substring.
// Returns an error if the substring is not found or if the value is not a string.
var AssureStringContains = func(substring string) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(string); ok {
			if !strings.Contains(v, substring) {
				return fmt.Errorf("string must contain %q, got %q", substring, v)
			}
			return nil
		}
		return fmt.Errorf("invalid type, got %T", value)
	}
}

// AssureStringNotContains ensures a string contains a specific substring.
// Returns an error if the substring is not found or if the value is not a string.
var AssureStringNotContains = func(substring string) ValidatorFunc {
	return func(value interface{}) error {
		v, ok := value.(string)
		if ok {
			if strings.Contains(v, substring) {
				return fmt.Errorf("string must not contain %q, got %q", substring, v)
			}
			return nil
		}
		v2, ok2 := value.(*string)
		if ok2 {
			if strings.Contains(*v2, substring) {
				return fmt.Errorf("string must not contain %q, got %q", substring, *v2)
			}
			return nil
		}
		return fmt.Errorf("invalid type, got %T", value)
	}
}

// AssureBoolTrue ensures a boolean value is true.
// Returns an error if the value is false or not a bool.
var AssureBoolTrue = func(value interface{}) error {
	if v, ok := value.(bool); ok {
		if !v {
			return fmt.Errorf("value must be true, got false")
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected bool, got %T", value)
}

// AssureBoolFalse ensures a boolean value is false.
// Returns an error if the value is true or not a bool.
var AssureBoolFalse = func(value interface{}) error {
	if v, ok := value.(bool); ok {
		if v {
			return fmt.Errorf("value must be false, got true")
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected bool, got %T", value)
}

// AssurePositiveInt ensures an integer is positive.
// Returns an error if the value is zero or negative, or not an int.
var AssurePositiveInt = func(value interface{}) error {
	if v, ok := value.(int); ok {
		if v <= 0 {
			return fmt.Errorf("value must be positive, got %d", v)
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected int, got %T", value)
}

// AssureNegativeInt ensures an integer is negative.
// Returns an error if the value is zero or positive, or not an int.
var AssureNegativeInt = func(value interface{}) error {
	if v, ok := value.(int); ok {
		if v >= 0 {
			return fmt.Errorf("value must be negative, got %d", v)
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected int, got %T", value)
}

// AssureIntGreaterThan ensures an integer is greater than (but not including) an int.
// Returns an error if the value is below, or not an int.
var AssureIntGreaterThan = func(above int) ValidatorFunc {
	return func(value interface{}) error {
		if v, err := toInt(value); err == nil {
			if v < above {
				return fmt.Errorf("value must be below %d", v)
			}
			return nil
		} else {
			return err
		}
	}
}

// AssureIntLessThan ensures an integer is less than (but not including) an int.
// Returns an error if the value is above, or not an int.
var AssureIntLessThan = func(below int) ValidatorFunc {
	return func(value interface{}) error {
		if v, err := toInt(value); err == nil {
			if v > below {
				return fmt.Errorf("value must be below %d", v)
			}
			return nil
		} else {
			return err
		}
	}
}

// AssureIntInRange ensures an integer is within a specified range (inclusive).
// Returns an error if the value is outside the range or not an int.
var AssureIntInRange = func(min, max int) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(int); ok {
			if v < min || v > max {
				return fmt.Errorf("value must be between %d and %d, got %d", min, max, v)
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected int, got %T", value)
	}
}

// AssureInt64GreaterThan ensures an integer is greater than (but not including) an int.
// Returns an error if the value is below, or not an int.
var AssureInt64GreaterThan = func(above int64) ValidatorFunc {
	return func(value interface{}) error {
		if v, err := toInt64(value); err == nil {
			if v < above {
				return fmt.Errorf("value must be below %d", v)
			}
			return nil
		} else {
			return err
		}
	}
}

// AssureInt64LessThan ensures an integer is less than (but not including) an int.
// Returns an error if the value is above, or not an int.
var AssureInt64LessThan = func(below int64) ValidatorFunc {
	return func(value interface{}) error {
		if v, err := toInt64(value); err == nil {
			if v > below {
				return fmt.Errorf("value must be below %d", v)
			}
			return nil
		} else {
			return err
		}
	}
}

// AssureInt64Positive ensures an int64 is positive.
// Returns an error if the value is zero or negative, or not an int64.
var AssureInt64Positive = func(value interface{}) error {
	if v, ok := value.(int64); ok {
		if v <= 0 {
			return fmt.Errorf("value must be positive, got %d", v)
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected int64, got %T", value)
}

// AssureInt64InRange ensures an int64 is within a specified range (inclusive).
// Returns a ValidatorFunc that checks the value against min and max.
var AssureInt64InRange = func(min, max int64) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(int64); ok {
			if v < min || v > max {
				return fmt.Errorf("value must be between %d and %d, got %d", min, max, v)
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected int64, got %T", value)
	}
}

// AssureFloat64Positive ensures a float64 is positive.
// Returns an error if the value is zero or negative, or not a float64.
var AssureFloat64Positive = func(value interface{}) error {
	if v, ok := value.(float64); ok {
		if v <= 0 {
			return fmt.Errorf("value must be positive, got %f", v)
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected float64, got %T", value)
}

// AssureFloat64InRange ensures a float64 is within a specified range (inclusive).
// Returns an error if the value is outside the range or not a float64.
var AssureFloat64InRange = func(min, max float64) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(float64); ok {
			if v < min || v > max {
				return fmt.Errorf("value must be between %f and %f, got %f", min, max, v)
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected float64, got %T", value)
	}
}

// AssureFloat64GreaterThan ensures an integer is greater than (but not including) an int.
// Returns an error if the value is below, or not an int.
var AssureFloat64GreaterThan = func(above float64) ValidatorFunc {
	return func(value interface{}) error {
		if v, err := toFloat64(value); err == nil {
			if v < above {
				return fmt.Errorf("value must be below %f", v)
			}
			return nil
		} else {
			return err
		}
	}
}

// AssureFloat64LessThan ensures an integer is less than (but not including) an int.
// Returns an error if the value is above, or not an int.
var AssureFloat64LessThan = func(below float64) ValidatorFunc {
	return func(value interface{}) error {
		if v, err := toFloat64(value); err == nil {
			if v > below {
				return fmt.Errorf("value must be below %f", v)
			}
			return nil
		} else {
			return err
		}
	}
}

// AssureFloat64NotNaN ensures a float64 is not NaN.
// Returns an error if the value is NaN or not a float64.
var AssureFloat64NotNaN = func(value interface{}) error {
	if v, ok := value.(float64); ok {
		if math.IsNaN(v) {
			return fmt.Errorf("value must not be NaN, got %f", v)
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected float64, got %T", value)
}

// AssureDurationGreaterThan ensures a time.Duration is greater than (but not including) a time.Duration.
// Returns an error if the value is below, or not an int.
var AssureDurationGreaterThan = func(above time.Duration) ValidatorFunc {
	return func(value interface{}) error {
		var t time.Duration
		if v, ok := value.(time.Duration); ok {
			t = v
		}
		if v, ok := value.(*time.Duration); ok {
			t = *v
		}
		if t < above {
			return fmt.Errorf("value must be above %v, got = %v", above, t)
		}
		return nil
	}
}

// AssureDurationLessThan ensures a time.Duration is less than (but not including) a time.Duration.
// Returns an error if the value is below, or not an int.
var AssureDurationLessThan = func(below time.Duration) ValidatorFunc {
	return func(value interface{}) error {
		var t time.Duration
		if v, ok := value.(time.Duration); ok {
			t = v
		}
		if v, ok := value.(*time.Duration); ok {
			t = *v
		}
		if t > below {
			return fmt.Errorf("value must be below %v, got = %v", below, t)
		}
		return nil
	}
}

// AssureDurationPositive ensures a time.Duration is positive.
// Returns an error if the value is zero or negative, or not a time.Duration.
var AssureDurationPositive = func(value interface{}) error {
	if v, ok := value.(time.Duration); ok {
		if v <= 0 {
			return fmt.Errorf("duration must be positive, got %s", v)
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected time.Duration, got %T", value)
}

// AssureDurationMin ensures a time.Duration is at least a minimum value.
// Returns a ValidatorFunc that checks the duration against the minimum.
var AssureDurationMin = func(min time.Duration) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(time.Duration); ok {
			if v < min {
				return fmt.Errorf("duration must be at least %s, got %s", min, v)
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected time.Duration, got %T", value)
	}
}

// AssureDurationMax ensures a time.Duration does not exceed a maximum value.
// Returns an error if the value exceeds the max or is not a time.Duration.
var AssureDurationMax = func(max time.Duration) ValidatorFunc {
	return func(value interface{}) error {
		if v, ok := value.(time.Duration); ok {
			if v > max {
				return fmt.Errorf("duration must not exceed %s, got %s", max, v)
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected time.Duration, got %T", value)
	}
}

// AssureListNotEmpty ensures a list is not empty.
// Returns an error if the list has no elements or is not a ListFlag.
var AssureListNotEmpty = func(value interface{}) error {
	switch v := value.(type) {
	case *ListFlag:
		if v == nil || len(*v.values) == 0 {
			return fmt.Errorf("list is empty")
		}
		return nil
	case *[]string:
		if len(*v) == 0 {
			return fmt.Errorf("list is empty")
		}
		return nil
	case []string:
		if len(v) == 0 {
			return fmt.Errorf("list is empty")
		}
		return nil
	default:
		return fmt.Errorf("invalid type, expected *ListFlag or []string, got %T", v)
	}
}

// AssureListMinLength ensures a list has at least a minimum number of elements.
// Returns an error if the list is too short or not a ListFlag.
var AssureListMinLength = func(min int) ValidatorFunc {
	return func(value interface{}) error {
		switch v := value.(type) {
		case *ListFlag:
			if len(*v.values) < min {
				return fmt.Errorf("list must have at least %d elements, got %d", min, len(*v.values))
			}
			return nil
		case *[]string:
			if len(*v) < min {
				return fmt.Errorf("list is empty")
			}
			return nil
		case []string:
			if len(v) < min {
				return fmt.Errorf("list is empty")
			}
			return nil
		default:
			return fmt.Errorf("invalid type, expected *ListFlag or []string, got %T", v)
		}
	}
}

// AssureListContains ensures a list contains a specific string value.
// Returns a ValidatorFunc that checks for the presence of the value.
var AssureListContains = func(value string) ValidatorFunc {
	return func(v interface{}) error {
		if list, ok := v.(*ListFlag); ok && list != nil {
			for _, item := range *list.values {
				if item == value {
					return nil
				}
			}
		}
		if list, ok := v.(*[]string); ok && list != nil {
			for _, item := range *list {
				if item == value {
					return nil
				}
			}
		}
		if list, ok := v.([]string); ok && list != nil {
			for _, item := range list {
				if item == value {
					return nil
				}
			}
			return fmt.Errorf("list must contain %q, got %v", value, list)
		}
		return fmt.Errorf("invalid type, expected *ListFlag, []string, or *[]string, got %T", v)
	}
}

// AssureListNotContains ensures a list contains a specific string value.
// Returns a ValidatorFunc that checks for the presence of the value.
var AssureListNotContains = func(value string) ValidatorFunc {
	return func(v interface{}) error {
		if list, ok := v.(*ListFlag); ok && list != nil {
			for _, item := range *list.values {
				if item == value {
					return fmt.Errorf("list cannot contain %s", item)
				}
			}
			return nil
		}
		if list, ok := v.(*[]string); ok && list != nil {
			for _, item := range *list {
				if item == value {
					return fmt.Errorf("list cannot contain %s", item)
				}
			}
			return nil
		}
		if list, ok := v.([]string); ok && list != nil {
			for _, item := range list {
				if item == value {
					return fmt.Errorf("list cannot contain %s", item)
				}
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected *ListFlag, []string, or *[]string, got %T", v)
	}
}

// AssureListContainsKey ensures a list contains a specific string.
// It accepts *ListFlag, *[]string, or []string and returns an error if the key string is not found
// or the type is invalid.
var AssureListContainsKey = func(key string) ValidatorFunc {
	return func(value interface{}) error {
		switch list := value.(type) {
		case *ListFlag:
			if list == nil {
				return fmt.Errorf("list is nil")
			}
			for _, item := range *list.values {
				if item == key {
					return nil
				}
			}
			return fmt.Errorf("list must contain %q, got %v", key, *list.values)
		case *[]string:
			if list == nil {
				return fmt.Errorf("list is nil")
			}
			for _, item := range *list {
				if item == key {
					return nil
				}
			}
			return fmt.Errorf("list must contain %q, got %v", key, *list)
		case []string:
			for _, item := range list {
				if item == key {
					return nil
				}
			}
			return fmt.Errorf("list must contain %q, got %v", key, list)
		default:
			return fmt.Errorf("invalid type, expected *ListFlag or []string, got %T", value)
		}
	}
}

// AssureListLength ensures a list has exactly the specified length.
// It accepts *ListFlag, *[]string, or []string and returns an error if the length differs
// or the type is invalid.
var AssureListLength = func(length int) ValidatorFunc {
	return func(value interface{}) error {
		switch list := value.(type) {
		case *ListFlag:
			if list == nil {
				return fmt.Errorf("list is nil")
			}
			if len(*list.values) != length {
				return fmt.Errorf("list must have length %d, got %d", length, len(*list.values))
			}
			return nil
		case *[]string:
			if list == nil {
				return fmt.Errorf("list is nil")
			}
			if len(*list) != length {
				return fmt.Errorf("list must have length %d, got %d", length, len(*list))
			}
			return nil
		case []string:
			if len(list) != length {
				return fmt.Errorf("list must have length %d, got %d", length, len(list))
			}
			return nil
		default:
			return fmt.Errorf("invalid type, expected *ListFlag or []string, got %T", value)
		}
	}
}

// AssureMapNotEmpty ensures a map is not empty.
// Returns an error if the map has no entries or is not a MapFlag.
var AssureMapNotEmpty = func(value interface{}) error {
	switch v := value.(type) {
	case *MapFlag:
		if len(*v.values) == 0 {
			return fmt.Errorf("map is empty")
		}
		return nil
	case *map[string]string:
		if len(*v) == 0 {
			return fmt.Errorf("list is empty")
		}
		return nil
	case map[string]string:
		if len(v) == 0 {
			return fmt.Errorf("list is empty")
		}
		return nil
	default:
		return fmt.Errorf("invalid type, expected *ListFlag or []string, got %T", v)
	}
}

// AssureMapHasKey ensures a map contains a specific key.
// Returns an error if the key is missing or the value is not a MapFlag.
var AssureMapHasKey = func(key string) ValidatorFunc {
	return func(value interface{}) error {
		switch v := value.(type) {
		case *MapFlag:
			if _, exists := (*v.values)[key]; !exists {
				return fmt.Errorf("map must contain key %q", key)
			}
			return nil
		case *map[string]string:
			if _, exists := (*v)[key]; !exists {
				return fmt.Errorf("map must contain key %q", key)
			}
			return nil
		case map[string]string:
			if _, exists := v[key]; !exists {
				return fmt.Errorf("map must contain key %q", key)
			}
			return nil
		default:
			return fmt.Errorf("invalid type, got %T", value)
		}
	}
}

// AssureMapHasNoKey ensures a map contains a specific key.
// Returns an error if the key is missing or the value is not a MapFlag.
var AssureMapHasNoKey = func(key string) ValidatorFunc {
	return func(value interface{}) error {
		switch v := value.(type) {
		case *MapFlag:
			if _, exists := (*v.values)[key]; exists {
				return fmt.Errorf("map must not contain key %q", key)
			}
			return nil
		case *map[string]string:
			if _, exists := (*v)[key]; exists {
				return fmt.Errorf("map must not contain key %q", key)
			}
			return nil
		case map[string]string:
			if _, exists := v[key]; exists {
				return fmt.Errorf("map must not contain key %q", key)
			}
			return nil
		default:
			return fmt.Errorf("invalid type, got %T", value)
		}
	}
}

// AssureMapValueMatches ensures a map has a specific key with a matching value.
// Returns a ValidatorFunc that checks for the key-value pair.
var AssureMapValueMatches = func(key, value string) ValidatorFunc {
	return func(v interface{}) error {
		switch m := v.(type) {
		case *MapFlag:
			if val, exists := (*m.values)[key]; exists {
				if val != value {
					return fmt.Errorf("map key %q must have value %q, got %q", key, value, val)
				}
				return nil
			}
			return fmt.Errorf("map key %q does not exist", key)
		case *map[string]string:
			if val, exists := (*m)[key]; exists {
				if val != value {
					return fmt.Errorf("map value %q must have value %q, got %q", key, value, val)
				}
				return nil
			}
			return fmt.Errorf("map key %q does not exist", key)
		case map[string]string:
			if val, exists := m[key]; exists {
				if val != value {
					return fmt.Errorf("map value %q must have value %q, got %q", key, value, val)
				}
				return nil
			}
			return fmt.Errorf("map key %q does not exist", key)
		default:
			return fmt.Errorf("invalid type, expected map[string]string, got %T", v)
		}
	}
}

// AssureMapHasKeys ensures a map contains all specified keys.
// Returns an error if any key is missing or the value is not a *MapFlag.
var AssureMapHasKeys = func(keys []string) ValidatorFunc {
	return func(value interface{}) error {
		missing := []string{}
		switch m := value.(type) {
		case *MapFlag:
			if m == nil {
				return fmt.Errorf("map is nil")
			}
			for _, key := range keys {
				if _, exists := (*m.values)[key]; !exists {
					missing = append(missing, key)
				}
			}
		case *map[string]string:
			if m == nil {
				return fmt.Errorf("map is nil")
			}
			for _, key := range keys {
				if _, exists := (*m)[key]; !exists {
					missing = append(missing, key)
				}
			}
		case map[string]string:
			for _, key := range keys {
				if _, exists := m[key]; !exists {
					missing = append(missing, key)
				}
			}
		default:
			return fmt.Errorf("invalid type, expected *MapFlag or map[string]string, got %T", value)
		}
		if len(missing) > 0 {
			return fmt.Errorf("map must contain keys %v, missing %v", keys, missing)
		}
		return nil
	}
}

// AssureMapLength ensures a map has exactly the specified length.
// It accepts *MapFlag, *map[string]string, or map[string]string and returns an error
// if the length differs or the type is invalid.
var AssureMapLength = func(length int) ValidatorFunc {
	return func(value interface{}) error {
		switch m := value.(type) {
		case *MapFlag:
			if m == nil {
				return fmt.Errorf("map is nil")
			}
			if len(*m.values) != length {
				return fmt.Errorf("map must have length %d, got %d", length, len(*m.values))
			}
			return nil
		case *map[string]string:
			if m == nil {
				return fmt.Errorf("map is nil")
			}
			if len(*m) != length {
				return fmt.Errorf("map must have length %d, got %d", length, len(*m))
			}
			return nil
		case map[string]string:
			if len(m) != length {
				return fmt.Errorf("map must have length %d, got %d", length, len(m))
			}
			return nil
		default:
			return fmt.Errorf("invalid type, expected *MapFlag or map[string]string, got %T", value)
		}
	}
}

// AssureMapNotLength ensures a map has exactly the specified length.
// It accepts *MapFlag, *map[string]string, or map[string]string and returns an error
// if the length differs or the type is invalid.
var AssureMapNotLength = func(length int) ValidatorFunc {
	return func(value interface{}) error {
		switch m := value.(type) {
		case *MapFlag:
			if m == nil {
				return fmt.Errorf("map is nil")
			}
			if len(*m.values) == length {
				return fmt.Errorf("map must not have length %d, got %d", length, len(*m.values))
			}
			return nil
		case *map[string]string:
			if m == nil {
				return fmt.Errorf("map is nil")
			}
			if len(*m) == length {
				return fmt.Errorf("map must not have length %d, got %d", length, len(*m))
			}
			return nil
		case map[string]string:
			if len(m) == length {
				return fmt.Errorf("map must not have length %d, got %d", length, len(m))
			}
			return nil
		default:
			return fmt.Errorf("invalid type, expected *MapFlag or map[string]string, got %T", value)
		}
	}
}
