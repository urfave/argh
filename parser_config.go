package argh

const (
	OneOrMoreValue  NValue = -2
	ZeroOrMoreValue NValue = -1
	ZeroValue       NValue = 0
)

var (
	POSIXyParserConfig = NewParserConfig()
)

type NValue int

func (nv NValue) Contains(i int) bool {
	tracef("NValue.Contains(%v)", i)

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

type ParserOption func(*ParserConfig)

func NewParserConfig(opts ...ParserOption) *ParserConfig {
	pCfg := &ParserConfig{}

	for _, opt := range opts {
		if opt != nil {
			opt(pCfg)
		}
	}

	if pCfg.Prog.IsZero() {
		pCfg.Prog = CommandConfig{}
		pCfg.Prog.init()
	}

	if pCfg.ScannerConfig == nil {
		pCfg.ScannerConfig = POSIXyScannerConfig
	}

	return pCfg
}

type CommandConfig struct {
	NValue     NValue
	ValueNames []string
	Flags      *Flags
	Commands   *Commands
}

func (cCfg *CommandConfig) IsZero() bool {
	return cCfg.NValue == NValue(0) &&
		cCfg.ValueNames == nil &&
		cCfg.Flags == nil &&
		cCfg.Commands == nil
}

func (cCfg *CommandConfig) init() {
	if cCfg.ValueNames == nil {
		cCfg.ValueNames = []string{}
	}

	if cCfg.Flags == nil {
		cCfg.Flags = &Flags{}
	}

	if cCfg.Commands == nil {
		cCfg.Commands = &Commands{}
	}
}

func (cCfg *CommandConfig) GetCommandConfig(name string) (CommandConfig, bool) {
	tracef("CommandConfig.GetCommandConfig(%q)", name)

	if cCfg.Commands == nil {
		cCfg.Commands = &Commands{Map: map[string]CommandConfig{}}
	}

	return cCfg.Commands.Get(name)
}

func (cCfg *CommandConfig) GetFlagConfig(name string) (FlagConfig, bool) {
	tracef("CommandConfig.GetFlagConfig(%q)", name)

	if cCfg.Flags == nil {
		cCfg.Flags = &Flags{Map: map[string]FlagConfig{}}
	}

	return cCfg.Flags.Get(name)
}

type FlagConfig struct {
	NValue     NValue
	Persist    bool
	ValueNames []string
}

type Flags struct {
	Parent *Flags
	Map    map[string]FlagConfig

	Automatic bool
}

func (fl *Flags) Get(name string) (FlagConfig, bool) {
	tracef("Flags.Get(%q)", name)

	if fl.Map == nil {
		fl.Map = map[string]FlagConfig{}
	}

	flCfg, ok := fl.Map[name]
	if !ok {
		if fl.Automatic {
			return FlagConfig{}, true
		}

		if fl.Parent != nil {
			flCfg, ok = fl.Parent.Get(name)
			return flCfg, ok && flCfg.Persist
		}
	}

	return flCfg, ok
}

type Commands struct {
	Map map[string]CommandConfig
}

func (cmd *Commands) Get(name string) (CommandConfig, bool) {
	tracef("Commands.Get(%q)", name)

	if cmd.Map == nil {
		cmd.Map = map[string]CommandConfig{}
	}

	cmdCfg, ok := cmd.Map[name]
	return cmdCfg, ok
}
