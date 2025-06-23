package figtree

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Value struct {
	Value      interface{}
	Mutagensis Mutagenesis
	Err        error
}

func (v *Value) Raw() interface{} {
	return v.Value
}

func (v *Value) Set(in string) error {
	switch v.Mutagensis {
	case tString:
		v.Value = in
	case tBool:
		if len(in) == 0 {
			in = "false"
		}
		val, err := toBool(in)
		if err != nil {
			v.Err = fmt.Errorf("failed to set bool value from string %q: %w", in, err)
			return v.Err
		}
		v.Value = val
	case tInt:
		if len(in) == 0 {
			in = "0"
		}
		val, err := toInt(in)
		if err != nil {
			v.Err = fmt.Errorf("failed to set int value from string %q: %w", in, err)
			return v.Err
		}
		v.Value = val
	case tInt64:
		if len(in) == 0 {
			in = "0"
		}
		val, err := toInt64(in)
		if err != nil {
			v.Err = fmt.Errorf("failed to set int64 value from string %q: %w", in, err)
			return v.Err
		}
		v.Value = val
	case tFloat64:
		if len(in) == 0 {
			in = "0.0"
		}
		val, err := toFloat64(in)
		if err != nil {
			v.Err = fmt.Errorf("failed to set float64 value from string %q: %w", in, err)
			return v.Err
		}
		v.Value = val
	case tDuration, tUnitDuration:
		if len(in) == 0 {
			err := v.Assign(zeroDuration)
			if err != nil {
				v.Err = fmt.Errorf("failed to set duration value from string %q: %w", in, err)
				return v.Err
			}
			return nil
		}
		va, er := time.ParseDuration(in)
		if er == nil {
			v.Value = va
			return nil
		}

		_val, err := ParseCustomDuration(in)
		if err != nil {
			if !strings.Contains(err.Error(), "invalid duration format") {
				v.Err = fmt.Errorf("failed to set duration value from string %q: %w", in, err)
				val, err := toInt64(in)
				if err != nil {
					v.Err = errors.Join(v.Err, fmt.Errorf("failed to set duration value from string %q", in))
					return v.Err
				}
				v.Value = time.Duration(val)
				return nil
			}
			_val, er := time.ParseDuration(in)
			if er != nil {
				v.Err = fmt.Errorf("failed to set duration value from string %q: %w", in, err)
				return v.Err
			}
			v.Value = _val
			return nil
		}
		v.Value = _val
		return nil

	case tList:
		if len(in) == 0 {
			err := v.Assign(zeroList)
			if err != nil {
				v.Err = fmt.Errorf("failed to set list value from string %q: %w", in, err)
				return v.Err
			}
			return nil
		}
		val, err := toStringSlice(in)
		if err != nil {
			v.Err = fmt.Errorf("failed to set list value from string %q: %w", in, err)
			return v.Err
		}
		if PolicyListAppend {
			vl, er := toStringSlice(v.Value)
			if er != nil {
				v.Err = fmt.Errorf("failed to set list value from string %q: %w", in, er)
				return v.Err
			}
			for _, x := range val {
				vl = append(vl, x)
			}
			vl = DeduplicateStrings(vl)
			v.Value = vl
		} else {
			v.Value = val
		}

	case tMap:
		if len(in) == 0 {
			err := v.Assign(zeroMap)
			if err != nil {
				v.Err = fmt.Errorf("failed to set map value from string %q: %w", in, err)
				return v.Err
			}
			return nil
		}
		val, err := toStringMap(in)
		if err != nil {
			v.Err = fmt.Errorf("failed to set map value from string %q: %w", in, err)
			return v.Err
		}
		if PolicyMapAppend {
			vm := v.Flesh().ToMap()
			for k, v := range val {
				vm[k] = v
			}
			v.Value = vm
		} else {
			v.Value = val
		}
	default:
		err := v.Assign(in)
		if err != nil {
			v.Err = fmt.Errorf("failed to set value from string %q: %w", in, err)
			return v.Err
		}
	}
	return nil
}

func (v *Value) Assign(as interface{}) error {
	switch as := as.(type) {
	case *ListFlag:
		v.Value = as.values
	case *MapFlag:
		v.Value = as.values
	case *Value:
		v.Value = as.Value
	default:
		v.Value = as
	}
	return nil
}

func (v *Value) Flesh() Flesh {
	return &figFlesh{v.Value, nil}
}

func (v *Value) String() string {
	vv := v.Value
	f := figFlesh{vv, nil}
	return f.ToString()
}
