package figtree

import (
	"fmt"
	"strings"
	"time"
)

// NewString with validator and withered support
func (tree *figTree) NewString(name string, value string, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	name = strings.ToLower(name)
	if _, exists := tree.figs[name]; exists {
		tree.problems = append(tree.problems, fmt.Errorf("name '%s' already exists", name))
		return tree
	}
	tree.activateFlagSet()
	vPtr := &Value{
		Value:      value,
		Mutagensis: tString,
	}
	tree.values.Store(name, vPtr)
	tree.flagSet.Var(vPtr, name, usage)
	def := &figFruit{
		name:        name,
		usage:       usage,
		Mutagenesis: tString,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	if _, exists := tree.withered[name]; !exists {
		tree.withered[name] = witheredFig{
			name:        name,
			Value:       *vPtr,
			Mutagenesis: tString,
		}
	}
	return tree
}

// NewBool with validator and withered support
func (tree *figTree) NewBool(name string, value bool, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	name = strings.ToLower(name)
	if _, exists := tree.figs[name]; exists {
		tree.problems = append(tree.problems, fmt.Errorf("name '%s' already exists", name))
		return tree
	}
	tree.activateFlagSet()
	v := &Value{
		Value:      value,
		Mutagensis: tBool,
	}
	tree.values.Store(name, v)
	tree.flagSet.Var(v, name, usage)
	def := &figFruit{
		name:        name,
		usage:       usage,
		Mutagenesis: tBool,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	if _, exists := tree.withered[name]; !exists {
		tree.withered[name] = witheredFig{
			name:        name,
			Value:       *v,
			Mutagenesis: tBool,
		}
	}
	return tree
}

// NewInt with validator and withered support
func (tree *figTree) NewInt(name string, value int, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	name = strings.ToLower(name)
	if _, exists := tree.figs[name]; exists {
		tree.problems = append(tree.problems, fmt.Errorf("name '%s' already exists", name))
		return tree
	}
	tree.activateFlagSet()
	v := &Value{
		Value:      value,
		Mutagensis: tInt,
	}
	tree.values.Store(name, v)
	tree.flagSet.Var(v, name, usage)
	def := &figFruit{
		name:        name,
		usage:       usage,
		Mutagenesis: tInt,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	if _, exists := tree.withered[name]; !exists {
		tree.withered[name] = witheredFig{
			name:        name,
			Value:       *v,
			Mutagenesis: tInt,
		} // Initialize withered with a copy
	}
	return tree
}

// NewInt64 with validator and withered support
func (tree *figTree) NewInt64(name string, value int64, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	name = strings.ToLower(name)
	if _, exists := tree.figs[name]; exists {
		tree.problems = append(tree.problems, fmt.Errorf("name '%s' already exists", name))
		return tree
	}
	tree.activateFlagSet()
	v := &Value{
		Value:      value,
		Mutagensis: tInt64,
	}
	tree.values.Store(name, v)
	tree.flagSet.Var(v, name, usage)
	def := &figFruit{
		name:        name,
		usage:       usage,
		Mutagenesis: tInt64,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	if _, exists := tree.withered[name]; !exists {
		tree.withered[name] = witheredFig{
			name:        name,
			Value:       *v,
			Mutagenesis: tInt64,
		}
	}
	return tree
}

// NewFloat64 with validator and withered support
func (tree *figTree) NewFloat64(name string, value float64, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	name = strings.ToLower(name)
	if _, exists := tree.figs[name]; exists {
		tree.problems = append(tree.problems, fmt.Errorf("name '%s' already exists", name))
		return tree
	}
	tree.activateFlagSet()
	v := &Value{
		Value:      value,
		Mutagensis: tFloat64,
	}
	tree.values.Store(name, v)
	tree.flagSet.Var(v, name, usage)
	def := &figFruit{
		name:        name,
		usage:       usage,
		Mutagenesis: tFloat64,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	if _, exists := tree.withered[name]; !exists {
		tree.withered[name] = witheredFig{
			name:        name,
			Value:       *v,
			Mutagenesis: tFloat64,
		}
	}
	return tree
}

// NewDuration with validator and withered support
func (tree *figTree) NewDuration(name string, value time.Duration, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	name = strings.ToLower(name)
	if _, exists := tree.figs[name]; exists {
		tree.problems = append(tree.problems, fmt.Errorf("name '%s' already exists", name))
		return tree
	}
	tree.activateFlagSet()
	v := &Value{
		Value:      value,
		Mutagensis: tDuration,
	}
	tree.values.Store(name, v)
	tree.flagSet.Var(v, name, usage)
	def := &figFruit{
		name:        name,
		usage:       usage,
		Mutagenesis: tDuration,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	if _, exists := tree.withered[name]; !exists {
		tree.withered[name] = witheredFig{
			name:        name,
			Value:       *v,
			Mutagenesis: tDuration,
		}
	}
	return tree
}

// NewUnitDuration registers a new time.Duration with a unit time.Duration against a name
func (tree *figTree) NewUnitDuration(name string, value, units time.Duration, usage string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	name = strings.ToLower(name)
	if _, exists := tree.figs[name]; exists {
		tree.problems = append(tree.problems, fmt.Errorf("name '%s' already exists", name))
		return tree
	}
	tree.activateFlagSet()
	v := &Value{
		Value:      value * units,
		Mutagensis: tUnitDuration,
	}
	tree.values.Store(name, v)
	tree.flagSet.Var(v, name, usage)
	def := &figFruit{
		name:        name,
		usage:       usage,
		Mutagenesis: tUnitDuration,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	if _, exists := tree.withered[name]; !exists {
		tree.withered[name] = witheredFig{
			name:        name,
			Value:       *v,
			Mutagenesis: tUnitDuration,
		}
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
	name = strings.ToLower(name)
	if _, exists := tree.figs[name]; exists {
		tree.problems = append(tree.problems, fmt.Errorf("name '%s' already exists", name))
		return tree
	}
	tree.activateFlagSet()
	v := &Value{
		Value:      ListFlag{values: value},
		Mutagensis: tList,
	}
	tree.values.Store(name, v)
	tree.flagSet.Var(v, name, usage)
	def := &figFruit{
		name:        name,
		usage:       usage,
		Mutagenesis: tList,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	if _, exists := tree.withered[name]; !exists {
		witheredVal := make([]string, len(value))
		copy(witheredVal, value)
		tree.withered[name] = witheredFig{
			name: name,
			Value: Value{
				Value:      witheredVal,
				Mutagensis: tList,
			},
			Mutagenesis: tList,
		}
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
	name = strings.ToLower(name)
	if _, exists := tree.figs[name]; exists {
		tree.problems = append(tree.problems, fmt.Errorf("name '%s' already exists", name))
		return tree
	}
	tree.activateFlagSet()
	v := &Value{
		Value:      MapFlag{values: value},
		Mutagensis: tMap,
	}
	tree.values.Store(name, v)
	tree.flagSet.Var(v, name, usage)
	def := &figFruit{
		name:        name,
		usage:       usage,
		Mutagenesis: tMap,
		Mutations:   make([]Mutation, 0),
		Validators:  make([]FigValidatorFunc, 0),
		Callbacks:   make([]Callback, 0),
		Rules:       make([]RuleKind, 0),
	}
	tree.figs[name] = def
	if _, exists := tree.withered[name]; !exists {
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
	}
	return tree
}
