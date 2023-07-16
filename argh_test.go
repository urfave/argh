package argh_test

import (
	"flag"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/urfave/argh"
)

func ptrTo[T any](v T) *T {
	return &v
}

func ptrFrom[T any](v *T) T {
	if v != nil {
		return *v
	}

	var zero T
	return zero
}

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

			_ = fl.Parse([]string{
				"-ok",
				"-dur", "42h42m10s",
				"-f64", "4242424242.42",
				"-i", "-42",
				"-i64", "-4242424242",
				"-s", "the answer",
				"-u", "42",
				"-u64", "4242424242",
			})
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
			var (
				okFlag  *bool
				durFlag *time.Duration
				f64Flag *float64
				iFlag   *int
				i64Flag *int64
				sFlag   *string
				uFlag   *uint
				u64Flag *uint64
			)

			pCfg := argh.NewParserConfig()
			pCfg.Prog = &argh.CommandConfig{
				Flags: &argh.Flags{
					Map: map[string]argh.FlagConfig{
						"ok": {
							On: func(fl argh.Flag) error {
								okFlag = ptrTo(true)
								return nil
							},
						},
						"dur": {
							NValue: 1,
							On: func(fl argh.Flag) error {
								if v, ok := argh.FirstValue(fl.Values, "0"); ok {
									if pt, err := time.ParseDuration(v); err != nil {
										durFlag = ptrTo(pt)
									}
								}
								return nil
							},
						},
						"f64": {
							NValue: 1,
							On: func(fl argh.Flag) error {
								if v, ok := argh.FirstValue(fl.Values, "0"); ok {
									if f, err := strconv.ParseFloat(v, 64); err == nil {
										f64Flag = ptrTo(f)
									}
								}
								return nil
							},
						},
						"i": {
							NValue: 1,
							On: func(fl argh.Flag) error {
								if v, ok := argh.FirstValue(fl.Values, "0"); ok {
									if i, err := strconv.ParseInt(v, 10, 64); err == nil {
										iFlag = ptrTo(int(i))
									}
								}
								return nil
							},
						},
						"i64": {
							NValue: 1,
							On: func(fl argh.Flag) error {
								if v, ok := argh.FirstValue(fl.Values, "0"); ok {
									if i, err := strconv.ParseInt(v, 10, 64); err == nil {
										i64Flag = ptrTo(i)
									}
								}
								return nil
							},
						},
						"s": {
							NValue: 1,
							On: func(fl argh.Flag) error {
								if v, ok := argh.FirstValue(fl.Values, "0"); ok {
									sFlag = ptrTo(v)
								}
								return nil
							},
						},
						"u": {
							NValue: 1,
							On: func(fl argh.Flag) error {
								if v, ok := argh.FirstValue(fl.Values, "0"); ok {
									if u, err := strconv.ParseUint(v, 10, 64); err == nil {
										uFlag = ptrTo(uint(u))
									}
								}
								return nil
							},
						},
						"u64": {
							NValue: 1,
							On: func(fl argh.Flag) error {
								if v, ok := argh.FirstValue(fl.Values, "0"); ok {
									if u, err := strconv.ParseUint(v, 10, 64); err == nil {
										u64Flag = ptrTo(u)
									}
								}
								return nil
							},
						},
					},
				},
			}

			_, _ = argh.ParseArgs([]string{
				"--ok",
				"--dur", "42h42m10s",
				"--f64", "4242424242.42",
				"-i", "-42",
				"--i64", "-4242424242",
				"-s", "the answer",
				"-u", "42",
				"--u64", "4242424242",
			}, pCfg)
			_ = fmt.Sprint(
				"okFlag", ptrFrom(okFlag),
				"durFlag", ptrFrom(durFlag),
				"f64Flag", ptrFrom(f64Flag),
				"iFlag", ptrFrom(iFlag),
				"i64Flag", ptrFrom(i64Flag),
				"sFlag", ptrFrom(sFlag),
				"uFlag", ptrFrom(uFlag),
				"u64Flag", ptrFrom(u64Flag),
			)
		}()
	}
}
