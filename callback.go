package figtree

import (
	"errors"
)

// WithCallback allows you to assign a slice of CallbackFunc to a Fig attached to a Tree.
// A callback is executed when the value of the configurable CHANGES
//
// Example:
//
//	figs := With(Options{Pollinate: true, Tracking, true, Harvest: 1776})
//	figs.NewString("domain", "", "domain name")
//	figs.WithCallback("domain", figtree.CallbackAfterVerify, func(value interface{}) error {
//		var sv string
//		switch v := value.(type) {
//			case *string:
//				sv = *v
//			case string:
//				sv = v
//			default:
//				return fmt.Errorf("invalid type want %T, got %T", sv, v)
//		}
//		// do something with the sv domain after its been verified
//	})
func (tree *Tree) WithCallback(name string, whenCallback CallbackAfter, runThis CallbackFunc) Fruit {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	fruit, exists := tree.figs[name]
	if !exists || fruit == nil {
		tree.mu.Unlock()
		tree.Resurrect(name)
		tree.mu.Lock()
		fruit = tree.figs[name]
	}
	if fruit == nil {
		return tree
	}
	fruit.Callbacks = append(fruit.Callbacks, Callback{
		CallbackAfter: whenCallback,
		CallbackFunc:  runThis,
	})
	tree.figs[name] = fruit
	return tree
}

// runCallbacks inspects each fig fruit on the tree and executes runCallbacks() against the fig fruit
func (tree *Tree) runCallbacks(callbackOn CallbackAfter) error {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	for _, fig := range tree.figs {
		if len(fig.Callbacks) == 0 {
			continue
		}
		err := fig.runCallbacks(callbackOn)
		if err != nil {
			return err
		}
	}
	return nil
}

// runCallbacks will take each registered callback and run it against the fig fruit
func (fig *Fig) runCallbacks(callbackOn CallbackAfter) error {
	if fig.Error != nil {
		return fig.Error
	}
	errs := make([]error, len(fig.Callbacks))
	for _, callback := range fig.Callbacks {
		if callback.CallbackAfter == callbackOn {
			err := callback.CallbackFunc(fig.Flesh)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
