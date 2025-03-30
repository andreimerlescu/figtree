package figtree

import (
	"fmt"
	"time"
)

// WithValidator adds a validator to an int flag
//
// Example:
//
//			figs := figtree.With(Options{Pollinate: true, Tracking: true, IgnoreEnvironment: true})
//			figs.NewString("domain", "", "domain name")
//	 	figs.WithValidator("domain", figtree.AssureStringHasPrefix("https://"))
//			err := figs.Parse() // if you're NOT using ./config.yaml
//			OR err := figs.Load() // if you're using ./config.yaml to populate domain
func (tree *Tree) WithValidator(name string, validator func(interface{}) error) Fruit {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if fig, ok := tree.figs[name]; ok {
		if fig.Validators == nil {
			fig.Validators = make([]ValidatorFunc, 0)
		}
		fig.Validators = append(fig.Validators, validator)
		tree.figs[name] = fig
	}
	return tree
}

// validateAll looks at Fig ValidatorFunc and returns the error if it fails otherwise it calls Tree.runCallbacks()
func (tree *Tree) validateAll() error {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	for name, fruit := range tree.figs {
		if fruit.Error != nil {
			return fruit.Error
		}
		for _, validator := range fruit.Validators {
			if fruit != nil && validator != nil {
				var val interface{}
				switch v := fruit.Flesh.(type) {
				case *int:
					val = *v
				case *int64:
					val = *v
				case *float64:
					val = *v
				case *string:
					val = *v
				case *bool:
					val = *v
				case *time.Duration:
					val = *v
				case *ListFlag:
					val = *v.values
				case *MapFlag:
					val = *v.values
				}
				if err := validator(val); err != nil {
					return fmt.Errorf("validation failed for %s: %v", name, err)
				}
			}
		}
	}
	return tree.runCallbacks(CallbackAfterVerify)
}
