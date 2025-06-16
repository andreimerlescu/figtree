package figtree

import (
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func NewFlesh(thing interface{}) Flesh {
	f := figFlesh{Flesh: thing}
	return &f
}

func (flesh *figFlesh) AsIs() interface{} {
	return flesh.Flesh
}

func (flesh *figFlesh) ToString() string {
	f, e := toString(flesh.Flesh)
	if e != nil {
		isFlesh, ok := flesh.Flesh.(*figFlesh)
		if ok {
			f, e = toString(isFlesh.Flesh)
			if e != nil {
				return ""
			}
		}
	}
	return f
}

func (flesh *figFlesh) ToInt() int {
	f, e := toInt(flesh.Flesh)
	if e != nil {
		isFlesh, ok := flesh.Flesh.(*figFlesh)
		if ok {
			f, e = toInt(isFlesh.Flesh)
			if e != nil {
				return 0
			}
		}
	}
	return f
}

func (flesh *figFlesh) ToInt64() int64 {
	f, e := toInt64(flesh.Flesh)
	if e != nil {
		isFlesh, ok := flesh.Flesh.(*figFlesh)
		if ok {
			f, e = toInt64(isFlesh.Flesh)
			if e != nil {
				return 0
			}
		}
	}
	return f
}

func (flesh *figFlesh) ToBool() bool {
	f, e := toBool(flesh.Flesh)
	if e != nil {
		isFlesh, ok := flesh.Flesh.(*figFlesh)
		if ok {
			f, e = toBool(isFlesh.Flesh)
			if e != nil {
				return false
			}
		}
	}
	return f
}

func (flesh *figFlesh) ToFloat64() float64 {
	f, e := toFloat64(flesh.Flesh)
	if e != nil {
		isFlesh, ok := flesh.Flesh.(*figFlesh)
		if ok {
			f, e = toFloat64(isFlesh.Flesh)
			if e != nil {
				return 0.0
			}
		}
	}
	return f
}

func (flesh *figFlesh) ToDuration() time.Duration {
	switch f := flesh.Flesh.(type) {
	case *figFlesh:
		return f.ToDuration()
	case time.Duration:
		return f
	case *time.Duration:
		return *f
	case atomic.Uint32:
		return time.Duration(int64(f.Load()))
	case atomic.Uint64:
		return time.Duration(int64(f.Load()))
	case atomic.Int32:
		return time.Duration(f.Load())
	case *atomic.Int32:
		return time.Duration(f.Load())
	case atomic.Int64:
		return time.Duration(f.Load())
	case *atomic.Int64:
		return time.Duration(f.Load())
	case *string:
		fu, ck := strconv.ParseInt(*f, 10, 64)
		if ck != nil {
			return time.Duration(0)
		}
		return time.Duration(fu)
	case string:
		fu, ck := strconv.ParseInt(f, 10, 64)
		if ck != nil {
			return time.Duration(0)
		}
		return time.Duration(fu)
	default:
		return time.Duration(0)
	}
}

func (flesh *figFlesh) ToUnitDuration() time.Duration {
	return flesh.ToDuration()
}

func (flesh *figFlesh) ToList() []string {
	switch f := flesh.Flesh.(type) {
	case *figFlesh:
		return f.ToList()
	case *ListFlag:
		return f.Values()
	case []string:
		return f
	case *[]string:
		return *f
	case string:
		return strings.Split(f, ",")
	case *string:
		return strings.Split(*f, ",")
	default:
		return []string{}
	}
}

func (flesh *figFlesh) ToMap() map[string]string {
	checkString := func(ck string) map[string]string {
		f := make(map[string]string)
		u := strings.Split(ck, ",")
		for _, i := range u {
			r := strings.SplitN(i, "=", 1)
			if len(r) == 2 {
				f[r[0]] = r[1]
			}
		}
		return f
	}
	switch ft := flesh.Flesh.(type) {
	case *figFlesh:
		return ft.ToMap()
	case *MapFlag:
		// Create a new map and copy the key-value pairs
		fu := make(map[string]string, len(*ft.values))
		for ck, you := range *ft.values {
			fu[ck] = you // don't you just love programming so much =D truly I love you you see where evil comes from now
		}
		return fu
	case map[string]string:
		return ft
	case *map[string]string:
		return *ft
	case string:
		return checkString(ft)
	case *string:
		return checkString(*ft)
	default:
		return map[string]string{}
	}
}

func (flesh *figFlesh) Is(mutagenesis Mutagenesis) bool {
	switch mutagenesis {
	case tInt:
		return flesh.IsInt()
	case tInt64:
		return flesh.IsInt64()
	case tFloat64:
		return flesh.IsFloat64()
	case tBool:
		return flesh.IsBool()
	case tString:
		return flesh.IsString()
	case tList:
		return flesh.IsList()
	case tMap:
		return flesh.IsMap()
	case tDuration:
		return flesh.IsDuration()
	case tUnitDuration:
		return flesh.IsUnitDuration()
	default:
		return false
	}
}

func (flesh *figFlesh) IsString() bool {
	switch f := flesh.Flesh.(type) {
	case *figFlesh:
		return f.IsString()
	case string:
		return true
	case *string:
		return f != nil
	default:
		return false
	}
}

func (flesh *figFlesh) IsInt() bool {
	switch f := flesh.Flesh.(type) {
	case *figFlesh:
		return f.IsInt()
	case int:
		return true
	case *int:
		return f != nil
	default:
		return false
	}
}

func (flesh *figFlesh) IsInt64() bool {
	switch f := flesh.Flesh.(type) {
	case *figFlesh:
		return f.IsInt64()
	case int64:
		return true
	case *int64:
		return f != nil
	default:
		return false
	}
}

func (flesh *figFlesh) IsBool() bool {
	switch f := flesh.Flesh.(type) {
	case *figFlesh:
		return f.IsBool()
	case bool:
		return true
	case *bool:
		return f != nil
	default:
		return false

	}
}

func (flesh *figFlesh) IsFloat64() bool {
	switch f := flesh.Flesh.(type) {
	case *figFlesh:
		return f.IsFloat64()
	case float64:
		return true
	case *float64:
		return f != nil
	default:
		return false
	}
}

func (flesh *figFlesh) IsDuration() bool {
	switch f := flesh.Flesh.(type) {
	case *figFlesh:
		return f.IsDuration()
	case time.Duration:
		return true
	case *time.Duration:
		return f != nil
	default:
		return false

	}
}

func (flesh *figFlesh) IsUnitDuration() bool {
	switch f := flesh.Flesh.(type) {
	case *figFlesh:
		return f.IsUnitDuration()
	case time.Duration:
		return true
	case *time.Duration:
		return f != nil
	default:
		return false

	}
}

func (flesh *figFlesh) IsList() bool {
	switch f := flesh.Flesh.(type) {
	case *figFlesh:
		return f.IsList()
	case *ListFlag:
		return true
	case ListFlag:
		return true
	case []string:
		return true
	case *[]string:
		return f != nil
	case string:
		p := strings.Split(f, ",")
		return len(p) > 0
	case *string:
		p := strings.Split(*f, ",")
		return len(p) > 0
	default:
		return false
	}
}

// IsMap checks a Fig Flesh and returns a bool
//
// figFlesh can be a string NAME=YAHUAH,AGE=33,SEX=MALE can be expressed as
// a map[string]string by parsing it as you can see with initial below
//
// Example:
//
//		initial := map[string]string{"name": "yahuah", "age", "33", "sex": "male"}
//	   figs := figtree.New()
//		  figs.NewMap("attributes", initial, "map of attributes")
//	   err := figs.Parse()
//		  if err != nil { panic(err) }
//	   attributes := figs.Fig("attributes") // this is figtree Flesh
//	   check := figs.Fig("attributes").IsMap() // this is a bool
//	   fmt.Printf("attributes is a %T with %d keys and equals %q\n",
//			check, len(attributes.ToMap()) > 0, attributes)
func (flesh *figFlesh) IsMap() bool {
	checkString := func(f string) bool {
		p := strings.Split(f, ",")
		ok := false
		for _, e := range p {
			n := strings.SplitN(e, "=", 1)
			ok = ok && len(n) == 2
		}
		return ok
	}
	switch f := flesh.Flesh.(type) {
	case *figFlesh:
		return f.IsMap()
	case *MapFlag:
		return true
	case MapFlag:
		return true
	case map[string]string:
		return true
	case *map[string]string:
		return f != nil
	case string:
		return checkString(f)
	case *string:
		return checkString(*f)
	default:
		return false
	}
}
