package figtree

import (
	"flag"
	"sync"
	"sync/atomic"
	"time"
)

// Fruit defines the interface for configuration management.
type Fruit interface {
	// WithValidator binds a ValidatorFunc to a Fig that returns Fruit
	WithValidator(name string, validator func(interface{}) error) Fruit
	// Parse can panic but interprets command line arguments defined with single dashes -example value -another sample
	Parse() error
	// ParseFile can panic but also can throw an error because it will attempt to load either JSON, YAML or INI file passed into it
	ParseFile(filename string) error

	// ErrorFor returns an error attached to a named Fig
	ErrorFor(name string) error

	// Recall allows you to unlock the Tree from changes and resume tracking
	Recall()
	// Curse allows you to lock the Tree from changes and stop tracking
	Curse()
	// Mutations receives Mutation data on a receiver channel
	Mutations() <-chan Mutation
	// Resurrect takes a nil Fig in the Tree.figs map and reloads it from ENV or the config file if available
	Resurrect(name string)

	// Load can panic but also can throw an error but will use the Environment Variable values if they are "EXAMPLE=value" or "ANOTHER=sample"
	Load() error
	// LoadFile accepts a path to a JSON, YAML or INI file to set values
	LoadFile(path string) error
	// Reload will refresh stored values of properties with their new Environment Variable values
	Reload() error

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

	// UnitDuration returns a pointer to stored time.Duration (-name=10 w/ units as time.Minute == 10 minutes time.Duration)
	UnitDuration(name string) *time.Duration
	// NewUnitDuration registers a new time.Duration by name and returns a pointer to it storing the initial value
	NewUnitDuration(name string, value, units time.Duration, usage string) *time.Duration

	// List returns a pointer to a []string containing strings
	List(name string) *[]string
	// NewList registers a new []string that can be assigned -name="ONE,TWO,THREE,FOUR"
	NewList(name string, value []string, usage string) *[]string

	// Map returns a pointer to a map[string]string containing strings
	Map(name string) *map[string]string
	// NewMap registers a new map[string]string that can be assigned -name="PROPERTY=VALUE,ANOTHER=VALUE"
	NewMap(name string, value map[string]string, usage string) *map[string]string

	// StoreInt replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreInt(name string, value int) Fruit
	// StoreInt64 replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreInt64(name string, value int64) Fruit
	// StoreFloat64 replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreFloat64(name string, value float64) Fruit
	// StoreString replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreString(name, value string) Fruit
	// StoreBool replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreBool(name string, value bool) Fruit
	// StoreDuration replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreDuration(name string, value time.Duration) Fruit
	// StoreUnitDuration replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreUnitDuration(name string, value, units time.Duration) Fruit
	// StoreList replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreList(name string, value []string) Fruit
	// StoreMap replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreMap(name string, value map[string]string) Fruit
}

// Tree stores figs that are defined by their name and Fig as well as a mutations channel and tracking bool for Options.Tracking
type Tree struct {
	ConfigFilePath string
	harvest        int
	pollinate      bool
	figs           map[string]*Fig
	withered       map[string]Fig
	mu             sync.RWMutex
	tracking       bool
	mutationsCh    chan Mutation
	flagSet        *flag.FlagSet
	filterTests    bool
	angel          *atomic.Bool
}

// Mutagenesis stores the type as a string like String, Bool, Float, etc to represent a supported Type
type Mutagenesis string

// Options allow you enable mutation tracking on your figs.Grow
type Options struct {
	ConfigFile string
	// Tracking creates a buffered channel that allows you to select { case mutation, ok := <-figs.Mutations(): }
	Tracking bool

	// Germinate enables the option to filter os.Args that begin with -test. prefix
	Germinate bool

	// Harvest allows you to set the buffer size of the Mutations channel
	Harvest int

	// Pollinate will enable Getters to lookup the environment for changes on every read
	Pollinate bool
}

type ValidatorFunc func(interface{}) error

type Fig struct {
	Flesh       interface{}
	Mutagenesis Mutagenesis
	Validator   ValidatorFunc
	Error       error
	Mutations   []Mutation
}

type Mutation struct {
	Property string
	Kind     string
	Way      string
	Old      interface{}
	New      interface{}
	When     time.Time
}
