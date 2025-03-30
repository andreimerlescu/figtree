package figtree

import (
	"os"
)

// Parsing Configuration

// Parse uses Tree.flagSet to run flag.Parse() on the registered figs and returns nil for validated results
func (tree *Tree) Parse() (err error) {
	tree.activateFlagSet()
	args := os.Args[1:]
	if tree.filterTests {
		args = filterTestFlags(args)
		err = tree.flagSet.Parse(args)
	} else {
		err = tree.flagSet.Parse(args)
	}
	if err != nil {
		return err
	}
	tree.readEnv()
	return tree.validateAll()
}

// ParseFile will check if filename is set and run loadFile on it.
func (tree *Tree) ParseFile(filename string) (err error) {
	tree.activateFlagSet()
	args := os.Args[1:]
	if tree.filterTests {
		args = filterTestFlags(args)
		err = tree.flagSet.Parse(args)
	} else {
		err = tree.flagSet.Parse(args)
	}
	if err != nil {
		return err
	}
	if filename != "" {
		return tree.loadFile(filename)
	}
	tree.readEnv()
	return tree.validateAll()
}
