package figtree

func (tree *figTree) WithAlias(name, alias string) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, exists := tree.aliases[alias]; exists {
		return
	}
	tree.aliases[alias] = name
	fig := tree.figs[name]
	if fig == nil {
		return
	}
	switch fig.Flesh.Flesh.(type) {
	case *ListFlag:

	}
	tree.flagSet.Var(&fig.Flesh, alias, fig.description)
}
