package argh

var (
	// POSIXyScannerConfig defines a scanner config that uses '-'
	// as the flag prefix, which also means that "--" is the "long
	// flag" prefix, a bare "--" is considered STOP_FLAG, and a
	// bare "-" is considered STDIN_FLAG.
	POSIXyScannerConfig = &ScannerConfig{
		AssignmentOperator: '=',
		FlagPrefix:         '-',
		MultiValueDelim:    ',',
	}
)

type ScannerConfig struct {
	AssignmentOperator rune
	FlagPrefix         rune
	MultiValueDelim    rune
}

func (cfg *ScannerConfig) IsFlagPrefix(ch rune) bool {
	return ch == cfg.FlagPrefix
}

func (cfg *ScannerConfig) IsMultiValueDelim(ch rune) bool {
	return ch == cfg.MultiValueDelim
}

func (cfg *ScannerConfig) IsAssignmentOperator(ch rune) bool {
	return ch == cfg.AssignmentOperator
}

func (cfg *ScannerConfig) IsBlankspace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (cfg *ScannerConfig) IsUnderscore(ch rune) bool {
	return ch == '_'
}
