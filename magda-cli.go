// Program to create, update & delete aspect schemas in Magda
package main

import (
	"fmt"
	"github.com/maxott/magda-cli/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

// var verbose

func main() {
	app := cmd.App()
	app.HelpFlag.Short('h')

	app.Flag("verbose", "Be chatty [MAGDA_VERBOSE]").Short('v').Envar("MAGDA_VERBOSE").Action(setVerbose).Bool()

	// app.PreAction(configLogger)
	command, err := app.Parse(os.Args[1:])
	kingpin.MustParse(command, err)
}

func setVerbose(c *kingpin.ParseContext) error {
	fmt.Printf(">>>>>>> SET VERBOSE\n")
	// //fmt.Printf(">>>>>>> VERBOSE %v\n\n", c.Elements)
	// for _, value := range c.Elements {
	// 	fmt.Printf(">>>>>>> VERBOSE %v - %s\n", value.Clause, value.Value)
	// }
	return nil
}
