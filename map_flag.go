package figtree

import (
	"fmt"
	"strings"
)

// MapFlag stores values in a map type configurable
type MapFlag struct {
	values *map[string]string
}

func (tree *figTree) MapKeys(name string) []string {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	fruit, exists := tree.figs[name]
	if !exists {
		return []string{}
	}
	switch v := fruit.Value.Value.(type) {
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
		keys := make([]string, 0, len(*v.values))
		for k := range *v.values {
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
	for key := range *m.values {
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
	for k, v := range *m.values {
		entries = append(entries, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(entries, ",")
}

var PolicyMapAppend = false

// Set accepts a value like KEY=VALUE,KEY=VALUE,KEY=VALUE to override map values
func (m *MapFlag) Set(value string) error {
	if m.values == nil {
		m.values = &map[string]string{}
	}
	if !PolicyMapAppend {
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
