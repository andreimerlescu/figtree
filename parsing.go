package figtree

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Parsing Configuration

func (tree *figTree) useValue(value *Value, err error) *Value {
	if err != nil {
		log.Printf("useValue caught err: %v", err)
	}
	return value
}

func (tree *figTree) from(name string) (*Value, error) {
	flagName := strings.ToLower(name)
	for alias, realname := range tree.aliases {
		if strings.EqualFold(alias, name) {
			flagName = realname
		}
	}
	valueAny, ok := tree.values.Load(flagName)
	if !ok {
		return nil, errors.New("no value for " + flagName)
	}
	value, ok := valueAny.(*Value)
	if !ok {
		return nil, errors.New("value for " + flagName + " is not a Value")
	}
	if value.Err != nil {
		return nil, value.Err
	}
	if value.Mutagensis == "" {
		value.Mutagensis = tree.MutagenesisOf(value)
	}
	return value, nil
}

var re *regexp.Regexp

func init() {
	re = regexp.MustCompile(`(\d+)([ylwdhms])`)
}

func ParseCustomDuration(input string) (time.Duration, error) {
	// Define the mapping of units to their respective durations
	unitMap := map[rune]time.Duration{
		'y': 365 * 24 * time.Hour, // Approximate year (non-leap)
		'l': 30 * 24 * time.Hour,  // Approximate month (30 days)
		'w': 7 * 24 * time.Hour,   // Week
		'd': 24 * time.Hour,       // Day
		'h': time.Hour,            // Hour
		'm': time.Minute,          // Minute
		's': time.Second,          // Second
	}

	// Regular expression to match number-unit pairs (e.g., "5h", "32m")
	matches := re.FindAllStringSubmatch(input, -1)

	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration format: %s", input)
	}

	var totalDuration time.Duration
	for _, match := range matches {
		// match[1] is the number, match[2] is the unit
		num, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, fmt.Errorf("invalid number in duration: %s", match[1])
		}

		unit := rune(match[2][0])
		duration, exists := unitMap[unit]
		if !exists {
			return 0, fmt.Errorf("invalid unit in duration: %s", match[2])
		}

		totalDuration += time.Duration(num) * duration
	}

	return totalDuration, nil
}

func (tree *figTree) checkFigErrors() error {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	for name, fig := range tree.figs {
		if fig.Error != nil {
			return fig.Error
		}
		value, err := tree.from(name)
		if err != nil {
			fig.Error = errors.Join(fig.Error, err)
			return fig.Error
		}
		if !strings.EqualFold(string(fig.Mutagenesis), string(value.Mutagensis)) {
			return fmt.Errorf("invalid Mutagenesis (Type) for flag -%s", name)
		}
		switch fig.Mutagenesis {
		case tString:
			_, e := toString(value)
			if e != nil {
				er := value.Assign(zeroString)
				if er != nil {
					e = errors.Join(e, er)
				}
				return ErrInvalidValue{name, e}
			}
		case tInt:
			_, e := toInt(value)
			if e != nil {
				er := value.Assign(zeroInt)
				if er != nil {
					e = errors.Join(e, er)
				}
				return ErrInvalidValue{name, e}
			}
		case tFloat64:
			_, e := toFloat64(value)
			if e != nil {
				er := value.Assign(zeroFloat64)
				if er != nil {
					e = errors.Join(e, er)
				}
				return ErrInvalidValue{name, e}
			}

		case tBool:
			_, e := toBool(value)
			if e != nil {
				er := value.Assign(zeroBool)
				if er != nil {
					e = errors.Join(e, er)
				}
				return ErrInvalidValue{name, e}
			}
		case tInt64, tUnitDuration, tDuration:
			if value.Mutagensis == tUnitDuration || value.Mutagensis == tDuration {
				vf := value.Flesh()
				var val time.Duration
				var err error
				if vf != nil {
					if _, ok := vf.AsIs().(time.Duration); !ok {
						val, err = ParseCustomDuration(vf.ToString())
						if err == nil {
							err = value.Assign(val)
							if err != nil {
								return ErrInvalidValue{name, err}
							}
							continue
						}
					} else {
						continue
					}
				}

			}
			_, e := toInt64(value)
			if e != nil {
				er := value.Assign(zeroInt64)
				if er != nil {
					e = errors.Join(e, er)
				}
				return ErrInvalidValue{name, e}
			}
		case tMap:
			_, e := toStringMap(value.Value)
			if e != nil {
				er := value.Assign(zeroMap)
				if er != nil {
					e = errors.Join(e, er)
				}
				return ErrInvalidValue{name, e}
			}
		case tList:
			_, e := toStringSlice(value.Value)
			if e != nil {
				er := value.Assign(zeroList)
				if er != nil {
					e = errors.Join(e, er)
				}
				return ErrInvalidValue{name, e}
			}
		default:
			return ErrInvalidValue{name, fmt.Errorf("unknown flag type")}
		}
	}
	return nil
}

// Parse uses figTree.flagSet to run flag.Parse() on the registered figs and returns nil for validated results
func (tree *figTree) Parse() (err error) {
	preloadErr := tree.preLoadOrParse()
	if preloadErr != nil {
		return preloadErr
	}
	if !tree.HasRule(RuleNoFlags) {
		tree.activateFlagSet()
		args := os.Args[1:]
		if tree.filterTests {
			args = filterTestFlags(args)
		}
		err = tree.flagSet.Parse(args)
		if err != nil {
			err2 := tree.checkFigErrors()
			if err2 != nil {
				err = errors.Join(err, err2)
			}
			return err
		}
		err = tree.loadFlagSet()
		if err != nil {
			return err
		}
		tree.readEnv()
		err = tree.applyWithered()
		if err != nil {
			return err
		}
		return tree.validateAll()
	}
	tree.readEnv()
	err = tree.applyWithered()
	if err != nil {
		return err
	}
	return tree.validateAll()
}

func (tree *figTree) applyWithered() error {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for name, fig := range tree.figs {
		if fig == nil {
			continue
		}
		value := tree.useValue(tree.from(name))
		if value.Mutagensis == tMap && PolicyMapAppend {
			vm := value.Flesh().ToMap()
			unique := make(map[string]string)
			withered := tree.withered[name]
			for k, v := range vm {
				unique[k] = v
			}
			for k, v := range withered.Value.Flesh().ToMap() {
				unique[k] = v
			}
			err := value.Assign(unique)
			if err != nil {
				return ErrInvalidValue{name, err}
			}
			tree.values.Store(name, value)
		}
		if value.Mutagensis == tList && PolicyListAppend {
			vl, e := toStringSlice(value.Value)
			if e != nil {
				return fmt.Errorf("failed toStringSlice: %w", e)
			}
			unique := make(map[string]struct{})
			for _, v := range vl {
				unique[v] = struct{}{}
			}
			withered := tree.withered[name]
			for _, w := range withered.Value.Flesh().ToList() {
				unique[w] = struct{}{}
			}
			var result []string
			for k, _ := range unique {
				result = append(result, k)
			}
			sort.Strings(result)
			err := value.Assign(result)
			if err != nil {
				return ErrInvalidValue{name, err}
			}
			tree.values.Store(name, value)
		}
	}
	return nil
}

// ParseFile will check if filename is set and run loadFile on it.
func (tree *figTree) ParseFile(filename string) (err error) {
	preloadErr := tree.preLoadOrParse()
	if preloadErr != nil {
		return preloadErr
	}
	if !tree.HasRule(RuleNoFlags) {
		tree.activateFlagSet()
		args := os.Args[1:]
		if tree.filterTests {
			args = filterTestFlags(args)
			err = tree.flagSet.Parse(args)
			if err != nil {
				err2 := tree.checkFigErrors()
				if err2 != nil {
					err = errors.Join(err, err2)
				}
			}
		} else {
			err = tree.flagSet.Parse(args)
		}
		if err != nil {
			return err
		}
	}
	err = tree.loadFlagSet()
	if err != nil {
		return err
	}
	if filename != "" {
		return tree.loadFile(filename)
	}
	tree.readEnv()
	return tree.validateAll()
}
