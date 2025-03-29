package figtree

import (
	"time"
)

// String with mutation tracking
func (fig *Tree) String(name string) *string {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.Lock()
		fig.Resurrect(name)
		fig.mu.Unlock()
		fruit = fig.figs[name]
	}
	s, _ := toString(fruit.Flesh)
	return &s
}

// Bool with mutation tracking
func (fig *Tree) Bool(name string) *bool {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.Unlock()
		fig.Resurrect(name)
		fig.mu.Lock()
		fruit = fig.figs[name]
	}
	s, _ := toBool(fruit.Flesh)
	return &s
}

// Int with mutation tracking
func (fig *Tree) Int(name string) *int {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.Unlock()
		fig.Resurrect(name)
		fig.mu.Lock()
		fruit = fig.figs[name]
	}
	s, _ := toInt(fruit.Flesh)
	return &s
}

// Int64 with mutation tracking
func (fig *Tree) Int64(name string) *int64 {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.Unlock()
		fig.Resurrect(name)
		fig.mu.Lock()
		fruit = fig.figs[name]
	}
	s, _ := toInt64(fruit.Flesh)
	return &s
}

// Float64 with mutation tracking
func (fig *Tree) Float64(name string) *float64 {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.Unlock()
		fig.Resurrect(name)
		fig.mu.Lock()
		fruit = fig.figs[name]
	}
	s, _ := toFloat64(fruit.Flesh)
	return &s
}

// Duration with mutation tracking
func (fig *Tree) Duration(name string) *time.Duration {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.Unlock()
		fig.Resurrect(name)
		fig.mu.Lock()
		fruit = fig.figs[name]
	}
	switch f := fruit.Flesh.(type) {
	case time.Duration:
		return &f
	case *time.Duration:
		return f
	default:
		return nil
	}
}

// UnitDuration with mutation tracking
func (fig *Tree) UnitDuration(name string) *time.Duration {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.Unlock()
		fig.Resurrect(name)
		fig.mu.Lock()
		fruit = fig.figs[name]
	}
	switch f := fruit.Flesh.(type) {
	case time.Duration:
		return &f
	case *time.Duration:
		return f
	default:
		return nil
	}
}

// List with mutation tracking
func (fig *Tree) List(name string) *[]string {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.Unlock()
		fig.Resurrect(name)
		fig.mu.Lock()
		fruit = fig.figs[name]
	}
	var v []string
	switch f := fruit.Flesh.(type) {
	case *ListFlag:
		return f.values
	case *[]string:
		v = *f
	case []string:
		v = f
	default:
	}
	return &v
}

// Map with mutation tracking
func (fig *Tree) Map(name string) *map[string]string {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.Unlock()
		fig.Resurrect(name)
		fig.mu.Lock()
		fruit = fig.figs[name]
	}
	switch f := fruit.Flesh.(type) {
	case *MapFlag:
		return f.values
	case *map[string]string:
		return f
	case map[string]string:
		return &f
	default:
		return nil
	}
}
