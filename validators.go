package figtree

import (
	"fmt"
	"time"
)

// WithValidator adds a validator to an int flag
func (fig *Tree) WithValidator(name string, validator func(interface{}) error) Fruit {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	if def, ok := fig.figs[name]; ok {
		if def.Validator != nil {
			fig.figs[name].Error = fmt.Errorf("validator for fig already exists, overwriting old validator")
		}
		def.Validator = validator
		fig.figs[name] = def
	}
	return fig
}

// validateAll looks at Fig ValidatorFunc and returns the error if it fails
func (fig *Tree) validateAll() error {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	for name, fruit := range fig.figs {
		if fruit.Error != nil {
			return fruit.Error
		}
		if fruit != nil && fruit.Validator != nil {
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
			if err := fruit.Validator(val); err != nil {
				return fmt.Errorf("validation failed for %s: %v", name, err)
			}
		}
	}
	return nil
}
