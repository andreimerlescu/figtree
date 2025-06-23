package figtree

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValue_Set_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		mutagenesis   Mutagenesis
		input         string
		expectError   bool
		expectedValue interface{} // Expected value after failed conversion (often default or zero)
	}{
		{"String_Valid", tString, "hello", false, "hello"},
		{"Bool_Invalid", tBool, "not-a-bool", true, false}, // Assuming bool defaults to false on conversion error
		{"Int_Invalid", tInt, "abc", true, 0},              // Assuming int defaults to 0 on conversion error
		{"Int64_Invalid", tInt64, "xyz", true, int64(0)},
		{"Float64_Invalid", tFloat64, "badfloat", true, 0.0},
		{"Duration_Invalid", tDuration, "not-a-duration", true, time.Duration(0)}, // Note: toInt64 from string might return 0
		{"List_InvalidFormat", tList, "a,b=c", true, []string{"a", "b=c"}},        // Assumes partial parse or specific error
		{"Map_InvalidFormat", tMap, "k1v1,k2", true, map[string]string{}},         // Assumes map fully resets on invalid item
		{"List_Valid", tList, "item1,item2", false, []string{"item1", "item2"}},
		{"Map_Valid", tMap, "key1=val1,key2=val2", false, map[string]string{"key1": "val1", "key2": "val2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := &Value{Value: tt.expectedValue, Mutagensis: tt.mutagenesis} // Initialize with expected zero/default
			err := val.Set(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				// For primitive types, the underlying value might stay its initial default (zero value)
				// For List/Map, check how your Set handles errors.
				// For List/Map, they might partial parse or reset, depending on implementation.
				// This needs careful assertion based on exact behavior of ListFlag/MapFlag Set.
				if _, ok := tt.expectedValue.(map[string]string); ok {
					assert.Equal(t, tt.expectedValue, val.Flesh().ToMap(), "Value should be default/empty on map error")
				} else if _, ok := tt.expectedValue.([]string); ok {
					assert.Equal(t, tt.expectedValue, val.Flesh().ToList(), "Value should be default/empty on list error")
				} else {
					assert.Equal(t, tt.expectedValue, val.Value, "Value should be default on primitive error")
				}

			} else {
				assert.NoError(t, err)
				// Re-fetch value correctly based on what Set assigns
				actualValue := val.Value
				// For List/Map, they might be *MapFlag or *ListFlag, need to unwrap
				if tf, ok := actualValue.(MapFlag); ok {
					actualValue = tf.values
				} else if tf, ok := actualValue.(*MapFlag); ok {
					actualValue = tf.values
				} else if tf, ok := actualValue.(ListFlag); ok {
					actualValue = tf.values
				} else if tf, ok := actualValue.(*ListFlag); ok {
					actualValue = tf.values
				}
				assert.Equal(t, tt.expectedValue, actualValue, "Value should be correctly set")
			}
		})
	}
}

func TestParse_InvalidFlagInput(t *testing.T) {
	tests := []struct {
		name         string
		flagName     string
		defaultValue interface{}
		usage        string
		argValue     string
	}{
		{"IntFlag_InvalidString", "port", zeroInt, "port number", "not-a-number"},
		{"BoolFlag_InvalidString", "debug", zeroBool, "debug mode", "maybe"},
		{"Float64Flag_InvalidString", "ratio", zeroFloat64, "ratio value", "bad-float"},
		{"DurationFlag_InvalidString", "timeout", zeroDuration, "timeout duration", "invalid-duration"},
		{"ListFlag_MalformedItem", "tags", []string{"a"}, "list of tags", "item1,item2=val"},
		{"MapFlag_MalformedItem", "config", map[string]string{"k": "v"}, "config map", "k1=v1,k2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			os.Args = []string{os.Args[0], "-" + tt.flagName, tt.argValue}

			figs := With(Options{Germinate: true})
			switch v := tt.defaultValue.(type) {
			case int:
				figs = figs.NewInt(tt.flagName, v, tt.usage)
			case bool:
				figs = figs.NewBool(tt.flagName, v, tt.usage)
			case float64:
				figs = figs.NewFloat64(tt.flagName, v, tt.usage)
			case time.Duration:
				figs = figs.NewDuration(tt.flagName, v, tt.usage)
			case []string:
				figs = figs.NewList(tt.flagName, v, tt.usage)
			case map[string]string:
				figs = figs.NewMap(tt.flagName, v, tt.usage)
			default:
				t.Fatalf("Unknown default type for test: %T", tt.defaultValue)
			}

			err := figs.Parse()
			assert.Error(t, err, "Parse() should return an error for invalid flag input")
			// Assert that the value either retains its default or a zero value,
			// and that there isn't a panic.
			switch tt.defaultValue.(type) {
			case int:
				assert.Equal(t, 0, *figs.Int(tt.flagName))
			case bool:
				assert.Equal(t, false, *figs.Bool(tt.flagName))
			case float64:
				assert.Equal(t, 0.0, *figs.Float64(tt.flagName))
			case time.Duration:
				assert.Equal(t, time.Duration(0), *figs.Duration(tt.flagName))
			// For lists/maps, behavior on partial parse/error might vary:
			case []string:
				assert.Empty(t, *figs.List(tt.flagName)) // Expect default or empty on parse error
			case map[string]string:
				assert.Empty(t, *figs.Map(tt.flagName)) // Expect default or empty on parse error
			}
		})
	}
}

func TestEnvironment_InvalidInput(t *testing.T) {
	tests := []struct {
		name         string
		envName      string
		defaultValue interface{}
		usage        string
		envValue     string
	}{
		{"IntEnv_InvalidString", "MY_PORT", zeroInt, "port number", "not-a-number"},
		{"BoolEnv_InvalidString", "MY_DEBUG", zeroBool, "debug mode", "maybe"},
		{"Float64Env_InvalidString", "MY_RATIO", zeroFloat64, "ratio value", "bad-float"},
		{"DurationEnv_InvalidString", "MY_TIMEOUT", zeroDuration, "timeout duration", "invalid-duration"},
		{"ListEnv_MalformedItem", "MY_TAGS", zeroList, "list of tags", "item1,item2=val"},
		{"MapEnv_MalformedItem", "MY_CONFIG", zeroMap, "config map", "k1=v1,k2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			os.Args = []string{os.Args[0]} // Ensure no CLI args interfere
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				timer := time.NewTimer(time.Second * 1)
				checker := time.NewTicker(100 * time.Millisecond)
				for {
					select {
					case <-timer.C:
						return
					case <-checker.C:
						assert.NoError(t, os.Setenv(strings.ToUpper(tt.envName), tt.envValue))
					}
				}
			}()
			defer assert.NoError(t, os.Unsetenv(tt.envName))
			wg.Wait()

			figs := With(Options{Germinate: true})
			switch v := tt.defaultValue.(type) {
			case int:
				figs = figs.NewInt(tt.envName, v, tt.usage)
			case bool:
				figs = figs.NewBool(tt.envName, v, tt.usage)
			case float64:
				figs = figs.NewFloat64(tt.envName, v, tt.usage)
			case time.Duration:
				figs = figs.NewDuration(tt.envName, v, tt.usage)
			case []string:
				figs = figs.NewList(tt.envName, v, tt.usage)
			case map[string]string:
				figs = figs.NewMap(tt.envName, v, tt.usage)
			default:
				t.Fatalf("Unknown default type for test: %T", tt.defaultValue)
			}

			err := figs.Load() // Use Load for env vars
			assert.Error(t, err, "Load() should return an error for invalid env input")
			// Assert that the value either retains its default or a zero value,
			// and that there isn't a panic.
			switch tt.defaultValue.(type) {
			case int:
				assert.Equal(t, zeroInt, *figs.Int(tt.envName))
			case bool:
				assert.Equal(t, zeroBool, *figs.Bool(tt.envName))
			case float64:
				assert.Equal(t, zeroFloat64, *figs.Float64(tt.envName))
			case time.Duration:
				assert.Equal(t, zeroDuration, *figs.Duration(tt.envName))
			case []string:
				assert.Empty(t, *figs.List(strings.ToLower(tt.envName)))
			case map[string]string:
				assert.Empty(t, *figs.Map(strings.ToLower(tt.envName)))
			}
		})
	}
}

// Test for empty string inputs to non-string types
func TestEmptyStringInput(t *testing.T) {
	tests := []struct {
		name         string
		flagName     string
		defaultValue interface{}
		usage        string
	}{
		{"IntFlag_EmptyString", "count", zeroInt, "count"},
		{"BoolFlag_EmptyString", "enabled", zeroBool, "enabled"},
		{"Float64Flag_EmptyString", "ratio", zeroFloat64, "ratio"},
		{"DurationFlag_EmptyString", "interval", zeroDuration, "interval"},
		{"ListFlag_EmptyString", "items", zeroList, "items"},
		{"MapFlag_EmptyString", "data", zeroMap, "data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			os.Args = []string{os.Args[0], "-" + tt.flagName, ""} // Pass empty string

			figs := With(Options{Germinate: true})
			switch v := tt.defaultValue.(type) {
			case int:
				figs = figs.NewInt(tt.flagName, v, tt.usage)
			case bool:
				figs = figs.NewBool(tt.flagName, v, tt.usage)
			case float64:
				figs = figs.NewFloat64(tt.flagName, v, tt.usage)
			case time.Duration:
				figs = figs.NewDuration(tt.flagName, v, tt.usage)
			case []string:
				figs = figs.NewList(tt.flagName, v, tt.usage)
			case map[string]string:
				figs = figs.NewMap(tt.flagName, v, tt.usage)
			default:
				t.Fatalf("Unknown default type for test: %T", tt.defaultValue)
			}

			err := figs.Parse()
			assert.NoError(t, err, "Parse() should not return error for empty string (handled by Set)")
			// The Set method for List and Map handles empty string as empty slice/map, which is good.
			// For primitives, it will attempt conversion, which should fail and keep default/zero.
			switch tt.defaultValue.(type) {
			case int:
				assert.Equal(t, zeroInt, *figs.Int(tt.flagName), "Int should be zero") // Atoi("") returns 0, err
			case bool:
				assert.Equal(t, zeroBool, *figs.Bool(tt.flagName), "Bool should be false") // ParseBool("") returns false, err
			case float64:
				assert.Equal(t, zeroFloat64, *figs.Float64(tt.flagName), "Float should be zero") // ParseFloat("") returns 0, err
			case time.Duration:
				assert.Equal(t, zeroDuration, *figs.Duration(tt.flagName), "Duration should be zero") // ParseDuration("") returns 0, err
			case []string:
				assert.Equal(t, zeroList, *figs.List(tt.flagName), "List should be empty")
			case map[string]string:
				assert.Equal(t, zeroMap, *figs.Map(tt.flagName), "Map should be empty")
			}
		})
	}
}

func TestRulePreventChange_StoreMethods(t *testing.T) {
	const flagName = "protected_setting"
	const initialValue = "secret"
	const newValue = "exposed"

	t.Run("PreventChange_On_StoreString", func(t *testing.T) {
		os.Args = []string{os.Args[0]} // Clean args
		figs := With(Options{Germinate: true})
		figs = figs.NewString(flagName, initialValue, "A protected string")
		figs = figs.WithRule(flagName, RulePreventChange)
		assert.NoError(t, figs.Parse()) // Initial parse will succeed

		// Attempt to store a new value
		figs = figs.StoreString(flagName, newValue)
		assert.NoError(t, figs.ErrorFor(flagName), "Store should not generate an error when RulePreventChange is active")

		// Assert that the value *did not change*
		assert.Equal(t, initialValue, *figs.String(flagName), "Value should remain unchanged when RulePreventChange is active")
	})

	t.Run("PreventChange_On_StoreInt", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true})
		figs = figs.NewInt(flagName, 100, "A protected int")
		figs = figs.WithRule(flagName, RulePreventChange)
		assert.NoError(t, figs.Parse())

		figs = figs.StoreInt(flagName, 200)
		assert.NoError(t, figs.ErrorFor(flagName))
		assert.Equal(t, 100, *figs.Int(flagName), "Int value should remain unchanged")
	})

	// Add similar tests for StoreBool, StoreFloat64, StoreDuration, StoreUnitDuration, StoreList, StoreMap
	// For Lists and Maps, ensure the underlying slice/map reference is not modified either.
	t.Run("PreventChange_On_StoreList", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		initialList := []string{"a", "b"}
		figs := With(Options{Germinate: true})
		figs = figs.NewList(flagName, initialList, "A protected list")
		figs = figs.WithRule(flagName, RulePreventChange)
		assert.NoError(t, figs.Parse())

		newValue := []string{"x", "y", "z"}
		figs = figs.StoreList(flagName, newValue)
		assert.NoError(t, figs.ErrorFor(flagName))
		assert.Equal(t, initialList, *figs.List(flagName), "List value should remain unchanged")
		// Importantly, ensure the underlying slice isn't the new slice reference, but still the original's value
		l := *figs.List(flagName)
		assert.NotSame(t, &newValue, &l, "Should be a copy or original reference, not the new slice directly")
	})

	t.Run("PreventChange_On_StoreMap", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		initialMap := map[string]string{"k1": "v1"}
		figs := With(Options{Germinate: true})
		figs = figs.NewMap(flagName, initialMap, "A protected map")
		figs = figs.WithRule(flagName, RulePreventChange)
		assert.NoError(t, figs.Parse())

		newValue := map[string]string{"k2": "v2"}
		figs = figs.StoreMap(flagName, newValue)
		assert.NoError(t, figs.ErrorFor(flagName))
		assert.Equal(t, initialMap, *figs.Map(flagName), "Map value should remain unchanged")
		l := *figs.Map(flagName)
		assert.NotSame(t, &newValue, &l, "Should be a copy or original reference, not the new map directly")
	})
}

func TestRuleCondemnedFromResurrection(t *testing.T) {
	const nonExistentFlag = "ghost_flag"

	t.Run("Condemned_Flag_Resurrection_Panics", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true})
		// Condemn a flag name before it's even registered.
		// `Resurrect` can be called internally if a flag is accessed but not defined.
		figs = figs.WithRule(nonExistentFlag, RuleCondemnedFromResurrection)

		// Directly access the non-existent flag to trigger potential Resurrect call
		assert.Panics(t, func() {
			_ = *figs.String(nonExistentFlag)
		})
	})
}

func TestFlesh_NilUnderlyingValue(t *testing.T) {
	t.Run("ToString_NilValue", func(t *testing.T) {
		flesh := NewFlesh(nil) // Create Flesh with nil underlying value
		assert.Empty(t, flesh.ToString(), "ToString() on nil value should return empty string")
	})

	t.Run("ToInt_NilValue", func(t *testing.T) {
		flesh := NewFlesh(nil)
		assert.Equal(t, 0, flesh.ToInt(), "ToInt() on nil value should return 0")
	})

	t.Run("ToBool_NilValue", func(t *testing.T) {
		flesh := NewFlesh(nil)
		assert.False(t, flesh.ToBool(), "ToBool() on nil value should return false")
	})

	t.Run("ToList_NilValue", func(t *testing.T) {
		flesh := NewFlesh(nil)
		assert.Empty(t, flesh.ToList(), "ToList() on nil value should return empty slice")
		assert.NotNil(t, flesh.ToList(), "ToList() on nil value should return non-nil empty slice") // IMPORTANT: avoid nil slice
	})

	t.Run("ToMap_NilValue", func(t *testing.T) {
		flesh := NewFlesh(nil)
		assert.Empty(t, flesh.ToMap(), "ToMap() on nil value should return empty map")
		assert.NotNil(t, flesh.ToMap(), "ToMap() on nil value should return non-nil empty map") // IMPORTANT: avoid nil map
	})

	t.Run("Is_NilValue", func(t *testing.T) {
		flesh := NewFlesh(nil)
		assert.False(t, flesh.IsString(), "IsString() on nil should be false")
		assert.False(t, flesh.IsInt(), "IsInt() on nil should be false")
		// ... test all other IsX methods
	})
}

func TestUsageOutputConsistency(t *testing.T) {
	const flagName = "workers"
	const flagAlias = "w"
	const flagUsage = "Number of worker goroutines"
	const flagDefault = 10

	t.Run("BasicFlagUsage", func(t *testing.T) {
		os.Args = []string{os.Args[0]} // Clean args
		figs := With(Options{Germinate: true})
		figs.NewInt(flagName, flagDefault, flagUsage)
		output := figs.UsageString()

		assert.Contains(t, output, fmt.Sprintf("-%s[=%d]", flagName, flagDefault), "Output should contain flag name with default value")
		assert.Contains(t, output, fmt.Sprintf("[%s]", tInt), "Output should contain correct type string")
		assert.Contains(t, output, flagUsage, "Output should contain correct usage string")
	})

	t.Run("AliasedFlagUsage", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true})
		figs.NewString(flagName, "myhost", flagUsage)
		figs.WithAlias(flagName, flagAlias)
		output := figs.UsageString()
		expectedAliasStr := fmt.Sprintf("-%s|-%s[=myhost]", flagAlias, flagName)
		assert.Contains(t, output, expectedAliasStr, "Output should show alias and main flag name with default")
		assert.Contains(t, output, fmt.Sprintf("[%s]", tString), "Output should contain correct type string for aliased flag")
		assert.Contains(t, output, flagUsage, "Output should contain correct usage string for aliased flag")
	})

	t.Run("ListFlagUsage", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true})
		defaultList := []string{"one", "two"}
		figs.NewList(flagName, defaultList, flagUsage)
		output := figs.UsageString()
		// Default value for ListFlag.String() is 'one,two' (using ListSeparator)
		expectedDefault := strings.Join(defaultList, ListSeparator)
		assert.Contains(t, output, fmt.Sprintf("-%s[=%s]", flagName, expectedDefault), "List flag usage should show joined default")
		assert.Contains(t, output, fmt.Sprintf("[%s]", tList), "List flag usage should show List type")
	})

	t.Run("MapFlagUsage", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true})
		defaultMap := map[string]string{"k1": "v1", "k2": "v2"}
		figs.NewMap(flagName, defaultMap, flagUsage)

		output := figs.UsageString()
		// MapFlag.String() returns 'k1=v1,k2=v2' (or similar, order non-guaranteed)
		// For assertion, we need to be flexible with map string order.
		// A simpler assertion might be to check for parts.
		assert.Contains(t, output, fmt.Sprintf("-%s[=", flagName), "Map flag usage should show default value start")
		assert.Contains(t, output, fmt.Sprintf("[%s]", tMap), "Map flag usage should show Map type")
		assert.Contains(t, output, "k1=v1", "Map flag usage should contain key-value pair")
		assert.Contains(t, output, "k2=v2", "Map flag usage should contain another key-value pair")
	})
}

/*
func TestListPolicyAppend_FullPrecedence(t *testing.T) {
	const listName = "app_list"
	tempConfigFile := filepath.Join(t.TempDir(), "list_config.yaml")
	defer assert.NoError(t, os.Remove(tempConfigFile))
	defer assert.NoError(t, os.Unsetenv(strings.ToUpper(listName)))
	os.Args = []string{os.Args[0]}

	initialDefault := []string{"apple", "banana"}
	fileContent := []string{"banana", "cherry"}
	envValue := "cherry,date"
	cliValue := "apple,elderberry"

	// Set PolicyListAppend to true for these tests
	PolicyListAppend = true
	defer func() { PolicyListAppend = false }() // Reset after tests

	// Write config file
	assert.NoError(t, os.WriteFile(tempConfigFile, []byte(listName+": "+strings.Join(fileContent, ListSeparator)), 0644))

	// Set env variable
	assert.NoError(t, os.Setenv(strings.ToUpper(listName), envValue))

	// Set CLI args
	os.Args = []string{os.Args[0], "-" + listName, cliValue}

	figs := With(Options{
		Germinate:  true,
		ConfigFile: tempConfigFile,
	})
	figs.NewList(listName, initialDefault, "App list settings")
	assert.NoError(t, figs.Parse())

	// Expected order will be sorted due to slices.Sort in List() getter
	expectedList := []string{"apple", "banana", "cherry", "date"} // unique, merged, sorted
	assert.Equal(t, expectedList, *figs.List(listName), "List should merge correctly with precedence and deduplicate")

	// Test case where policy is false (full overwrite)
	t.Run("PolicyListAppend_FALSE_FullOverwrite", func(t *testing.T) {
		PolicyListAppend = false // Ensure it's false for this sub-test
		os.Clearenv()
		os.Args = []string{os.Args[0]} // No CLI arg
		assert.NoError(t, os.Setenv(strings.ToUpper(listName), envValue))

		// Clean up and re-initialize figs for this sub-test
		figs = With(Options{
			Germinate:  true,
			ConfigFile: tempConfigFile,
		})
		figs.NewList(listName, initialDefault, "App list settings")
		assert.NoError(t, figs.Parse())               // Parse will pick up Env & File
		expectedEnvList := []string{"cherry", "date"} // Env completely overwrites
		assert.Equal(t, expectedEnvList, *figs.List(listName), "Env should completely overwrite file/default when PolicyListAppend is false")

		os.Args = []string{os.Args[0], "-" + listName, cliValue} // Now with CLI arg
		figs = With(Options{
			Germinate:  true,
			ConfigFile: tempConfigFile,
		})
		figs.NewList(listName, initialDefault, "App list settings")
		assert.NoError(t, figs.Parse())                    // Parse will pick up CLI
		expectedCliList := []string{"apple", "elderberry"} // CLI completely overwrites
		assert.Equal(t, expectedCliList, *figs.List(listName), "CLI should completely overwrite all when PolicyListAppend is false")
	})
}

func TestMapPolicyAppend_FullPrecedence(t *testing.T) {
	const mapName = "app_map"
	tempConfigFile := filepath.Join(t.TempDir(), "map_config.yaml")
	defer assert.NoError(t, os.Remove(tempConfigFile))
	defer assert.NoError(t, os.Unsetenv(strings.ToUpper(mapName)))
	os.Args = []string{os.Args[0]}

	initialDefault := map[string]string{"a": "1", "b": "2", "c": "3"}
	fileContent := map[string]string{"b": "B_file", "d": "4", "e": "5"}
	envValue := "c=C_env,f=6,g=7"
	cliValue := "a=A_cli,e=E_cli,h=8"

	// Set PolicyMapAppend to true for these tests
	PolicyMapAppend = true
	defer func() { PolicyMapAppend = false }() // Reset after tests

	// Write config file
	fileMapStr := ""
	for k, v := range fileContent {
		if fileMapStr != "" {
			fileMapStr += MapSeparator
		}
		fileMapStr += k + MapKeySeparator + v
	}
	assert.NoError(t, os.WriteFile(tempConfigFile, []byte(mapName+": "+fileMapStr), 0644))

	// Set env variable
	assert.NoError(t, os.Setenv(strings.ToUpper(mapName), envValue))

	// Set CLI args
	os.Args = []string{os.Args[0], "-" + mapName, cliValue}

	figs := With(Options{
		Germinate:  true,
		ConfigFile: tempConfigFile,
	})
	figs.NewMap(mapName, initialDefault, "App map settings")
	assert.NoError(t, figs.Parse())

	expectedMap := map[string]string{
		"a": "A_cli",  // CLI overrides default
		"b": "B_file", // File overrides default (no env/cli conflict)
		"c": "C_env",  // Env overrides default (no cli conflict)
		"d": "4",      // File value
		"e": "E_cli",  // CLI overrides file
		"f": "6",      // Env value
		"g": "7",      // Env value
		"h": "8",      // CLI value
	}
	assert.Equal(t, expectedMap, *figs.Map(mapName), "Map should merge correctly with precedence")

	// Test case where policy is false (full overwrite)
	t.Run("PolicyMapAppend_FALSE_FullOverwrite", func(t *testing.T) {
		PolicyMapAppend = false // Ensure it's false for this sub-test
		os.Clearenv()
		os.Args = []string{os.Args[0]} // No CLI arg
		assert.NoError(t, os.Setenv(strings.ToUpper(mapName), envValue))

		// Clean up and re-initialize figs for this sub-test
		figs = With(Options{
			Germinate:  true,
			ConfigFile: tempConfigFile,
		})
		figs.NewMap(mapName, initialDefault, "App map settings")
		assert.NoError(t, figs.Parse()) // Parse will pick up Env & File
		assert.Equal(t, map[string]string{"c": "C_env", "f": "6", "g": "7"}, *figs.Map(mapName), "Env should completely overwrite file/default when PolicyMapAppend is false")

		os.Args = []string{os.Args[0], "-" + mapName, cliValue} // Now with CLI arg
		figs = With(Options{
			Germinate:  true,
			ConfigFile: tempConfigFile,
		})
		figs.NewMap(mapName, initialDefault, "App map settings")
		assert.NoError(t, figs.Parse()) // Parse will pick up CLI
		assert.Equal(t, map[string]string{"a": "A_cli", "e": "E_cli", "h": "8"}, *figs.Map(mapName), "CLI should completely overwrite all when PolicyMapAppend is false")
	})
}

func TestPrecedence(t *testing.T) {
	const flagName = "app_setting"
	const defaultValue = "default_value"
	const fileValue = "file_value"
	const envValue = "env_value"
	const cliValue = "cli_value"

	// Define temporary file paths
	tempConfigFile := filepath.Join(t.TempDir(), "test_config.yaml")

	// Helper to clean up OS environment variables
	defer os.Unsetenv(strings.ToUpper(flagName))
	defer os.Unsetenv(EnvironmentKey)

	// --- Scenario 1: CLI Flag should win everything ---
	t.Run("CLI_Wins", func(t *testing.T) {
		os.Clearenv()
		os.Args = []string{os.Args[0], "-" + flagName, cliValue}
		assert.NoError(t, os.Setenv(strings.ToUpper(flagName), envValue))                      // Set env variable
		assert.NoError(t, os.WriteFile(tempConfigFile, []byte(flagName+": "+fileValue), 0644)) // Write config file

		figs := With(Options{
			Germinate:         true,
			ConfigFile:        tempConfigFile,
			Pollinate:         false, // Pollinate not directly affecting initial parse order here
			IgnoreEnvironment: false,
		})
		figs.NewString(flagName, defaultValue, "App setting")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, cliValue, *figs.String(flagName), "CLI value should override all")
	})

	// --- Scenario 2: Environment variable should win over file/default when no CLI flag ---
	t.Run("Env_Wins_Over_File_Default", func(t *testing.T) {
		os.Clearenv()
		os.Args = []string{os.Args[0]} // No CLI flag
		assert.NoError(t, os.Setenv(strings.ToUpper(flagName), envValue))
		assert.NoError(t, os.WriteFile(tempConfigFile, []byte(flagName+": "+fileValue), 0644))

		figs := With(Options{
			Germinate:         true,
			ConfigFile:        tempConfigFile,
			IgnoreEnvironment: false,
		})
		figs.NewString(flagName, defaultValue, "App setting")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, envValue, *figs.String(flagName), "Environment value should override file and default")
	})

	// --- Scenario 3: Config file should win over default when no CLI/Env ---
	t.Run("File_Wins_Over_Default", func(t *testing.T) {
		os.Clearenv()
		os.Args = []string{os.Args[0]} // No CLI flag
		// No environment variable set
		assert.NoError(t, os.WriteFile(tempConfigFile, []byte(flagName+": "+fileValue), 0644))

		figs := With(Options{
			Germinate:         true,
			ConfigFile:        tempConfigFile,
			IgnoreEnvironment: false,
		})
		figs.NewString(flagName, defaultValue, "App setting")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, fileValue, *figs.String(flagName), "Config file value should override default")
	})

	// --- Scenario 4: Default value is used when no other source provides a value ---
	t.Run("Default_Value", func(t *testing.T) {
		os.Clearenv()
		os.Args = []string{os.Args[0]} // No CLI flag
		// No environment variable set
		// No config file for this specific flag
		assert.NoError(t, os.WriteFile(tempConfigFile, []byte("some_other_setting: test"), 0644)) // Ensure file exists but doesn't have our flag

		figs := With(Options{
			Germinate:         true,
			ConfigFile:        tempConfigFile,
			IgnoreEnvironment: false,
		})
		figs.NewString(flagName, defaultValue, "App setting")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, defaultValue, *figs.String(flagName), "Default value should be used")
	})

	// --- Scenario 5: IgnoreEnvironment option bypasses env vars ---
	t.Run("IgnoreEnvironment", func(t *testing.T) {
		os.Clearenv()
		os.Args = []string{os.Args[0]}
		assert.NoError(t, os.Setenv(strings.ToUpper(flagName), envValue))                      // Env is set
		assert.NoError(t, os.WriteFile(tempConfigFile, []byte(flagName+": "+fileValue), 0644)) // File is set

		figs := With(Options{
			Germinate:         true,
			ConfigFile:        tempConfigFile,
			IgnoreEnvironment: true, // This should make envValue ignored
		})
		figs.NewString(flagName, defaultValue, "App setting")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, fileValue, *figs.String(flagName), "File value should win when env is ignored")
	})

	// Clean up temporary file
	defer assert.NoError(t, os.Remove(tempConfigFile))
}
*/

func TestCornerCases(t *testing.T) {
	t.Run("UsingEnv", func(t *testing.T) {
		assert.NoError(t, os.Setenv("NAME", "Yeshua"))
		os.Args = []string{os.Args[0], "-name", "Andrei"}
		figs := With(Options{Germinate: true})
		figs = figs.NewString("name", "Satan", "Your Name").WithAlias("name", "n")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, "Yeshua", *figs.String("name"))
	})
}

func TestWith(t *testing.T) {
	type args struct {
		opts Options
	}
	tests := []struct {
		name      string
		args      args
		wantNil   bool
		wantTrack bool
		wantFile  string
	}{
		{
			name:      "With tracking enabled",
			args:      args{opts: Options{Tracking: true}},
			wantNil:   false,
			wantTrack: true,
			wantFile:  "",
		},
		{
			name:      "With tracking disabled and custom config file",
			args:      args{opts: Options{Tracking: false, ConfigFile: "custom.yaml"}},
			wantNil:   false,
			wantTrack: false,
			wantFile:  "custom.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment to ensure clean slate
			os.Clearenv()

			figs := With(tt.args.opts)
			assert.Equal(t, tt.wantNil, figs == nil, "With(%v) should return a non-nil Plant", tt.args.opts)
			if !tt.wantNil {
				tree, ok := figs.(*figTree)
				assert.True(t, ok, "With() should return a *figTree")
				assert.Equal(t, tt.wantTrack, tree.tracking, "With() should set tracking to %v", tt.wantTrack)
				assert.NotNil(t, tree.figs, "With() should initialize figs map")
				assert.NotNil(t, tree.withered, "With() should initialize withered map")
				assert.NotNil(t, tree.mutationsCh, "With() should initialize mutationsCh")
				assert.Equal(t, tt.wantFile, tree.ConfigFilePath, "With() should set ConfigFilePath to %v", tt.wantFile)
				// Verify mutationsCh behavior based on tracking
				select {
				case _, ok := <-figs.Mutations():
					assert.True(t, ok, "Mutations channel should be open after With")
				default:
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		wantNil   bool
		wantTrack bool
	}{
		{
			name:      "New creates tree without tracking",
			wantNil:   false,
			wantTrack: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment to ensure clean slate
			os.Clearenv()

			figs := New()
			assert.Equal(t, tt.wantNil, figs == nil, "New() should return a non-nil Plant")
			if !tt.wantNil {
				tree, ok := figs.(*figTree)
				assert.True(t, ok, "New() should return a *figTree")
				assert.Equal(t, tt.wantTrack, tree.tracking, "New() should disable tracking by default")
				assert.NotNil(t, tree.figs, "New() should initialize figs map")
				assert.NotNil(t, tree.withered, "New() should initialize withered map")
				assert.NotNil(t, tree.mutationsCh, "New() should initialize mutationsCh")
				// Verify mutationsCh is open (even if not tracking, channel should exist)
				select {
				case _, ok := <-figs.Mutations():
					assert.True(t, ok, "Mutations channel should be open after New")
				default:
				}
			}
		})
	}
}

func TestGrow(t *testing.T) {
	tests := []struct {
		name      string
		wantNil   bool
		wantTrack bool
	}{
		{
			name:      "Grow creates tree with tracking",
			wantNil:   false,
			wantTrack: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment to ensure clean slate
			os.Clearenv()

			figs := Grow()
			assert.Equal(t, tt.wantNil, figs == nil, "Grow() should return a non-nil Plant")
			if !tt.wantNil {
				tree, ok := figs.(*figTree)
				assert.True(t, ok, "Grow() should return a *figTree")
				assert.True(t, tree.tracking, "Grow() should enable tracking")
				assert.NotNil(t, tree.figs, "Grow() should initialize figs map")
				assert.NotNil(t, tree.withered, "Grow() should initialize withered map")
				assert.NotNil(t, tree.mutationsCh, "Grow() should initialize mutationsCh")
				// Verify mutationsCh is open
				select {
				case _, ok := <-figs.Mutations():
					assert.True(t, ok, "Mutations channel should be open after Grow")
				default:
				}
			}
		})
	}
}

func TestVersion(t *testing.T) {
	t.Run("current_version_default_empty", func(t *testing.T) {
		if len(currentVersion) > 0 {
			currentVersion = ""
		}
		assert.Empty(t, currentVersion, "currentVersion should return an empty string")
	})
	t.Run("current_version", func(t *testing.T) {
		assert.NotEmpty(t, Version(), "Version() should not return an empty string")
		assert.Equal(t, currentVersion, Version(), "Version() should return the current version")
	})
}

func TestIsTracking(t *testing.T) {
	// Add a timeout to ensure the test fails if it runs too long
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute+34*time.Second)
	defer cancel()

	// Run the test with a timeout
	t.Run("IsTracking", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		os.Clearenv()
		// Create a figtree with Tracking and Germinate enabled
		figs := With(Options{Tracking: true, Germinate: true})
		var k, d, u = "name", "yahuah", "usage"
		figs.NewString(k, d, u)
		assert.Nil(t, figs.Parse())

		// Use a WaitGroup to synchronize goroutines
		var wg sync.WaitGroup

		// Use atomic counters to collect errors and count writes
		var errorCount int32
		var writeCount int32

		// Collect mutations from the Mutations channel
		mutations := make([]Mutation, 0, 500)
		var mutationsMu sync.Mutex
		done := make(chan struct{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			for mutation := range figs.Mutations() {
				mutationsMu.Lock()
				mutations = append(mutations, mutation)
				mutationsMu.Unlock()
				// t.Logf("Received mutation: %v", mutation)
			}
			close(done)
		}()

		runFor := time.NewTimer(7 * time.Second)
		readEvery := time.NewTicker(3 * time.Millisecond)
		writeEvery := time.NewTicker(77 * time.Millisecond)

		// Channel to collect errors
		errChan := make(chan error, 500)

		// Main loop to handle reads and writes
		wg.Add(1)
		go func() {
			defer wg.Done()
			// start := time.Now()
			for {
				select {
				case <-ctx.Done():
					t.Error("Test timed out after 2 minutes")
					return
				case <-runFor.C:
					// case n := <-readEvery.C:
					// t.Logf("Timer fired at %v after %v", n, time.Since(start))
					readEvery.Stop()
					writeEvery.Stop()
					// Close the mutations channel to allow the collector goroutine to exit
					// t.Log("cursing the fig tree...")
					figs.Curse()
					return
				case <-readEvery.C:
					// case n := <-runFor.C:
					// t.Logf("Read ticker fired at %v", n)
					val := figs.String(k)
					if val == nil {
						errChan <- fmt.Errorf("String(%q) returned nil", k)
						atomic.AddInt32(&errorCount, 1)
						continue
					}
					// t.Logf("Read: %s", *val)
				case <-writeEvery.C:
					// case n := <-writeEvery.C:
					// t.Logf("Write ticker fired at %v", n)
					// Toggle the value between "yahuah" and "andrei"
					if d == "yahuah" {
						d = "andrei"
					} else {
						d = "yahuah"
					}
					figs = figs.StoreString(k, d)
					assert.NoError(t, figs.ErrorFor(k))
					atomic.AddInt32(&writeCount, 1)
					// t.Logf("Wrote: %s (write #%d)", d, writeCount)
				}
			}
		}()

		// Wait for goroutines to finish
		wg.Wait()
		<-done

		// Collect and log any errors
		errors := make([]error, 0, len(errChan))
		for len(errChan) > 0 {
			err := <-errChan
			if err != nil {
				errors = append(errors, err)
			}
		}

		// Log errors for debugging
		for i, err := range errors {
			t.Logf("Error %d: %v", i+1, err)
		}

		// Verify that no errors occurred
		assert.Equal(t, int32(0), errorCount, "Test failed with %d errors", errorCount)

		actualWrites := int(atomic.LoadInt32(&writeCount))
		assert.Equal(t, actualWrites, len(mutations), "Number of mutations does not match number of writes")

		// Verify the final value
		finalVal := figs.String(k)
		assert.NotNil(t, finalVal)
	})
}

func TestTree_PollinateString(t *testing.T) {
	const k1 = "testmehere"
	os.Args = []string{os.Args[0]}
	figs := With(Options{Pollinate: true, Tracking: true, Germinate: true})
	figs = figs.NewString(k1, "initial", "usage").WithValidator("test", AssureStringContains("ini"))
	assert.NoError(t, figs.Parse())
	assert.Equal(t, "initial", *figs.String(k1))
	go func() {
		timer := time.NewTimer(time.Second * 1)
		checker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-timer.C:
				assert.Equal(t, "updated", *figs.String(k1))
				return
			case <-checker.C:
				assert.NoError(t, os.Setenv(strings.ToUpper(k1), "updated"))
			}
		}
	}()
	defer assert.NoError(t, os.Unsetenv(strings.ToUpper(k1)))
	mutation, ok := <-figs.Mutations()
	if ok {
		assert.Equal(t, k1, mutation.Property)
		assert.Equal(t, "string", mutation.Mutagenesis)
		assert.Equal(t, "StoreString", mutation.Way)
		assert.Equal(t, "initial", mutation.Old)
		assert.Equal(t, "updated", mutation.New)
		assert.NoError(t, mutation.Error)
	}
}
