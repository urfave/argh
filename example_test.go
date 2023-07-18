package argh_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/urfave/argh"
)

func ExampleParserConfig() {
	cmdState := map[string]argh.Command{}
	flagState := map[string]argh.Flag{}

	pCfg := argh.NewParserConfig()
	pCfg.Prog.NValue = argh.OneOrMoreValue
	pCfg.Prog.ValueNames = []string{"val"}
	pCfg.Prog.On = func(cf argh.Command) error {
		cmdState["prog"] = cf

		fmt.Printf("prog Name: %[1]q\n", cf.Name)
		fmt.Printf("prog Values: %[1]q\n", cf.Values)
		fmt.Printf("prog len(Nodes): %[1]v\n", len(cf.Nodes))

		return nil
	}

	pCfg.Prog.SetFlagConfig("a", &argh.FlagConfig{
		NValue: 2,
		On: func(cf argh.Flag) error {
			flagState["a"] = cf

			fmt.Printf("prog -a Name: %[1]q\n", cf.Name)
			fmt.Printf("prog -a Values: %[1]q\n", cf.Values)
			fmt.Printf("prog -a len(Nodes): %[1]v\n", len(cf.Nodes))

			return nil
		},
	})

	pCfg.Prog.SetFlagConfig("b", &argh.FlagConfig{
		Persist: true,
		On: func(cf argh.Flag) error {
			flagState["b"] = cf

			fmt.Printf("prog -b Name: %[1]q\n", cf.Name)
			fmt.Printf("prog -b Values: %[1]q\n", cf.Values)
			fmt.Printf("prog -b len(Nodes): %[1]v\n", len(cf.Nodes))

			return nil
		},
	})

	sub := &argh.CommandConfig{
		NValue:     3,
		ValueNames: []string{"pilot", "navigator", "comms"},
		On: func(cf argh.Command) error {
			cmdState["sub"] = cf

			fmt.Printf("prog sub Name: %[1]q\n", cf.Name)
			fmt.Printf("prog sub Values: %[1]q\n", cf.Values)
			fmt.Printf("prog sub len(Nodes): %[1]v\n", len(cf.Nodes))

			return nil
		},
	}

	sub.SetFlagConfig("c", &argh.FlagConfig{
		NValue:     1,
		ValueNames: []string{"feels"},
		On: func(cf argh.Flag) error {
			flagState["c"] = cf

			fmt.Printf("prog sub -c Name: %[1]q\n", cf.Name)
			fmt.Printf("prog sub -c Values: %[1]q\n", cf.Values)
			fmt.Printf("prog sub -c len(Nodes): %[1]v\n", len(cf.Nodes))

			return nil
		},
	})

	pCfg.Prog.SetCommandConfig("sub", sub)

	// simulate command line args
	os.Args = []string{
		"hello", "-a=from", "the", "ether", "sub", "marge", "-c=hurlish", "patty", "selma", "-b",
	}

	if _, err := json.Marshal(pCfg.Prog); err != nil {
		log.Fatal(err)
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

	fmt.Printf("command state: ")

	if err := enc.Encode(cmdState); err != nil {
		log.Fatalf("failed to jsonify: %v", err)
	}

	fmt.Printf("flag state: ")

	if err := enc.Encode(flagState); err != nil {
		log.Fatalf("failed to jsonify: %v", err)
	}

	// Output:
	// prog -a Name: "a"
	// prog -a Values: [{"0" "from"} {"1" "the"}]
	// prog -a len(Nodes): 4
	// prog sub -c Name: "c"
	// prog sub -c Values: [{"feels" "hurlish"}]
	// prog sub -c len(Nodes): 2
	// prog -b Name: "b"
	// prog -b Values: []
	// prog -b len(Nodes): 0
	// prog sub Name: "sub"
	// prog sub Values: [{"pilot" "marge"} {"navigator" "patty"} {"comms" "selma"}]
	// prog sub len(Nodes): 10
	// prog Name: "hello"
	// prog Values: [{"val" "ether"}]
	// prog len(Nodes): 6
	// command state: {
	//   "prog": {
	//     "Name": "hello",
	//     "Values": [
	//       {
	//         "Key": "val",
	//         "Value": "ether"
	//       }
	//     ],
	//     "Nodes": [
	//       {},
	//       {
	//         "Name": "a",
	//         "Values": [
	//           {
	//             "Key": "0",
	//             "Value": "from"
	//           },
	//           {
	//             "Key": "1",
	//             "Value": "the"
	//           }
	//         ],
	//         "Nodes": [
	//           {},
	//           {
	//             "Value": "from"
	//           },
	//           {},
	//           {
	//             "Value": "the"
	//           }
	//         ]
	//       },
	//       {},
	//       {
	//         "Value": "ether"
	//       },
	//       {},
	//       {
	//         "Name": "sub",
	//         "Values": [
	//           {
	//             "Key": "pilot",
	//             "Value": "marge"
	//           },
	//           {
	//             "Key": "navigator",
	//             "Value": "patty"
	//           },
	//           {
	//             "Key": "comms",
	//             "Value": "selma"
	//           }
	//         ],
	//         "Nodes": [
	//           {},
	//           {
	//             "Value": "marge"
	//           },
	//           {},
	//           {
	//             "Name": "c",
	//             "Values": [
	//               {
	//                 "Key": "feels",
	//                 "Value": "hurlish"
	//               }
	//             ],
	//             "Nodes": [
	//               {},
	//               {
	//                 "Value": "hurlish"
	//               }
	//             ]
	//           },
	//           {},
	//           {
	//             "Value": "patty"
	//           },
	//           {},
	//           {
	//             "Value": "selma"
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
	//     "Values": [
	//       {
	//         "Key": "pilot",
	//         "Value": "marge"
	//       },
	//       {
	//         "Key": "navigator",
	//         "Value": "patty"
	//       },
	//       {
	//         "Key": "comms",
	//         "Value": "selma"
	//       }
	//     ],
	//     "Nodes": [
	//       {},
	//       {
	//         "Value": "marge"
	//       },
	//       {},
	//       {
	//         "Name": "c",
	//         "Values": [
	//           {
	//             "Key": "feels",
	//             "Value": "hurlish"
	//           }
	//         ],
	//         "Nodes": [
	//           {},
	//           {
	//             "Value": "hurlish"
	//           }
	//         ]
	//       },
	//       {},
	//       {
	//         "Value": "patty"
	//       },
	//       {},
	//       {
	//         "Value": "selma"
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
	// flag state: {
	//   "a": {
	//     "Name": "a",
	//     "Values": [
	//       {
	//         "Key": "0",
	//         "Value": "from"
	//       },
	//       {
	//         "Key": "1",
	//         "Value": "the"
	//       }
	//     ],
	//     "Nodes": [
	//       {},
	//       {
	//         "Value": "from"
	//       },
	//       {},
	//       {
	//         "Value": "the"
	//       }
	//     ]
	//   },
	//   "b": {
	//     "Name": "b",
	//     "Values": null,
	//     "Nodes": null
	//   },
	//   "c": {
	//     "Name": "c",
	//     "Values": [
	//       {
	//         "Key": "feels",
	//         "Value": "hurlish"
	//       }
	//     ],
	//     "Nodes": [
	//       {},
	//       {
	//         "Value": "hurlish"
	//       }
	//     ]
	//   }
	// }
}
