// Program to create, update & delete aspect schemas in Magda
package main

import (
	"os"

	"github.com/maxott/magda-cli/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := cmd.App()
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
