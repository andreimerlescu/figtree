package figtree

import (
	"embed"
	"path/filepath"
	"strings"
)

//go:embed VERSION
var versionBytes embed.FS

var currentVersion string

func Version() string {
	if len(currentVersion) == 0 {
		versionBytes, err := versionBytes.ReadFile("VERSION")
		if err != nil {
			return ""
		}
		currentVersion = strings.TrimSpace(string(versionBytes))
	}
	return currentVersion
}

const (
	DefaultYAMLFile string = "config.yml"  // Default filename for a YAML configuration file
	DefaultJSONFile string = "config.json" // Default filename for a JSON configuration file
	DefaultINIFile  string = "config.ini"  // Default filename for a INI configuration file

	tString       Mutagenesis = "String"
	tBool         Mutagenesis = "Bool"
	tInt          Mutagenesis = "Int"
	tInt64        Mutagenesis = "Int64"
	tFloat64      Mutagenesis = "Float64"
	tDuration     Mutagenesis = "Duration"
	tUnitDuration Mutagenesis = "UnitDuration"
	tList         Mutagenesis = "List"
	tMap          Mutagenesis = "Map"

	CallbackAfterChange  CallbackWhen = "CallbackAfterChange"
	CallbackAfterRead    CallbackWhen = "CallbackAfterRead"
	CallbackAfterVerify  CallbackWhen = "CallbackAfterVerify"
	CallbackBeforeChange CallbackWhen = "CallbackBeforeChange"
	CallbackBeforeRead   CallbackWhen = "CallbackBeforeRead"
	CallbackBeforeVerify CallbackWhen = "CallbackBeforeVerify"
)

// Mutageneses is the plural form of Mutagenesis and this is a slice of Mutagenesis
var Mutageneses = []Mutagenesis{tString, tBool, tInt, tInt64, tFloat64, tDuration, tUnitDuration, tList, tMap}

// EnvironmentKey stores the preferred ENV that contains the path to your configuration file (.ini, .json or .yaml)
var EnvironmentKey string = "CONFIG_FILE"

// ConfigFilePath stores the path to the configuration file of choice
var ConfigFilePath string = filepath.Join(".", DefaultYAMLFile)
