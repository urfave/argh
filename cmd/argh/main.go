package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"git.meatballhat.com/x/box-o-sand/argh"
)

func main() {
	ast, err := argh.ParseArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(ast, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
