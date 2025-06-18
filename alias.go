package figtree

import "flag"

func (tree *figTree) WithAlias(name, alias string) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, exists := tree.aliases[alias]; exists {
		return tree
	}
	tree.aliases[alias] = name
	fig, ok := tree.figs[name]
	if !ok {
		return tree
	}
	flag.Var(&fig.Value, alias, "Alias of -"+name)
	return tree
}
