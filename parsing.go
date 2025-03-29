package figtree

import (
	"os"
)

// Parsing Configuration

// Parse uses Tree.flagSet to run flag.Parse() on the registered figs and returns nil for validated results
func (fig *Tree) Parse() (err error) {
	fig.activateFlagSet()
	args := os.Args[1:]
	if fig.filterTests {
		args = filterTestFlags(args)
		err = fig.flagSet.Parse(args)
	} else {
		err = fig.flagSet.Parse(args)
	}
	if err != nil {
		return err
	}
	fig.readEnv()
	return fig.validateAll()
}

// ParseFile will check if filename is set and run loadFile on it.
func (fig *Tree) ParseFile(filename string) (err error) {
	fig.activateFlagSet()
	args := os.Args[1:]
	if fig.filterTests {
		args = filterTestFlags(args)
		err = fig.flagSet.Parse(args)
	} else {
		err = fig.flagSet.Parse(args)
	}
	if err != nil {
		return err
	}
	if filename != "" {
		return fig.loadFile(filename)
	}
	fig.readEnv()
	return fig.validateAll()
}
