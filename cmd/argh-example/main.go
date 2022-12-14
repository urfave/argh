package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"git.meatballhat.com/x/argh"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	asJSON := os.Getenv("ARGH_OUTPUT_JSON") == "enabled"

	log.SetFlags(0)

	pCfg := argh.NewParserConfig()
	pCfg.Prog = &argh.CommandConfig{
		NValue:     argh.OneOrMoreValue,
		ValueNames: []string{"val"},
		Flags: &argh.Flags{
			Automatic: true,
		},
	}

	pt, err := argh.ParseArgs(os.Args, pCfg)
	if err != nil {
		log.Fatal(err)
	}

	ast := argh.ToAST(pt.Nodes)

	if asJSON {
		b, err := json.MarshalIndent(ast, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))

		return
	}

	spew.Dump(ast)
}
