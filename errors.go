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
