package figtree

import (
	"encoding/json"
	"flag"
	check "github.com/andreimerlescu/checkfs"
	"github.com/andreimerlescu/checkfs/file"
	"github.com/go-ini/ini"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// Mutations returns a receiver channel of Mutation data
func (fig *Tree) Mutations() <-chan Mutation {
	return fig.mutationsCh
}

// Recall is when you bring the mutations channel back to life and you unlock making further changes to the fig *Tree
func (fig *Tree) Recall() {
	fig.angel.Store(false)
	fig.mutationsCh = make(chan Mutation, fig.harvest)
	fig.tracking = true
}

// Curse is when you lock the fig *Tree from further changes, stop tracking and close the channel
func (fig *Tree) Curse() {
	fig.angel.Store(true)
	fig.tracking = false
	close(fig.mutationsCh)
}

// Resurrect revives a missing or nil definition, checking env and config files first
func (fig *Tree) Resurrect(name string) {
	fig.mu.Lock()
	defer fig.mu.Unlock()
	if _, exists := fig.figs[name]; !exists {
		// Check environment first
		if val, ok := os.LookupEnv(name); ok {
			ptr := new(string)
			*ptr = strings.Clone(val) // Use strings.Clone as requested
			fig.figs[name] = &Fig{Flesh: ptr}
			flag.String(name, val, "Resurrected from environment")
			return
		}

		// Check config files with traditional for loop
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
								fig.figs[name] = &Fig{Flesh: ptr}
								flag.String(name, strVal, "Resurrected from JSON")
								return
							}
						}
					case ".yaml", ".yml":
						if yaml.Unmarshal(data, &m) == nil && m[name] != nil {
							if strVal, err := toString(m[name]); err == nil {
								ptr := new(string)
								*ptr = strings.Clone(strVal)
								fig.figs[name] = &Fig{Flesh: ptr}
								flag.String(name, strVal, "Resurrected from YAML")
								return
							}
						}
					case ".ini":
						if cfg, err := ini.Load(data); err == nil {
							if val := cfg.Section("").Key(name).String(); val != "" {
								ptr := new(string)
								*ptr = strings.Clone(val)
								fig.figs[name] = &Fig{Flesh: ptr}
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
		fig.figs[name] = &Fig{Flesh: ptr}
		flag.String(name, "", "Resurrected configuration")
	}
}
