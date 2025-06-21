package figtree

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
)

// Parsing Configuration

func (tree *figTree) useValue(value *Value, err error) *Value {
	if err != nil {
		log.Printf("useValue caught err: %v", err)
	}
	return value
}

func (tree *figTree) from(name string) (*Value, error) {
	valueAny, ok := tree.values.Load(name)
	if !ok {
		return nil, errors.New("no value for " + name)
	}
	value, ok := valueAny.(*Value)
	if !ok {
		return nil, errors.New("value for " + name + " is not a Value")
	}
	if value.Mutagensis == "" {
		value.Mutagensis = tree.MutagenesisOf(value)
	}
	return value, nil
}

// Parse uses figTree.flagSet to run flag.Parse() on the registered figs and returns nil for validated results
func (tree *figTree) Parse() (err error) {
	if !tree.HasRule(RuleNoFlags) {
		tree.activateFlagSet()
		args := os.Args[1:]
		if tree.filterTests {
			args = filterTestFlags(args)
		}
		err = tree.flagSet.Parse(args)
		if err != nil {
			return err
		}
		tree.mu.Lock()
		for name, fig := range tree.figs {
			if fig == nil {
				continue
			}
			value := tree.useValue(tree.from(name))
			if value.Mutagensis == tMap && PolicyMapAppend {
				vm := value.Flesh().ToMap()
				unique := make(map[string]string)
				withered := tree.withered[name]
				for k, v := range vm {
					unique[k] = v
				}
				for k, v := range withered.Value.Flesh().ToMap() {
					unique[k] = v
				}
				err = value.Assign(unique)
				if err != nil {
					return fmt.Errorf("failed to assign %s due to %w", name, err)
				}
				tree.values.Store(name, value)
			}
			if value.Mutagensis == tList && PolicyListAppend {
				vl := value.Flesh().ToList()
				w := tree.withered[name]
				wl := w.Value.Flesh().ToList()
				unique := make(map[string]struct{})
				for _, v := range wl {
					unique[v] = struct{}{}
				}
				for _, v := range vl {
					unique[v] = struct{}{}
				}
				withered := tree.withered[name]
				for _, w := range withered.Value.Flesh().ToList() {
					unique[w] = struct{}{}
				}
				var result []string
				for k, _ := range unique {
					result = append(result, k)
				}
				sort.Strings(result)
				err = value.Assign(result)
				if err != nil {
					return fmt.Errorf("failed assign to %s: %w", name, err)
				}
				tree.values.Store(name, value)
			}
		}
		tree.mu.Unlock()
		err = tree.loadFlagSet()
		if err != nil {
			return err
		}
		tree.readEnv()
		return tree.validateAll()
	}
	tree.readEnv()
	tree.mu.Lock()
	for name, fig := range tree.figs {
		if fig == nil {
			continue
		}
		value := tree.useValue(tree.from(name))
		if value.Mutagensis == tMap && PolicyMapAppend {
			vm := value.Flesh().ToMap()
			unique := make(map[string]string)
			withered := tree.withered[name]
			for k, v := range vm {
				unique[k] = v
			}
			for k, v := range withered.Value.Flesh().ToMap() {
				unique[k] = v
			}
			err = value.Assign(unique)
			if err != nil {
				return fmt.Errorf("failed to assign %s due to %w", name, err)
			}
			tree.values.Store(name, value)
		}
		if value.Mutagensis == tList && PolicyListAppend {
			vl, e := toStringSlice(value.Value)
			if e != nil {
				return fmt.Errorf("failed toStringSlice: %w", e)
			}
			unique := make(map[string]struct{})
			for _, v := range vl {
				unique[v] = struct{}{}
			}
			withered := tree.withered[name]
			for _, w := range withered.Value.Flesh().ToList() {
				unique[w] = struct{}{}
			}
			var result []string
			for k, _ := range unique {
				result = append(result, k)
			}
			sort.Strings(result)
			err = value.Assign(result)
			if err != nil {
				return fmt.Errorf("failed assign to %s: %w", name, err)
			}
			tree.values.Store(name, value)
		}
	}
	tree.mu.Unlock()
	return tree.validateAll()
}

// ParseFile will check if filename is set and run loadFile on it.
func (tree *figTree) ParseFile(filename string) (err error) {
	if !tree.HasRule(RuleNoFlags) {
		tree.activateFlagSet()
		args := os.Args[1:]
		if tree.filterTests {
			args = filterTestFlags(args)
			err = tree.flagSet.Parse(args)
		} else {
			err = tree.flagSet.Parse(args)
		}
		if err != nil {
			return err
		}
	}
	err = tree.loadFlagSet()
	if err != nil {
		return err
	}
	if filename != "" {
		return tree.loadFile(filename)
	}
	tree.readEnv()
	return tree.validateAll()
}
