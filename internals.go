package figtree

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"gopkg.in/yaml.v3"
)

// CONFIGURABLE INTERNAL FUNCTIONS

// loadFile will parse the filename for .yaml or .ini or .json and run the related loadJSON, loadYAML, or loadINI on it
func (tree *Tree) loadFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return tree.loadJSON(data)
	case ".yaml", ".yml":
		return tree.loadYAML(data)
	case ".ini":
		return tree.loadINI(data)
	default:
		return errors.New("unsupported file extension")
	}
}

// loadJSON parses the DefaultJSONFile or the value of the EnvironmentKey or ConfigFilePath into json.Unmarshal
func (tree *Tree) loadJSON(data []byte) error {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	return tree.setValuesFromMap(jsonData)
}

// loadYAML parses the DefaultYAMLFile or the value of the EnvironmentKey or ConfigFilePath into yaml.Unmarshal
func (tree *Tree) loadYAML(data []byte) error {
	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return err
	}
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	for n, d := range yamlData {
		t := tree.MutagensisOf(d)
		var fruit *Fig
		var exists bool
		if fruit, exists = tree.figs[n]; exists && fruit != nil {
			tf := tree.MutagensisOf(fruit.Flesh)
			if strings.EqualFold(string(tf), string(t)) {
				tree.figs[n].Flesh = d
				continue
			}
		}

		switch d.(type) {
		case *string:
			tree.figs[n] = &Fig{
				Flesh:         d,
				Mutagenesis:   tString,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
			tree.withered[n] = Fig{
				Flesh:         d,
				Mutagenesis:   tString,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
		case *bool:
			tree.figs[n] = &Fig{
				Flesh:         d,
				Mutagenesis:   tBool,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
			tree.withered[n] = Fig{
				Flesh:         d,
				Mutagenesis:   tBool,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
		case *int:
			tree.figs[n] = &Fig{
				Flesh:         d,
				Mutagenesis:   tInt,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
			tree.withered[n] = Fig{
				Flesh:         d,
				Mutagenesis:   tInt,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
		case *int64:
			tree.figs[n] = &Fig{
				Flesh:         d,
				Mutagenesis:   tInt64,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
			tree.withered[n] = Fig{
				Flesh:         d,
				Mutagenesis:   tInt64,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
		case *float64:
			tree.figs[n] = &Fig{
				Flesh:         d,
				Mutagenesis:   tFloat64,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
			tree.withered[n] = Fig{
				Flesh:         d,
				Mutagenesis:   tFloat64,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
		case *time.Duration:
			tree.figs[n] = &Fig{
				Flesh:         d,
				Mutagenesis:   tDuration,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
			tree.withered[n] = Fig{
				Flesh:         d,
				Mutagenesis:   tDuration,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
		case *[]string:
			tree.figs[n] = &Fig{
				Flesh:         d,
				Mutagenesis:   tList,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
			tree.withered[n] = Fig{
				Flesh:         d,
				Mutagenesis:   tList,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
		case *map[string]string:
			tree.figs[n] = &Fig{
				Flesh:         d,
				Mutagenesis:   tMap,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
			tree.withered[n] = Fig{
				Flesh:         d,
				Mutagenesis:   tMap,
				Mutations:     make([]Mutation, 0),
				Validators:    make([]ValidatorFunc, 0),
				Callbacks:     make([]Callback, 0),
				CallbackAfter: CallbackAfterVerify,
			}
		}
	}

	return tree.setValuesFromMap(yamlData)
}

// loadINI parses the DefaultINIFile or the value of the EnvironmentKey or ConfigFilePath into ini.Load()
func (tree *Tree) loadINI(data []byte) error {
	cfg, err := ini.Load(data)
	if err != nil {
		return err
	}
	iniData := make(map[string]interface{})
	for key := range tree.figs {
		if val := cfg.Section("").Key(key).String(); val != "" {
			iniData[key] = val
		}
	}
	return tree.setValuesFromMap(iniData)
}

// setValuesFromMap uses the data map to store the configurable figs
func (tree *Tree) setValuesFromMap(data map[string]interface{}) error {
	for key, value := range data {
		if _, exists := tree.figs[key]; exists {
			if err := tree.mutateFig(key, value); err != nil {
				return fmt.Errorf("error setting key %s: %w", key, err)
			}
		}
	}
	return nil
}

func (tree *Tree) setValue(flagVal interface{}, value interface{}) error {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	switch ptr := flagVal.(type) {
	case *int:
		if v, ok := value.(int); ok {
			*ptr = v
			return nil
		}
		intVal, err := toInt(value)
		if err != nil {
			return err
		}
		*ptr = intVal
	case *int64:
		if v, ok := value.(int64); ok {
			*ptr = v
			return nil
		}
		int64Val, err := toInt64(value)
		if err != nil {
			return err
		}
		*ptr = int64Val
	case *float64:
		if v, ok := value.(float64); ok {
			*ptr = v
			return nil
		}
		floatVal, err := toFloat64(value)
		if err != nil {
			return err
		}
		*ptr = floatVal
	case *string:
		if v, ok := value.(string); ok {
			*ptr = v
			return nil
		}
		strVal, err := toString(value)
		if err != nil {
			return err
		}
		*ptr = strVal
	case *bool:
		if v, ok := value.(bool); ok {
			*ptr = v
			return nil
		}
		boolVal, err := toBool(value)
		if err != nil {
			return err
		}
		*ptr = boolVal
	case *time.Duration:
		if v, ok := value.(time.Duration); ok {
			*ptr = v
			return nil
		}
		strVal, err := toString(value)
		if err != nil {
			return err
		}
		duration, err := time.ParseDuration(strVal)
		if err != nil {
			return err
		}
		*ptr = duration
	case *ListFlag:
		listVal, err := toStringSlice(value)
		if err != nil {
			return err
		}
		*ptr.values = listVal
	case *MapFlag:
		mapVal, err := toStringMap(value)
		if err != nil {
			return err
		}
		*ptr.values = mapVal
	default:
		return fmt.Errorf("unsupported flag type %T", ptr)
	}
	return nil
}

// readEnv checks the os.LookupEnv on each Fig in the Tree
func (tree *Tree) readEnv() {
	for name := range tree.figs {
		tree.checkAndSetFromEnv(name)
	}
}

// checkAndSetFromEnv uses os.LookupEnv and assigns it to the figs name value
func (tree *Tree) checkAndSetFromEnv(name string) {
	if !tree.ignoreEnv {
		if val, exists := os.LookupEnv(name); exists {
			_ = tree.mutateFig(name, val)
		}
	}
}

// mutateFig replaces the value interface{} and sends a Mutation into Mutations
func (tree *Tree) mutateFig(name string, value interface{}) error {
	def, ok := tree.figs[name]
	if !ok || def == nil {
		tree.Resurrect(name)
		def = tree.figs[name]
	}
	if def == nil {
		return fmt.Errorf("no definition for key %s", name)
	}
	var old interface{}
	var dead interface{}
	witheredFig, ok := tree.withered[name]
	dead = witheredFig.Flesh
	old = def.Flesh
	def.Flesh = value
	t1 := tree.MutagensisOf(&old)
	t2 := tree.MutagensisOf(value)
	if t1 == "" && t2 != "" {
		t1 = t2
	}
	if !strings.EqualFold(string(t1), string(t2)) {
		return fmt.Errorf("type mismatch for key %s", name)
	}
	// if tree.tracking && old != dead && dead != value
	if tree.tracking && !reflect.DeepEqual(old, dead) && !reflect.DeepEqual(dead, value) {
		tree.mutationsCh <- Mutation{
			Property: name,
			Kind:     fmt.Sprintf("%T", value),
			Way:      "mutateFig",
			Old:      old,
			New:      value,
			When:     time.Now(),
		}
	}
	return nil
}

// activateFlagSet sets flag.CommandLine to Tree.flagSet
func (tree *Tree) activateFlagSet() Fruit {
	flag.CommandLine = tree.flagSet
	return tree
}

// assignFlagSet assigns a new *flag.FlagSet to Tree.flagSet
func (tree *Tree) assignFlagSet(newSet *flag.FlagSet) Fruit {
	tree.flagSet = newSet
	return tree
}

// filterTestFlags removes test-specific flags (e.g., -test.v) from the args slice
func filterTestFlags(args []string) []string {
	var filtered []string
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-test.") {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}
