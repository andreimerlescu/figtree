package figtree

import (
	"strings"
)

// ListFlag stores values in a list type configurable
type ListFlag struct {
	values *[]string
}

// String returns the slice of strings using strings.Join
func (l *ListFlag) String() string {
	if l.values == nil {
		return ""
	}
	return strings.Join(*l.values, ",")
}

// Set unpacks a comma separated value argument and appends items to the list of []string
func (l *ListFlag) Set(value string) error {
	if l.values == nil {
		l.values = &[]string{}
	}
	items := strings.Split(value, ",")
	*l.values = append(*l.values, items...)
	return nil
}
