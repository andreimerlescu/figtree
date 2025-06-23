package figtree

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	zeroString   string
	zeroBool     bool
	zeroFloat64  float64
	zeroInt      int
	zeroInt64    int64
	zeroDuration time.Duration
	zeroList     []string
	zeroMap      map[string]string // Will be nil map by default, should be make(map[string]string) if always non-nil
)

func init() {
	zeroList = make([]string, 0)
	zeroMap = make(map[string]string)
}

// String with mutation tracking
func (tree *figTree) String(name string) *string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	originalName := strings.Clone(name) // in case we need it
	defer func() {
		name = originalName // restore after we're done
	}()
	name = strings.ToLower(name)
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		return nil
	}
	err := fruit.runCallbacks(tree, CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroString
	}
	value, err := tree.from(name)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroString
	}
	s := value.Flesh().ToString()
	if !tree.HasRule(RuleNoEnv) && !fruit.HasRule(RuleNoEnv) && !tree.ignoreEnv && tree.pollinate {
		envs := os.Environ()
		var e string
		var ok bool
		for _, env := range envs {
			if strings.EqualFold(env, strings.ToUpper(name)) {
				v := strings.Split(env, "=")
				e = v[1]
			}
		}
		if len(e) == 0 {
			e, ok = os.LookupEnv(strings.ToUpper(name))
		}
		if ok && len(e) > 0 {
			if !strings.EqualFold(strings.ToLower(e), strings.ToLower(s)) {
				s = strings.Clone(e)
				tree.mu.RUnlock()
				tree.Store(fruit.Mutagenesis, name, e)
				tree.mu.RLock()
				fruit = tree.figs[name]
			}
		}
	}
	err = fruit.runCallbacks(tree, CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroString
	}
	return &s
}

// Bool with mutation tracking
func (tree *figTree) Bool(name string) *bool {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	originalName := strings.Clone(name) // in case we need it
	defer func() {
		name = originalName // restore after we're done
	}()
	name = strings.ToLower(name)
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		return nil
	}
	err := fruit.runCallbacks(tree, CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroBool
	}
	value, err := tree.from(name)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroBool
	}
	s := value.Flesh().ToBool()
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
	err = fruit.runCallbacks(tree, CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroBool
	}
	return &s
}

// Int with mutation tracking
func (tree *figTree) Int(name string) *int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	originalName := strings.Clone(name) // in case we need it
	defer func() {
		name = originalName // restore after we're done
	}()
	name = strings.ToLower(name)
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		return nil
	}
	err := fruit.runCallbacks(tree, CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroInt
	}
	value, err := tree.from(name)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroInt
	}
	if value.Err != nil {
		fmt.Printf("value(%s).Err = %v\n", name, value.Err.Error())
		return &zeroInt
	}
	s := value.Flesh().ToInt()
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
	err = fruit.runCallbacks(tree, CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroInt
	}
	return &s
}

// Int64 with mutation tracking
func (tree *figTree) Int64(name string) *int64 {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	originalName := strings.Clone(name) // in case we need it
	defer func() {
		name = originalName // restore after we're done
	}()
	name = strings.ToLower(name)
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		return nil
	}
	err := fruit.runCallbacks(tree, CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroInt64
	}
	value, err := tree.from(name)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroInt64
	}
	s := value.Flesh().ToInt64()
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
	err = fruit.runCallbacks(tree, CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroInt64
	}
	return &s
}

// Float64 with mutation tracking
func (tree *figTree) Float64(name string) *float64 {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	originalName := strings.Clone(name) // in case we need it
	defer func() {
		name = originalName // restore after we're done
	}()
	name = strings.ToLower(name)
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		return nil
	}
	err := fruit.runCallbacks(tree, CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroFloat64
	}
	value, err := tree.from(name)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroFloat64
	}
	s := value.Flesh().ToFloat64()
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
	err = fruit.runCallbacks(tree, CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroFloat64
	}
	return &s
}

// Duration with mutation tracking
func (tree *figTree) Duration(name string) *time.Duration {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	originalName := strings.Clone(name) // in case we need it
	defer func() {
		name = originalName // restore after we're done
	}()
	name = strings.ToLower(name)
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		return nil
	}
	err := fruit.runCallbacks(tree, CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroDuration
	}
	var d time.Duration
	value, err := tree.from(name)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroDuration
	}
	switch f := value.Value.(type) {
	case int64:
		d = time.Duration(f)
	case *int64:
		d = time.Duration(*f)
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
	err = fruit.runCallbacks(tree, CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroDuration
	}
	return &d
}

// UnitDuration with mutation tracking
func (tree *figTree) UnitDuration(name string) *time.Duration {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	originalName := strings.Clone(name) // in case we need it
	defer func() {
		name = originalName // restore after we're done
	}()
	name = strings.ToLower(name)
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		return nil
	}
	err := fruit.runCallbacks(tree, CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroDuration
	}
	var d time.Duration
	value, err := tree.from(name)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroDuration
	}
	switch f := value.Value.(type) {
	case int64:
		d = time.Duration(f)
	case *int64:
		d = time.Duration(*f)
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
	err = fruit.runCallbacks(tree, CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroDuration
	}
	return &d
}

// List returns a copy of the []string stored inside the Value of the figFruit
func (tree *figTree) List(name string) *[]string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	originalName := strings.Clone(name) // in case we need it
	defer func() {
		name = originalName // restore after we're done
	}()
	name = strings.ToLower(name)
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		return nil
	}
	value, err := tree.from(name)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroList
	}
	err = fruit.runCallbacks(tree, CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroList
	}
	var v []string
	switch f := value.Value.(type) {
	case string:
		v = []string{f}
	case *string:
		v = []string{*f}
	case ListFlag:
		v = make([]string, len(f.values))
		copy(v, f.values)
	case *ListFlag:
		v = make([]string, len(f.values))
		copy(v, f.values)
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
			i := strings.Split(e, ListSeparator)
			if len(i) == 0 {
				v = []string{}
			} else if !slices.Equal(v, i) {
				tree.mu.RUnlock()
				tree.Store(fruit.Mutagenesis, name, i)
				tree.mu.RLock()
				value, err := tree.from(name)
				if err != nil {
					fruit.Error = errors.Join(fruit.Error, err)
					tree.figs[name] = fruit
					return &zeroList
				}
				switch f := value.Value.(type) {
				case string:
					v = []string{f}
				case *string:
					v = []string{*f}
				case ListFlag:
					v = make([]string, len(f.values))
					copy(v, f.values)
				case *ListFlag:
					v = make([]string, len(f.values))
					copy(v, f.values)
				case *[]string:
					v = make([]string, len(*f))
					copy(v, *f)
				case []string:
					v = make([]string, len(f))
					copy(v, f)
				default:
					panic("unreachable")
				}
			}
		}
	}
	err = fruit.runCallbacks(tree, CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroList
	}
	slices.Sort(v)
	return &v
}

// Map with mutation tracking
func (tree *figTree) Map(name string) *map[string]string {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	originalName := strings.Clone(name) // in case we need it
	defer func() {
		name = originalName // restore after we're done
	}()
	name = strings.ToLower(name)
	if _, exists := tree.aliases[name]; exists {
		name = tree.aliases[name]
	}
	fruit, ok := tree.figs[name]
	if !ok || fruit == nil {
		return nil
	}
	value, err := tree.from(name)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroMap
	}
	err = fruit.runCallbacks(tree, CallbackBeforeRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroMap
	}
	var v map[string]string
	switch f := value.Value.(type) {
	case string:
		v = map[string]string{}
		for _, x := range strings.Split(f, MapSeparator) {
			parts := strings.SplitN(x, MapKeySeparator, 2)
			if len(parts) != 2 {
				continue
			}
			v[parts[0]] = parts[1]
		}
	case *string:
		v = map[string]string{}
		for _, x := range strings.Split(*f, MapSeparator) {
			parts := strings.SplitN(x, MapKeySeparator, 2)
			if len(parts) != 2 {
				continue
			}
			v[parts[0]] = parts[1]
		}
	case MapFlag:
		v = make(map[string]string, len(f.values))
		for k, val := range f.values {
			v[k] = val
		}
	case *MapFlag:
		// Create a new map and copy the key-value pairs
		v = make(map[string]string, len(f.values))
		for k, val := range f.values {
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
					parts := strings.Split(iv, MapKeySeparator)
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
					value, err := tree.from(name)
					if err != nil {
						fruit.Error = errors.Join(fruit.Error, err)
						tree.figs[name] = fruit
						return &zeroMap
					}
					switch f := value.Value.(type) {
					case string:
						v = map[string]string{}
						for _, x := range strings.Split(f, MapSeparator) {
							parts := strings.SplitN(x, MapKeySeparator, 2)
							if len(parts) != 2 {
								continue
							}
							v[parts[0]] = parts[1]
						}
					case *string:
						v = map[string]string{}
						for _, x := range strings.Split(*f, MapSeparator) {
							parts := strings.SplitN(x, MapKeySeparator, 2)
							if len(parts) != 2 {
								continue
							}
							v[parts[0]] = parts[1]
						}
					case MapFlag:
						v = make(map[string]string, len(f.values))
						for k, val := range f.values {
							v[k] = val
						}
					case *MapFlag:
						v = make(map[string]string, len(f.values))
						for k, val := range f.values {
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
	err = fruit.runCallbacks(tree, CallbackAfterRead)
	if err != nil {
		fruit.Error = errors.Join(fruit.Error, err)
		tree.figs[name] = fruit
		return &zeroMap
	}
	return &v
}
