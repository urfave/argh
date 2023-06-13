package argh_test

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/argh"
)

func ExampleParserConfig_simple() {
	pCfg := argh.NewParserConfig()
	pCfg.Prog = &argh.CommandConfig{
		NValue:     argh.OneOrMoreValue,
		ValueNames: []string{"val"},
		Flags: &argh.Flags{
			Automatic: true,
		},
	}

	// simulate command line args
	os.Args = []string{"hello", "there"}

	pt, err := argh.ParseArgs(os.Args, pCfg)
	if err != nil {
		argh.PrintParserError(os.Stderr, err)
		os.Exit(86)
		return
	}

	if pt == nil {
		log.Fatal("no parse tree?")
	}

	fmt.Printf("parsed!\n")

	// Output:
	// parsed!
}
