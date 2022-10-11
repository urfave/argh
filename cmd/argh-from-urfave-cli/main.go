package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"git.meatballhat.com/x/argh"
	"github.com/urfave/cli/v2"
)

func main() {
	pc := argh.NewParserConfig()

	app := &cli.App{
		Name: "argh-from-urfave-cli",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
			&cli.BoolFlag{Name: "happy", Aliases: []string{"H"}},
		},
		Action: func(cCtx *cli.Context) error {
			return dumpAST(cCtx.App.Writer, os.Args, pc)
		},
		Commands: []*cli.Command{
			{
				Name: "poof",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "like", Aliases: []string{"L"}},
					&cli.BoolFlag{Name: "loudly", Aliases: []string{"l"}},
				},
				Action: func(cCtx *cli.Context) error {
					return dumpAST(cCtx.App.Writer, os.Args, pc)
				},
				Subcommands: []*cli.Command{
					{
						Name: "immediately",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "really", Aliases: []string{"X"}},
						},
						Action: func(cCtx *cli.Context) error {
							return dumpAST(cCtx.App.Writer, os.Args, pc)
						},
					},
				},
			},
		},
	}

	mapFlags(pc.Prog.Flags, app.Flags, true)
	mapCommands(pc.Prog, app.Commands)

	app.RunContext(context.Background(), os.Args)
}

func dumpAST(w io.Writer, args []string, pc *argh.ParserConfig) error {
	pt, err := argh.ParseArgs(args, pc)
	if err != nil {
		return err
	}

	q := argh.NewQuerier(pt.Nodes)

	jsonBytes, err := json.MarshalIndent(
		map[string]any{
			"ast":    q.AST(),
			"config": pc,
		},
		"", "  ",
	)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, string(jsonBytes))
	return err
}

func mapFlags(cmdFlags *argh.Flags, flags []cli.Flag, persist bool) {
	for _, fl := range flags {
		nValue := argh.ZeroValue

		if df, ok := fl.(cli.DocGenerationFlag); ok {
			if df.TakesValue() {
				if dsf, ok := df.(cli.DocGenerationSliceFlag); ok && dsf.IsSliceFlag() {
					nValue = argh.OneOrMoreValue
				} else {
					nValue = argh.NValue(1)
				}
			}
		}

		flCfg := argh.FlagConfig{
			NValue:     nValue,
			Persist:    persist,
			ValueNames: []string{fl.Names()[0]},
		}

		for _, flAlias := range fl.Names() {
			cmdFlags.Set(flAlias, flCfg)
		}
	}
}

func mapCommands(cCfg *argh.CommandConfig, cmds []*cli.Command) {
	for _, cmd := range cmds {
		// TODO: vary nValue if/when cli.Command accepts positional args?
		cmdCfg := argh.NewCommandConfig()

		mapFlags(cmdCfg.Flags, cmd.Flags, false)

		cCfg.Commands.Set(cmd.Name, *cmdCfg)

		if len(cmd.Subcommands) == 0 {
			return
		}

		mapCommands(cmdCfg, cmd.Subcommands)
	}
}
