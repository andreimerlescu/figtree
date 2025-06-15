package figtree

func (tree *figTree) WithAlias(name, alias string) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if _, exists := tree.aliases[alias]; exists {
		return
	}
	tree.aliases[alias] = name
}
