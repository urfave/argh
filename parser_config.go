package argh

const (
	OneOrMoreValue  NValue = -2
	ZeroOrMoreValue NValue = -1
	ZeroValue       NValue = 0
)

var (
	DefaultParserConfig = &ParserConfig{
		Commands:      map[string]CommandConfig{},
		Flags:         map[string]FlagConfig{},
		ScannerConfig: DefaultScannerConfig,
	}
)

type NValue int

type ParserConfig struct {
	Prog     CommandConfig
	Commands map[string]CommandConfig
	Flags    map[string]FlagConfig

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
