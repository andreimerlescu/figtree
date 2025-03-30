package figtree

import (
	"fmt"
	"os"
	"path/filepath"

	check "github.com/andreimerlescu/checkfs"
	"github.com/andreimerlescu/checkfs/file"
)

// Reload will readEnv on each flag in the configurable package
func (tree *Tree) Reload() error {
	tree.readEnv()
	return tree.validateAll()
}

// Load uses the EnvironmentKey and the DefaultJSONFile, DefaultYAMLFile, and DefaultINIFile to run ParseFile if it exists
func (tree *Tree) Load() (err error) {
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

	files := []string{
		os.Getenv(EnvironmentKey),
		tree.ConfigFilePath,
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
			if err := tree.loadFile(f); err != nil {
				return fmt.Errorf("failed to load %s: %w", f, err)
			}
		}
	}

	tree.readEnv()
	return tree.validateAll()
}

// LoadFile accepts a path and uses it to populate the Tree
func (tree *Tree) LoadFile(path string) (err error) {
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
	var loadErr error
	if loadErr = check.File(path, file.Options{Exists: true}); loadErr == nil {
		if err2 := tree.loadFile(path); err2 != nil {
			return fmt.Errorf("failed to loadFile %s: %w", path, err2)
		}
		tree.readEnv()
		err3 := tree.validateAll()
		if err3 != nil {
			return fmt.Errorf("failed to validateAll: %w", err3)
		}
		return nil
	}
	tree.readEnv()
	err3 := tree.validateAll()
	if err3 != nil {
		return fmt.Errorf("failed to validateAll: %w", err3)
	}
	return fmt.Errorf("failed to LoadFile %s due to err %v", path, loadErr)
}
