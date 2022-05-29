package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"git.meatballhat.com/x/box-o-sand/argh"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	asJSON := os.Getenv("ARGH_OUTPUT_JSON") == "enabled"

	log.SetFlags(0)

	pt, err := argh.ParseArgs(os.Args, nil)
	if err != nil {
		log.Fatal(err)
	}

	ast := argh.NewQuerier(pt.Nodes).AST()

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
