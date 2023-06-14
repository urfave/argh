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
			Map: map[string]argh.FlagConfig{
				"a": {
					NValue: 2,
					On: func(cf argh.CommandFlag) {
						fmt.Printf("prog -a Name: %[1]q\n", cf.Name)
						fmt.Printf("prog -a Values: %[1]q\n", cf.Values)
						fmt.Printf("prog -a len(Nodes): %[1]v\n", len(cf.Nodes))
					},
				},
			},
		},
		On: func(cf argh.CommandFlag) {
			fmt.Printf("prog Name: %[1]q\n", cf.Name)
			fmt.Printf("prog Values: %[1]q\n", cf.Values)
			fmt.Printf("prog len(Nodes): %[1]v\n", len(cf.Nodes))
		},
	}

	// simulate command line args
	os.Args = []string{"hello", "-a=from", "the", "ether"}

	pt, err := argh.ParseArgs(os.Args, pCfg)
	if err != nil {
		argh.PrintParserError(os.Stderr, err)
		os.Exit(86)
		return
	}

	if pt == nil {
		log.Fatal("no parse tree?")
	}

	// Output:
	// prog -a Name: "a"
	// prog -a Values: map["0":"from" "1":"the"]
	// prog -a len(Nodes): 4
	// prog Name: "hello"
	// prog Values: map["val":"ether"]
	// prog len(Nodes): 4
}
