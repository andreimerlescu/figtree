package figtree

import (
	"errors"
)

// WithCallback allows you to assign a slice of CallbackFunc to a figFruit attached to a figTree.
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
func (tree *figTree) WithCallback(name string, whenCallback CallbackWhen, runThis CallbackFunc) Plant {
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
		CallbackWhen: whenCallback,
		CallbackFunc: runThis,
	})
	tree.figs[name] = fruit
	return tree
}

// runCallbacks inspects each fig fruit on the tree and executes runCallbacks() against the fig fruit
func (tree *figTree) runCallbacks(callbackOn CallbackWhen) error {
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
func (fig *figFruit) runCallbacks(callbackOn CallbackWhen) error {
	if fig.Error != nil {
		return fig.Error
	}
	errs := make([]error, len(fig.Callbacks))
	for _, callback := range fig.Callbacks {
		if callback.CallbackWhen == callbackOn {
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
