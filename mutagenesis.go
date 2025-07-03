package figtree

import (
	"flag"
	"time"
)

func (m Mutagenesis) Kind() string {
	switch m {
	case tString:
		return "string|*string"
	case tBool:
		return "bool|*bool"
	case tInt:
		return "int|*int"
	case tInt64:
		return "int64|*int64"
	case tFloat64:
		return "float64|*float64"
	case tDuration, tUnitDuration:
		return "time.Duration|*time.Duration"
	case tList:
		return "ListFlag|*ListFlag|[]string|*[]string"
	case tMap:
		return "MapFlag|*MapFlag|map[string]string|*map[string]string"
	default:
		return string(m)
	}
}

// MutagenesisOfFig returns the Mutagensis of the name
func (tree *figTree) MutagenesisOfFig(name string) Mutagenesis {
	fruit, ok := tree.figs[name]
	if !ok {
		return ""
	}
	return fruit.Mutagenesis
}

func MutagenesisOf(what interface{}) Mutagenesis {
	switch x := what.(type) {
	case Value:
		return x.Mutagensis
	case flag.Value:
		fv, e := toFloat64(x.String())
		if e == nil {
			return MutagenesisOf(fv)
		}
		i64v, e := toInt64(x.String())
		if e == nil {
			return MutagenesisOf(i64v)
		}
		iv, e := toInt(x.String())
		if e == nil {
			return MutagenesisOf(iv)
		}
		bv, e := toBool(x.String())
		if e == nil {
			return MutagenesisOf(bv)
		}
		sv, e := toStringSlice(x.String())
		if e == nil {
			return MutagenesisOf(sv)
		}
		mv, e := toStringMap(x.String())
		if e == nil {
			return MutagenesisOf(mv)
		}
		return ""

	case int:
		return tInt
	case *int:
		return tInt
	case *int64:
		return tInt64
	case int64:
		return tInt64
	case string:
		return tString
	case *string:
		return tString
	case bool:
		return tBool
	case *bool:
		return tBool
	case *float64:
		return tFloat64
	case float64:
		return tFloat64
	case time.Duration:
		return tDuration
	case *time.Duration:
		return tDuration
	case []string:
		return tList
	case *[]string:
		return tList
	case map[string]string:
		return tMap
	case *map[string]string:
		return tMap
	default:
		return ""
	}
}

// MutagenesisOf accepts anything and allows you to determine the Mutagensis of the type of from what
// Example:
//
//	tree.MutagenesisOf("hello") // Returns tString
//	tree.MutagenesisOf(42)      // Returns tInt
func (tree *figTree) MutagenesisOf(what interface{}) Mutagenesis {
	return MutagenesisOf(what)
}
