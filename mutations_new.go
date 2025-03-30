package figtree

import (
	"flag"
	"time"
)

// MutagensisOfFig returns the Mutagensis of the name
func (tree *Tree) MutagensisOfFig(name string) Mutagenesis {
	fruit, ok := tree.figs[name]
	if !ok {
		return ""
	}
	return fruit.Mutagenesis
}

// MutagensisOf accepts anything and allows you to determine the Mutagensis of the type of from what
func (tree *Tree) MutagensisOf(what interface{}) Mutagenesis {
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
func (tree *Tree) NewString(name string, value string, usage string) *string {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.String(name, value, usage)
	def := &Fig{
		Flesh:         ptr,
		Mutagenesis:   tString,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	tree.figs[name] = def
	tree.withered[name] = Fig{
		Flesh:         new(string),
		Mutagenesis:   tString,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	*tree.withered[name].Flesh.(*string) = value
	return ptr
}

// NewBool with validator and withered support
func (tree *Tree) NewBool(name string, value bool, usage string) *bool {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.Bool(name, value, usage)
	def := &Fig{
		Flesh:         ptr,
		Mutagenesis:   tBool,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	tree.figs[name] = def
	tree.withered[name] = Fig{
		Flesh:         new(bool),
		Mutagenesis:   tBool,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	*tree.withered[name].Flesh.(*bool) = value
	return ptr
}

// NewInt with validator and withered support
func (tree *Tree) NewInt(name string, value int, usage string) *int {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if flesh, exists := tree.figs[name]; exists {
		return flesh.Flesh.(*int)
	}
	tree.activateFlagSet()
	ptr := flag.Int(name, value, usage)
	def := &Fig{
		Flesh:         ptr,
		Mutagenesis:   tInt,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	tree.figs[name] = def
	tree.withered[name] = Fig{
		Flesh:         new(int),
		Mutagenesis:   tInt,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	} // Initialize withered with a copy
	*tree.withered[name].Flesh.(*int) = value
	return ptr
}

// NewInt64 with validator and withered support
func (tree *Tree) NewInt64(name string, value int64, usage string) *int64 {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.Int64(name, value, usage)
	def := &Fig{
		Flesh:         ptr,
		Mutagenesis:   tInt64,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	tree.figs[name] = def
	tree.withered[name] = Fig{
		Flesh:         new(int64),
		Mutagenesis:   tInt64,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	*tree.withered[name].Flesh.(*int64) = value
	return ptr
}

// NewFloat64 with validator and withered support
func (tree *Tree) NewFloat64(name string, value float64, usage string) *float64 {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.Float64(name, value, usage)
	def := &Fig{
		Flesh:         ptr,
		Mutagenesis:   tFloat64,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	tree.figs[name] = def
	tree.withered[name] = Fig{
		Flesh:         new(float64),
		Mutagenesis:   tFloat64,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	*tree.withered[name].Flesh.(*float64) = value
	return ptr
}

// NewDuration with validator and withered support
func (tree *Tree) NewDuration(name string, value time.Duration, usage string) *time.Duration {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.Duration(name, value, usage)
	def := &Fig{
		Flesh:         ptr,
		Mutagenesis:   tDuration,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	tree.figs[name] = def
	tree.withered[name] = Fig{
		Flesh:         new(time.Duration),
		Mutagenesis:   tDuration,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	*tree.withered[name].Flesh.(*time.Duration) = value
	return ptr
}

// NewUnitDuration registers a new time.Duration with a unit time.Duration against a name
func (tree *Tree) NewUnitDuration(name string, value, units time.Duration, usage string) *time.Duration {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	ptr := flag.Duration(name, value*units, usage)
	def := &Fig{
		Flesh:         ptr,
		Mutagenesis:   tUnitDuration,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	tree.figs[name] = def
	tree.withered[name] = Fig{
		Flesh:         new(time.Duration),
		Mutagenesis:   tUnitDuration,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	*tree.withered[name].Flesh.(*time.Duration) = value * units
	return ptr
}

// NewList with validator and withered support
func (tree *Tree) NewList(name string, value []string, usage string) *[]string {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	l := &ListFlag{values: &value}
	tree.activateFlagSet()
	flag.Var(l, name, usage)
	def := &Fig{
		Flesh:         l,
		Mutagenesis:   tList,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	tree.figs[name] = def
	witheredVal := make([]string, len(value))
	copy(witheredVal, value)
	tree.withered[name] = Fig{
		Flesh:         &ListFlag{values: &witheredVal},
		Mutagenesis:   tList,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	return l.values
}

// NewMap with validator and withered support
func (tree *Tree) NewMap(name string, value map[string]string, usage string) *map[string]string {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	m := &MapFlag{values: &value}
	tree.activateFlagSet()
	flag.Var(m, name, usage)
	def := &Fig{
		Flesh:         m,
		Mutagenesis:   tMap,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	tree.figs[name] = def
	witheredVal := make(map[string]string)
	for k, v := range value {
		witheredVal[k] = v
	}
	tree.withered[name] = Fig{
		Flesh: &MapFlag{
			values: &witheredVal,
		},
		Mutagenesis:   tMap,
		Mutations:     make([]Mutation, 0),
		Validators:    make([]ValidatorFunc, 0),
		Callbacks:     make([]Callback, 0),
		CallbackAfter: CallbackAfterVerify,
	}
	return m.values
}
