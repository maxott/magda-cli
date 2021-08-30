// Program to create, update & delete aspect schemas in Magda
package main

import (
	"github.com/maxott/magda-cli/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	lgrus "github.com/sirupsen/logrus"

	"github.com/maxott/magda-cli/pkg/log/logrus"
)

func main() {
	app := cmd.App()
	app.HelpFlag.Short('h')

	app.Flag("verbose", "Be chatty [MAGDA_VERBOSE]").Short('v').Envar("MAGDA_VERBOSE").Action(setVerbose).Bool()
	app.Flag("debug", "Be very chatty [MAGDA_DEBUG]").Short('d').Envar("MAGDA_DEBUG").Action(setDebug).Bool()

	// app.PreAction(configLogger)
	command, err := app.Parse(os.Args[1:])
	kingpin.MustParse(command, err)
}

func setVerbose(c *kingpin.ParseContext) error {
	cmd.SetLogger(logrus.NewSimpleLogger(lgrus.InfoLevel))
	return nil
}

func setDebug(c *kingpin.ParseContext) error {
	cmd.SetLogger(logrus.NewSimpleLogger(lgrus.DebugLevel))
	return nil
}
