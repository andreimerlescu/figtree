package figtree

import (
	"fmt"
	"strings"
)

// resolveName returns the canonical fig name for a given name or alias.
// Callers must hold tree.mu (read or write) before calling this.
func (tree *figTree) resolveName(name string) string {
	name = strings.ToLower(name)
	if canonical, exists := tree.aliases[name]; exists {
		return canonical
	}
	return name
}

func (tree *figTree) Problems() []error {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return append([]error(nil), tree.problems...)
}

func (tree *figTree) WithAlias(name, alias string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	name = strings.ToLower(name)
	alias = strings.ToLower(alias)
	if _, exists := tree.aliases[alias]; exists {
		tree.problems = append(tree.problems, fmt.Errorf("WithAlias: alias exists -%s", name))
		return tree
	}
	if _, exists := tree.figs[name]; !exists {
		tree.problems = append(tree.problems, fmt.Errorf("WithAlias: no fig named -%s", name))
		return tree
	}
	ptr, ok := tree.values.Load(name)
	if !ok {
		tree.problems = append(tree.problems, fmt.Errorf("WithAlias: no value found for -%s", name))
		return tree
	}
	value, ok := ptr.(*Value)
	if !ok {
		tree.problems = append(tree.problems, fmt.Errorf("WithAlias: failed to cast value for -%s", name))
		return tree
	}
	if _, exists := tree.figs[alias]; exists {
		tree.problems = append(tree.problems, fmt.Errorf("WithAlias: alias -%s conflicts with existing fig name", alias))
		return tree
	}
	if tree.flagSet.Lookup(alias) != nil {
		tree.problems = append(tree.problems, fmt.Errorf("WithAlias: alias -%s conflicts with existing flag", alias))
		return tree
	}
	tree.aliases[alias] = name // only register after all validations pass
	tree.flagSet.Var(value, alias, "Alias of -"+name)
	return tree
}
