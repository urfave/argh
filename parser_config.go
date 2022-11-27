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
	tracef(2, "NValue.Contains(%v)", i)

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
		pCfg.Prog = NewCommandConfig()
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

	callbacks []func(*CommandFlag) error
}

func NewCommandConfig() *CommandConfig {
	cCfg := &CommandConfig{}
	cCfg.init()

	return cCfg
}

func (cCfg *CommandConfig) On(f func(*CommandFlag) error) {
	cCfg.init()

	cCfg.callbacks = append(cCfg.callbacks, f)
}

func (cCfg *CommandConfig) applyCallbacks(cf *CommandFlag) error {
	for _, callback := range cCfg.callbacks {
		if err := callback(cf); err != nil {
			return err
		}
	}

	return nil
}

func (cCfg *CommandConfig) Child() *CommandConfig {
	child := NewCommandConfig()
	child.Flags.Parent = cCfg.Flags
	return child
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

	if cCfg.callbacks == nil {
		cCfg.callbacks = []func(*CommandFlag) error{}
	}
}

// GetParsedCommand returns the CommandConfig for the command that was
// parsed or nil if none
func (cCfg *CommandConfig) GetParsedCommand() (string, *CommandConfig) {
	return "", nil
	/*
		tracef(2, "%[1]p CommandConfig.GetParsedCommand()", cCfg)
		if cCfg.Node == nil {
			tracef(2, "%[1]p CommandConfig.GetParsedCommand() -> \"\", nil", cCfg)
			return "", nil
		}

		name, subCfg := cCfg.Commands.GetParsed()
		tracef(2, "%[1]p CommandConfig.GetParsedCommand() -> %[2]q, %[3]p", cCfg, name, subCfg)
		return name, subCfg
	*/
}

// NFlag works like flag.FlagSet.NFlag
func (cCfg *CommandConfig) NFlag() int {
	return 0
	/*
		n := 0

		for _, flCfg := range cCfg.Flags.Map {
			if flCfg.Node != nil {
				n++
			}
		}

		return n
	*/
}

// Visit works like flag.FlagSet.Visit
func (cCfg *CommandConfig) Visit(f func(*CommandFlag)) {
	_ = f
	/*
		names := make([]string, len(cCfg.Flags.Map))

		i := 0
		for name := range cCfg.Flags.Map {
			names[i] = name
			i++
		}

		sort.Strings(names)

		for _, name := range names {
			if flCfg, ok := cCfg.Flags.GetShallow(name); ok && flCfg.Node != nil {
				f(flCfg.Node)
			}
		}
	*/
}

// Args is like flag.FlagSet.Args
func (cCfg *CommandConfig) Args() []string {
	return []string{}
	/*
		tracef(2, "%[1]p CommandConfig.Args()", cCfg)
		ret := []string{}

		if cCfg.Node == nil {
			tracef(2, "%[1]p CommandConfig.Args() -> %+[2]v", cCfg, ret)
			return ret
		}

		for _, node := range cCfg.Node.Nodes {
			if idn, ok := node.(*Ident); ok {
				ret = append(ret, idn.Literal)
			}

			if p, ok := node.(*PassthroughArgs); ok {
				for _, ptNode := range p.Nodes {
					if idn, ok := ptNode.(*Ident); ok {
						ret = append(ret, idn.Literal)
					}
				}
			}
		}

		tracef(2, "%[1]p CommandConfig.Args() -> %+[2]v", cCfg, ret)
		return ret
	*/
}

// Lookup is like flag.FlagSet.Lookup
func (cCfg *CommandConfig) Lookup(name string) *FlagConfig {
	if flCfg, ok := cCfg.Flags.Get(name); ok {
		return flCfg
	}

	return nil
}

// Set is like flag.FlagSet.Set
func (cCfg *CommandConfig) Set(name, value string) error {
	return nil
	/*
		cCfg.SetFlagConfig(
			name,
			&FlagConfig{
				Node: &CommandFlag{
					Name:   name,
					Values: map[string]string{name: value},
				},
			},
		)

		return nil
	*/
}

func (cCfg *CommandConfig) GetCommandConfig(name string) (*CommandConfig, bool) {
	tracef(2, "%[1]p CommandConfig.GetCommandConfig(%[2]q)", cCfg, name)

	cCfg.init()

	return cCfg.Commands.Get(name)
}

func (cCfg *CommandConfig) GetFlagConfig(name string) (*FlagConfig, bool) {
	tracef(2, "%[1]p CommandConfig.GetFlagConfig(%[2]q)", cCfg, name)

	cCfg.init()

	return cCfg.Flags.Get(name)
}

func (cCfg *CommandConfig) SetFlagConfig(name string, flCfg *FlagConfig) {
	tracef(2, "%[1]p CommandConfig.SetFlagConfig(%+[2]q, %[3]p)", cCfg, name, flCfg)

	cCfg.init()

	cCfg.Flags.Set(name, flCfg)
}

func (cCfg *CommandConfig) SetDefaultFlagConfig(name string, flCfg *FlagConfig) {
	tracef(2, "%[1]p CommandConfig.SetDefaultFlagConfig(%[2]q, %[3]p)", cCfg, name, flCfg)

	cCfg.init()

	cCfg.Flags.SetDefault(name, flCfg)
}

type FlagConfig struct {
	NValue     NValue
	Persist    bool
	ValueNames []string

	callbacks []func(*CommandFlag) error
}

func (flCfg *FlagConfig) On(f func(*CommandFlag) error) {
	if flCfg.callbacks == nil {
		flCfg.callbacks = []func(*CommandFlag) error{}
	}

	flCfg.callbacks = append(flCfg.callbacks, f)
}

func (flCfg *FlagConfig) applyCallbacks(cf *CommandFlag) error {
	for _, callback := range flCfg.callbacks {
		if err := callback(cf); err != nil {
			return err
		}
	}

	return nil
}

func (flCfg *FlagConfig) Set(value string) error {
	return nil

	/*
		if flCfg.Node == nil {
			if len(flCfg.ValueNames) == 0 {
				return fmt.Errorf("cannot set value of flag without name: %w", Error)
			}

			flCfg.Node = &CommandFlag{
				Name: flCfg.ValueNames[0],
			}
		}

		flCfg.Node.Values[flCfg.Node.Name] = value

		return nil
	*/
}

func (flCfg *FlagConfig) Name() string {
	return ""
	/*
		if flCfg.Node == nil {
			return ""
		}

		return flCfg.Node.Name
	*/
}

// Value is like flag.Value.String (see alias)
func (flCfg *FlagConfig) Value() string {
	if vals := flCfg.Values(); len(vals) > 0 {
		return vals[0]
	}

	return ""
}

func (flCfg *FlagConfig) LookupValue() (string, bool) {
	if vals := flCfg.Values(); len(vals) > 0 {
		return vals[0], true
	}

	return "", false
}

func (flCfg *FlagConfig) String() string { return flCfg.Value() }

func (flCfg *FlagConfig) Values() []string {
	return []string{}

	/*
		if flCfg.Node == nil {
			return []string{}
		}

		values := make([]string, len(flCfg.Node.Values))

		i := 0
		for _, value := range flCfg.Node.Values {
			values[i] = value
			i++
		}

		return values
	*/
}

type Flags struct {
	Parent *Flags
	Map    map[string]*FlagConfig

	Automatic bool
}

func (fl *Flags) GetShallow(name string) (*FlagConfig, bool) {
	fl.ensureMap()

	flCfg, ok := fl.Map[name]
	return flCfg, ok
}

func (fl *Flags) Get(name string) (*FlagConfig, bool) {
	tracef(2, "%[1]p Flags.Get(%[2]q)", fl, name)

	fl.ensureMap()

	flCfg, ok := fl.Map[name]
	if !ok {
		if fl.Automatic {
			tracef(2, "%[1]p Flags.Get(%[1]q) -> {} (automatic)", fl, name)
			return &FlagConfig{}, true
		}

		if fl.Parent != nil {
			v, ok := fl.Parent.Get(name)
			tracef(2, "%[1]p Flags.Get(%[2]q) -> %[3]p", fl, name, v)
			return v, ok && v.Persist
		}
	}

	tracef(2, "%[1]p Flags.Get(%[2]q) -> %[3]p", fl, name, flCfg)
	return flCfg, ok
}

func (fl *Flags) Set(name string, flCfg *FlagConfig) {
	tracef(2, "%[1]p Flags.Set(%[2]q, %[3]p)", fl, name, flCfg)

	fl.ensureMap()

	fl.Map[name] = flCfg
}

func (fl *Flags) SetDefault(name string, flCfg *FlagConfig) {
	tracef(2, "%[1]p Flags.SetDefault(%[2]q, %[3]p)", fl, name, flCfg)

	fl.ensureMap()

	if _, ok := fl.Map[name]; !ok {
		fl.Map[name] = flCfg
	}
}

func (fl *Flags) ensureMap() {
	if fl.Map == nil {
		fl.Map = map[string]*FlagConfig{}
	}
}

type Commands struct {
	Map map[string]*CommandConfig
}

func (cmd *Commands) GetParsed() (string, *CommandConfig) {
	return "", nil

	/*
		for name, loopCcfg := range cmd.Map {
			cCfg := loopCcfg

			tracef(2, "%[1]p Commands.GetParsed() check (%[2]q, %[3]p)", cmd, name, cCfg)

			if cCfg.Node != nil {
				tracef(2, "%[1]p Commands.GetParsed() -> (%[2]q, %[3]p)", cmd, name, cCfg)
				return name, cCfg
			}
		}

		tracef(2, "%[1]p Commands.GetParsed() -> (\"\", nil)", cmd)
		return "", nil
	*/
}

func (cmd *Commands) Get(name string) (*CommandConfig, bool) {
	tracef(2, "%[1]p Commands.Get(%[2]q)", cmd, name)

	cmd.ensureMap()

	cmdCfg, ok := cmd.Map[name]
	if !ok {
		tracef(2, "%[1]p Commands.Get(%[2]q) -> nil", cmd, name)
		return nil, ok
	}

	tracef(2, "%[1]p Commands.Get(%[2]q) -> %[3]p", cmd, name, cmdCfg)
	return cmdCfg, ok
}

func (cmd *Commands) Set(name string, cmdCfg *CommandConfig) {
	tracef(2, "%[1]p Commands.Set(%[2]q, %[3]p)", cmd, name, cmdCfg)

	cmd.ensureMap()

	cmd.Map[name] = cmdCfg
}

func (cmd *Commands) SetDefault(name string, cmdCfg *CommandConfig) {
	tracef(2, "%[1]p Commands.SetDefault(%[2]q, %[3]p)", cmd, name, cmdCfg)

	cmd.ensureMap()

	if _, ok := cmd.Map[name]; ok {
		cmd.Map[name] = cmdCfg
	}
}

func (cmd *Commands) ensureMap() {
	if cmd.Map == nil {
		cmd.Map = map[string]*CommandConfig{}
	}
}
