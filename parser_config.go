package argh

const (
	OneOrMoreValue  NValue = -2
	ZeroOrMoreValue NValue = -1
	ZeroValue       NValue = 0
)

var (
	POSIXyParserConfig = &ParserConfig{
		Prog:          CommandConfig{},
		ScannerConfig: POSIXyScannerConfig,
	}
)

type NValue int

func (nv NValue) Contains(i int) bool {
	if i < int(ZeroValue) {
		return false
	}

	if nv == OneOrMoreValue || nv == ZeroOrMoreValue {
		return true
	}

	return int(nv) > i
}

type ParserConfig struct {
	Prog CommandConfig

	ScannerConfig *ScannerConfig
}

type CommandConfig struct {
	NValue     NValue
	ValueNames []string
	Flags      map[string]FlagConfig
	Commands   map[string]CommandConfig
}

type FlagConfig struct {
	NValue     NValue
	ValueNames []string
}
