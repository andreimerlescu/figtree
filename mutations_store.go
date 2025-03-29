package figtree

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func (fig *Tree) Store(mut Mutagenesis, name string, value interface{}) Fruit {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		if !fig.angel.Load() {
			fig.mu.Unlock()
			fig.Resurrect(name)
			fig.mu.Lock()
			fruit = fig.figs[name]
		}
	}
	if fruit == nil {
		if !fig.angel.Load() {
			fig.mu.Unlock()
			fig.Resurrect(name)
			fig.mu.Lock()
			fruit = fig.figs[name]
		}
	}
	if fruit == nil {
		return fig
	}
	if fig.angel.Load() {
		fig.figs[name].Error = errors.Join(fig.figs[name].Error, fmt.Errorf("fig fruit is an angel so we cannot store %s inside %s", fig.MutagensisOf(value), fig.MutagensisOf(fruit.Flesh)))
		return fig
	}
	mv := fig.MutagensisOf(value)
	if mv == tDuration && mut == tUnitDuration {
		mv = tUnitDuration
	}
	if !strings.EqualFold(string(mv), string(fruit.Mutagenesis)) {
		fig.figs[name].Error = errors.Join(fig.figs[name].Error, fmt.Errorf("will not store %s inside %s", fig.MutagensisOf(value), fig.MutagensisOf(fruit.Flesh)))
		return fig
	}
	if _, exists := fig.withered[name]; !exists {
		fig.withered[name] = Fig{Flesh: fruit.Flesh, Mutagenesis: tString, Error: fmt.Errorf("missing withered value for %s", name)}
	}
	changed, previous, current := fig.persist(fruit, mut, name, value)
	if fig.tracking && changed {
		fig.mutationsCh <- Mutation{
			Property: name,
			Kind:     strings.ToLower(string(mut)),
			Way:      "Store" + string(mut),
			Old:      previous,
			New:      current,
			When:     time.Now(),
		}
	}
	return fig
}

// StoreString replaces the name with the new value while issuing a Mutation if Tree.tracking is true
func (fig *Tree) StoreString(name, value string) Fruit {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		if !fig.angel.Load() {
			fig.Resurrect(name)
			fruit = fig.figs[name]
		}
	}
	if fruit == nil {
		if !fig.angel.Load() {
			fig.Resurrect(name)
			fruit = fig.figs[name]
		}
	}
	if fruit == nil {
		return fig
	}
	if fig.angel.Load() {
		fig.figs[name].Error = errors.Join(fig.figs[name].Error, fmt.Errorf("fig fruit is an angel so we cannot store %s inside %s", fig.MutagensisOf(value), fig.MutagensisOf(fruit.Flesh)))
		return fig
	}
	mv := fig.MutagensisOf(value)
	if !strings.EqualFold(string(mv), string(fruit.Mutagenesis)) {
		fig.figs[name].Error = errors.Join(fig.figs[name].Error, fmt.Errorf("will not store %s inside %s", fig.MutagensisOf(value), fig.MutagensisOf(fruit.Flesh)))
		return fig
	}
	if _, exists := fig.withered[name]; !exists {
		fig.withered[name] = Fig{Flesh: fruit.Flesh, Mutagenesis: tString, Error: fmt.Errorf("missing withered value for %s", name)}
	}
	old, err := toString(fruit.Flesh)
	if err != nil {
		fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
		return fig
	}
	current, err := toString(&value)
	if err != nil {
		fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
		return fig
	}
	fruit.Flesh = current
	fig.figs[name] = fruit
	changed := !strings.EqualFold(old, current) // string
	if fig.tracking && changed {
		fig.mutationsCh <- Mutation{
			Property: name,
			Kind:     "string",
			Way:      "StoreString",
			Old:      old,
			New:      current,
			When:     time.Now(),
		}
	}
	return fig
}

// StoreBool replaces the name with the new value while issuing a Mutation if Tree.tracking is true
func (fig *Tree) StoreBool(name string, value bool) Fruit {
	return fig.Store(tBool, name, value)
}

// StoreInt replaces the name with the new value while issuing a Mutation if Tree.tracking is true
func (fig *Tree) StoreInt(name string, value int) Fruit {
	return fig.Store(tInt, name, value)
}

// StoreInt64 replaces the name with the new value while issuing a Mutation if Tree.tracking is true
func (fig *Tree) StoreInt64(name string, value int64) Fruit {
	return fig.Store(tInt64, name, value)
}

// StoreFloat64 replaces the name with the new value while issuing a Mutation if Tree.tracking is true
func (fig *Tree) StoreFloat64(name string, value float64) Fruit {
	return fig.Store(tFloat64, name, value)
}

// StoreDuration replaces the name with the new value while issuing a Mutation if Tree.tracking is true
func (fig *Tree) StoreDuration(name string, value time.Duration) Fruit {
	return fig.Store(tDuration, name, value)
}

// StoreUnitDuration replaces the name with the new value while issuing a Mutation if Tree.tracking is true
func (fig *Tree) StoreUnitDuration(name string, value, units time.Duration) Fruit {
	return fig.Store(tUnitDuration, name, value*units)
}

// StoreList replaces the name with the new value while issuing a Mutation if Tree.tracking is true
func (fig *Tree) StoreList(name string, value []string) Fruit {
	return fig.Store(tList, name, value)
}

// StoreMap replaces the name with the new value while issuing a Mutation if Tree.tracking is true
func (fig *Tree) StoreMap(name string, value map[string]string) Fruit {
	return fig.Store(tMap, name, value)
}

// persist requires the Tree.mu to be locked before using this func and is an internal func
func (fig *Tree) persist(fruit *Fig, mut Mutagenesis, name string, value interface{}) (changed bool, previous, current interface{}) {
	flesh := fruit.Flesh
	switch mut {
	case tMap:
		var old *map[string]string
		var err error
		switch f := flesh.(type) {
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
			return false, flesh, value
		}
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
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
			return false, old, flesh
		}
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, old, value
		}
		fruit.Flesh = current
		fig.figs[name] = fruit
		return old != current, old, current
	case tList:
		var old *[]string
		var err error
		switch v := flesh.(type) {
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
			return false, flesh, value
		}
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
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
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, old, value
		}
		fruit.Flesh = current
		fig.figs[name] = fruit
		return old != current, old, current
	case tUnitDuration:
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
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
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
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, old, value
		}
		fruit.Flesh = current
		fig.figs[name] = fruit
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
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
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
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, old, value
		}
		fruit.Flesh = current
		fig.figs[name] = fruit
		return old != current, old, current
	case tFloat64:
		old, err := toFloat64(flesh)
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, flesh, value
		}
		current, err := toFloat64(value)
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, old, value
		}
		fruit.Flesh = current
		fig.figs[name] = fruit
		return old != current, old, current
	case tInt64:
		old, err := toInt64(flesh)
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, flesh, value
		}
		current, err := toInt64(value)
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, old, value
		}
		fruit.Flesh = current
		fig.figs[name] = fruit
		return old != current, old, current
	case tInt:
		old, err := toInt(flesh)
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, flesh, value
		}
		current, err := toInt(value)
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, old, value
		}
		fruit.Flesh = current
		fig.figs[name] = fruit
		return old != current, old, current
	case tString:
		old, err := toString(flesh)
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, flesh, value
		}
		current, err := toString(value)
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, old, value
		}
		fruit.Flesh = current
		fig.figs[name] = fruit
		return !strings.EqualFold(old, current), old, current
	case tBool:
		old, err := toBool(flesh)
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, flesh, value
		}
		current, err := toBool(value)
		if err != nil {
			fig.figs[name].Error = errors.Join(fig.figs[name].Error, err)
			return false, old, value
		}
		fruit.Flesh = current
		fig.figs[name] = fruit
		return old != current, old, current
	default:
		return false, flesh, value
	}
}
