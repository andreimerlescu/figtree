package figtree

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// AssureStringHasSuffix ensures a string ends with the given suffix.
var AssureStringHasSuffix = func(suffix string) FigValidatorFunc {
	return makeStringValidator(
		func(s string) bool { return strings.HasSuffix(s, suffix) },
		"string must have suffix %q, got %q",
	)
}

// AssureStringHasPrefix ensures a string starts with the given prefix.
var AssureStringHasPrefix = func(prefix string) FigValidatorFunc {
	return makeStringValidator(
		func(s string) bool { return strings.HasPrefix(s, prefix) },
		"string must have prefix %q, got %q",
	)
}

// AssureStringNoSuffixes ensures a string ends with a suffix
// Returns a figValidatorFunc that checks for the substring (case-sensitive).
var AssureStringNoSuffixes = func(suffixes []string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			for _, suffix := range suffixes {
				if strings.HasSuffix(v.ToString(), suffix) {
					return fmt.Errorf("string must not have suffix substring %q, got %q", suffix, v)
				}
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringNoPrefixes ensures a string begins with a prefix
// Returns a figValidatorFunc that checks for the substring (case-sensitive).
var AssureStringNoPrefixes = func(prefixes []string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(v.ToString(), prefix) {
					return fmt.Errorf("string must not have prefix %q, got %q", prefix, v)
				}
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringHasSuffixes ensures a string ends with a suffix
// Returns a figValidatorFunc that checks for the substring (case-sensitive).
var AssureStringHasSuffixes = func(suffixes []string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			for _, suffix := range suffixes {
				if !strings.HasSuffix(v.ToString(), suffix) {
					return fmt.Errorf("string must have suffix %q, got %q", suffix, v)
				}
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringHasPrefixes ensures a string begins with a prefix
// Returns a figValidatorFunc that checks for the substring (case-sensitive).
var AssureStringHasPrefixes = func(prefixes []string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			for _, prefix := range prefixes {
				if !strings.HasPrefix(v.ToString(), prefix) {
					return fmt.Errorf("string must have prefix %q, got %q", prefix, v)
				}
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringNoSuffix ensures a string ends with a suffix
// Returns a figValidatorFunc that checks for the substring (case-sensitive).
var AssureStringNoSuffix = func(suffix string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			vs := v.ToString()
			if strings.HasSuffix(vs, suffix) {
				return fmt.Errorf("string must not have suffix substring %q, got %q", suffix, vs)
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringNoPrefix ensures a string begins with a prefix
// Returns a figValidatorFunc that checks for the substring (case-sensitive).
var AssureStringNoPrefix = func(prefix string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			vs := v.ToString()
			if strings.HasPrefix(vs, prefix) {
				return fmt.Errorf("string must not have prefix substring %q, got %q", prefix, vs)
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringLengthLessThan ensures the string is less than an int
// Returns a figValidatorFunc that checks for the substring (case-sensitive).
var AssureStringLengthLessThan = func(length int) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			vs := v.ToString()
			if len(vs) > length {
				return fmt.Errorf("string must be less than %d chars, got %d", length, len(vs))
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringLengthGreaterThan ensures the string is greater than an int
// Returns a figValidatorFunc that checks for the substring (case-sensitive).
var AssureStringLengthGreaterThan = func(length int) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			vs := v.ToString()
			if len(vs) < length {
				return fmt.Errorf("string must be greater than %d chars, got %d", length, len(vs))
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringSubstring ensures a string contains a specific substring.
// Returns a figValidatorFunc that checks for the substring (case-sensitive).
var AssureStringSubstring = func(sub string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			vs := v.ToString()
			if !strings.Contains(vs, sub) {
				return fmt.Errorf("string must contain substring %q, got %q", sub, vs)
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringLength ensures a string contains a specific substring.
// Returns a figValidatorFunc that checks for the substring (case-sensitive).
var AssureStringLength = func(length int) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			vs := v.ToString()
			if len(vs) < length {
				return fmt.Errorf("string must be at least %d chars, got %q", length, len(vs))
			}
			return nil

		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringNotLength ensures a string contains a specific substring.
// Returns a figValidatorFunc that checks for the substring (case-sensitive).
var AssureStringNotLength = func(length int) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			vs := v.ToString()
			if len(vs) == length {
				return fmt.Errorf("string must not be %d chars, got %q", length, len(vs))
			}
			return nil
		}
		return fmt.Errorf("invalid type, expected string, got %T", value)
	}
}

// AssureStringNotEmpty ensures a string is not empty.
// Returns an error if the value is an empty string or not a string.
var AssureStringNotEmpty = func(value interface{}) error {
	v := figFlesh{value}
	if v.IsString() {
		if len(v.ToString()) == 0 {
			return fmt.Errorf("empty string")
		}
		return nil
	}
	return fmt.Errorf("invalid type, got %T", value)
}

// AssureStringContains ensures a string contains a specific substring.
// Returns an error if the substring is not found or if the value is not a string.
var AssureStringContains = func(substring string) FigValidatorFunc {
	return makeStringValidator(
		func(s string) bool { return strings.Contains(s, substring) },
		"string must contain %q, got %q",
	)
}

// AssureStringNotContains ensures a string contains a specific substring.
// Returns an error if the substring is not found or if the value is not a string.
var AssureStringNotContains = func(substring string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if v.IsString() {
			vs := v.ToString()
			if strings.Contains(vs, substring) {
				return fmt.Errorf("string must not contain %q, got %q", substring, vs)
			}
			return nil
		}
		return fmt.Errorf("invalid type, got %T", value)
	}
}

// AssureBoolTrue ensures a boolean value is true.
// Returns an error if the value is false or not a bool.
var AssureBoolTrue = func(value interface{}) error {
	v := figFlesh{value}
	if v.IsBool() {
		if !v.ToBool() {
			return fmt.Errorf("value must be true, got false")
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected bool, got %T", value)
}

// AssureBoolFalse ensures a boolean value is false.
// Returns an error if the value is true or not a bool.
var AssureBoolFalse = func(value interface{}) error {
	v := figFlesh{value}
	if v.IsBool() {
		if v.ToBool() {
			return fmt.Errorf("value must be false, got true")
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected bool, got %T", value)
}

// AssureIntPositive ensures an integer is positive.
// Returns an error if the value is zero or negative, or not an int.
var AssureIntPositive = func(value interface{}) error {
	v := figFlesh{value}
	if v.IsInt() {
		if v.ToInt() < 0 {
			return fmt.Errorf("value must be positive, got %d", v.ToInt())
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected int, got %T", value)
}

// AssureIntNegative ensures an integer is negative.
// Returns an error if the value is zero or positive, or not an int.
var AssureIntNegative = func(value interface{}) error {
	v := figFlesh{value}
	if v.IsInt() {
		if v.ToInt() > 0 {
			return fmt.Errorf("value must be negative, got %d", v.ToInt())
		}
		return nil
	}
	return fmt.Errorf("invalid type, expected int, got %T", value)
}

// AssureIntGreaterThan ensures an integer is greater than (but not including) an int.
// Returns an error if the value is below, or not an int.
var AssureIntGreaterThan = func(above int) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsInt() {
			return fmt.Errorf("invalid type, expected int, got %T", value)
		}
		i := v.ToInt()
		if i < above {
			return fmt.Errorf("value must be below %d", i)
		}
		return nil
	}
}

// AssureIntLessThan ensures an integer is less than (but not including) an int.
// Returns an error if the value is above, or not an int.
var AssureIntLessThan = func(below int) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsInt() {
			return fmt.Errorf("invalid type, expected int, got %T", value)
		}
		i := v.ToInt()
		if i > below {
			return fmt.Errorf("value must be below %d", i)
		}
		return nil
	}
}

// AssureIntInRange ensures an integer is within a specified range (inclusive).
// Returns an error if the value is outside the range or not an int.
var AssureIntInRange = func(min, max int) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsInt() {
			return fmt.Errorf("invalid type, expected int, got %T", value)
		}
		i := v.ToInt()
		if i < min || i > max {
			return fmt.Errorf("value must be between %d and %d, got %d", min, max, i)
		}
		return nil
	}
}

// AssureInt64GreaterThan ensures an integer is greater than (but not including) an int.
// Returns an error if the value is below, or not an int.
var AssureInt64GreaterThan = func(above int64) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsInt64() {
			return fmt.Errorf("invalid type, expected int64, got %T", value)
		}
		i := v.ToInt64()
		if i < above {
			return fmt.Errorf("value must be below %d", i)
		}
		return nil
	}
}

// AssureInt64LessThan ensures an integer is less than (but not including) an int.
// Returns an error if the value is above, or not an int.
var AssureInt64LessThan = func(below int64) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsInt64() {
			return fmt.Errorf("value must be int64, got %d", value)
		}
		i := v.ToInt64()
		if i > below {
			return fmt.Errorf("value must be below %d", i)
		}
		return nil
	}
}

// AssureInt64Positive ensures an int64 is positive.
// Returns an error if the value is zero or negative, or not an int64.
var AssureInt64Positive = func(value interface{}) error {
	v := figFlesh{value}
	if !v.IsInt64() {
		return fmt.Errorf("invalid type, expected int64, got %T", value)
	}
	i := v.ToInt64()
	if i <= 0 {
		return fmt.Errorf("value must be positive, got %d", i)
	}
	return nil
}

// AssureInt64InRange ensures an int64 is within a specified range (inclusive).
// Returns a figValidatorFunc that checks the value against min and max.
var AssureInt64InRange = func(min, max int64) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsInt64() {
			return fmt.Errorf("invalid type, expected int64, got %T", value)
		}
		i := v.ToInt64()
		if i < min || i > max {
			return fmt.Errorf("value must be between %d and %d, got %d", min, max, i)
		}
		return nil
	}
}

// AssureFloat64Positive ensures a float64 is positive.
// Returns an error if the value is zero or negative, or not a float64.
var AssureFloat64Positive = func(value interface{}) error {
	v := figFlesh{value}
	if !v.IsFloat64() {
		return fmt.Errorf("invalid type, expected float64, got %T", value)
	}
	f := v.ToFloat64()
	if f <= 0 {
		return fmt.Errorf("value must be positive, got %f", f)
	}
	return nil
}

// AssureFloat64InRange ensures a float64 is within a specified range (inclusive).
// Returns an error if the value is outside the range or not a float64.
var AssureFloat64InRange = func(min, max float64) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsFloat64() {
			return fmt.Errorf("invalid type, expected float64, got %T", value)
		}
		f := v.ToFloat64()
		if f < min || f > max {
			return fmt.Errorf("value must be between %f and %f, got %f", min, max, f)
		}
		return nil
	}
}

// AssureFloat64GreaterThan ensures an integer is greater than (but not including) an int.
// Returns an error if the value is below, or not an int.
var AssureFloat64GreaterThan = func(above float64) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsFloat64() {
			return fmt.Errorf("invalid type, expected float64, got %T", value)
		}
		f := v.ToFloat64()
		if f < above {
			return fmt.Errorf("value must be below %f", f)
		}
		return nil
	}
}

// AssureFloat64LessThan ensures an integer is less than (but not including) an int.
// Returns an error if the value is above, or not an int.
var AssureFloat64LessThan = func(below float64) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsFloat64() {
			return fmt.Errorf("invalid type, expected float64, got %T", value)
		}
		f := v.ToFloat64()
		if f > below {
			return fmt.Errorf("value must be below %f", f)
		}
		return nil
	}
}

// AssureFloat64NotNaN ensures a float64 is not NaN.
// Returns an error if the value is NaN or not a float64.
var AssureFloat64NotNaN = func(value interface{}) error {
	v := figFlesh{value}
	if !v.IsFloat64() {
		return fmt.Errorf("invalid type, expected float64, got %T", value)
	}
	n := v.ToFloat64()
	if math.IsNaN(n) {
		return fmt.Errorf("value must not be NaN, got %f", n)
	}
	return nil
}

// AssureDurationGreaterThan ensures a time.Duration is greater than (but not including) a time.Duration.
// Returns an error if the value is below, or not an int.
var AssureDurationGreaterThan = func(above time.Duration) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsDuration() {
			return fmt.Errorf("value must be a duration, got %v", value)
		}
		t := v.ToDuration()
		if t < above {
			return fmt.Errorf("value must be above %v, got = %v", above, t)
		}
		return nil
	}
}

// AssureDurationLessThan ensures a time.Duration is less than (but not including) a time.Duration.
// Returns an error if the value is below, or not an int.
var AssureDurationLessThan = func(below time.Duration) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsDuration() {
			return fmt.Errorf("value must be a duration, got %v", value)
		}
		t := v.ToDuration()
		if t > below {
			return fmt.Errorf("value must be below %v, got = %v", below, t)
		}
		return nil
	}
}

// AssureDurationPositive ensures a time.Duration is positive.
// Returns an error if the value is zero or negative, or not a time.Duration.
var AssureDurationPositive = func(value interface{}) error {
	v := figFlesh{value}
	if !v.IsDuration() {
		return fmt.Errorf("invalid type, expected time.Duration, got %T", value)
	}
	d := v.ToDuration()
	if d <= 0 {
		return fmt.Errorf("duration must be positive, got %s", d)
	}
	return nil
}

// AssureDurationMin ensures a time.Duration is at least a minimum value.
// Returns a figValidatorFunc that checks the duration against the minimum.
var AssureDurationMin = func(min time.Duration) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsDuration() {
			return fmt.Errorf("value must be a duration, got %s", v)
		}
		d := v.ToDuration()
		if d < min {
			return fmt.Errorf("duration must be at least %s, got %s", min, d)
		}
		return nil
	}
}

// AssureDurationMax ensures a time.Duration does not exceed a maximum value.
// Returns an error if the value exceeds the max or is not a time.Duration.
var AssureDurationMax = func(max time.Duration) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsDuration() {
			return fmt.Errorf("invalid type, expected time.Duration, got %T", value)
		}
		d := v.ToDuration()
		if d > max {
			return fmt.Errorf("duration must not exceed %s, got %s", max, d)
		}
		return nil
	}
}

// AssureListNotEmpty ensures a list is not empty.
// Returns an error if the list has no elements or is not a ListFlag.
var AssureListNotEmpty = func(value interface{}) error {
	v := figFlesh{value}
	if !v.IsList() {
		return fmt.Errorf("invalid type, expected *ListFlag or []string, got %T", v)
	}
	l := v.ToList()
	if len(l) == 0 {
		return fmt.Errorf("list is empty")
	}
	return nil
}

// AssureListMinLength ensures a list has at least a minimum number of elements.
// Returns an error if the list is too short or not a ListFlag.
var AssureListMinLength = func(min int) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsList() {
			return fmt.Errorf("invalid type, expected *ListFlag or []string, got %T", v)
		}
		l := v.ToList()
		if len(l) < min {
			return fmt.Errorf("list is empty")
		}
		return nil
	}
}

// AssureListContains ensures a list contains a specific string value.
// Returns a figValidatorFunc that checks for the presence of the value.
var AssureListContains = func(inside string) FigValidatorFunc {
	return func(value interface{}) error {
		v := NewFlesh(value)
		if !v.IsList() {
			return fmt.Errorf("invalid type, expected ListFlag or []string, got %T", value)
		}
		l := v.ToList()
		for _, item := range l {
			if item == inside {
				return nil
			}
		}
		return fmt.Errorf("list must contain %q, got %v", inside, l)
	}
}

// AssureListNotContains ensures a list contains a specific string value.
// Returns a figValidatorFunc that checks for the presence of the value.
var AssureListNotContains = func(not string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsList() {
			return fmt.Errorf("invalid type, expected *ListFlag, []string, or *[]string, got %T", v)
		}
		l := v.ToList()
		for _, item := range l {
			if item == not {
				return fmt.Errorf("list cannot contain %s", item)
			}
		}
		return nil
	}
}

// AssureListContainsKey ensures a list contains a specific string.
// It accepts *ListFlag, *[]string, or []string and returns an error if the key string is not found
// or the type is invalid.
var AssureListContainsKey = func(key string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsList() {
			return fmt.Errorf("invalid type, expected *ListFlag or []string, got %T", value)
		}
		l := v.ToList()
		for _, item := range l {
			if item == key {
				return nil
			}
		}
		return fmt.Errorf("list must contain %q, got %v", key, l)
	}
}

// AssureListLength ensures a list has exactly the specified length.
// It accepts *ListFlag, *[]string, or []string and returns an error if the length differs
// or the type is invalid.
var AssureListLength = func(length int) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsList() {
			return fmt.Errorf("invalid type, expected *ListFlag or []string, got %T", value)
		}
		l := v.ToList()
		if len(l) != length {
			return fmt.Errorf("list must have length %d, got %d", length, len(l))
		}
		return nil
	}
}

// AssureMapNotEmpty ensures a map is not empty.
// Returns an error if the map has no entries or is not a MapFlag.
var AssureMapNotEmpty = func(value interface{}) error {
	v := figFlesh{value}
	if !v.IsMap() {
		return fmt.Errorf("invalid type, expected *ListFlag or []string, got %T", v)
	}
	m := v.ToMap()
	if len(m) == 0 {
		return fmt.Errorf("list is empty")
	}
	return nil
}

// AssureMapHasKey ensures a map contains a specific key.
// Returns an error if the key is missing or the value is not a MapFlag.
var AssureMapHasKey = func(key string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsMap() {
			return fmt.Errorf("invalid type, got %T", value)
		}
		m := v.ToMap()
		if _, exists := m[key]; !exists {
			return fmt.Errorf("map must contain key %q", key)
		}
		return nil
	}
}

// AssureMapHasNoKey ensures a map contains a specific key.
// Returns an error if the key is missing or the value is not a MapFlag.
var AssureMapHasNoKey = func(key string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsMap() {
			return fmt.Errorf("invalid type, got %T", value)
		}
		m := v.ToMap()
		if _, exists := m[key]; exists {
			return fmt.Errorf("map must not contain key %q", key)
		}
		return nil
	}
}

// AssureMapValueMatches ensures a map has a specific key with a matching value.
// Returns a figValidatorFunc that checks for the key-value pair.
var AssureMapValueMatches = func(key, match string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsMap() {
			return fmt.Errorf("%s is not a map", key)
		}
		m := v.ToMap()
		if val, exists := m[key]; exists {
			if val != match {
				return fmt.Errorf("map value %q must have value %q, got %q", key, match, val)
			}
			return nil
		}
		return fmt.Errorf("map key %q does not exist", key)
	}
}

// AssureMapHasKeys ensures a map contains all specified keys.
// Returns an error if any key is missing or the value is not a *MapFlag.
var AssureMapHasKeys = func(keys []string) FigValidatorFunc {
	return func(value interface{}) error {
		var missing []string
		v := figFlesh{value}
		if !v.IsMap() {
			return fmt.Errorf("invalid type, expected map[string]string, got %T", v)
		}
		m := v.ToMap()
		for _, key := range keys {
			if _, exists := m[key]; !exists {
				missing = append(missing, key)
			}
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
var AssureMapLength = func(length int) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsMap() {
			return fmt.Errorf("invalid type, expected *MapFlag or map[string]string, got %T", value)
		}
		m := v.ToMap()
		if len(m) != length {
			return fmt.Errorf("map must have length %d, got %d", length, len(m))
		}
		return nil
	}
}

// AssureMapNotLength ensures a map has exactly the specified length.
// It accepts *MapFlag, *map[string]string, or map[string]string and returns an error
// if the length differs or the type is invalid.
var AssureMapNotLength = func(length int) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsMap() {
			return fmt.Errorf("invalid type, got %T", value)
		}
		m := v.ToMap()
		if len(m) == length {
			return fmt.Errorf("map must not have length %d, got %d", length, len(m))
		}
		return nil
	}
}
