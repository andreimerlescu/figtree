package figtree

import (
	"path/filepath"
)

const (
	DefaultYAMLFile string = "config.yml"  // Default filename for a YAML configuration file
	DefaultJSONFile string = "config.json" // Default filename for a JSON configuration file
	DefaultINIFile  string = "config.ini"  // Default filename for a INI configuration file

	VERSION = "v2.0.0"

	tString       Mutagenesis = "String"
	tBool         Mutagenesis = "Bool"
	tInt          Mutagenesis = "Int"
	tInt64        Mutagenesis = "Int64"
	tFloat64      Mutagenesis = "Float64"
	tDuration     Mutagenesis = "Duration"
	tUnitDuration Mutagenesis = "UnitDuration"
	tList         Mutagenesis = "List"
	tMap          Mutagenesis = "Map"
)

const CallbackAfterChange CallbackAfter = "CallbackAfterChange"
const CallbackAfterRead CallbackAfter = "CallbackAfterRead"
const CallbackAfterVerify CallbackAfter = "CallbackAfterVerify"

// Mutageneses is the plural form of Mutagenesis and this is a slice of Mutagenesis
var Mutageneses = []Mutagenesis{tString, tBool, tInt, tInt64, tFloat64, tDuration, tUnitDuration, tList, tMap}

// EnvironmentKey stores the preferred ENV that contains the path to your configuration file (.ini, .json or .yaml)
var EnvironmentKey string = "CONFIG_FILE"

// ConfigFilePath stores the path to the configuration file of choice
var ConfigFilePath string = filepath.Join(".", DefaultYAMLFile)
