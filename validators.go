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
func (tree *figTree) WithValidator(name string, validator func(interface{}) error) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if fig, ok := tree.figs[name]; ok {
		if fig.Validators == nil {
			fig.Validators = make([]FigValidatorFunc, 0)
		}
		fig.Validators = append(fig.Validators, validator)
		tree.figs[name] = fig
	}
	return tree
}

// WithValidators uses WithValidator to pass multiple Assure into a type
// Example:
//
//			figs := figtree.Grow()
//			figs.NewString("name", "", "Your name")
//			figs.WithValidators("name",
//				figtree.AssureStringNotEmpty,
//				figtree.AssureStringNotContains("god"),
//	         figtree.AssureStringLengthGreaterThan(2))
func (tree *figTree) WithValidators(name string, validators ...func(interface{}) error) Plant {
	for _, v := range validators {
		tree.WithValidator(name, v)
	}
	return tree
}

// validateAll looks at figFruit FigValidatorFunc and returns the error if it fails otherwise it calls figTree.runCallbacks()
func (tree *figTree) validateAll() error {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	err := tree.runCallbacks(CallbackBeforeVerify)
	if err != nil {
		return err
	}
	for name, fruit := range tree.figs {
		if fruit.Error != nil {
			return fruit.Error
		}
		for _, validator := range fruit.Validators {
			if fruit != nil && validator != nil {
				var val interface{}
				switch v := fruit.Flesh.Flesh.(type) {
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

// makeStringValidator creates a validator for string-based checks.
func makeStringValidator(check func(string) bool, errFormat string) FigValidatorFunc {
	return func(value interface{}) error {
		v := figFlesh{value}
		if !v.IsString() {
			return fmt.Errorf("expected string, got %T", value)
		}
		s := v.ToString()
		if !check(s) {
			return fmt.Errorf(errFormat, s)
		}
		return nil
	}
}
