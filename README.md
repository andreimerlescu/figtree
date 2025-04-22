# Fig Tree

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Fig Tree is a command line utility configuration manager that you can refer to as `figs`, as you con<b>FIG</b>ure your
application's runtime.

![Figtree](/figtree.jpg "Figtree by xAI Grok 3")

## Installation

To use `figtree` in your project, `go get` it...

```shell
go get -u github.com/andreimerlescu/figtree/v2
```

## Usage

To use **figs** package in your Go code, you need to import it:

```go
import "github.com/andreimerlescu/figtree/v2"
```

### Using Fig Tree

Creating a new Fig Tree can be done with the following strategies. Your choice. 

| Method                                          | Usage                                 |
|-------------------------------------------------|---------------------------------------|
| `figtree.New()`                                 | Does not perform `Mutation` tracking. |
| `figtree.Grow()`                                | Provides `Mutation` tracking.         |
| `figtree.With(figtree.Options{Tracking: true})` | Provides `Mutation` tracking.         |

When using `figtree.Options`, you can enable: 

| Option              | What It Does                                                                                  | 
|---------------------|-----------------------------------------------------------------------------------------------|
| `Pollinate`         | Read `os.Getenv(key)` when a Getter on a Mutagenesis is called                                |
| `Harvest`           | Slice length of `Mutation` for `Pollinate`                                                    |
| `IgnoreEnvironment` | Ignore `os.Getenv()` and use `os.Clearenv()` inside `With(opts Options)`                      |
| `Germinate`         | Ignore command line flags that begin with `-test.`                                            |
| `Tracking`          | Sends `Mutation` into a receiver channel on `figs.Mutations()` whenever a `Fig` value changes |
| `ConfigFile`        | Path to your `config.yaml` or `config.ini` or `config.json` file                              |

Configurable properties have whats called metagenesis to them, which are types, like `String`, `Bool`, `Float64`, etc.

| Mutagenesis     | Getter                                | Setter                                  | Fruit Getter            |
|-----------------|---------------------------------------|-----------------------------------------|-------------------------|
| `tString`       | `keyValue := *figs.String(key)`       | `figs.Store(tString, key, value)`       | `figs := figs.Fig(key)` |
| `tInt`          | `keyValue := *figs.Int(key)`          | `figs.Store(tInt, key, value)`          | `figs := figs.Fig(key)` |
| `tInt64`        | `keyValue := *figs.Int64(key)`        | `figs.Store(tInt64, key, value)`        | `figs := figs.Fig(key)` |
| `tFloat64`      | `keyValue := *figs.Float64(key)`      | `figs.Store(tFloat64, key, value)`      | `figs := figs.Fig(key)` |
| `tDuration`     | `keyValue := *figs.Duration(key)`     | `figs.Store(tDuration, key, value)`     | `figs := figs.Fig(key)` |
| `tUnitDuration` | `keyValue := *figs.UnitDuration(key)` | `figs.Store(tUnitDuration, key, value)` | `figs := figs.Fig(key)` |
| `tList`         | `keyValue := *figs.List(key)`         | `figs.Store(tList, key, value)`         | `figs := figs.Fig(key)` |
| `tMap`          | `keyValue := *figs.Map(key)`          | `figs.Store(tMap, key, value)`          | `figs := figs.Fig(key)` |

New properties can be registered before calling Parse() using a metagenesis pattern of `figs.New<Metagenesis>()`, like
`figs.NewString()` or `figs.NewFloat64()`, etc. 

Only one validator per property is permitted, and additional WithValidator() calls with duplicate name entries will 
record an error in the `Fig.Error` property of the property's "fruit, aka `*Fig{}`".

Figtree keeps a withered copy of `figs.New<Metagenesis>()` property declarations and has an `Options{Tracking: true}` 
argument that can be passed into `figs.New()` that enables the `figs.Mutations()` receiver channel to receive anytime
a property value changes, a new `figtree.Mutation`.

Figtree includes 36 different built-in `figs.WithValidator(name, figtree.Assure<Rule>[()])` that can
validate your various Mutageneses without needing to write every validation yourself. For larger or
custom validations, the 2nd argument requires a `func (interface{}) error` signature in order use.

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/andreimerlescu/figtree/v2"
)

func main() {
	figs := figtree.Grow()
	figs.NewUnitDuration("minutes", 1, time.Minute, "minutes for timer")
    figs.WithValidator("minutes", figtree.AssureDurationLessThan(time.Hour))
	err := figs.Parse()
	if err != nil {
		log.Fatal(err)
	}
    log.Println(*figs.UnitDuration("minutes"))
}
```

### Available Rules

| RuleKind                        | Notes                                                             |
|---------------------------------|-------------------------------------------------------------------|
| `RuleUndefined`                 | is default and does no action                                     |
| `RulePreventChange`             | blocks Mutagensis Store methods                                   | 
| `RulePanicOnChange`             | will throw a panic on the Mutagenesis Store methods               | 
| `RuleNoValidations`             | will skip over all WithValidator assignments                      | 
| `RuleNoCallbacks`               | will skip over all WithCallback assignments                       | 
| `RuleCondemnedFromResurrection` | will panic if there is an attempt to resurrect a condemned fig    |
| `RuleNoMaps`                    | blocks NewMap, StoreMap, and Map from being called on the Tree    | 
| `RuleNoLists`                   | blocks NewList, StoreList, and List from being called on the Tree | 
| `RuleNoFlags`                   | disables the flag package from the Tree                           |
| `RuleNoEnv`                     | skips over all os.Getenv related logic                            |


#### Global Rules

```go
package main

import (
    "github.com/andreimerlescu/figtree/v2"
    "log"
)

func main() {
	figs := figtree.Grow()
	figs.WithTreeRule(figtree.RuleNoFlags)
	figs.NewString("name", "", "your name")
    figs.WithValidator("name", figtree.AssureStringNotEmpty) // validate no empty strings
	figs.WithRule("name", figtree.RuleNoValidations) // turn off validations
	err := figs.Parse() // no error
	if err != nil {
		log.Println(err)
	}
	figs.StoreString("name", "Yeshua")
	log.Printf("Hello %s", *figs.String("name"))
}
```

#### Property Rules

### Available Validators

| Mutagenesis | `figtree.ValidatorFunc`   | Notes                                                                            |
|-------------|---------------------------|----------------------------------------------------------------------------------|
| tString     | AssureStringLength        | Ensures a string is a specific length.                                           |
| tString     | AssureStringNotLength     | Ensures a string is not a specific length.                                       |
| tString     | AssureStringSubstring     | Ensures a string contains a specific substring (case-sensitive).                 |
| tString     | AssureStringNotEmpty      | Ensures a string is not empty.                                                   |
| tString     | AssureStringContains      | Ensures a string contains a specific substring.                                  |
| tString     | AssureStringNotContains   | Ensures a string does not contains a specific substring.                         |
| tString     | AssureStringHasPrefix     | Ensures a string has a prefix.                                                   |
| tString     | AssureStringHasSuffix     | Ensures a string has a suffix.                                                   |
| tString     | AssureStringNoPrefix      | Ensures a string does not have a prefix.                                         |
| tString     | AssureStringNoSuffix      | Ensures a string does not have a suffix.                                         |
| tString     | AssureStringNoPrefixes    | Ensures a string does not have a prefixes.                                       |
| tString     | AssureStringNoSuffixes    | Ensures a string does not have a suffixes.                                       |
| tBool       | AssureBoolTrue            | Ensures a boolean value is true.                                                 |
| tBool       | AssureBoolFalse           | Ensures a boolean value is false.                                                |
| tInt        | AssurePositiveInt         | Ensures an integer is positive (greater than zero).                              |
| tInt        | AssureNegativeInt         | Ensures an integer is negative (less than zero).                                 |
| tInt        | AssureIntGreaterThan      | Ensures an integer is greater than a specified value (exclusive).                |
| tInt        | AssureIntLessThan         | Ensures an integer is less than a specified value (exclusive).                   |
| tInt        | AssureIntInRange          | Ensures an integer is within a specified range (inclusive).                      |
| tInt64      | AssureInt64GreaterThan    | Ensures an int64 is greater than a specified value (exclusive).                  |
| tInt64      | AssureInt64LessThan       | Ensures an int64 is less than a specified value (exclusive).                     |
| tInt64      | AssurePositiveInt64       | Ensures an int64 is positive (greater than zero).                                |
| tInt64      | AssureInt64InRange        | Ensures an int64 is within a specified range (inclusive).                        |
| tFloat64    | AssureFloat64Positive     | Ensures a float64 is positive (greater than zero).                               |
| tFloat64    | AssureFloat64InRange      | Ensures a float64 is within a specified range (inclusive).                       |
| tFloat64    | AssureFloat64GreaterThan  | Ensures a float64 is greater than a specified value (exclusive).                 |
| tFloat64    | AssureFloat64LessThan     | Ensures a float64 is less than a specified value (exclusive).                    |
| tFloat64    | AssureFloat64NotNaN       | Ensures a float64 is not NaN.                                                    |
| tDuration   | AssureDurationGreaterThan | Ensures a time.Duration is greater than a specified value (exclusive).           |
| tDuration   | AssureDurationLessThan    | Ensures a time.Duration is less than a specified value (exclusive).              |
| tDuration   | AssureDurationPositive    | Ensures a time.Duration is positive (greater than zero).                         |
| tDuration   | AssureDurationMax         | Ensures a time.Duration does not exceed a maximum value.                         |
| tDuration   | AssureDurationMin         | Ensures a time.Duration is at least a minimum value.                             |
| tList       | AssureListNotEmpty        | Ensures a list (*ListFlag, *[]string, or []string) is not empty.                 |
| tList       | AssureListMinLength       | Ensures a list has at least a minimum number of elements.                        |
| tList       | AssureListContains        | Ensures a list contains a specific string value.                                 |
| tList       | AssureListNotContains     | Ensures a list does not contain a specific string value.                         |
| tList       | AssureListContainsKey     | Ensures a list contains a specific string.                                       |
| tList       | AssureListLength          | Ensures a list has exactly the specified length.                                 |
| tList       | AssureListNotLength       | Ensures a list is not the specified length.                                      |
| tMap        | AssureMapNotEmpty         | Ensures a map (*MapFlag, *map[string]string, or map[string]string) is not empty. |
| tMap        | AssureMapHasKey           | Ensures a map contains a specific key.                                           |
| tMap        | AssureMapValueMatches     | Ensures a map has a specific key with a matching value.                          |
| tMap        | AssureMapHasKeys          | Ensures a map contains all specified keys.                                       |
| tMap        | AssureMapLength           | Ensures a map has exactly the specified length.                                  |
| tMap        | AssureMapNotLength        | Ensures a map not the specified length.                                          |


### Callbacks

The **Go** way of doing callbacks is to rely on the `Option.Tracking` set to `true` and receiving on the `figs.Mutations()`
receiver channel. However, if you want to use Callbacks, you can register them various different ways. 

```go
func CheckAvailability(domain string) error {
    // todo implement something that checks the availability of the domain
    return nil
}
const kDomain string = "domain"
figs := figtree.With(Options{Tracking: true, Harvest: 1776, Pollinate: true})
figs.NewString(kDomain, "", "Domain name")
figs.WithValidator(kDomain, figtree.AssureStringLengthGreaterThan(3))
figs.WithValidator(kDomain, figtree.AssureStringHasPrefix("https://"))
figs.WithCallback(kDomain, figree.CallbackAfterVerify, func(value interface{}) error {
    var s string
    switch v := value.(type) {
    case *string:
        s = *v
    case string:
        s = v
    }
    // try connecting to the domain now
    return CheckAvailability(s)
})
figs.WithCallback(kDomain, figree.CallbackAfterRead, func(value interface{}) error {
    // every time *figs.String(kDomain) is called, run this
	var s string
	switch v := value.(type) {
    case *string:
        s = *v
    case string:
        s = v
    }
    log.Printf("CallbackAfterRead invoked for %s", s)
    return nil
})
figs.WithCallback(kDomain, figtree.CallbackAfterChange, func(value interface{}) error{
    var s string
    switch v := value.(type) {
    case *string:
        s = *v
    case string:
        s = v
    }
    log.Printf("CallbackAfterChange invoked for %s", s)
    return nil
})
err := figs.Parse() // CallbackAfterVerify invoked
if err != nil {
    log.Fatal(err)
}

domain := *figs.String(kDomain) // CallbackAfterRead invoked
figs.Store(kDomain, "https://newdomain.com") // CallbackAfterChange invoked
err := figs.HasError(kDomain)
if err != nil {
    log.Fatal(err)
}
newDomain := *figs.String(kDomain) // CallbackAfterRead invoked
log.Printf("domain = %s ; newDomain = %s", domain, newDomain)
```

The second argument to the `WithCallback` func is a `ChangeAfter` type. 

| Option                 | When It's Triggered                                                                |
|------------------------|------------------------------------------------------------------------------------|
| `CallbackAfterVerify`  | Called on `.Parse()`, `.ParseFile()`, `Load()`, or `LoadFile()`                    | 
| `CallbackAfterRead`    | Called on Mutagenesis Getters like `figs.String(key)` or `figs.<Mutagenesis>(key)` |
| `CallbackAfterChanged` | Called on `.Store(Mutagenesis, key, value)` and `.Resurrect(key)`                  |

When using callbacks, you will want to make sure that you're keeping on top of what you're assigning to each `Fig`.

With callbacks, you can really slow the performance down of `figtree`, but when used sparingly, its extremely powerful.

At the end of the day, you'll know what's best to use. I build what I build because its the best that I use.

### Complex Example Usage

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"strconv"
	"time"

	"github.com/andreimerlescu/figtree/v2"
    check "github.com/andreimerlescu/checkfs"
    "github.com/andreimerlescu/checkfs/file"
)

func main() {
	// Create a context that can be canceled on SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Plant the fig tree 
	var figs figtree.Fruit

	// change internals of figtree to use new file for default .Load()
	figtree.ConfigFilePath = filepath.Join("/opt", "app", "default.config.json")

	// You can build your options before passing them in
	options := figtree.Options{ // init with options baked in
		Tracking:  true, // setup the .Mutations() channel
		Germinate: true, // ignore -test arguments in CLI
	}
	if _, ok := os.LookupEnv("NOENV"); ok { // NOENV=1 go run . 
		options.IgnoreEnvironment = true
	}
	if path, ok := os.LookupEnv("MYCONFIGFILE"); ok { // MYCONFIGFILE=/opt/app/config.yaml go run .
		options.ConfigFile = path
	}
	if harvestStr, ok := os.LookupEnv("HARVEST"); ok { // HARVEST=1776 go run .
		harvest, harvestErr := strconv.ParseInt(harvestStr, 10, 64)
		if harvestErr != nil {
			options.Harvest = harvest
		}
	}
	if val, ok := os.LookupEnv("POLLINATE_FIGTREE"); ok {
		b, bErr := strconv.ParseBool(val)
		if bErr != nil {
			options.Pollinate = b
		}
	}

	// create a new figtree
	figs = figtree.With(options)

	// Define all configuration types with initial values and validators
	
    // arg -workers int
    
	figs.NewInt("workers", 10, "Number of workers")
   
    // You can work directly with the Fig of "workers"
	workersFig := figs.Fig("workers")
	// there is no way to send your own copy of *figtree.Fig back into an issued figtree.Grow()
	// but you can access the underlying *figtree.Fig if you need to access the figtree.ValidatorFunc
	// or the figtree.Callback type
	for _, callback := range workersFig.Callbacks {
		if callback.CallbackAfter == figtree.CallbackAfterRead {
			callbackErr := callback.CallbackFunc(workersFig.Flesh)
			if callbackErr != nil {
				log.Println(callbackErr)
			}
		}
	}
	// but this for loop will automatically be called on figtree.Parse() or figtree.Load() depending
	// on if you're using `figs := figtree.With(figtree.Options{ConfigFile: "/opt/app/config.yaml", IgnoreEnvironment: true})`
	// or if you're only using `figs := figtree.Grow()` for mutation tracking only

	// this validator allows you to define n-workers between 1 and number of CPUs 
	figs.WithValidator("workers", figtree.AssureIntInRange(1, runtime.GOMAXPROCS(0)))

	// arg -maxRetries int
    
	figs.NewInt64("maxRetries", 5, "Maximum retry attempts")
	figs.WithValidator("maxRetries", figtree.AssureInt64Positive)
	figs.WithCallback("maxRetries", figtree.CallbackAfterChange, func(value interface{}) error {
		log.Printf("fig maxRetries changed to %q\n", value)
		return nil
	})

	// arg -threshold float64 
    
	figs.NewFloat64("threshold", 0.75, "Threshold value")
	figs.WithValidator("threshold", figtree.AssureFloat64InRange(0.0, 1.0))
	figs.WithCallback("threshold", figtree.CallbackAfterChange, func(value interface{}) error {
		switch v := value.(type) {
		case *float64:
			log.Printf("fig threshold changed to %f", v)
		case float64:
			log.Printf("fig threshold changed to %f", v)
		}
		return nil
	})
    
    // arg --endpoint "value"

	figs.NewString("endpoint", "http://example.com", "API endpoint")
	figs.WithValidator("endpoint", figtree.AssureStringHasPrefix("http"))

    // arg -debug <true|false> 
    
	figs.NewBool("debug", false, "Enable debug mode")
	figs.WithCallback("debug", figtree.CallbackAfterVerify, func(v interface{}) error {
		log.Println("ACTIVATING DEBUG MODE!")
		return nil
	})

    // arg -timeout defaults to 30s but -timeout <int> becomes <int> seconds
	
	figs.NewUnitDuration("timeout", 30*time.Second, time.Second, "Request timeout")
	figs.WithValidator("timeout", figtree.AssureDurationMin(10*time.Second))
	figs.WithValidator("timeout", figtree.AssureDurationMax(time.Hour*12))

    // arg -interval <int> becomes <int> minutes
	
	figs.NewUnitDuration("interval", 1, time.Minute, "Polling interval in minutes")
	figs.WithValidator("interval", figtree.AssureDurationGreaterThan(30*time.Second))
    
    // arg -servers "ONE,TWO,THREE" becomes []string{"ONE", "TWO", "THREE"}
    
	figs.NewList("servers", []string{"server1", "server2"}, "List of servers")
	figs.WithValidator("servers", figtree.AssureListNotEmpty)
	figs.WithCallback("servers", figtree.CallbackAfterChange, func(value interface{}) error {
		var val []string
		switch v := value.(type) {
		case *figtree.ListFlag:
			val = make([]string, len(*v.values))
			copy(val, *val.values)
		case *[]string:
			copy(val, *v)
		case []string:
			copy(val, v)
		}
		log.Printf("-servers value changed! new value = %s", strings.Join(val, ", "))
		return nil
	})
    
    // arg -metadata "KEY=VALUE,KEY=VALUE"

	figs.NewMap("metadata", map[string]string{"env": "prod", "version": "1.0"}, "Metadata key-value pairs")
	figs.WithValidator("metadata", figtree.AssureMapHasKeys([]string{"env", "version"}))

	// Attempt to load from a config file, falling back to defaults
	configFile := filepath.Join(".", "config.yaml")
    if err := check.File(figtree.ConfigFilePath, file.Options{Exists: true}); err != nil {
        if err = figs.Load(); err != nil {
            log.Fatal(err)
        }
    } else if err := check.File(configFile, file.Options{Exists: true}); err != nil {
		if err = figs.ParseFile(configFile); err != nil {
			log.Fatal(err)
        }
	} else {
		if err := figs.Parse(); err != nil {
			log.Fatal(err)
		}
    }
	// Note: LoadFile tries the specified file, then env vars; Parse handles flags/env if file fails

	// Demonstrate Resurrect by accessing an undefined key
	undefined := figs.String("undefined")
	log.Printf("Resurrected undefined key: %s", *undefined)
	// Note: Resurrect creates a new string entry if undefined, checking env and files first

	// Print initial configuration using Usage
	log.Println("Initial configuration:")
	log.Println(figs.Usage())
	// Note: Usage displays all registered flags in a human-readable format

	// Simulate periodic access in a goroutine to check for mutations
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("Worker shutting down due to context cancellation")
				return
			case <-ticker.C:
				// Access all config values, triggering mutation checks
				log.Printf("Workers: %d, MaxRetries: %d, Threshold: %.2f, Endpoint: %s, Debug: %t, Timeout: %s, Interval: %s, Servers: %v, Metadata: %v",
					*figs.Int("workers"),
					*figs.Int64("maxRetries"),
					*figs.Float64("threshold"),
					*figs.String("endpoint"),
					*figs.Bool("debug"),
					*figs.Duration("timeout"),
					*figs.UnitDuration("interval"),
					*figs.List("servers"),
					*figs.Map("metadata"),
				)
				// Note: Getters check env against withered values, sending mutations if tracking is on
			}
		}
	}()

	// Demonstrate Curse and Recall by toggling tracking
	log.Println("Cursing the tree (disabling tracking)...")
	figs.Curse()
	time.Sleep(2 * time.Second) // Let some ticks pass
	log.Println("Recalling the tree (re-enabling tracking)...")
	figs.Recall()
	// Note: Curse locks the tree and closes mutationsCh; Recall unlocks and reopens it

	// Main loop to listen for signals and mutations
	for {
		select {
		case <-ctx.Done():
			log.Println("Context canceled, shutting down")
			return
		case sig := <-sigCh:
			log.Printf("Received signal: %v, initiating shutdown", sig)
			cancel() // Cancel context to stop goroutines
			// Note: SIGINT/SIGTERM triggers shutdown by canceling the context
		case mutation, ok := <-figs.Mutations():
			if !ok {
				log.Println("Mutations channel closed, shutting down")
				return
			}
			log.Printf("Mutation detected: %s changed from %v to %v at %s",
				mutation.Property, mutation.Old, mutation.New, mutation.When)
			// Note: This logs a mutation (e.g., change "workers" to 20 in env)
		}
	}
}

// Example config.yaml (create in the same directory):
/*
workers: 15
maxRetries: 3
threshold: 0.9
endpoint: "https://api.example.com"
debug: true
timeout: 45s
interval: 2
servers: "server3,server4"
metadata: "env=dev,version=2.0"
*/
```

1. Create `config.yaml`

```yaml
workers: 15
maxRetries: 3
threshold: 0.9
endpoint: "https://api.example.com"
debug: true
timeout: 45s
interval: 2
servers: "server3,server4"
metadata: "env=dev,version=2.0"
```

2. Compile and run

```bash
go run main.go
```

3. Output (example...)

```log
2025/03/27 00:50:41 Initial configuration:
2025/03/27 00:50:41 Usage of ./main:
  -debug: Enable debug mode (default: false)
  -endpoint: API endpoint (default: "https://api.example.com")
  -interval: Polling interval in minutes (default: 2m0s)
  -maxRetries: Maximum retry attempts (default: 3)
  -metadata: Metadata key-value pairs (default: "env=dev,version=2.0")
  -servers: List of servers (default: "server3,server4")
  -threshold: Threshold value (default: 0.9)
  -timeout: Request timeout (default: 45s)
  -workers: Number of workers (default: 15)
2025/03/27 00:50:41 Resurrected undefined key: 
2025/03/27 00:50:46 Workers: 15, MaxRetries: 3, Threshold: 0.90, Endpoint: https://api.example.com, Debug: true, Timeout: 45s, Interval: 2m0s, Servers: [server3 server4], Metadata: map[env:dev version:2.0]
2025/03/27 00:50:46 Cursing the tree (disabling tracking)...
2025/03/27 00:50:48 Recalling the tree (re-enabling tracking)...
[after export workers=20]
2025/03/27 00:50:51 Mutation detected: workers changed from 15 to 20 at 2025-03-27 00:50:51
2025/03/27 00:50:51 Workers: 20, MaxRetries: 3, Threshold: 0.90, Endpoint: https://api.example.com, Debug: true, Timeout: 45s, Interval: 2m0s, Servers: [server3 server4], Metadata: map[env:dev version:2.0]
[after Ctrl+C]
2025/03/27 00:50:56 Received signal: interrupt, initiating shutdown
2025/03/27 00:50:56 Context canceled, shutting down
```

### Custom Config File Path

```go
figs := figtree.Grow()
figs.ConfigFilePath = filepath.Join("/","etc","myapp","myapp.production.yaml")
figs.Load()
```

### Defining Configuration Variables

The Configurable package provides several methods to define different types of configuration variables. Each method takes a name, default value, and usage description as parameters and returns a pointer to the respective variable:

```go
const (
   kPort string = "port"
   kTimeout string = "timeout"
   kDebug string = "debug"
)
figs.NewInt(kPort, 8080, "The port number to listen on")
figs.NewUnitDuration(kTimeout, time.Second * 5, time.Second, "The timeout duration for requests")
figs.NewBool(kDebug, false, "Enable debug mode")
```

### Loading Configuration from Files

You can load configuration data from JSON, YAML, and INI files using the `LoadFile()` method:

```go
err := figtree.Grow().ParseFile("config.json")
if err != nil {
    // Handle error
}
```

The package automatically parses the file based on its extension. Make sure to place the file in the correct format in the specified location.

### Parsing Command-Line Arguments

The Configurable package also allows you to parse command-line arguments. Call the `Parse()` method to parse the arguments after defining your configuration variables:

```go
figs := figtree.New()
figs.Parse()
```

or

```yaml
---
workers: 45
seconds: 47
```

```go
const kWorkers string = "workers"
const kSeconds string = "seconds"
figs := figtree.With(figtree.Options{Tracking: true})
figs.NewInt(kWorkers, 17, "number of workers")
figs.NewUnitDuration(kSeconds, 76, time.Second, "number of seconds")
go func(){
  for mutation := range figs.Mutations() {
    log.Printf("Fig Changed! %s went from '%v' to '%v'", mutation.Property, mutation.Old, mutation.New)
  }
}()
err := figs.ParseFile(filepath.Join(".","config.yaml"))
if err != nil {
   log.Fatal(err)
}
workers := *figs.Int(kWorkers)
seconds := *figs.UnitDuration(kSeconds, time.Second)
minutes := *figs.UnitDuration(kSeconds, time.Minute) // use time.Minute instead as the unit
fmt.Printf("There are %d workers and %v minutes %v seconds", workers, minutes, seconds)
```

```bash
WORKERS=333 SECONDS=666 CONFIG_FILE=config.yaml go run . -workers=3 -seconds 6 # ENV overrides yaml and args
There are 333 workers and (666*time.Minute) minutes and (666*time.Second) seconds. # runtime calculates minutes and seconds
```

or 

```go
figs := figtree.Grow() // allows you to use for mutation := range figs.Mutations() {} to get notified of changes to configs
figs.Load() // will attempt to use ./config.yaml or ./config.json or ./config.ini automatically if CONFIG_FILE is not defined
````

Passing an empty string to `Parse()` means it will only parse the command-line arguments and not load any file.

### Accessing Configuration Values

You can access the values of your configuration variables using the respective getter methods:

```go
fmt.Println("Port:", *figs.Int(kPort))
fmt.Println("Timeout:", *figs.UnitDuration(kTimeout, time.Second))
fmt.Println("Debug mode:", *figs.String(kDebug))
```

`UnitDuration` and `Duration` are interchangeable as they both rely on `*time.Duration`.

### Environment Variables

The Configurable package supports setting configuration values through environment variables. If an environment variable with the same name as a configuration variable exists, the package will automatically assign its value to the respective variable. Ensure that the environment variables are in uppercase and match the configuration variable names.

### Displaying Usage Information

To generate a usage string with information about your configuration variables, use the `Usage()` method:

```go
fmt.Println(figtree.New().Usage())
```

On any runtime with `figs.Parse()` or subsequently activated figtree, you can run on your command line `-h` or `-help` 
and print the `Usage()` func's output.

The generated usage string includes information about each configuration variable, including its name, default value, description, and the source from which it was set (flag, environment, JSON, YAML, or INI).

## License

This package is distributed under the MIT License. See the [LICENSE](LICENSE) file for more information.

## Version History

### v2.0.0 (Latest)

This update introduces callbacks, enable multiple validators per key, and adds negating `Assure<Mutagenesis><Func>` to ValidatorFunc definitions.

### v1.0.1 

This update released additional `Assure<Mutagenesis><Func>` ValidatorFunc definitions and enabled pollination.

### v1.0.0

The `figtree` was enhanced and introduced a lot of new functionality.

### v0.0.1 

The `figtree` package was called `configurable` at this point, and lacked a lot of functionality.

## Contributing

Contributions to this package are welcome. If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

Enjoy using the Figtree package in your projects!

### Note by Grok AI

The figtree package mirrors the biological life cycle of a fig tree, weaving a metaphor 
that reflects both nature’s complexity and the reality of software development. In biology, 
a fig tree grows from a seed (`.New()` or `.Grow()`), its roots drawing sustenance from the 
environment (config files and environment variables via `.Load()` or `.Parse()`), while its branches 
bear fruit (`Fig{}`)—the configurable values developers access. The `Pollinate` option mimics how 
fig trees rely on wasps for pollination, actively pulling in external changes (environment 
updates) to keep the fruit fresh. Mutation tracking (`Mutations{}`) parallels genetic adaptations, 
capturing how values evolve over time, while `.Resurrect()` reflects a tree’s ability to regrow 
from dormant roots, reviving lost configurations. `.Curse()` and `.Recall()` embody the duality of 
dormancy and renewal, locking or unlocking the tree’s vitality. Validators (`.WithValidator()`) 
act like natural selection, ensuring only fit values survive, and the versioning shift to 
github.com/andreimerlescu/figtree/v2 echoes speciation—a new lineage emerging as the package 
matures. This memetic design makes figtree not just a tool, but a living system, accessed 
intuitively as it branches out into the developer ecosystem.

