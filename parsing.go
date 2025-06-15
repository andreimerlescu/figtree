package figtree

import (
	"os"
)

// Parsing Configuration

// Parse uses figTree.flagSet to run flag.Parse() on the registered figs and returns nil for validated results
func (tree *figTree) Parse() (err error) {
	if !tree.HasRule(RuleNoFlags) {
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
	}
	tree.readEnv()
	return tree.validateAll()
}

// ParseFile will check if filename is set and run loadFile on it.
func (tree *figTree) ParseFile(filename string) (err error) {
	if !tree.HasRule(RuleNoFlags) {
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
	}
	if filename != "" {
		return tree.loadFile(filename)
	}
	tree.readEnv()
	return tree.validateAll()
}
