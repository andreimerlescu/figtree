package figtree

import (
	"fmt"
	"strings"
)

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
