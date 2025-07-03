package figtree

import (
	"fmt"
)

// ErrorFor returns an error on a given name if one exists
func (tree *figTree) ErrorFor(name string) error {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	fruit, exists := tree.figs[name]
	if !exists || fruit == nil {
		return fmt.Errorf("no tree named %s", name)
	}
	return fruit.Error
}

type ErrInvalidType struct {
	Wanted Mutagenesis
	Got    any
}

func (e ErrInvalidType) Error() string {
	return fmt.Sprintf("invalid type ; got %s ; wanted %s", e.Got, e.Wanted.Kind())
}

type ErrConversion struct {
	From Mutagenesis
	To   Mutagenesis
	Got  any
}

func (e ErrConversion) Error() string {
	return fmt.Sprintf("failed to convert %v (type %T) type %s into %s", e.Got, e.Got, e.From.Kind(), e.To.Kind())
}

type ErrInvalidValue struct {
	Name string
	Err  error
}

func (e ErrInvalidValue) Error() string {
	return fmt.Sprintf("invalid value for flag -%s: %s", e.Name, e.Err.Error())
}
