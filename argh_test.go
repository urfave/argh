package argh_test

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"git.meatballhat.com/x/argh"
)

func BenchmarkStdlibFlag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			fl := flag.NewFlagSet("bench", flag.PanicOnError)
			okFlag := fl.Bool("ok", false, "")
			durFlag := fl.Duration("dur", time.Second, "")
			f64Flag := fl.Float64("f64", float64(42.0), "")
			iFlag := fl.Int("i", -11, "")
			i64Flag := fl.Int64("i64", -111111111111, "")
			sFlag := fl.String("s", "hello", "")
			uFlag := fl.Uint("u", 11, "")
			u64Flag := fl.Uint64("u64", 11111111111111111111, "")

			_ = fl.Parse([]string{})
			_ = fmt.Sprint(
				"fl", fl,
				"okFlag", *okFlag,
				"durFlag", *durFlag,
				"f64Flag", *f64Flag,
				"iFlag", *iFlag,
				"i64Flag", *i64Flag,
				"sFlag", *sFlag,
				"uFlag", *uFlag,
				"u64Flag", *u64Flag,
			)
		}()
	}
}

func BenchmarkArgh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			pCfg := argh.NewParserConfig()
			pCfg.Prog = &argh.CommandConfig{
				Flags: &argh.Flags{
					Map: map[string]argh.FlagConfig{
						"ok":  {},
						"dur": {},
						"f64": {},
						"i":   {},
						"i64": {},
						"s":   {},
						"u":   {},
						"u64": {},
					},
				},
			}

			_, _ = argh.ParseArgs([]string{}, pCfg)
		}()
	}
}
