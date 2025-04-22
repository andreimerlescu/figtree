package figtree

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func (tree *figTree) ReadFrom(path string) error {
	_, fileErr := os.Stat(path)
	if os.IsNotExist(fileErr) || os.IsPermission(fileErr) {
		return fileErr
	}
	return tree.loadFile(path)
}

func (tree *figTree) SaveTo(path string) error {
	var properties = make(map[string]interface{})
	tree.mu.Lock()
	defer tree.mu.Unlock()
	for name, fig := range tree.figs {
		properties[name] = fig.Flesh.Flesh
	}
	formatValue := func(val interface{}) string {
		return fmt.Sprintf("%v", val)
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		yamlBytes, yamlErr := yaml.Marshal(properties)
		if yamlErr != nil {
			return yamlErr
		}
		return os.WriteFile(path, yamlBytes, 0644)
	case ".json":
		jsonBytes, jsonErr := json.MarshalIndent(properties, "", "  ")
		if jsonErr != nil {
			return jsonErr
		}
		return os.WriteFile(path, jsonBytes, 0644)
	case ".ini":
		cfg := ini.Empty()
		for key, value := range properties {
			switch v := value.(type) {
			case map[string]interface{}:
				section, err := cfg.NewSection(key)
				if err != nil {
					return err
				}
				for sk, sv := range v {
					section.Key(sk).SetValue(formatValue(sv))
				}
			default:
				cfg.Section("").Key(key).SetValue(formatValue(value))
			}
		}
		return cfg.SaveTo(path)
	default:
		return errors.New("invalid file extension provided")
	}
}
