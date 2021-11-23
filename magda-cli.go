// Program to create, update & delete aspect schemas in Magda
package main

import (
	"fmt"
	"os"

	"github.com/maxott/magda-cli/cmd"
	"github.com/maxott/magda-cli/pkg/adapter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/alecthomas/kingpin.v2"
)

var Version string

func main() {
	app := cmd.App()
	app.HelpFlag.Short('h')

	setLogger(zapcore.ErrorLevel)
	app.Flag("verbose", "Be chatty [MAGDA_VERBOSE]").Short('v').Envar("MAGDA_VERBOSE").Action(setVerbose).Bool()
	app.Flag("debug", "Be very chatty [MAGDA_DEBUG]").Short('d').Envar("MAGDA_DEBUG").Action(setDebug).Bool()
	app.Flag("version", "Print out version").Action(printVersion).Bool()

	// app.PreAction(configLogger)
	_, err := app.Parse(os.Args[1:])
	if err != nil {
		//var e interface{} = err
		if aerr, ok := err.(adapter.IAdapterError); ok {
			fmt.Printf("ERROR: %s\n", aerr.Error())
		}
	}
	cmd.Logger().Sync()
}

func setLogger(logLevel zapcore.Level) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(logLevel)
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	cmd.SetLogger(logger)
}

func setVerbose(c *kingpin.ParseContext) error {
	setLogger(zapcore.InfoLevel)
	return nil
}

func setDebug(c *kingpin.ParseContext) error {
	setLogger(zapcore.DebugLevel)
	return nil
}

func printVersion(c *kingpin.ParseContext) error {
	fmt.Printf("Version: %s\n", Version)
	os.Exit(0)
	return nil
}
