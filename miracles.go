package figtree

// Mutations returns a receiver channel of Mutation data
func (tree *figTree) Mutations() <-chan Mutation {
	return tree.mutationsCh
}

// Recall is when you bring the mutations channel back to life and you unlock making further changes to the fig *figTree
func (tree *figTree) Recall() {
	tree.angel.Store(false)
	tree.mutationsCh = make(chan Mutation, tree.harvest)
	tree.tracking = true
}

// Curse is when you lock the fig *figTree from further changes, stop tracking and close the channel
func (tree *figTree) Curse() {
	tree.angel.Store(true)
	tree.tracking = false
	close(tree.mutationsCh)
}

// FigFlesh returns a Flesh interface to the Value on the figTree
func (tree *figTree) FigFlesh(name string) Flesh {
	valueAny, exists := tree.values.Load(name)
	if !exists {
		return nil
	}
	value, ok := valueAny.(*Value)
	if !ok {
		return nil
	}
	return value.Flesh()
}
