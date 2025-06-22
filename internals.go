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
func (tree *figTree) loadFile(filename string) error {
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
func (tree *figTree) loadJSON(data []byte) error {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	return tree.setValuesFromMap(jsonData)
}

// loadYAML parses the DefaultYAMLFile or the value of the EnvironmentKey or ConfigFilePath into yaml.Unmarshal
func (tree *figTree) loadYAML(data []byte) error {
	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return err
	}
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.activateFlagSet()
	for n, d := range yamlData {
		var fruit *figFruit
		var exists bool
		if fruit, exists = tree.figs[n]; exists && fruit != nil {
			value := tree.useValue(tree.from(fruit.name))
			ds, err := toString(d)
			if err != nil {
				return fmt.Errorf("unable toString value for %s: %w", n, err)
			}
			err = value.Set(ds)
			if err != nil {
				return fmt.Errorf("unable to Set value for %s: %w", n, err)
			}
			tree.values.Store(fruit.name, value)
			continue
		}
		mut := tree.MutagenesisOf(d)
		vf, er := tree.from(n)
		if er == nil && vf != nil && strings.EqualFold(string(vf.Mutagensis), string(tree.MutagenesisOf(d))) {
			err := vf.Assign(d)
			if err != nil {
				return fmt.Errorf("unable to assign value for %s: %w", n, err)
			}
			tree.values.Store(n, vf)
			mut = vf.Mutagensis
		} else {
			value := &Value{
				Value:      d,
				Mutagensis: mut,
			}
			tree.values.Store(n, value)
		}
		fruit = &figFruit{
			name:        n,
			Mutagenesis: mut,
			usage:       "Unknown, loaded from config file",
			Mutations:   make([]Mutation, 0),
			Validators:  make([]FigValidatorFunc, 0),
			Callbacks:   make([]Callback, 0),
			Rules:       make([]RuleKind, 0),
		}
		withered := witheredFig{
			name: n,
			Value: Value{
				Mutagensis: mut,
				Value:      d,
			},
		}
		tree.figs[n] = fruit
		tree.withered[n] = withered
	}

	return nil
}

// loadINI parses the DefaultINIFile or the value of the EnvironmentKey or ConfigFilePath into ini.Load()
func (tree *figTree) loadINI(data []byte) error {
	cfg, err := ini.Load(data)
	if err != nil {
		return err
	}
	iniData := make(map[string]interface{})
	validKeys := make(map[string]struct{}, len(tree.figs))
	for key := range tree.figs {
		validKeys[key] = struct{}{}
	}
	for _, section := range cfg.Sections() {
		sectionName := section.Name()
		prefix := ""
		if sectionName != ini.DefaultSection {
			prefix = sectionName + "."
		}
		for _, key := range section.Keys() {
			keyName := prefix + key.Name()
			if _, exists := validKeys[keyName]; exists {
				if val, err := key.Int(); err == nil {
					iniData[keyName] = val
				} else if val, err := key.Bool(); err == nil {
					iniData[keyName] = val
				} else if val, err := key.Float64(); err == nil {
					iniData[keyName] = val
				} else {
					iniData[keyName] = key.String()
				}
			}
		}
	}
	return tree.setValuesFromMap(iniData)
}

// setValuesFromMap uses the data map to store the configurable figs
func (tree *figTree) setValuesFromMap(data map[string]interface{}) error {
	for key, value := range data {
		if _, exists := tree.figs[key]; exists {
			if err := tree.mutateFig(key, value); err != nil {
				return fmt.Errorf("error setting key %s: %w", key, err)
			}
		}
	}
	return nil
}

func (tree *figTree) setValue(flagVal interface{}, value interface{}) error {
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
		ptr.values = listVal
	case *MapFlag:
		mapVal, err := toStringMap(value)
		if err != nil {
			return err
		}
		ptr.values = mapVal
	default:
		return fmt.Errorf("unsupported flag type %T", ptr)
	}
	return nil
}

// readEnv checks the os.LookupEnv on each figFruit in the figTree
func (tree *figTree) readEnv() {
	if tree.HasRule(RuleNoEnv) {
		return
	}
	for name := range tree.figs {
		tree.checkAndSetFromEnv(name)
	}
}

// checkAndSetFromEnv uses os.LookupEnv and assigns it to the figs name value
func (tree *figTree) checkAndSetFromEnv(name string) {
	if tree.HasRule(RuleNoEnv) {
		return
	}
	if !tree.ignoreEnv {
		if val, exists := os.LookupEnv(name); exists {
			_ = tree.mutateFig(name, val)
		}
	}
}

// mutateFig replaces the value interface{} and sends a Mutation into Mutations
func (tree *figTree) mutateFig(name string, value interface{}) error {
	def, ok := tree.figs[name]
	if !ok || def == nil {
		return fmt.Errorf("no such fig: %s", name)
	}
	var old interface{}
	var dead interface{}
	witheredFig, ok := tree.withered[name]
	dead = witheredFig.Value.Value
	valueAny, ok := tree.values.Load(name)
	if !ok {
		return nil
	}
	_value, ok := valueAny.(*Value)
	if !ok {
		return nil
	}
	old = _value.Flesh()
	err := _value.Assign(value)
	if err != nil {
		return err
	}
	tree.values.Store(name, _value)
	t1 := tree.MutagenesisOf(&old)
	t2 := tree.MutagenesisOf(_value.Value)
	if t1 == "" && t2 != "" {
		t1 = t2
	}
	if !strings.EqualFold(string(t1), string(t2)) {
		return fmt.Errorf("type mismatch for key %s", name)
	}
	// if tree.tracking && old != dead && dead != value
	if tree.tracking && !reflect.DeepEqual(old, dead) && !reflect.DeepEqual(dead, value) {
		tree.mutationsCh <- Mutation{
			Property:    name,
			Mutagenesis: fmt.Sprintf("%T", value),
			Way:         "mutateFig",
			Old:         old,
			New:         value,
			When:        time.Now(),
		}
	}
	return nil
}

// activateFlagSet sets flag.CommandLine to figTree.flagSet
func (tree *figTree) activateFlagSet() Plant {
	if tree.HasRule(RuleNoFlags) {
		return tree
	}
	flag.CommandLine = tree.flagSet
	return tree
}

// assignFlagSet assigns a new *flag.FlagSet to figTree.flagSet
func (tree *figTree) assignFlagSet(newSet *flag.FlagSet) Plant {
	if tree.HasRule(RuleNoFlags) {
		return tree
	}
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

func DeduplicateStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}
