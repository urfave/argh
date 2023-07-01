package argh

const (
	OneOrMoreValue  NValue = -2
	ZeroOrMoreValue NValue = -1
	ZeroValue       NValue = 0
)

var (
	zeroValuePtr = func() *NValue {
		v := ZeroValue
		return &v
	}()
)

type NValue int

func (nv NValue) Required() bool {
	if nv == OneOrMoreValue {
		return true
	}

	return int(nv) >= 1
}

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
	Prog *CommandConfig

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

	if pCfg.Prog == nil {
		pCfg.Prog = &CommandConfig{}
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

	On func(CommandFlag)
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

func (cCfg *CommandConfig) SetCommandConfig(name string, sCfg *CommandConfig) {
	tracef("CommandConfig.SetCommandConfig(%q, ...)", name)

	if cCfg.Commands == nil {
		cCfg.Commands = &Commands{Map: map[string]CommandConfig{}}
	}

	sCfg.init()
	sCfg.Flags.Parent = cCfg.Flags

	cCfg.Commands.Set(name, sCfg)
}

func (cCfg *CommandConfig) GetFlagConfig(name string) (FlagConfig, bool) {
	tracef("CommandConfig.GetFlagConfig(%q)", name)

	if cCfg.Flags == nil {
		cCfg.Flags = &Flags{Map: map[string]FlagConfig{}}
	}

	return cCfg.Flags.Get(name)
}

func (cCfg *CommandConfig) SetFlagConfig(name string, flCfg *FlagConfig) {
	tracef("CommandConfig.SetFlagConfig(%q, ...)", name)

	if cCfg.Flags == nil {
		cCfg.Flags = &Flags{Map: map[string]FlagConfig{}}
	}

	cCfg.Flags.Set(name, flCfg)
}

type FlagConfig struct {
	NValue     NValue
	Persist    bool
	ValueNames []string

	On func(CommandFlag)
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

func (fl *Flags) Set(name string, flCfg *FlagConfig) {
	tracef("Flags.Get(%q)", name)

	if fl.Map == nil {
		fl.Map = map[string]FlagConfig{}
	}

	fl.Map[name] = *flCfg
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

func (cmd *Commands) Set(name string, cCfg *CommandConfig) {
	tracef("Commands.Set(%q, ...)", name)

	if cmd.Map == nil {
		cmd.Map = map[string]CommandConfig{}
	}

	cCfg.init()

	cmd.Map[name] = *cCfg
}
