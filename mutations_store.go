package figtree

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func (tree *figTree) Store(mut Mutagenesis, name string, value interface{}) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		if !tree.angel.Load() {
			tree.mu.Unlock()
			tree.Resurrect(name)
			tree.mu.Lock()
			fruit = tree.figs[name]
		}
	}
	if fruit == nil {
		if !tree.angel.Load() {
			tree.mu.Unlock()
			tree.Resurrect(name)
			tree.mu.Lock()
			fruit = tree.figs[name]
		}
	}
	if fruit == nil {
		return tree
	}
	if tree.HasRule(RulePreventChange) || fruit.HasRule(RulePreventChange) {
		return tree
	}
	if tree.angel.Load() {
		tree.figs[name].Error = errors.Join(tree.figs[name].Error, fmt.Errorf("tree fruit is an angel so we cannot store %s inside %s", tree.MutagenesisOf(value), tree.MutagenesisOf(fruit.Value.Value)))
		return tree
	}
	mv := tree.MutagenesisOf(value)
	if mv == tDuration && mut == tUnitDuration {
		mv = tUnitDuration
	}
	if !strings.EqualFold(string(mv), string(fruit.Mutagenesis)) {
		tree.figs[name].Error = errors.Join(tree.figs[name].Error, fmt.Errorf("will not store %s inside %s", tree.MutagenesisOf(value), tree.MutagenesisOf(fruit.Value.Value)))
		return tree
	}
	if _, exists := tree.withered[name]; !exists {
		tree.withered[name] = witheredFig{
			Value:       fruit.Value,
			Mutagenesis: tString,
			Error:       fmt.Errorf("missing withered value for %s", name),
		}
	}
	err := fruit.runCallbacks(CallbackBeforeChange)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
	}
	changed, previous, current := tree.persist(fruit, mut, name, value)
	if !changed {
		return tree
	}
	err = fruit.runCallbacks(CallbackAfterChange)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
	}
	tree.figs[name] = fruit
	if tree.tracking && !tree.angel.Load() {
		tree.mutationsCh <- Mutation{
			Property:    name,
			Mutagenesis: strings.ToLower(string(mut)),
			Way:         "Store" + string(mut),
			Old:         previous,
			New:         current,
			When:        time.Now(),
			Error:       err,
		}
	}
	return tree
}

// StoreString replaces the name with the new value while issuing a Mutation if figTree.tracking is true
func (tree *figTree) StoreString(name, value string) Plant {
	return tree.Store(tString, name, value)
}

// StoreBool replaces the name with the new value while issuing a Mutation if figTree.tracking is true
func (tree *figTree) StoreBool(name string, value bool) Plant {
	return tree.Store(tBool, name, value)
}

// StoreInt replaces the name with the new value while issuing a Mutation if figTree.tracking is true
func (tree *figTree) StoreInt(name string, value int) Plant {
	return tree.Store(tInt, name, value)
}

// StoreInt64 replaces the name with the new value while issuing a Mutation if figTree.tracking is true
func (tree *figTree) StoreInt64(name string, value int64) Plant {
	return tree.Store(tInt64, name, value)
}

// StoreFloat64 replaces the name with the new value while issuing a Mutation if figTree.tracking is true
func (tree *figTree) StoreFloat64(name string, value float64) Plant {
	return tree.Store(tFloat64, name, value)
}

// StoreDuration replaces the name with the new value while issuing a Mutation if figTree.tracking is true
func (tree *figTree) StoreDuration(name string, value time.Duration) Plant {
	return tree.Store(tDuration, name, value)
}

// StoreUnitDuration replaces the name with the new value while issuing a Mutation if figTree.tracking is true
func (tree *figTree) StoreUnitDuration(name string, value, units time.Duration) Plant {
	return tree.Store(tUnitDuration, name, value*units)
}

// StoreList replaces the name with the new value while issuing a Mutation if figTree.tracking is true
func (tree *figTree) StoreList(name string, value []string) Plant {
	return tree.Store(tList, name, value)
}

// StoreMap replaces the name with the new value while issuing a Mutation if figTree.tracking is true
func (tree *figTree) StoreMap(name string, value map[string]string) Plant {
	return tree.Store(tMap, name, value)
}

// persist requires the figTree.mu to be locked before using this func and is an internal func
func (tree *figTree) persist(fruit *figFruit, mut Mutagenesis, name string, value interface{}) (changed bool, previous, current interface{}) {
	if fruit.HasRule(RulePreventChange) {
		return false, fruit, fruit
	}
	if fruit.HasRule(RulePanicOnChange) {
		panic("RuleOnPanicChange triggered for " + fruit.name)
	}
	flesh := fruit.Value.Value
	switch mut {
	case tMap:
		var old *map[string]string
		var err error
		switch f := flesh.(type) {
		case *Value:
			old = f.Value.(*map[string]string)
		case *MapFlag:
			old = f.values
		case *map[string]string:
			old = f
		case map[string]string:
			old = &f
		case string:
			m := map[string]string{}
			for _, p := range strings.Split(f, ",") {
				z := strings.SplitN(p, "=", 1)
				if len(z) == 2 {
					m[z[0]] = z[1]
				}
			}
			old = &m
		case *string:
			m := map[string]string{}
			for _, p := range strings.Split(*f, ",") {
				z := strings.SplitN(p, "=", 1)
				if len(z) == 2 {
					m[z[0]] = z[1]
				}
			}
			old = &m
		default:
			return false, f, value
		}
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, flesh, value
		}
		var current *map[string]string
		switch f := value.(type) {
		case *MapFlag:
			current = f.values
		case *map[string]string:
			current = f
		case map[string]string:
			current = &f
		case string:
			m := map[string]string{}
			for _, p := range strings.Split(f, ",") {
				z := strings.SplitN(p, "=", 1)
				if len(z) == 2 {
					m[z[0]] = z[1]
				}
			}
			current = &m
		case *string:
			m := map[string]string{}
			for _, p := range strings.Split(*f, ",") {
				z := strings.SplitN(p, "=", 1)
				if len(z) == 2 {
					m[z[0]] = z[1]
				}
			}
			current = &m
		default:
			return false, old, f
		}
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, old, value
		}
		fruit.Value.Value = current
		tree.figs[name] = fruit
		return old != current, old, current
	case tList:
		var old *[]string
		var err error
		switch v := flesh.(type) {
		case *Value:
			old = v.Value.(*[]string)
		case *ListFlag:
			old = v.values
		case *[]string:
			old = v
		case []string:
			old = &v
		case string:
			x := strings.Split(v, ",")
			old = &x
		case *string:
			x := strings.Split(*v, ",")
			old = &x
		default:
			return false, v, value
		}
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, flesh, value
		}
		var current *[]string
		switch v := value.(type) {
		case *ListFlag:
			current = v.values
		case []string:
			current = &v
		case *[]string:
			current = v
		case string:
			x := strings.Split(v, ",")
			current = &x
		case *string:
			x := strings.Split(*v, ",")
			current = &x
		default:
			return false, old, flesh
		}
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, old, value
		}
		fruit.Value.Value = current
		tree.figs[name] = fruit
		return old != current, old, current
	case tUnitDuration:
		var old time.Duration
		var err error
		switch v := flesh.(type) {
		case *Value:
			old = v.Value.(time.Duration)
		case time.Duration:
			old = v
		case *time.Duration:
			old = *v
		case string:
			old, err = time.ParseDuration(v)
		case *string:
			old, err = time.ParseDuration(*v)
		default:
			return false, flesh, value
		}
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, flesh, value
		}
		var current time.Duration
		switch v := value.(type) {
		case time.Duration:
			current = v
		case *time.Duration:
			current = *v
		case string:
			current, err = time.ParseDuration(v)
		case *string:
			current, err = time.ParseDuration(*v)
		default:
			return false, flesh, value

		}
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, old, value
		}
		fruit.Value.Value = current
		tree.figs[name] = fruit
		return old != current, old, current
	case tDuration:
		var old time.Duration
		var err error
		switch v := flesh.(type) {
		case time.Duration:
			old = v
		case *time.Duration:
			old = *v
		case string:
			old, err = time.ParseDuration(v)
		case *string:
			old, err = time.ParseDuration(*v)
		default:
			return false, flesh, value

		}
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, flesh, value
		}
		var current time.Duration
		switch v := value.(type) {
		case time.Duration:
			current = v
		case *time.Duration:
			current = *v
		case string:
			current, err = time.ParseDuration(v)
		case *string:
			current, err = time.ParseDuration(*v)
		default:
			return false, flesh, value

		}
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, old, value
		}
		fruit.Value.Value = current
		tree.figs[name] = fruit
		return old != current, old, current
	case tFloat64:
		old, err := toFloat64(flesh)
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, flesh, value
		}
		current, err := toFloat64(value)
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, old, value
		}
		fruit.Value.Value = current
		tree.figs[name] = fruit
		return old != current, old, current
	case tInt64:
		old, err := toInt64(flesh)
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, flesh, value
		}
		current, err := toInt64(value)
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, old, value
		}
		fruit.Value.Value = current
		tree.figs[name] = fruit
		return old != current, old, current
	case tInt:
		old, err := toInt(flesh)
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, flesh, value
		}
		current, err := toInt(value)
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, old, value
		}
		fruit.Value.Value = current
		tree.figs[name] = fruit
		return old != current, old, current
	case tString:
		old, err := toString(flesh)
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, flesh, value
		}
		current, err := toString(value)
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, old, value
		}
		fruit.Value.Value = current
		tree.figs[name] = fruit
		return !strings.EqualFold(old, current), old, current
	case tBool:
		old, err := toBool(flesh)
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, flesh, value
		}
		current, err := toBool(value)
		if err != nil {
			tree.figs[name].Error = errors.Join(tree.figs[name].Error, err)
			return false, old, value
		}
		fruit.Value.Value = current
		tree.figs[name] = fruit
		return old != current, old, current
	default:
		return false, flesh, value
	}
}
