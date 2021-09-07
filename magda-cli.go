// Program to create, update & delete aspect schemas in Magda
package main

import (
	"os"

	"github.com/maxott/magda-cli/cmd"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := cmd.App()
	app.HelpFlag.Short('h')

	setLogger(zapcore.WarnLevel) 
	app.Flag("verbose", "Be chatty [MAGDA_VERBOSE]").Short('v').Envar("MAGDA_VERBOSE").Action(setVerbose).Bool()
	app.Flag("debug", "Be very chatty [MAGDA_DEBUG]").Short('d').Envar("MAGDA_DEBUG").Action(setDebug).Bool()
	

	// app.PreAction(configLogger)
	command, err := app.Parse(os.Args[1:])
	kingpin.MustParse(command, err)
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
