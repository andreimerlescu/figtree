package figtree

import (
	"fmt"
	"time"
)

type Value struct {
	Value      interface{}
	Mutagensis Mutagenesis
}

func (v *Value) Raw() interface{} {
	return v.Value
}

func (v *Value) Set(in string) error {
	switch v.Mutagensis {
	case tString:
		v.Value = in
	case tBool:
		val, err := toBool(in)
		if err != nil {
			return fmt.Errorf("failed to set bool value from string %q: %w", in, err)
		}
		v.Value = val
	case tInt:
		val, err := toInt(in)
		if err != nil {
			return fmt.Errorf("failed to set int value from string %q: %w", in, err)
		}
		v.Value = val
	case tInt64:
		val, err := toInt64(in)
		if err != nil {
			return fmt.Errorf("failed to set int64 value from string %q: %w", in, err)
		}
		v.Value = val
	case tFloat64:
		val, err := toFloat64(in)
		if err != nil {
			return fmt.Errorf("failed to set float64 value from string %q: %w", in, err)
		}
		v.Value = val
	case tDuration, tUnitDuration:
		val, err := toInt64(in)
		if err != nil {
			return fmt.Errorf("failed to set duration value from string %q: %w", in, err)
		}
		v.Value = time.Duration(val)
	case tList:
		val, err := toStringSlice(in)
		if err != nil {
			return fmt.Errorf("failed to set list value from string %q: %w", in, err)
		}
		if PolicyListAppend {
			vl := v.Flesh().ToList()
			for _, x := range val {
				vl = append(vl, x)
			}
			vl = DeduplicateStrings(vl)
			v.Value = vl
		} else {
			v.Value = val
		}

	case tMap:
		val, err := toStringMap(in)
		if err != nil {
			return fmt.Errorf("failed to set map value from string %q: %w", in, err)
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
		v.Value = in
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
	return &figFlesh{v.Value}
}

func (v *Value) String() string {
	vv := v.Value
	f := figFlesh{vv}
	return f.ToString()
}
