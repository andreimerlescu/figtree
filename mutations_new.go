package figtree

import (
	"flag"
	"time"
)

// MutagensisOfFig returns the Mutagensis of the name
func (fig *Tree) MutagensisOfFig(name string) Mutagenesis {
	fruit, ok := fig.figs[name]
	if !ok {
		return ""
	}
	return fruit.Mutagenesis
}

// MutagensisOf accepts anything and allows you to determine the Mutagensis of the type of from what
func (fig *Tree) MutagensisOf(what interface{}) Mutagenesis {
	switch what.(type) {
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

// NewString with validator and withered support
func (fig *Tree) NewString(name string, value string, usage string) *string {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fig.activateFlagSet()
	ptr := flag.String(name, value, usage)
	def := &Fig{Flesh: ptr, Mutagenesis: tString}
	fig.figs[name] = def
	fig.withered[name] = Fig{Flesh: new(string), Mutagenesis: tString}
	*fig.withered[name].Flesh.(*string) = value
	return ptr
}

// NewBool with validator and withered support
func (fig *Tree) NewBool(name string, value bool, usage string) *bool {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fig.activateFlagSet()
	ptr := flag.Bool(name, value, usage)
	def := &Fig{Flesh: ptr, Mutagenesis: tBool}
	fig.figs[name] = def
	fig.withered[name] = Fig{Flesh: new(bool), Mutagenesis: tBool}
	*fig.withered[name].Flesh.(*bool) = value
	return ptr
}

// NewInt with validator and withered support
func (fig *Tree) NewInt(name string, value int, usage string) *int {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	if flesh, exists := fig.figs[name]; exists {
		return flesh.Flesh.(*int)
	}
	fig.activateFlagSet()
	ptr := flag.Int(name, value, usage)
	def := &Fig{Flesh: ptr, Mutagenesis: tInt}
	fig.figs[name] = def
	fig.withered[name] = Fig{Flesh: new(int), Mutagenesis: tInt} // Initialize withered with a copy
	*fig.withered[name].Flesh.(*int) = value
	return ptr
}

// NewInt64 with validator and withered support
func (fig *Tree) NewInt64(name string, value int64, usage string) *int64 {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fig.activateFlagSet()
	ptr := flag.Int64(name, value, usage)
	def := &Fig{Flesh: ptr, Mutagenesis: tInt64}
	fig.figs[name] = def
	fig.withered[name] = Fig{Flesh: new(int64), Mutagenesis: tInt64}
	*fig.withered[name].Flesh.(*int64) = value
	return ptr
}

// NewFloat64 with validator and withered support
func (fig *Tree) NewFloat64(name string, value float64, usage string) *float64 {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fig.activateFlagSet()
	ptr := flag.Float64(name, value, usage)
	def := &Fig{Flesh: ptr, Mutagenesis: tFloat64}
	fig.figs[name] = def
	fig.withered[name] = Fig{Flesh: new(float64), Mutagenesis: tFloat64}
	*fig.withered[name].Flesh.(*float64) = value
	return ptr
}

// NewDuration with validator and withered support
func (fig *Tree) NewDuration(name string, value time.Duration, usage string) *time.Duration {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fig.activateFlagSet()
	ptr := flag.Duration(name, value, usage)
	def := &Fig{Flesh: ptr, Mutagenesis: tDuration}
	fig.figs[name] = def
	fig.withered[name] = Fig{Flesh: new(time.Duration), Mutagenesis: tDuration}
	*fig.withered[name].Flesh.(*time.Duration) = value
	return ptr
}

// NewUnitDuration registers a new time.Duration with a unit time.Duration against a name
func (fig *Tree) NewUnitDuration(name string, value, units time.Duration, usage string) *time.Duration {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fig.activateFlagSet()
	ptr := flag.Duration(name, value*units, usage)
	def := &Fig{Flesh: ptr, Mutagenesis: tUnitDuration}
	fig.figs[name] = def
	fig.withered[name] = Fig{Flesh: new(time.Duration), Mutagenesis: tUnitDuration}
	*fig.withered[name].Flesh.(*time.Duration) = value * units
	return ptr
}

// NewList with validator and withered support
func (fig *Tree) NewList(name string, value []string, usage string) *[]string {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	l := &ListFlag{values: &value}
	fig.activateFlagSet()
	flag.Var(l, name, usage)
	def := &Fig{Flesh: l, Mutagenesis: tList}
	fig.figs[name] = def
	witheredVal := make([]string, len(value))
	copy(witheredVal, value)
	fig.withered[name] = Fig{Flesh: &ListFlag{values: &witheredVal}, Mutagenesis: tList}
	return l.values
}

// NewMap with validator and withered support
func (fig *Tree) NewMap(name string, value map[string]string, usage string) *map[string]string {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	m := &MapFlag{values: &value}
	fig.activateFlagSet()
	flag.Var(m, name, usage)
	def := &Fig{Flesh: m, Mutagenesis: tMap}
	fig.figs[name] = def
	witheredVal := make(map[string]string)
	for k, v := range value {
		witheredVal[k] = v
	}
	fig.withered[name] = Fig{Flesh: &MapFlag{values: &witheredVal}, Mutagenesis: tMap}
	return m.values
}
