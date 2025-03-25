package figs

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	check "github.com/andreimerlescu/checkfs"
	"github.com/andreimerlescu/checkfs/file"
	"github.com/go-ini/ini"
	"gopkg.in/yaml.v3"
)

const (
	DefaultYAMLFile string = "config.yml"  // Default filename for a YAML configuration file
	DefaultJSONFile string = "config.json" // Default filename for a JSON configuration file
	DefaultINIFile  string = "config.ini"  // Default filename for a INI configuration file
)

// EnvironmentKey stores the preferred ENV that contains the path to your configuration file (.ini, .json or .yaml)
var EnvironmentKey string = "CONFIG_FILE"

// ConfigFilePath stores the path to the configuration file of choice
var ConfigFilePath string = filepath.Join(".", DefaultYAMLFile)

// Fruit defines the interface for configuration management.
type Fruit interface {
	// Parse can panic but interprets command line arguments defined with single dashes -example value -another sample
	Parse()
	// ParseFile can panic but also can throw an error because it will attempt to load either JSON, YAML or INI file passed into it
	ParseFile(filename string) error

	// Load can panic but also can throw an error but will use the Environment Variable values if they are "EXAMPLE=value" or "ANOTHER=sample"
	Load() error
	// Reload will refresh stored values of properties with their new Environment Variable values
	Reload()

	// Usage displays the helpful menu of figs registered using -h or -help
	Usage() string

	// Int returns a pointer to a registered int32 by name as -name=1 a pointer to 1 is returned
	Int(name string) *int
	// NewInt registers a new int32 flag by name and returns a pointer to the int32 storing the initial value
	NewInt(name string, value int, usage string) *int

	// Int64 returns a pointer to a registered int64 by name as -name=1 a pointer to 1 is returned
	Int64(name string) *int64
	// NewInt64 registers a new int32 flag by name and returns a pointer to the int64 storing the initial value
	NewInt64(name string, value int64, usage string) *int64

	// Float64 returns a pointer to a registered float64 by name as -name=1.0 a pointer to 1.0 is returned
	Float64(name string) *float64
	// NewFloat64 registers a new float64 flag by name and returns a pointer to the float64 storing the initial value
	NewFloat64(name string, value float64, usage string) *float64

	// String returns a pointer to stored string by -name=value
	String(name string) *string
	// NewString registers a new string flag by name and returns a pointer to the string storing the initial value
	NewString(name, value, usage string) *string

	// Bool returns a pointer to stored bool by -name=true
	Bool(name string) *bool
	// NewBool registers a new bool flag by name and returns a pointer to the bool storing the initial value
	NewBool(name string, value bool, usage string) *bool

	// Duration returns a pointer to stored time.Duration (unitless) by name like -minutes=10 (requires multiplication of * time.Minute to match memetics of "minutes" flag name and human interpretation of this)
	Duration(name string) *time.Duration
	// NewDuration registers a new time.Duration by name and returns a pointer to it storing the initial value
	NewDuration(name string, value time.Duration, usage string) *time.Duration

	// List returns a pointer to a []string containing strings
	List(name string) *[]string
	// NewList registers a new []string that can be assigned -name="ONE,TWO,THREE,FOUR"
	NewList(name string, value []string, usage string) *[]string

	// Map returns a pointer to a map[string]string containing strings
	Map(name string) *map[string]string
	// NewMap registers a new map[string]string that can be assigned -name="PROPERTY=VALUE,ANOTHER=VALUE"
	NewMap(name string, value map[string]string, usage string) *map[string]string
}

type Tree struct {
	figs map[string]interface{}
	mu   sync.RWMutex
}

// New will initialize the Tree package
// Usage:
//
//				When defining:
//				    cfg := configurable.New()
//			     cfg.NewInt("workers", 10, "number of workers")
//			     cfg.Parse()
//		      OR err := cfg.Load()
//		      OR err := cfg.ParseFile("path/to/file.json")
//	       THEN workers := *cfg.Int("workers") // workers is a regular int type but could be 0
func New() Fruit {
	cfg := &Tree{figs: make(map[string]interface{}), mu: sync.RWMutex{}}
	return cfg
}

// Parsing Configuration

// Parse runs flag.Parse() on the registered figs and uses defer recover() if an error occurs to return it
func (c *Tree) Parse() {
	flag.Parse()
	c.readEnv()
}

// ParseFile will check if filename is set and run loadFile on it.
func (c *Tree) ParseFile(filename string) error {
	flag.Parse()
	if filename != "" {
		return c.loadFile(filename)
	}
	c.readEnv()
	return nil
}

// Reload will readEnv on each flag in the configurable package
func (c *Tree) Reload() {
	c.readEnv()
}

// Load uses the EnvironmentKey and the DefaultJSONFile, DefaultYAMLFile, and DefaultINIFile to run ParseFile if it exists
func (c *Tree) Load() error {
	flag.Parse()

	files := []string{
		os.Getenv(EnvironmentKey),
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
			if err := c.loadFile(f); err != nil {
				return fmt.Errorf("failed to load %s: %w", f, err)
			}
			c.readEnv()
			return nil
		}
	}

	c.readEnv()
	return fmt.Errorf("no valid configuration file found among %v", files)
}

var ChainLoaded = false

// Usage prints a helpful table of figs in a human-readable format
func (c *Tree) Usage() string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "Usage of %s:\n", os.Args[0])
	flag.VisitAll(func(f *flag.Flag) {
		_, _ = fmt.Fprintf(&sb, "  -%s: %s (default: %s)\n", f.Name, f.Usage, f.DefValue)
	})
	return sb.String()
}

// INT

// NewInt registers a new integer 32 bit in the configurable package
func (c *Tree) NewInt(name string, value int, usage string) *int {
	c.mu.Lock()
	defer c.mu.Unlock()
	ptr := flag.Int(name, value, usage)
	c.figs[name] = ptr
	return ptr
}

// INT64

// Int returns an int32 from the configuration by name
func (c *Tree) Int(name string) *int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if ptr, ok := c.figs[name].(*int); ok {
		return ptr
	}
	return nil
}

// NewInt64 registers a new integer 64 bit in the configurable package
func (c *Tree) NewInt64(name string, value int64, usage string) *int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	var i = flag.Int64(name, value, usage)
	c.figs[name] = i
	return i
}

// FLOAT64

// Int64 returns the pointer to the int64 from the configuration by name
func (c *Tree) Int64(name string) *int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, _ := c.figs[name].(*int64)
	return val
}

// NewFloat64 registers a new float64 bit in the configurable package
func (c *Tree) NewFloat64(name string, value float64, usage string) *float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	var i = flag.Float64(name, value, usage)
	c.figs[name] = i
	return i
}

// DURATION

// Float64 returns the pointer to the float64 from the configuration by name
func (c *Tree) Float64(name string) *float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, _ := c.figs[name].(*float64)
	return val
}

// NewDuration registers a time.Duration in the configurable package
func (c *Tree) NewDuration(name string, value time.Duration, usage string) *time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	var i = flag.Duration(name, value, usage)
	c.figs[name] = i
	return i
}

// STRING

// Duration returns the pointer to the time.Duration from the configuration by name
func (c *Tree) Duration(name string) *time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, _ := c.figs[name].(*time.Duration)
	return val
}

// NewString registers a string in the configurable package
func (c *Tree) NewString(name string, value string, usage string) *string {
	c.mu.Lock()
	defer c.mu.Unlock()
	var s = flag.String(name, value, usage)
	c.figs[name] = s
	return s
}

// BOOL

// String returns the pointer to the string from the configuration by name
func (c *Tree) String(name string) *string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, _ := c.figs[name].(*string)
	return val
}

// NewBool registers a bool in the configurable package
func (c *Tree) NewBool(name string, value bool, usage string) *bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	var b = flag.Bool(name, value, usage)
	c.figs[name] = b
	return b
}

// LIST

// Bool returns the pointer to the string from the configuration by name
func (c *Tree) Bool(name string) *bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, _ := c.figs[name].(*bool)
	return val
}

// ListFlag stores values in a list type configurable
type ListFlag struct {
	values *[]string
}

// String returns the slice of strings using strings.Join
func (l *ListFlag) String() string {
	if l.values == nil {
		return ""
	}
	return strings.Join(*l.values, ",")
}

// Set unpacks a comma separated value argument and appends items to the list of []string
func (l *ListFlag) Set(value string) error {
	if l.values == nil {
		l.values = &[]string{}
	}
	items := strings.Split(value, ",")
	*l.values = append(*l.values, items...)
	return nil
}

// NewList registers a new list in the configurable package
func (c *Tree) NewList(name string, value []string, usage string) *[]string {
	c.mu.Lock()
	defer c.mu.Unlock()
	l := &ListFlag{values: &value}
	flag.Var(l, name, usage)
	c.figs[name] = l
	return l.values
}

// MAP

// List returns the pointer to the []string list in the configurable package
func (c *Tree) List(name string) *[]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if ptr, ok := c.figs[name].(*ListFlag); ok {
		return ptr.values
	}
	return nil
}

// MapFlag stores values in a map type configurable
type MapFlag struct {
	values *map[string]string
}

// String returns the map[string]string as string=string,string=string,...
func (m *MapFlag) String() string {
	if m.values == nil {
		return ""
	}
	var entries []string
	for k, v := range *m.values {
		entries = append(entries, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(entries, ",")
}

// Set accepts a value like KEY=VALUE,KEY=VALUE,KEY=VALUE to override map values
func (m *MapFlag) Set(value string) error {
	if m.values == nil {
		m.values = &map[string]string{}
	}
	pairs := strings.Split(value, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid map item: %s", pair)
		}
		(*m.values)[kv[0]] = kv[1]
	}
	return nil
}

// NewMap registers a new map[string]string in the configurable package
func (c *Tree) NewMap(name string, value map[string]string, usage string) *map[string]string {
	c.mu.Lock()
	defer c.mu.Unlock()
	m := &MapFlag{values: &value}
	flag.Var(m, name, usage)
	c.figs[name] = m
	return m.values
}

// Map returns the pointer to the map[string]string in the configurable package
func (c *Tree) Map(name string) *map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if ptr, ok := c.figs[name].(*MapFlag); ok {
		return ptr.values
	}
	return nil
}

// CONFIGURABLE INTERNAL FUNCTIONS

// loadFile will parse the filename for .yaml or .ini or .json and run the related loadJSON, loadYAML, or loadINI on it
func (c *Tree) loadFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return c.loadJSON(data)
	case ".yaml", ".yml":
		return c.loadYAML(data)
	case ".ini":
		return c.loadINI(data)
	default:
		return errors.New("unsupported file extension")
	}
}

// loadJSON parses the DefaultJSONFile or the value of the EnvironmentKey or ConfigFilePath into json.Unmarshal
func (c *Tree) loadJSON(data []byte) error {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	return c.setValuesFromMap(jsonData)
}

// loadYAML parses the DefaultYAMLFile or the value of the EnvironmentKey or ConfigFilePath into yaml.Unmarshal
func (c *Tree) loadYAML(data []byte) error {
	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return err
	}
	return c.setValuesFromMap(yamlData)
}

// loadINI parses the DefaultINIFile or the value of the EnvironmentKey or ConfigFilePath into ini.Load()
func (c *Tree) loadINI(data []byte) error {
	cfg, err := ini.Load(data)
	if err != nil {
		return err
	}
	iniData := make(map[string]interface{})
	for key := range c.figs {
		if val := cfg.Section("").Key(key).String(); val != "" {
			iniData[key] = val
		}
	}
	return c.setValuesFromMap(iniData)
}

// setValuesFromMap uses the data map to store the configurable figs
func (c *Tree) setValuesFromMap(data map[string]interface{}) error {
	for key, value := range data {
		if flagVal, exists := c.figs[key]; exists {
			if err := c.setValue(flagVal, value); err != nil {
				return fmt.Errorf("error setting key %s: %w", key, err)
			}
		}
	}
	return nil
}

// setValue is a generic method on the Tree interface to set a value to a flag type
func (c *Tree) setValue(flagVal interface{}, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch ptr := flagVal.(type) {
	case *int:
		intVal, err := toInt(value)
		if err != nil {
			return err
		}
		*ptr = intVal
	case *int64:
		int64Val, err := toInt64(value)
		if err != nil {
			return err
		}
		*ptr = int64Val
	case *float64:
		floatVal, err := toFloat64(value)
		if err != nil {
			return err
		}
		*ptr = floatVal
	case *string:
		strVal, err := toString(value)
		if err != nil {
			return err
		}
		*ptr = strVal
	case *bool:
		boolVal, err := toBool(value)
		if err != nil {
			return err
		}
		*ptr = boolVal
	case *time.Duration:
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
		*ptr.values = append(*ptr.values, listVal...)
	case *MapFlag:
		mapVal, err := toStringMap(value)
		if err != nil {
			return err
		}
		for k, v := range mapVal {
			(*ptr.values)[k] = v
		}
	default:
		return fmt.Errorf("unsupported flag type for key %v", ptr)
	}
	return nil
}

// ReadEnv uses the os.LookupEnv to store the value for each property key
func (c *Tree) readEnv() {
	for name, _ := range c.figs {
		c.checkAndSetFromEnv(name)
	}
}

// checkAndSetFromEnv uses os.LookupEnv and assigns it to the figs name value
func (c *Tree) checkAndSetFromEnv(name string) {
	if val, exists := os.LookupEnv(name); exists {
		if flagVal, exists := c.figs[name]; exists {
			_ = c.setValue(flagVal, val)
		}
	}
}

// INTERFACE TYPE CONVERSIONS

// toInt returns an interface{} as an int or returns an error
func toInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %v to int", value)
	}
}

// toInt64 returns an interface{} as an int64 or returns an error
func toInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %v to int64", value)
	}
}

// toFloat64 returns an interface{} as an float64 or returns an error
func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %v to float64", value)
	}
}

// toString returns an interface{} as a string or returns an error
func toString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", fmt.Errorf("cannot convert %v to string", value)
	}
}

// toBool returns an interface{} as a bool or returns an error
func toBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("cannot convert %v to bool", value)
	}
}

// toStringSlice returns an interface{} as a []string{} or returns an error
func toStringSlice(value interface{}) ([]string, error) {
	switch v := value.(type) {
	case []interface{}:
		var result []string
		for _, item := range v {
			str, err := toString(item)
			if err != nil {
				return nil, err
			}
			result = append(result, str)
		}
		return result, nil
	case string:
		if v == "" {
			return []string{}, nil
		}
		return strings.Split(v, ","), nil
	default:
		return nil, fmt.Errorf("cannot convert %v to []string", value)
	}
}

// toStringMap returns an interface{} as a map[string]string or returns an error
func toStringMap(value interface{}) (map[string]string, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		result := make(map[string]string)
		for key, val := range v {
			strVal, err := toString(val)
			if err != nil {
				return nil, err
			}
			result[key] = strVal
		}
		return result, nil
	case string:
		if v == "" {
			return map[string]string{}, nil
		}
		pairs := strings.Split(v, ",")
		result := make(map[string]string)
		for _, pair := range pairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid map item: %s", pair)
			}
			result[kv[0]] = kv[1]
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %v to map[string]string", value)
	}
}
