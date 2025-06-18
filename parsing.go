package figtree

import (
	"flag"
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
		for _, fig := range tree.figs {
			fv := tree.flagSet.Lookup(fig.name)
			v := fv.Value.String()
			e := fig.Value.Set(v)
			if e != nil {
				err = e
			}
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
	tree.flagSet.VisitAll(func(f *flag.Flag) {
		if fig, exists := tree.figs[f.Name]; exists {
			fig.Value = Value{
				Value:      f.Value,
				Mutagensis: tree.MutagenesisOf(fig.Value),
			}
		}
	})
	if filename != "" {
		return tree.loadFile(filename)
	}
	tree.readEnv()
	return tree.validateAll()
}
