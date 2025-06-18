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
func (tree *figTree) String(name string) *string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		tree.mu.RUnlock()
		tree.Resurrect(name)
		tree.mu.RLock()
		fruit = tree.figs[name]
	}
	err := fruit.runCallbacks(CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return nil
	}
	s, err := toString(fruit.Value.Value)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return nil
	}
	if !tree.HasRule(RuleNoEnv) && !fruit.HasRule(RuleNoEnv) && !tree.ignoreEnv && tree.pollinate {
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
				tree.mu.RUnlock()
				tree.Store(fruit.Mutagenesis, name, e)
				tree.mu.RLock()
				fruit = tree.figs[name]
			}
		}
	}
	err = fruit.runCallbacks(CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
	}
	return &s
}

// Bool with mutation tracking
func (tree *figTree) Bool(name string) *bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		tree.mu.RUnlock()
		tree.Resurrect(name)
		tree.mu.RLock()
		fruit = tree.figs[name]
	}
	err := fruit.runCallbacks(CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return nil
	}
	s := fruit.Value.Flesh().ToBool()
	if !tree.HasRule(RuleNoEnv) && !fruit.HasRule(RuleNoEnv) && !tree.ignoreEnv && tree.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			pb, err := strconv.ParseBool(e)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				tree.mu.RUnlock()
				tree.mu.Lock()
				tree.figs[name] = fruit
				tree.mu.Unlock()
				tree.mu.RLock()
			}
			if pb != s {
				s = pb
				tree.mu.RUnlock()
				tree.Store(fruit.Mutagenesis, name, pb)
				tree.mu.RLock()
			}
		}
	}
	err = fruit.runCallbacks(CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
	}
	return &s
}

// Int with mutation tracking
func (tree *figTree) Int(name string) *int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		tree.mu.RUnlock()
		tree.Resurrect(name)
		tree.mu.RLock()
		fruit = tree.figs[name]
	}
	err := fruit.runCallbacks(CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return nil
	}
	s := fruit.Value.Flesh().ToInt()
	if !tree.HasRule(RuleNoEnv) && !fruit.HasRule(RuleNoEnv) && !tree.ignoreEnv && tree.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			h, err := strconv.Atoi(e)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				tree.mu.RUnlock()
				tree.mu.Lock()
				tree.figs[name] = fruit
				tree.mu.Unlock()
				tree.mu.RLock()
			}
			if s != h {
				s = h
				tree.mu.RUnlock()
				tree.Store(fruit.Mutagenesis, name, h)
				tree.mu.RLock()
			}

		}
	}
	err = fruit.runCallbacks(CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
	}
	return &s
}

// Int64 with mutation tracking
func (tree *figTree) Int64(name string) *int64 {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		tree.mu.RUnlock()
		tree.Resurrect(name)
		tree.mu.RLock()
		fruit = tree.figs[name]
	}
	err := fruit.runCallbacks(CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return nil
	}
	s := fruit.Value.Flesh().ToInt64()
	if !tree.HasRule(RuleNoEnv) && !fruit.HasRule(RuleNoEnv) && !tree.ignoreEnv && tree.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			h, err := strconv.ParseInt(e, 10, 64)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				tree.mu.RUnlock()
				tree.mu.Lock()
				tree.figs[name] = fruit
				tree.mu.Unlock()
				tree.mu.RLock()
			}
			if s != h {
				s = h
				tree.mu.RUnlock()
				tree.Store(fruit.Mutagenesis, name, h)
				tree.mu.RLock()
			}

		}
	}
	err = fruit.runCallbacks(CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
	}
	return &s
}

// Float64 with mutation tracking
func (tree *figTree) Float64(name string) *float64 {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		tree.mu.RUnlock()
		tree.Resurrect(name)
		tree.mu.RLock()
		fruit = tree.figs[name]
	}
	err := fruit.runCallbacks(CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return nil
	}
	s := fruit.Value.Flesh().ToFloat64()
	if !tree.HasRule(RuleNoEnv) && !fruit.HasRule(RuleNoEnv) && !tree.ignoreEnv && tree.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			h, err := strconv.ParseFloat(e, 64)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				tree.mu.RUnlock()
				tree.mu.Lock()
				tree.figs[name] = fruit
				tree.mu.Unlock()
				tree.mu.RLock()
			}
			if s != h {
				s = h
				tree.mu.RUnlock()
				tree.Store(fruit.Mutagenesis, name, h)
				tree.mu.RLock()
			}

		}
	}
	err = fruit.runCallbacks(CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
	}
	return &s
}

// Duration with mutation tracking
func (tree *figTree) Duration(name string) *time.Duration {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		tree.mu.RUnlock()
		tree.Resurrect(name)
		tree.mu.RLock()
		fruit = tree.figs[name]
	}
	err := fruit.runCallbacks(CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return nil
	}
	var d time.Duration
	switch f := fruit.Value.Value.(type) {
	case time.Duration:
		d = f
	case *time.Duration:
		d = *f
	default:
		return nil
	}
	if !tree.HasRule(RuleNoEnv) && !fruit.HasRule(RuleNoEnv) && !tree.ignoreEnv && tree.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			h, err := time.ParseDuration(e)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				tree.mu.RUnlock()
				tree.mu.Lock()
				tree.figs[name] = fruit
				tree.mu.Unlock()
				tree.mu.RLock()
			}
			if h != d {
				d = h
				tree.mu.RUnlock()
				tree.Store(fruit.Mutagenesis, name, h)
				tree.mu.RLock()
			}
		}
	}
	err = fruit.runCallbacks(CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
	}
	return &d
}

// UnitDuration with mutation tracking
func (tree *figTree) UnitDuration(name string) *time.Duration {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		tree.mu.RUnlock()
		tree.Resurrect(name)
		tree.mu.RLock()
		fruit = tree.figs[name]
	}
	err := fruit.runCallbacks(CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return nil
	}
	var d time.Duration
	switch f := fruit.Value.Value.(type) {
	case time.Duration:
		d = f
	case *time.Duration:
		d = *f
	default:
		return nil
	}
	if !tree.HasRule(RuleNoEnv) && !fruit.HasRule(RuleNoEnv) && !tree.ignoreEnv && tree.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			h, err := time.ParseDuration(e)
			if err != nil {
				fruit.Error = errors.Join(fruit.Error, err)
				tree.mu.RUnlock()
				tree.mu.Lock()
				tree.figs[name] = fruit
				tree.mu.Unlock()
				tree.mu.RLock()
			}
			if h != d {
				d = h
				tree.mu.RUnlock()
				tree.Store(fruit.Mutagenesis, name, h)
				tree.mu.RLock()
			}
		}
	}
	err = fruit.runCallbacks(CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
	}
	return &d
}

// List with mutation tracking
func (tree *figTree) List(name string) *[]string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		tree.mu.RUnlock()
		tree.Resurrect(name)
		tree.mu.RLock()
		fruit = tree.figs[name]
	}
	err := fruit.runCallbacks(CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return nil
	}
	var v []string
	switch f := fruit.Value.Value.(type) {
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
	if !tree.HasRule(RuleNoEnv) && !fruit.HasRule(RuleNoEnv) && !tree.ignoreEnv && tree.pollinate {
		e := os.Getenv(name)
		if len(e) > 0 {
			i := strings.Split(e, ",")
			if len(i) == 0 {
				v = []string{}
			} else if !slices.Equal(v, i) {
				tree.mu.RUnlock()
				tree.Store(fruit.Mutagenesis, name, i)
				tree.mu.RLock()
				fruit = tree.figs[name]
				switch f := fruit.Value.Value.(type) {
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
	err = fruit.runCallbacks(CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
	}
	return &v
}

// Map with mutation tracking
func (tree *figTree) Map(name string) *map[string]string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		tree.mu.RUnlock()
		tree.Resurrect(name)
		tree.mu.RLock()
		fruit = tree.figs[name]
	}
	err := fruit.runCallbacks(CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return nil
	}
	var v map[string]string
	switch f := fruit.Value.Value.(type) {
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
	if !tree.HasRule(RuleNoEnv) && !fruit.HasRule(RuleNoEnv) && !tree.ignoreEnv && tree.pollinate {
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
					tree.mu.RUnlock()
					tree.Store(fruit.Mutagenesis, name, newMap)
					tree.mu.RLock()
					fruit = tree.figs[name]
					switch f := fruit.Value.Value.(type) {
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
	err = fruit.runCallbacks(CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
	}
	return &v
}
