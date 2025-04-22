package figtree

import (
	"flag"
	"time"
)

// MutagenesisOfFig returns the Mutagensis of the name
func (tree *figTree) MutagenesisOfFig(name string) Mutagenesis {
	fruit, ok := tree.figs[name]
	if !ok {
		return ""
	}
	return fruit.Mutagenesis
}

// MutagenesisOf accepts anything and allows you to determine the Mutagensis of the type of from what
// Example:
//
//	tree.MutagenesisOf("hello") // Returns tString
//	tree.MutagenesisOf(42)      // Returns tInt
func (tree *figTree) MutagenesisOf(what interface{}) Mutagenesis {
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
func (tree *figTree) NewString(name string, value string, usage string) *string {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.String(name, value, usage)
	def := &figFruit{
		name:        name,
		Flesh:       figFlesh{ptr},
		Mutagenesis: tString,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Flesh:       figFlesh{},
		Mutagenesis: tString,
	}
	witheredFig, exists := tree.withered[name]
	if !exists {
		witheredFig.Flesh = figFlesh{value}
	}
	tree.withered[name] = witheredFig
	return ptr
}

// NewBool with validator and withered support
func (tree *figTree) NewBool(name string, value bool, usage string) *bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.Bool(name, value, usage)
	def := &figFruit{
		name:        name,
		Flesh:       figFlesh{ptr},
		Mutagenesis: tBool,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Flesh:       figFlesh{new(bool)},
		Mutagenesis: tBool,
	}
	*tree.withered[name].Flesh.Flesh.(*bool) = value
	return ptr
}

// NewInt with validator and withered support
func (tree *figTree) NewInt(name string, value int, usage string) *int {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if flesh, exists := tree.figs[name]; exists {
		if flesh.Flesh.Is(flesh.Mutagenesis) {
			return flesh.Flesh.Flesh.(*int)
		}
	}
	tree.activateFlagSet()
	ptr := flag.Int(name, value, usage)
	def := &figFruit{
		name:        name,
		Flesh:       figFlesh{ptr},
		Mutagenesis: tInt,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Flesh:       figFlesh{new(int)},
		Mutagenesis: tInt,
	} // Initialize withered with a copy
	*tree.withered[name].Flesh.Flesh.(*int) = value
	return ptr
}

// NewInt64 with validator and withered support
func (tree *figTree) NewInt64(name string, value int64, usage string) *int64 {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.Int64(name, value, usage)
	def := &figFruit{
		name:        name,
		Flesh:       figFlesh{ptr},
		Mutagenesis: tInt64,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Flesh:       figFlesh{new(int64)},
		Mutagenesis: tInt64,
	}
	*tree.withered[name].Flesh.Flesh.(*int64) = value
	return ptr
}

// NewFloat64 with validator and withered support
func (tree *figTree) NewFloat64(name string, value float64, usage string) *float64 {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.Float64(name, value, usage)
	def := &figFruit{
		name:        name,
		Flesh:       figFlesh{ptr},
		Mutagenesis: tFloat64,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Flesh:       figFlesh{new(float64)},
		Mutagenesis: tFloat64,
	}
	*tree.withered[name].Flesh.Flesh.(*float64) = value
	return ptr
}

// NewDuration with validator and withered support
func (tree *figTree) NewDuration(name string, value time.Duration, usage string) *time.Duration {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.Duration(name, value, usage)
	def := &figFruit{
		name:        name,
		Flesh:       figFlesh{ptr},
		Mutagenesis: tDuration,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Flesh:       figFlesh{new(time.Duration)},
		Mutagenesis: tDuration,
	}
	*tree.withered[name].Flesh.Flesh.(*time.Duration) = value
	return ptr
}

// NewUnitDuration registers a new time.Duration with a unit time.Duration against a name
func (tree *figTree) NewUnitDuration(name string, value, units time.Duration, usage string) *time.Duration {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.Duration(name, value*units, usage)
	def := &figFruit{
		name:        name,
		Flesh:       figFlesh{ptr},
		Mutagenesis: tUnitDuration,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Flesh:       figFlesh{new(time.Duration)},
		Mutagenesis: tUnitDuration,
	}
	*tree.withered[name].Flesh.Flesh.(*time.Duration) = value * units
	return ptr
}

// NewList with validator and withered support
func (tree *figTree) NewList(name string, value []string, usage string) *[]string {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.HasRule(RuleNoLists) {
		return nil
	}
	ptr := &ListFlag{values: &value}
	tree.activateFlagSet()
	flag.Var(ptr, name, usage)
	def := &figFruit{
		name:        name,
		Flesh:       figFlesh{ptr},
		Mutagenesis: tList,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	witheredVal := make([]string, len(value))
	copy(witheredVal, value)
	tree.withered[name] = witheredFig{
		name:        name,
		Flesh:       figFlesh{&ListFlag{values: &witheredVal}},
		Mutagenesis: tList,
	}
	return ptr.values
}

// NewMap with validator and withered support
func (tree *figTree) NewMap(name string, value map[string]string, usage string) *map[string]string {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.HasRule(RuleNoMaps) {
		return nil
	}
	ptr := &MapFlag{values: &value}
	tree.activateFlagSet()
	flag.Var(ptr, name, usage)
	def := &figFruit{
		name:        name,
		Flesh:       figFlesh{ptr},
		Mutagenesis: tMap,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	witheredVal := make(map[string]string)
	for k, v := range value {
		witheredVal[k] = v
	}
	tree.withered[name] = witheredFig{
		name: name,
		Flesh: figFlesh{&MapFlag{
			values: &witheredVal,
		}},
		Mutagenesis: tMap,
	}
	return ptr.values
}
