# Figtree vs Viper: A Practical Comparison

This document is written for Go developers and engineering organizations evaluating
configuration management packages. It assumes you are **not** using remote configuration
sources — which represents the majority of real-world Go projects.

---

## TL;DR

| Capability | Viper | Figtree |
|---|---|---|
| CLI flags | ✅ (via pflag) | ✅ (via stdlib flag) |
| Environment variables | ✅ | ✅ |
| Config files (YAML/JSON/INI) | ✅ | ✅ |
| File watching | ✅ | 🔜 planned v2.1.1+ |
| Struct unmarshaling | ✅ | 🔜 planned v2.2.0+ |
| Per-property validators | ❌ | ✅ 36 built-in |
| Per-property callbacks | ❌ | ✅ |
| Mutation tracking channel | ❌ | ✅ |
| Property aliases | ⚠️ shallow | ✅ full propagation |
| Property rules | ❌ | ✅ |
| Struct tag validation (assure:) | ❌ | 🔜 planned v2.2.0+ |
| Organizational branches | ❌ | 🔜 planned v2.3.0+ |
| Known race conditions | ⚠️ yes | ✅ fixed |
| Remote config sources | ✅ | 🔜 planned |
| stdlib flag compatibility | ❌ | ✅ |
| Zero dependencies (core) | ❌ | ✅ |

---

## Philosophy

**Viper** was designed to be the swiss army knife of Go configuration. It integrates
with `pflag`, supports remote backends like etcd and Consul, and has accumulated years
of community contributions. Its breadth is its strength and its weakness — the API
surface is large, the concurrency model has known issues, and per-property behavior
is not possible without wrapping viper yourself.

**Figtree** was designed around a single premise: every configurable property in your
application is a first-class citizen with its own validators, callbacks, aliases, and
rules. The tree is the unit of organization. Properties are the unit of behavior.
Mutations are observable. The API is opinionated so your application does not have to be.

---

## Feature Breakdown

### CLI Flags

Viper delegates flag parsing to `pflag`, a POSIX-compliant flag package that is itself
a dependency. Figtree uses the Go standard library `flag` package directly, introducing
zero additional dependencies for flag parsing.

```go
// Viper (requires pflag)
pflag.String("host", "localhost", "database host")
viper.BindPFlag("host", pflag.Lookup("host"))

// Figtree
figs.NewString("host", "localhost", "database host")
```

### Environment Variables

Both packages support environment variable binding. Figtree resolves environment
variables automatically by uppercasing the property name and does not require explicit
binding calls.

```go
// Viper
viper.SetEnvPrefix("APP")
viper.AutomaticEnv()
viper.BindEnv("host")

// Figtree — automatic, no binding required
figs.NewString("host", "localhost", "database host")
// HOST=db.example.com ./myapp works automatically
```

### Config Files

Both packages support YAML, JSON, and INI formats. Figtree resolves the config file
path through a priority chain: `CONFIG_FILE` environment variable, `Options.ConfigFile`,
package-level `ConfigFilePath`, then conventional filenames in the working directory.

```go
// Viper
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.ReadInConfig()

// Figtree
figs := figtree.With(figtree.Options{ConfigFile: "config.yaml"})
figs.Load()
```

### File Watching

Viper provides `viper.WatchConfig()` with an `OnConfigChange` callback. Figtree provides
`Options{Watch: true}` with mutations flowing through the existing `Mutations()` channel
— no separate callback registration required. File-driven changes flow through the same
validators, callbacks, and rules as programmatic changes.

```go
// Viper
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    // raw event, no validation, no type safety
})

// Figtree
figs := figtree.With(figtree.Options{
    Watch:     true,
    Tracking:  true,
    ConfigFile: "config.yaml",
})
figs.Load()
for mutation := range figs.Mutations() {
    log.Printf("%s changed: %v → %v", mutation.Property, mutation.Old, mutation.New)
}
```

### Per-Property Validators

Viper has no built-in validation. Developers must validate values after retrieval,
scattering validation logic across the codebase. Figtree ships 36 built-in validators
and supports custom `func(interface{}) error` validators registered per property.

```go
// Viper — validation is your problem
host := viper.GetString("host")
if host == "" {
    return errors.New("host is required")
}

// Figtree — validation is declared at registration
figs.NewString("host", "", "database host")
figs.WithValidator("host", figtree.AssureStringNotEmpty)
figs.WithValidator("host", figtree.AssureStringHasPrefix("postgres://"))
// Parse() or Load() returns error if validation fails
```

The full validator table covers strings, booleans, integers, int64, float64, durations,
lists, and maps. See the [Available Validators](README.md#available-validators) section
for the complete reference.

### Per-Property Callbacks

Viper has no per-property callback system. Figtree supports three callback hooks per
property: `CallbackAfterVerify` (on Parse/Load), `CallbackAfterRead` (on every getter
call), and `CallbackAfterChange` (on every Store call and file-driven reload).

```go
// Viper — no equivalent

// Figtree
figs.WithCallback("host", figtree.CallbackAfterChange, func(value interface{}) error {
    log.Printf("host changed to %v — reconnecting", value)
    return reconnect(value.(string))
})
```

### Mutation Tracking

Viper provides a single `OnConfigChange` hook for file-driven changes only. Programmatic
changes via `viper.Set()` produce no observable event. Figtree provides a buffered
channel that receives a `Mutation` for every value change regardless of source — file,
environment variable, flag, or programmatic `Store()` call.

```go
// Viper — only file changes, no channel
viper.OnConfigChange(func(e fsnotify.Event) { ... })

// Figtree — all changes, channel-based, select-compatible
figs := figtree.With(figtree.Options{Tracking: true, Harvest: 100})
go func() {
    for mutation := range figs.Mutations() {
        log.Printf("%s: %v → %v at %s",
            mutation.Property, mutation.Old, mutation.New, mutation.When)
    }
}()
```

### Property Aliases

Viper supports `viper.RegisterAlias("host", "h")` but the alias only applies at the
getter level. Validators, callbacks, and setters do not propagate through aliases.
Figtree aliases are full citizens — everything that works on the canonical name works
identically on the alias.

```go
// Viper — alias is getter-only
viper.RegisterAlias("verbose", "v")
// viper.Set("v", true) does NOT trigger OnConfigChange for "verbose"

// Figtree — alias propagates everywhere
figs.NewBool("verbose", false, "enable verbose output")
figs.WithAlias("verbose", "v")
figs.WithValidator("verbose", figtree.AssureBoolTrue)
figs.StoreString("v", "true")    // validator fires
*figs.Bool("verbose")            // true
*figs.Bool("v")                  // true
```

### Property Rules

Viper has no concept of property-level rules. Figtree provides rules that govern how
a property behaves at runtime, applied per-property or tree-wide.

```go
// Figtree rules — no Viper equivalent
figs.WithRule("db-password", figtree.RulePreventChange)   // immutable after Parse
figs.WithRule("debug", figtree.RulePanicOnChange)         // panic on change
figs.WithTreeRule(figtree.RuleNoFlags)                    // disable all CLI flags
figs.WithTreeRule(figtree.RuleNoEnv)                      // disable all env vars
```

Full rule reference:

| Rule | Behavior |
|---|---|
| `RulePreventChange` | Blocks all Store calls after initial value is set |
| `RulePanicOnChange` | Panics on any Store call |
| `RuleNoValidations` | Skips all WithValidator assignments |
| `RuleNoCallbacks` | Skips all WithCallback assignments |
| `RuleNoFlags` | Disables CLI flag parsing for the tree |
| `RuleNoEnv` | Skips all os.Getenv logic |
| `RuleNoMaps` | Blocks NewMap, StoreMap, and Map |
| `RuleNoLists` | Blocks NewList, StoreList, and List |
| `RuleCondemnedFromResurrection` | Panics on Resurrect attempts |

### Struct Unmarshaling with Validation

Viper provides `viper.Unmarshal(&cfg)` which populates a struct from the config store.
It has no validation layer — you get the values, validation is your responsibility.
Figtree's `Unmarshal` populates structs and runs inline validators declared via the
`assure:` struct tag, covering all 36 built-in validators.

```go
// Viper
var cfg Config
viper.Unmarshal(&cfg)
// validate cfg yourself

// Figtree
type DatabaseConfig struct {
    Host     string        `fig:"host"    assure:"notEmpty|hasPrefix=postgres://"`
    Port     int           `fig:"port"    assure:"inRange=1024,65535"`
    Timeout  time.Duration `fig:"timeout" assure:"min=5s|max=2m"`
}
var cfg DatabaseConfig
err := figs.Unmarshal(&cfg)
// validation runs inline, UnmarshalError carries field, fig key, and failing token
```

### Concurrency

Viper has well-documented race conditions that have been open issues for years. Safe
concurrent use of viper requires external locking that the developer must implement
and maintain. Figtree was designed with concurrency as a first-class concern — all
internal operations use `sync.RWMutex` correctly, the mutations channel is safely
buffered, and `Store()` releases locks before channel sends to prevent deadlock.

---

## Migration from Viper to Figtree

The conceptual mapping is straightforward for projects not using remote config sources.

### Installation

```bash
go get -u github.com/andreimerlescu/figtree/v2
```

### Import

```go
// Before
import "github.com/spf13/viper"

// After
import "github.com/andreimerlescu/figtree/v2"
```

### Initialization

```go
// Viper
viper.SetDefault("host", "localhost")
viper.SetDefault("port", 5432)
viper.AutomaticEnv()

// Figtree
figs := figtree.Grow()
figs.NewString("host", "localhost", "database host")
figs.NewInt("port", 5432, "database port")
```

### Reading Values

```go
// Viper
host := viper.GetString("host")
port := viper.GetInt("port")

// Figtree
host := *figs.String("host")
port := *figs.Int("port")
```

### Writing Values

```go
// Viper
viper.Set("host", "db.example.com")

// Figtree
figs.StoreString("host", "db.example.com")
```

### Config File Loading

```go
// Viper
viper.SetConfigFile("config.yaml")
viper.ReadInConfig()

// Figtree
figs := figtree.With(figtree.Options{ConfigFile: "config.yaml"})
figs.Load()
```

### Struct Unmarshaling

```go
// Viper
var cfg Config
viper.Unmarshal(&cfg)

// Figtree
var cfg Config
figs.Unmarshal(&cfg)
```

---

## When to Choose Viper

- Your project depends on `pflag` and you cannot migrate
- You require remote configuration sources (etcd, Consul, Vault) today
- You have deep existing viper integration with significant migration cost

## When to Choose Figtree

- You want per-property validation without writing it yourself
- You want observable mutations via a channel
- You want immutability rules on specific properties
- You want struct unmarshaling with inline validation
- You want a concurrency-safe configuration package
- You are not using viper's remote configuration capabilities
- You want zero non-stdlib dependencies in your configuration layer
- You want organizational structure via branches (v2.3.0+)

---

## Summary

For the majority of Go projects that use viper for flags, environment variables, and
config files — and nothing more — figtree offers a more expressive, safer, and more
maintainable alternative. The per-property model gives you validation, callbacks, rules,
and mutation tracking that viper requires you to build yourself.

If your project does not use viper's remote configuration sources, figtree is worth
evaluating as a direct replacement.

---

*Figtree is maintained by [@andreimerlescu](https://github.com/andreimerlescu).*
*Current stable release: v2.1.0*
*License: MIT*
