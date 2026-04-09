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

	// Guard: alias already registered
	if existing, exists := tree.aliases[alias]; exists {
		if existing != name {
			tree.problems = append(tree.problems,
				fmt.Errorf("WithAlias: alias -%s already maps to -%s, cannot remap to -%s", alias, existing, name))
		}
		// idempotent: same alias→name pair is a no-op, not an error
		return tree
	}

	// Guard: canonical fig must exist
	if _, exists := tree.figs[name]; !exists {
		tree.problems = append(tree.problems,
			fmt.Errorf("WithAlias: no fig named -%s", name))
		return tree
	}

	// Guard: alias must not shadow an existing fig name
	if _, exists := tree.figs[alias]; exists {
		tree.problems = append(tree.problems,
			fmt.Errorf("WithAlias: alias -%s conflicts with existing fig name", alias))
		return tree
	}

	// Guard: alias must not shadow an existing flag (covers both figs and
	// any flags registered outside of figtree, e.g. via flagSet.Var directly)
	if tree.flagSet.Lookup(alias) != nil {
		tree.problems = append(tree.problems,
			fmt.Errorf("WithAlias: alias -%s conflicts with existing flag", alias))
		return tree
	}

	// Guard: underlying value must be retrievable and correctly typed
	ptr, ok := tree.values.Load(name)
	if !ok {
		tree.problems = append(tree.problems,
			fmt.Errorf("WithAlias: no value found for -%s", name))
		return tree
	}
	value, ok := ptr.(*Value)
	if !ok {
		tree.problems = append(tree.problems,
			fmt.Errorf("WithAlias: value for -%s is %T, expected *Value", name, ptr))
		return tree
	}

	// All validations passed — register the alias
	tree.aliases[alias] = name
	tree.flagSet.Var(value, alias, "Alias of -"+name)
	return tree
}
