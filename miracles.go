package figtree

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	check "github.com/andreimerlescu/checkfs"
	"github.com/andreimerlescu/checkfs/file"
	"github.com/go-ini/ini"
	"gopkg.in/yaml.v3"
)

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

// Resurrect revives a missing or nil definition, checking env and config files first
func (tree *figTree) Resurrect(name string) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	if tree.HasRule(RuleCondemnedFromResurrection) {
		log.Fatalf("resurrection %q condemned", name)
	}
	if _, exists := tree.figs[name]; !exists {
		// Check environment first
		if !tree.ignoreEnv {
			if val, ok := os.LookupEnv(name); ok {
				ptr := new(string)
				*ptr = strings.Clone(val) // Use strings.Clone as requested
				tree.figs[name] = &figFruit{
					Flesh:       figFlesh{ptr},
					Mutagenesis: tree.MutagenesisOf(val),
					Validators:  make([]FigValidatorFunc, 0),
					Callbacks:   make([]Callback, 0),
					Mutations:   make([]Mutation, 0),
				}
				flag.String(name, val, "Resurrected from environment")
				return
			}
		}

		// Check config files with traditional for loop
		envVal := ""
		if !tree.ignoreEnv {
			envVal = os.Getenv(EnvironmentKey)
		}
		files := []string{
			envVal,
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
				data, err := os.ReadFile(f)
				if err == nil {
					var m map[string]interface{}
					ext := strings.ToLower(filepath.Ext(f))
					switch ext {
					case ".json":
						if json.Unmarshal(data, &m) == nil && m[name] != nil {
							if strVal, err := toString(m[name]); err == nil {
								ptr := new(string)
								*ptr = strings.Clone(strVal)
								tree.figs[name] = &figFruit{
									Flesh:       figFlesh{ptr},
									Mutagenesis: tree.MutagenesisOf(ptr),
									Validators:  make([]FigValidatorFunc, 0),
									Callbacks:   make([]Callback, 0),
									Mutations:   make([]Mutation, 0),
								}
								flag.String(name, strVal, "Resurrected from JSON")
								return
							}
						}
					case ".yaml", ".yml":
						if yaml.Unmarshal(data, &m) == nil && m[name] != nil {
							if strVal, err := toString(m[name]); err == nil {
								ptr := new(string)
								*ptr = strings.Clone(strVal)
								tree.figs[name] = &figFruit{
									Flesh:       figFlesh{ptr},
									Mutagenesis: tree.MutagenesisOf(ptr),
									Validators:  make([]FigValidatorFunc, 0),
									Callbacks:   make([]Callback, 0),
									Mutations:   make([]Mutation, 0),
								}
								flag.String(name, strVal, "Resurrected from YAML")
								return
							}
						}
					case ".ini":
						if cfg, err := ini.Load(data); err == nil {
							if val := cfg.Section("").Key(name).String(); val != "" {
								ptr := new(string)
								*ptr = strings.Clone(val)
								tree.figs[name] = &figFruit{
									Flesh:       figFlesh{ptr},
									Mutagenesis: tree.MutagenesisOf(ptr),
									Validators:  make([]FigValidatorFunc, 0),
									Callbacks:   make([]Callback, 0),
									Mutations:   make([]Mutation, 0),
								}
								flag.String(name, val, "Resurrected from INI")
								return
							}
						}
					}
				}
			}
		}

		// Default to empty string if no value found
		ptr := new(string)
		*ptr = ""
		tree.figs[name] = &figFruit{
			Flesh:       figFlesh{ptr},
			Mutagenesis: tree.MutagenesisOf(ptr),
			Validators:  make([]FigValidatorFunc, 0),
			Callbacks:   make([]Callback, 0),
			Mutations:   make([]Mutation, 0),
		}
		flag.String(name, "", "Resurrected configuration")
	}
}

// Fig returns a figFruit on the fig figTree
func (tree *figTree) Fig(name string) Flesh {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	fruit, exists := tree.figs[name]
	if !exists {
		return nil
	}
	return &fruit.Flesh
}
