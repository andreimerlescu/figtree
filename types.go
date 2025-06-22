package figtree

import (
	"flag"
	"sync"
	"sync/atomic"
	"time"
)

type Withables interface {
	// WithCallback registers a new CallbackWhen with a CallbackFunc on a figFruit on the figTree by its name
	WithCallback(name string, whenCallback CallbackWhen, runThis CallbackFunc) Plant
	// WithAlias registers a short form of the name of a figFruit on the figTree
	WithAlias(name, alias string) Plant
	// WithRule attaches a RuleKind to a figFruit
	WithRule(name string, rule RuleKind) Plant
	// WithTreeRule assigns a global rule on the Tree
	WithTreeRule(rule RuleKind) Plant
	// WithValidator binds a figValidatorFunc to a figFruit that returns Plant
	WithValidator(name string, validator func(interface{}) error) Plant
}

type Savable interface {
	// SaveTo will store the Tree in a path file
	SaveTo(path string) error
}

type Readable interface {
	// ReadFrom will attempt to load the file into the Tree
	ReadFrom(path string) error
}

type Parsable interface {
	// Parse can panic but interprets command line arguments defined with single dashes -example value -another sample
	Parse() error
	// ParseFile can panic but also can throw an error because it will attempt to load either JSON, YAML or INI file passed into it
	ParseFile(filename string) error
}

type Mutable interface {
	// Mutations receives Mutation data on a receiver channel
	Mutations() <-chan Mutation
	// MutagenesisOfFig will look up a Fruit by name and return the Metagenesis of it
	MutagenesisOfFig(name string) Mutagenesis
	// MutagenesisOf takes anything and returns the Mutagenesis of it
	MutagenesisOf(what interface{}) Mutagenesis
}

type Loadable interface {
	// Load can panic but also can throw an error but will use the Environment Variable values if they are "EXAMPLE=value" or "ANOTHER=sample"
	Load() error
	// LoadFile accepts a path to a JSON, YAML or INI file to set values
	LoadFile(path string) error
	// Reload will refresh stored values of properties with their new Environment Variable values
	Reload() error
}

type Divine interface {
	// Recall allows you to unlock the figTree from changes and resume tracking
	Recall()
	// Curse allows you to lock the figTree from changes and stop tracking
	Curse()
}

type Intable interface {
	// Int returns a pointer to a registered int32 by name as -name=1 a pointer to 1 is returned
	Int(name string) *int
	// NewInt registers a new int32 flag by name and returns a pointer to the int32 storing the initial value
	NewInt(name string, value int, usage string) Plant
	// StoreInt replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreInt(name string, value int) Plant
}

type Intable64 interface {
	// Int64 returns a pointer to a registered int64 by name as -name=1 a pointer to 1 is returned
	Int64(name string) *int64
	// NewInt64 registers a new int32 flag by name and returns a pointer to the int64 storing the initial value
	NewInt64(name string, value int64, usage string) Plant
	// StoreInt64 replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreInt64(name string, value int64) Plant
}

type Floatable interface {
	// Float64 returns a pointer to a registered float64 by name as -name=1.0 a pointer to 1.0 is returned
	Float64(name string) *float64
	// NewFloat64 registers a new float64 flag by name and returns a pointer to the float64 storing the initial value
	NewFloat64(name string, value float64, usage string) Plant
	// StoreFloat64 replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreFloat64(name string, value float64) Plant
}

type String interface {
	// String returns a pointer to stored string by -name=value
	String(name string) *string
	// NewString registers a new string flag by name and returns a pointer to the string storing the initial value
	NewString(name, value, usage string) Plant
	// StoreString replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreString(name, value string) Plant
}

type Flaggable interface {
	// Bool returns a pointer to stored bool by -name=true
	Bool(name string) *bool
	// NewBool registers a new bool flag by name and returns a pointer to the bool storing the initial value
	NewBool(name string, value bool, usage string) Plant
	// StoreBool replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreBool(name string, value bool) Plant
}

type Durable interface {
	// Duration returns a pointer to stored time.Duration (unitless) by name like -minutes=10 (requires multiplication of * time.Minute to match memetics of "minutes" flag name and human interpretation of this)
	Duration(name string) *time.Duration
	// NewDuration registers a new time.Duration by name and returns a pointer to it storing the initial value
	NewDuration(name string, value time.Duration, usage string) Plant
	// StoreDuration replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreDuration(name string, value time.Duration) Plant
	// UnitDuration returns a pointer to stored time.Duration (-name=10 w/ units as time.Minute == 10 minutes time.Duration)
	UnitDuration(name string) *time.Duration
	// NewUnitDuration registers a new time.Duration by name and returns a pointer to it storing the initial value
	NewUnitDuration(name string, value, units time.Duration, usage string) Plant
	// StoreUnitDuration replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreUnitDuration(name string, value, units time.Duration) Plant
}

type Listable interface {
	// List returns a pointer to a []string containing strings
	List(name string) *[]string
	// NewList registers a new []string that can be assigned -name="ONE,TWO,THREE,FOUR"
	NewList(name string, value []string, usage string) Plant
	// StoreList replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreList(name string, value []string) Plant
}

type Mappable interface {
	// Map returns a pointer to a map[string]string containing strings
	Map(name string) *map[string]string
	// NewMap registers a new map[string]string that can be assigned -name="PROPERTY=VALUE,ANOTHER=VALUE"
	NewMap(name string, value map[string]string, usage string) Plant
	// MapKeys returns the keys of the map[string]string as a []string
	MapKeys(name string) []string
	// StoreMap replaces name with value and can issue a Mutation when receiving on Mutations()
	StoreMap(name string, value map[string]string) Plant
}

type CoreAbilities interface {
	Withables
	Savable
	Readable
	Parsable
	Mutable
	Loadable
	Divine
}

type CoreMutations interface {
	Intable
	Intable64
	Floatable
	String
	Flaggable
	Durable
	Listable
	Mappable
}

type Core interface {
	// FigFlesh returns a figFruit from the figTree by its name
	FigFlesh(name string) Flesh

	// ErrorFor returns an error attached to a named figFruit
	ErrorFor(name string) error

	// Usage displays the helpful menu of figs registered using -h or -help
	Usage()
}

// Plant defines the interface for configuration management.
type Plant interface {
	Core
	CoreAbilities
	CoreMutations
}

// figTree stores figs that are defined by their name and figFruit as well as a mutations channel and tracking bool for Options.Tracking
type figTree struct {
	ConfigFilePath string
	GlobalRules    []RuleKind
	harvest        int
	pollinate      bool
	figs           map[string]*figFruit
	values         *sync.Map
	withered       map[string]witheredFig
	aliases        map[string]string
	sourceLocker   sync.RWMutex
	mu             sync.RWMutex
	tracking       bool
	problems       []error
	mutationsCh    chan Mutation
	flagSet        *flag.FlagSet
	filterTests    bool
	angel          *atomic.Bool
	ignoreEnv      bool
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

	// IgnoreEnvironment is a part of free will, it lets us disregard our environment (ENV vars)
	IgnoreEnvironment bool
}

type FigValidatorFunc func(interface{}) error

type witheredFig struct {
	Error       error
	Mutagenesis Mutagenesis
	Value       Value
	name        string
}

type figFruit struct {
	Validators  []FigValidatorFunc
	Mutations   []Mutation
	Callbacks   []Callback
	Rules       []RuleKind
	Locker      *sync.RWMutex
	Error       error
	Mutagenesis Mutagenesis
	name        string
	usage       string
}

type figFlesh struct {
	Flesh interface{}
}

type Flesh interface {
	Is(mutagenesis Mutagenesis) bool
	AsIs() interface{}
	IsString() bool
	IsInt() bool
	IsInt64() bool
	IsBool() bool
	IsFloat64() bool
	IsDuration() bool
	IsUnitDuration() bool
	IsList() bool
	IsMap() bool

	ToString() string
	ToInt() int
	ToInt64() int64
	ToBool() bool
	ToFloat64() float64
	ToDuration() time.Duration
	ToUnitDuration() time.Duration
	ToList() []string
	ToMap() map[string]string
}

type Callback struct {
	CallbackWhen CallbackWhen
	CallbackFunc CallbackFunc
}

type CallbackWhen string

type CallbackFunc func(interface{}) error

type Mutation struct {
	Property    string
	Mutagenesis string
	Way         string
	Old         interface{}
	New         interface{}
	When        time.Time
	Error       error
}

var ListSeparator = ","
var MapSeparator = ","
var MapKeySeparator = "="
