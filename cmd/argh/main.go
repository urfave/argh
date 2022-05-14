package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"git.meatballhat.com/x/box-o-sand/argh"
)

func main() {
	log.SetFlags(0)

	ast, err := argh.ParseArgs(os.Args, nil)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(ast, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
