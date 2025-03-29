package figtree

import (
	"fmt"
)

// ErrorFor returns an error on a given name if one exists
func (fig *Tree) ErrorFor(name string) error {
	fig.mu.RLock()
	defer fig.mu.RUnlock()
	fruit, exists := fig.figs[name]
	if !exists || fruit == nil {
		return fmt.Errorf("no fig named %s", name)
	}
	return fruit.Error
}
