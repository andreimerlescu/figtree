package figtree

import (
	"fmt"
	"strings"
)

func (tree *figTree) WithAlias(name, alias string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	name = strings.ToLower(name)
	alias = strings.ToLower(alias)
	if _, exists := tree.aliases[alias]; exists {
		return tree
	}
	tree.aliases[alias] = name
	ptr, ok := tree.values.Load(name)
	if !ok {
		fmt.Println("failed to load -" + name + " value")
		return tree
	}
	value, ok := ptr.(*Value)
	if !ok {
		fmt.Println("failed to cast -" + name + " value")
		return tree
	}
	tree.flagSet.Var(value, alias, "Alias of -"+name)
	return tree
}
