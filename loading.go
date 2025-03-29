package figtree

import (
	"fmt"
	"os"
	"path/filepath"

	check "github.com/andreimerlescu/checkfs"
	"github.com/andreimerlescu/checkfs/file"
)

// Reload will readEnv on each flag in the configurable package
func (fig *Tree) Reload() error {
	fig.readEnv()
	return fig.validateAll()
}

// Load uses the EnvironmentKey and the DefaultJSONFile, DefaultYAMLFile, and DefaultINIFile to run ParseFile if it exists
func (fig *Tree) Load() (err error) {
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

	files := []string{
		os.Getenv(EnvironmentKey),
		fig.ConfigFilePath,
		ConfigFilePath,
		filepath.Join(".", DefaultJSONFile),
		filepath.Join(".", DefaultINIFile),
	}
	for i := 0; i < len(files); i++ {
		f := files[i]
		if f == "" {
			continue
		}
		if err := check.File(f, file.Options{Exists: true}); err == nil {
			if err := fig.loadFile(f); err != nil {
				return fmt.Errorf("failed to load %s: %w", f, err)
			}
		}
	}

	fig.readEnv()
	return fig.validateAll()
}

// LoadFile accepts a path and uses it to populate the Tree
func (fig *Tree) LoadFile(path string) (err error) {
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
	var loadErr error
	if loadErr = check.File(path, file.Options{Exists: true}); loadErr == nil {
		if err2 := fig.loadFile(path); err2 != nil {
			return fmt.Errorf("failed to loadFile %s: %w", path, err2)
		}
		fig.readEnv()
		err3 := fig.validateAll()
		if err3 != nil {
			return fmt.Errorf("failed to validateAll: %w", err3)
		}
		return nil
	}
	fig.readEnv()
	err3 := fig.validateAll()
	if err3 != nil {
		return fmt.Errorf("failed to validateAll: %w", err3)
	}
	return fmt.Errorf("failed to LoadFile %s due to err %v", path, loadErr)
}
