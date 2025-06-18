package figtree

import (
	"strings"
)

// ListFlag stores values in a list type configurable
type ListFlag struct {
	values *[]string
}

func (tree *figTree) ListValues(name string) []string {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	fruit, exists := tree.figs[name]
	if !exists {
		return []string{}
	}
	return fruit.Value.Flesh().ToList()
}

func (l *ListFlag) Values() []string {
	if l.values == nil {
		return []string{}
	}
	return *l.values
}

// String returns the slice of strings using strings.Join
func (l *ListFlag) String() string {
	if l.values == nil {
		return ""
	}
	return strings.Join(*l.values, ",")
}

// PolicyListAppend will apply ListFlag.Set to the list of values and not append to any existing values in the ListFlag
var PolicyListAppend bool = false

// Set unpacks a comma separated value argument and appends items to the list of []string
func (l *ListFlag) Set(value string) error {
	if l.values == nil {
		l.values = &[]string{}
	}
	items := strings.Split(value, ",")
	if PolicyListAppend {
		*l.values = append(*l.values, items...)
	} else {
		*l.values = items
	}
	return nil
}
