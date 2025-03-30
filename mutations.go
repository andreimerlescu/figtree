package figtree

import (
	"errors"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// String with mutation tracking
func (fig *Tree) String(name string) *string {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.RUnlock()
		fig.Resurrect(name)
		fig.mu.RLock()
		fruit = fig.figs[name]
	}
	s, _ := toString(fruit.Flesh)
	if fig.pollinate {
		envs := os.Environ()
		var e string
		var ok bool
		for _, env := range envs {
			if strings.EqualFold(env, name) {
				v := strings.Split(env, "=")
				e = v[1]
			}
		}
		if len(e) == 0 {
			e, ok = os.LookupEnv(name)
		}
		if ok && len(e) > 0 {
			if !strings.EqualFold(e, s) {
				s = strings.Clone(e)
				fig.mu.RUnlock()
				fig.Store(fruit.Mutagenesis, name, e)
				fig.mu.RLock()
				fruit = fig.figs[name]
			}
		}
	}
	return &s
}

// Bool with mutation tracking
func (fig *Tree) Bool(name string) *bool {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.RUnlock()
		fig.Resurrect(name)
		fig.mu.RLock()
		fruit = fig.figs[name]
	}
	s, _ := toBool(fruit.Flesh)
	if fig.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			pb, err := strconv.ParseBool(e)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				fig.mu.RUnlock()
				fig.mu.Lock()
				fig.figs[name] = fruit
				fig.mu.Unlock()
				fig.mu.RLock()
			}
			if pb != s {
				s = pb
				fig.mu.RUnlock()
				fig.Store(fruit.Mutagenesis, name, pb)
				fig.mu.RLock()
			}
		}
	}
	return &s
}

// Int with mutation tracking
func (fig *Tree) Int(name string) *int {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.RUnlock()
		fig.Resurrect(name)
		fig.mu.RLock()
		fruit = fig.figs[name]
	}
	s, _ := toInt(fruit.Flesh)
	if fig.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			h, err := strconv.Atoi(e)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				fig.mu.RUnlock()
				fig.mu.Lock()
				fig.figs[name] = fruit
				fig.mu.Unlock()
				fig.mu.RLock()
			}
			if s != h {
				s = h
				fig.mu.RUnlock()
				fig.Store(fruit.Mutagenesis, name, h)
				fig.mu.RLock()
			}

		}
	}
	return &s
}

// Int64 with mutation tracking
func (fig *Tree) Int64(name string) *int64 {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.RUnlock()
		fig.Resurrect(name)
		fig.mu.RLock()
		fruit = fig.figs[name]
	}
	s, _ := toInt64(fruit.Flesh)
	if fig.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			h, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				fig.mu.RUnlock()
				fig.mu.Lock()
				fig.figs[name] = fruit
				fig.mu.Unlock()
				fig.mu.RLock()
			}
			if s != h {
				s = h
				fig.mu.RUnlock()
				fig.Store(fruit.Mutagenesis, name, h)
				fig.mu.RLock()
			}

		}
	}
	return &s
}

// Float64 with mutation tracking
func (fig *Tree) Float64(name string) *float64 {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.RUnlock()
		fig.Resurrect(name)
		fig.mu.RLock()
		fruit = fig.figs[name]
	}
	s, _ := toFloat64(fruit.Flesh)
	if fig.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			h, err := strconv.ParseFloat(e, 64)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				fig.mu.RUnlock()
				fig.mu.Lock()
				fig.figs[name] = fruit
				fig.mu.Unlock()
				fig.mu.RLock()
			}
			if s != h {
				s = h
				fig.mu.RUnlock()
				fig.Store(fruit.Mutagenesis, name, h)
				fig.mu.RLock()
			}

		}
	}
	return &s
}

// Duration with mutation tracking
func (fig *Tree) Duration(name string) *time.Duration {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.RUnlock()
		fig.Resurrect(name)
		fig.mu.RLock()
		fruit = fig.figs[name]
	}
	var d time.Duration
	switch f := fruit.Flesh.(type) {
	case time.Duration:
		d = f
	case *time.Duration:
		d = *f
	default:
		return nil
	}
	if fig.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			h, err := time.ParseDuration(e)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				fig.mu.RUnlock()
				fig.mu.Lock()
				fig.figs[name] = fruit
				fig.mu.Unlock()
				fig.mu.RLock()
			}
			if h != d {
				d = h
				fig.mu.RUnlock()
				fig.Store(fruit.Mutagenesis, name, h)
				fig.mu.RLock()
			}
		}
	}
	return &d
}

// UnitDuration with mutation tracking
func (fig *Tree) UnitDuration(name string) *time.Duration {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.RUnlock()
		fig.Resurrect(name)
		fig.mu.RLock()
		fruit = fig.figs[name]
	}
	var d time.Duration
	switch f := fruit.Flesh.(type) {
	case time.Duration:
		d = f
	case *time.Duration:
		d = *f
	default:
		return nil
	}
	if fig.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			h, err := time.ParseDuration(e)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				fig.mu.RUnlock()
				fig.mu.Lock()
				fig.figs[name] = fruit
				fig.mu.Unlock()
				fig.mu.RLock()
			}
			if h != d {
				d = h
				fig.mu.RUnlock()
				fig.Store(fruit.Mutagenesis, name, h)
				fig.mu.RLock()
			}
		}
	}
	return &d
}

// List with mutation tracking
func (fig *Tree) List(name string) *[]string {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.RUnlock()
		fig.Resurrect(name)
		fig.mu.RLock()
		fruit = fig.figs[name]
	}
	var v []string
	switch f := fruit.Flesh.(type) {
	case *ListFlag:
		v = make([]string, len(*f.values))
		copy(v, *f.values)
	case *[]string:
		v = make([]string, len(*f))
		copy(v, *f)
	case []string:
		v = make([]string, len(f))
		copy(v, f)
	default:
		return nil
	}
	if fig.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			i := strings.Split(e, ",")
			if len(i) == 0 {
				v = []string{}
			} else if !slices.Equal(v, i) {
				fig.mu.RUnlock()
				fig.Store(fruit.Mutagenesis, name, i)
				fig.mu.RLock()
				fruit = fig.figs[name]
				switch f := fruit.Flesh.(type) {
				case *ListFlag:
					v = make([]string, len(*f.values))
					copy(v, *f.values)
				case *[]string:
					v = make([]string, len(*f))
					copy(v, *f)
				case []string:
					v = make([]string, len(f))
					copy(v, f)
				}
			}
		}
	}
	return &v
}

// Map with mutation tracking
func (fig *Tree) Map(name string) *map[string]string {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	fruit, ok := fig.figs[name]
	if !ok || fruit == nil {
		fig.mu.RUnlock()
		fig.Resurrect(name)
		fig.mu.RLock()
		fruit = fig.figs[name]
	}
	var v map[string]string
	switch f := fruit.Flesh.(type) {
	case *MapFlag:
		// Create a new map and copy the key-value pairs
		v = make(map[string]string, len(*f.values))
		for k, val := range *f.values {
			v[k] = val
		}
	case *map[string]string:
		v = make(map[string]string, len(*f))
		for k, val := range *f {
			v[k] = val
		}
	case map[string]string:
		v = make(map[string]string, len(f))
		for k, val := range f {
			v[k] = val
		}
	default:
		return nil
	}
	if fig.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			i := strings.Split(e, ",")
			if len(i) == 0 {
				v = map[string]string{}
			} else {
				newMap := make(map[string]string)
				for _, iv := range i {
					parts := strings.Split(iv, "=")
					if len(parts) == 2 {
						newMap[parts[0]] = parts[1]
					}
				}
				equal := true
				if len(v) == len(newMap) {
					for k, val := range v {
						if newVal, exists := newMap[k]; !exists || newVal != val {
							equal = false
							break
						}
					}
				} else {
					equal = false
				}
				if !equal {
					fig.mu.RUnlock()
					fig.Store(fruit.Mutagenesis, name, newMap)
					fig.mu.RLock()
					fruit = fig.figs[name]
					switch f := fruit.Flesh.(type) {
					case *MapFlag:
						v = make(map[string]string, len(*f.values))
						for k, val := range *f.values {
							v[k] = val
						}
					case *map[string]string:
						v = make(map[string]string, len(*f))
						for k, val := range *f {
							v[k] = val
						}
					case map[string]string:
						v = make(map[string]string, len(f))
						for k, val := range f {
							v[k] = val
						}
					}
				}
			}
		}
	}
	return &v
}
