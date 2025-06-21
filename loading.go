package figtree

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	check "github.com/andreimerlescu/checkfs"
	"github.com/andreimerlescu/checkfs/file"
)

// Reload will readEnv on each flag in the configurable package
func (tree *figTree) Reload() error {
	tree.readEnv()
	return tree.validateAll()
}

// Load uses the EnvironmentKey and the DefaultJSONFile, DefaultYAMLFile, and DefaultINIFile to run ParseFile if it exists
func (tree *figTree) Load() (err error) {
	if !tree.HasRule(RuleNoFlags) {
		tree.activateFlagSet()
		args := os.Args[1:]
		if tree.filterTests {
			args = filterTestFlags(args)
			err = tree.flagSet.Parse(args)
		} else {
			err = tree.flagSet.Parse(args)
		}
		err = tree.loadFlagSet()
		if err != nil {
			return err
		}
	}
	first := ""
	if !tree.HasRule(RuleNoEnv) {
		first = os.Getenv(EnvironmentKey)
	}
	files := []string{
		first,
		tree.ConfigFilePath,
		ConfigFilePath,
		filepath.Join(".", DefaultJSONFile),
		filepath.Join(".", DefaultINIFile),
	}
	for i := 0; i < len(files); i++ {
		f := files[i]
		if f == "" {
			continue
		}
		if err := check.File(f, file.Options{Exists: true}); err == nil {
			if err := tree.loadFile(f); err != nil {
				return fmt.Errorf("failed to load %s: %w", f, err)
			}
		}
	}

	tree.readEnv()
	return tree.validateAll()
}

// LoadFile accepts a path and uses it to populate the figTree
func (tree *figTree) LoadFile(path string) (err error) {
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
	var loadErr error
	if loadErr = check.File(path, file.Options{Exists: true}); loadErr == nil {
		if err2 := tree.loadFile(path); err2 != nil {
			return fmt.Errorf("failed to loadFile %s: %w", path, err2)
		}
		tree.readEnv()
		err3 := tree.loadFlagSet()
		if err3 != nil {
			return err3
		}
		err4 := tree.validateAll()
		if err4 != nil {
			return fmt.Errorf("failed to validateAll: %w", err4)
		}
		return nil
	}
	err3 := tree.loadFlagSet()
	if err3 != nil {
		return err3
	}
	tree.readEnv()
	err4 := tree.validateAll()
	if err4 != nil {
		return fmt.Errorf("failed to validateAll: %w", err4)
	}
	return fmt.Errorf("failed to LoadFile %s due to err %v", path, loadErr)
}

func (tree *figTree) loadFlagSet() (e error) {
	defer func() {
		if e != nil {
			_, _ = fmt.Fprintln(os.Stderr, e)
		}
		if r := recover(); r != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v", r)
		}
	}()
	tree.flagSet.VisitAll(func(f *flag.Flag) {
		value := tree.useValue(tree.from(f.Name))
		if value.Mutagensis == tMap && tree.MutagenesisOf(f.Value) == tMap {
			merged := value.Flesh().ToMap()
			withered := tree.withered[f.Name]
			witheredValue := withered.Value.Flesh().ToMap()
			flagged, err := toStringMap(f.Value)
			if err != nil {
				e = fmt.Errorf("failed to load %s: %w", f.Name, err)
				return
			}
			result := make(map[string]string)
			for k, v := range witheredValue {
				result[k] = v
			}
			for k, v := range merged {
				result[k] = v
			}
			for k, v := range flagged {
				result[k] = v
			}
			err = value.Assign(result)
			if err != nil {
				e = fmt.Errorf("failed to load %s: %w", f.Name, err)
				return
			}
			tree.values.Store(f.Name, value)
			return
		}
		if value.Mutagensis == tList && tree.MutagenesisOf(f.Value) == tList {
			merged, err := toStringSlice(value.Value)
			if err != nil {
				e = fmt.Errorf("failed to load %s: %w", f.Name, err)
				return
			}
			flagged, err := toStringSlice(f.Value)
			if err != nil {
				e = fmt.Errorf("failed to load %s: %w", f.Name, err)
				return
			}
			unique := make(map[string]bool)
			for _, v := range merged {
				unique[v] = true
			}
			for _, v := range flagged {
				unique[v] = true
			}
			var newValue []string
			for k, _ := range unique {
				newValue = append(newValue, k)
			}
			err = value.Assign(newValue)
			if err != nil {
				e = fmt.Errorf("failed to load %s: %w", f.Name, err)
			}
			tree.values.Store(f.Name, value)
			return

		}
		err := value.Set(f.Value.String())
		if err != nil {
			e = fmt.Errorf("failed to value.Set(%s) due to err: %w", f.Value.String(), err)
			return
		}
		tree.values.Store(f.Name, value)
	})
	return nil
}
