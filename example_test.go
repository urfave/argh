package argh_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/urfave/argh"
)

func ExampleParserConfig_simple() {
	state := map[string]argh.CommandFlag{}

	pCfg := argh.NewParserConfig()
	pCfg.Prog.NValue = argh.OneOrMoreValue
	pCfg.Prog.ValueNames = []string{"val"}
	pCfg.Prog.On = func(cf argh.CommandFlag) {
		state["prog"] = cf

		fmt.Printf("prog Name: %[1]q\n", cf.Name)
		fmt.Printf("prog Values: %[1]q\n", cf.Values)
		fmt.Printf("prog len(Nodes): %[1]v\n", len(cf.Nodes))
	}

	pCfg.Prog.SetFlagConfig("a", &argh.FlagConfig{
		NValue: 2,
		On: func(cf argh.CommandFlag) {
			state["a"] = cf

			fmt.Printf("prog -a Name: %[1]q\n", cf.Name)
			fmt.Printf("prog -a Values: %[1]q\n", cf.Values)
			fmt.Printf("prog -a len(Nodes): %[1]v\n", len(cf.Nodes))
		},
	})

	pCfg.Prog.SetFlagConfig("b", &argh.FlagConfig{
		Persist: true,
		On: func(cf argh.CommandFlag) {
			state["b"] = cf

			fmt.Printf("prog -b Name: %[1]q\n", cf.Name)
			fmt.Printf("prog -b Values: %[1]q\n", cf.Values)
			fmt.Printf("prog -b len(Nodes): %[1]v\n", len(cf.Nodes))
		},
	})

	pCfg.Prog.SetCommandConfig("sub", &argh.CommandConfig{
		NValue:     3,
		ValueNames: []string{"pilot", "navigator", "comms"},
		On: func(cf argh.CommandFlag) {
			state["sub"] = cf

			fmt.Printf("prog sub Name: %[1]q\n", cf.Name)
			fmt.Printf("prog sub Values: %[1]q\n", cf.Values)
			fmt.Printf("prog sub len(Nodes): %[1]v\n", len(cf.Nodes))
		},
	})

	// simulate command line args
	os.Args = []string{
		"hello", "-a=from", "the", "ether", "sub", "marge", "patty", "selma", "-b",
	}

	pt, err := argh.ParseArgs(os.Args, pCfg)
	if err != nil {
		argh.PrintParserError(os.Stderr, err)
		os.Exit(86)
		return
	}

	if pt == nil {
		log.Fatal("no parse tree?")
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	fmt.Printf("state: ")

	if err := enc.Encode(state); err != nil {
		log.Fatalf("failed to jsonify: %v", err)
	}

	// Output:
	// prog -a Name: "a"
	// prog -a Values: map["0":"from" "1":"the"]
	// prog -a len(Nodes): 4
	// prog -b Name: "b"
	// prog -b Values: map[]
	// prog -b len(Nodes): 0
	// prog sub Name: "sub"
	// prog sub Values: map["comms":"selma" "navigator":"patty" "pilot":"marge"]
	// prog sub len(Nodes): 8
	// prog Name: "hello"
	// prog Values: map["val":"ether"]
	// prog len(Nodes): 6
	// state: {
	//   "a": {
	//     "Name": "a",
	//     "Values": {
	//       "0": "from",
	//       "1": "the"
	//     },
	//     "Nodes": [
	//       {},
	//       {
	//         "Literal": "from"
	//       },
	//       {},
	//       {
	//         "Literal": "the"
	//       }
	//     ]
	//   },
	//   "b": {
	//     "Name": "b",
	//     "Values": null,
	//     "Nodes": null
	//   },
	//   "prog": {
	//     "Name": "hello",
	//     "Values": {
	//       "val": "ether"
	//     },
	//     "Nodes": [
	//       {},
	//       {
	//         "Name": "a",
	//         "Values": {
	//           "0": "from",
	//           "1": "the"
	//         },
	//         "Nodes": [
	//           {},
	//           {
	//             "Literal": "from"
	//           },
	//           {},
	//           {
	//             "Literal": "the"
	//           }
	//         ]
	//       },
	//       {},
	//       {
	//         "Literal": "ether"
	//       },
	//       {},
	//       {
	//         "Name": "sub",
	//         "Values": {
	//           "comms": "selma",
	//           "navigator": "patty",
	//           "pilot": "marge"
	//         },
	//         "Nodes": [
	//           {},
	//           {
	//             "Literal": "marge"
	//           },
	//           {},
	//           {
	//             "Literal": "patty"
	//           },
	//           {},
	//           {
	//             "Literal": "selma"
	//           },
	//           {},
	//           {
	//             "Name": "b",
	//             "Values": null,
	//             "Nodes": null
	//           }
	//         ]
	//       }
	//     ]
	//   },
	//   "sub": {
	//     "Name": "sub",
	//     "Values": {
	//       "comms": "selma",
	//       "navigator": "patty",
	//       "pilot": "marge"
	//     },
	//     "Nodes": [
	//       {},
	//       {
	//         "Literal": "marge"
	//       },
	//       {},
	//       {
	//         "Literal": "patty"
	//       },
	//       {},
	//       {
	//         "Literal": "selma"
	//       },
	//       {},
	//       {
	//         "Name": "b",
	//         "Values": null,
	//         "Nodes": null
	//       }
	//     ]
	//   }
	// }
}
