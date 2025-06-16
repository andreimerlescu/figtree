package figtree

import (
	"fmt"
)

// ErrorFor returns an error on a given name if one exists
func (tree *figTree) ErrorFor(name string) error {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	name = tree.resolveName(name)
	fruit, exists := tree.figs[name]
	if !exists || fruit == nil {
		return fmt.Errorf("no tree named %s", name)
	}
	return fruit.Error
}

func (fig *figFruit) Unwrap() error {
	return fig.Error
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

func (e ErrInvalidValue) Unwrap() error {
	return e.Err
}

const (
	ErrWayBeBelow      string = "be below"
	ErrWayBeAbove      string = "be above"
	ErrWayBeBetweenFmt string = "be between %v and %v"
	ErrWayBePositive   string = "be positive"
	ErrWayBeNegative   string = "be negative"
	ErrWayBeNotNaN     string = "not be NaN"
)

type ErrValue struct {
	Way   string
	Value any
	Than  any
}

func (e ErrValue) Error() string {
	if e.Than != nil {
		return fmt.Sprintf("invalid value ; must be %s than %v ; got %v", e.Way, e.Value, e.Than)
	}
	return fmt.Sprintf("invalid value ; must be %s ; got %v", e.Way, e.Value)
}

type ErrLoadFailure struct {
	What string
	Err  error
}

func (e ErrLoadFailure) Error() string {
	return fmt.Sprintf("failed to load %s: %s", e.What, e.Err.Error())
}

func (e ErrLoadFailure) Unwrap() error {
	return e.Err
}

type ErrValidationFailure struct {
	Err error
}

func (e ErrValidationFailure) Error() string {
	return fmt.Sprintf("failed to validateAll with err: %v", e.Err)
}

func (e ErrValidationFailure) Unwrap() error {
	return e.Err
}
