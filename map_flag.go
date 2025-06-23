package figtree

import (
	"fmt"
	"maps"
	"os"
	"strings"
)

// MapFlag stores values in a map type configurable
type MapFlag struct {
	values map[string]string
}

func (tree *figTree) MapKeys(name string) []string {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	originalName := strings.Clone(name) // capture original
	defer func() {
		name = originalName // return value to original
	}()
	name = strings.ToLower(name)
	fruit, exists := tree.figs[name]
	if !exists {
		return []string{}
	}
	if fruit == nil {
		return []string{}
	}
	_value, err := tree.from(name)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "tree.from(%s) return err: %v", name, err.Error())
		return []string{}
	}
	switch v := _value.Value.(type) {
	case nil:
		return []string{}
	case map[string]string:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		return keys
	case *map[string]string:
		keys := make([]string, 0, len(*v))
		for k := range *v {
			keys = append(keys, k)
		}
		return keys
	case Value:
		keys := make([]string, 0, len(*v.Value.(*map[string]string)))
		for k := range *v.Value.(*map[string]string) {
			keys = append(keys, k)
		}
		return keys
	case *MapFlag:
		keys := make([]string, 0, len(v.values))
		for k := range v.values {
			keys = append(keys, k)
		}
		return keys
	default:
		return []string{}
	}
}

func (m *MapFlag) Keys() []string {
	if m.values == nil {
		return []string{}
	}
	var keys []string
	for key := range m.values {
		keys = append(keys, key)
	}
	return keys
}

// String returns the map[string]string as string=string,string=string,...
func (m *MapFlag) String() string {
	if m.values == nil {
		return ""
	}
	var entries []string
	for k, v := range m.values {
		entries = append(entries, fmt.Sprintf("%s%s%s", k, MapKeySeparator, v))
	}
	return strings.Join(entries, ",")
}

var PolicyMapAppend = false

// Set accepts a value like KEY=VALUE,KEY=VALUE,KEY=VALUE to override map values
func (m *MapFlag) Set(value string) error {
	if m.values == nil || !PolicyMapAppend {
		m.values = map[string]string{}
	}
	existing := maps.Clone(m.values)
	if PolicyMapAppend {
		for k, v := range existing {
			m.values[k] = v
		}
	}
	adding := make(map[string]string)
	pairs := strings.Split(value, MapSeparator)
	for _, pair := range pairs {
		kv := strings.SplitN(pair, MapKeySeparator, 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid map item: %s", pair)
		}
		adding[kv[0]] = kv[1]
	}
	for k, v := range adding {
		m.values[k] = v
	}
	return nil
}
