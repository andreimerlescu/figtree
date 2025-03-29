# Fig Tree

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Fig Tree is a command line utility configuration manager that you can refer to as `figs`, as you con<b>FIG</b>ure your
application's runtime.

## Installation

To use `figtree` in your project, `go get` it...

```shell
go get -u github.com/andreimerlescu/figtree
```

## Usage

To use **figs** package in your Go code, you need to import it:

```go
import "github.com/andreimerlescu/figtree"
```

### Using Fig Tree

| Method                                          | Usage                                 |
|-------------------------------------------------|---------------------------------------|
| `figtree.New()`                                 | Does not perform `Mutation` tracking. |
| `figtree.Grow()`                                | Provides `Mutation` tracking.         |
| `figtree.With(figtree.Options{Tracking: true})` | Provides `Mutation` tracking.         |

Configurable properties have whats called metagenesis to them, which are types, like `String`, `Bool`, `Float64`, etc.

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
    "github.com/andreimerlescu/figtree"
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

### Available Validators

| Mutagenesis | `figtree.ValidatorFunc`   | Notes                                                                            |
|-------------|---------------------------|----------------------------------------------------------------------------------|
| tString     | AssureStringLength        | Ensures a string is a specific length.                                           |
| tString     | AssureStringSubstring     | Ensures a string contains a specific substring (case-sensitive).                 |
| tString     | AssureStringNotEmpty      | Ensures a string is not empty.                                                   |
| tString     | AssureStringContains      | Ensures a string contains a specific substring.                                  |
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
| tList       | AssureListContainsKey     | Ensures a list contains a specific string.                                       |
| tList       | AssureListLength          | Ensures a list has exactly the specified length.                                 |
| tMap        | AssureMapNotEmpty         | Ensures a map (*MapFlag, *map[string]string, or map[string]string) is not empty. |
| tMap        | AssureMapHasKey           | Ensures a map contains a specific key.                                           |
| tMap        | AssureMapValueMatches     | Ensures a map has a specific key with a matching value.                          |
| tMap        | AssureMapHasKeys          | Ensures a map contains all specified keys.                                       |
| tMap        | AssureMapLength           | Ensures a map has exactly the specified length.                                  |



### Complex Example Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "path/filepath"
    "strings"
    "syscall"
    "time"

    "github.com/andreimerlescu/figtree"
)

func main() {
    // Create a context that can be canceled on SIGINT/SIGTERM
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Set up signal handling for graceful shutdown
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    // Plant the fig tree with mutation tracking enabled
    figs := figtree.Grow()
    // Note: Grow() sets Options{Tracking: true}, enabling mutation tracking via mutationsCh

    // Define all configuration types with initial values and validators
    figs.NewInt("workers", 10, "Number of workers")
    figs.WithValidator("workers", func(v interface{}) error {
        if val, ok := v.(int); ok && val < 1 {
            return fmt.Errorf("workers must be at least 1")
        }
        return nil
    })
    figs.NewInt64("maxRetries", 5, "Maximum retry attempts")
    figs.WithValidator("maxRetries", figtree.AssurePositiveInt64)
    figs.NewFloat64("threshold", 0.75, "Threshold value")
    figs.WithValidator("threshold", func(v interface{}) error {
        if val, ok := v.(float64); ok && (val < 0 || val > 1) {
            return fmt.Errorf("threshold must be between 0 and 1")
        }
        return nil
    })
    figs.NewString("endpoint", "http://example.com", "API endpoint")
    figs.WithValidator("endpoint", func(v interface{}) error {
        if val, ok := v.(string); ok && !strings.HasPrefix(val, "http") {
            return fmt.Errorf("endpoint must start with http")
        }
        return nil
    })
    figs.NewBool("debug", false, "Enable debug mode")
    figs.NewDuration("timeout", 30*time.Second, "Request timeout")
    figs.NewUnitDuration("interval", 1, time.Minute, "Polling interval in minutes")
    figs.WithValidator("interval", figtree.AssureDurationGreaterThan(30*time.Second))
    figs.NewList("servers", []string{"server1", "server2"}, "List of servers")
    figs.NewMap("metadata", map[string]string{"env": "prod", "version": "1.0"}, "Metadata key-value pairs")
    // Note: Each New* method registers a flag and initializes withered; WithValidator adds validation logic

    // Attempt to load from a config file, falling back to defaults
    configFile := filepath.Join(".", "config.yaml")
    if err := figs.ParseFile(configFile); err != nil {
        log.Printf("No config file at %s, using defaults: %v", configFile, err)
        err := figs.Parse() // Parse command-line flags and environment variables
        if err != nil {
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
fmt.Println("Port:", *cfigs.Int(kPort))
fmt.Println("Timeout:", *cfigs.UnitDuration(kTimeout, time.Second))
fmt.Println("Debug mode:", *cfigs.String(kDebug))
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

## Contributing

Contributions to this package are welcome. If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

Enjoy using the Figtree package in your projects!


