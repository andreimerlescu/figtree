package figtree

import (
	"flag"
	"fmt"
)

func (tree *figTree) WithAlias(name, alias string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
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
	flag.Var(value, alias, "Alias of -"+name)
	return tree
}
