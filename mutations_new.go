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
func (tree *figTree) NewString(name string, value string, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	v := Value{
		Value:      value,
		Mutagensis: tree.MutagenesisOf(value),
	}
	flag.Var(&v, name, usage)
	def := &figFruit{
		name:        name,
		Value:       v,
		Mutagenesis: tString,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	theWitheredFig, exists := tree.withered[name]
	if !exists {
		tree.withered[name] = witheredFig{
			name:        name,
			Value:       v,
			Mutagenesis: tString,
		}
	}
	tree.withered[name] = theWitheredFig
	return tree
}

// NewBool with validator and withered support
func (tree *figTree) NewBool(name string, value bool, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	v := Value{
		Value:      value,
		Mutagensis: tree.MutagenesisOf(value),
	}
	flag.Var(&v, name, usage)
	def := &figFruit{
		name:        name,
		Value:       v,
		Mutagenesis: tBool,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Value:       v,
		Mutagenesis: tBool,
	}
	return tree
}

// NewInt with validator and withered support
func (tree *figTree) NewInt(name string, value int, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	v := Value{
		Value:      value,
		Mutagensis: tree.MutagenesisOf(value),
	}
	flag.Var(&v, name, usage)
	def := &figFruit{
		name:        name,
		Value:       v,
		Mutagenesis: tInt,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Value:       v,
		Mutagenesis: tInt,
	} // Initialize withered with a copy
	return tree
}

// NewInt64 with validator and withered support
func (tree *figTree) NewInt64(name string, value int64, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	v := Value{
		Value:      value,
		Mutagensis: tree.MutagenesisOf(value),
	}
	flag.Var(&v, name, usage)
	def := &figFruit{
		name:        name,
		Value:       v,
		Mutagenesis: tInt64,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Value:       v,
		Mutagenesis: tInt64,
	}
	return tree
}

// NewFloat64 with validator and withered support
func (tree *figTree) NewFloat64(name string, value float64, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	v := Value{
		Value:      value,
		Mutagensis: tree.MutagenesisOf(value),
	}
	flag.Var(&v, name, usage)
	def := &figFruit{
		name:        name,
		Value:       v,
		Mutagenesis: tFloat64,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Value:       v,
		Mutagenesis: tFloat64,
	}
	return tree
}

// NewDuration with validator and withered support
func (tree *figTree) NewDuration(name string, value time.Duration, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	v := Value{
		Value:      value,
		Mutagensis: tree.MutagenesisOf(value),
	}
	flag.Var(&v, name, usage)
	def := &figFruit{
		name:        name,
		Value:       v,
		Mutagenesis: tDuration,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Value:       v,
		Mutagenesis: tDuration,
	}
	return tree
}

// NewUnitDuration registers a new time.Duration with a unit time.Duration against a name
func (tree *figTree) NewUnitDuration(name string, value, units time.Duration, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	v := Value{
		Value:      value * units,
		Mutagensis: tree.MutagenesisOf(value),
	}
	flag.Var(&v, name, usage)
	def := &figFruit{
		name:        name,
		Value:       v,
		Mutagenesis: tUnitDuration,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	tree.withered[name] = witheredFig{
		name:        name,
		Value:       v,
		Mutagenesis: tUnitDuration,
	}
	return tree
}

// NewList with validator and withered support
func (tree *figTree) NewList(name string, value []string, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.HasRule(RuleNoLists) {
		return tree
	}
	ptr := &ListFlag{values: &value}
	tree.activateFlagSet()
	v := Value{
		Value:      ptr,
		Mutagensis: tree.MutagenesisOf(ptr),
	}
	flag.Var(&v, name, usage)
	def := &figFruit{
		name:        name,
		Value:       v,
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
		Value:       v,
		Mutagenesis: tList,
	}
	return tree
}

// NewMap with validator and withered support
func (tree *figTree) NewMap(name string, value map[string]string, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.HasRule(RuleNoMaps) {
		return tree
	}
	ptr := &MapFlag{values: &value}
	tree.activateFlagSet()
	v := Value{
		Value:      ptr,
		Mutagensis: tree.MutagenesisOf(ptr),
	}
	flag.Var(&v, name, usage)
	def := &figFruit{
		name:        name,
		Value:       v,
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
		Value: Value{
			Value:      witheredVal,
			Mutagensis: tree.MutagenesisOf(witheredVal),
		},
		Mutagenesis: tMap,
	}
	return tree
}
