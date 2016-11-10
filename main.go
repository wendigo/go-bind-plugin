package main

import (
	"flag"
	"os"

	"log"

	"github.com/wendigo/go-bind-plugin/cli"
)

func main() {
	logger := log.New(os.Stderr, "go-bind-plugin ", log.Ltime)
	config := cli.Config{}

	flagset := flag.NewFlagSet("go-bind-plugin", flag.ContinueOnError)

	flagset.StringVar(&config.PluginPath, "plugin-path", "", "Path to plugin (.so file)")
	flagset.StringVar(&config.PluginPackage, "plugin-package", "", "Plugin package url (as accepted by go get)")
	flagset.StringVar(&config.OutputName, "output-name", "PluginAPI", "Output struct name")
	flagset.StringVar(&config.OutputPath, "output-path", "plugin_api.go", "Output file path")
	flagset.BoolVar(&config.DereferenceVariables, "dereference-vars", false, "Dereference plugin variables")
	flagset.BoolVar(&config.CheckSha256, "sha256", false, "Write plugin's sha256 checksum to wrapper and validate it when loading it")
	flagset.BoolVar(&config.FormatCode, "format", true, "Format generated output file with gofmt")
	flagset.BoolVar(&config.ForcePluginRebuild, "rebuild", false, "Rebuild plugin on every run")
	flagset.StringVar(&config.OutputPackage, "output-package", "main", "Output package (can be derived from output-path)")
	flagset.BoolVar(&config.HideVariables, "hide-vars", false, "Do not export plugin variables")
	flagset.BoolVar(&config.AsInterface, "interface", false, "Generate and return interface instead of struct (turns on -hide-vars)")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		logger.Fatal(err)
	}

	cli, err := cli.New(config, logger)

	if err != nil {
		logger.Fatal(err)
	}

	if err := cli.GenerateFile(); err != nil {
		logger.Fatal(err)
	}
}
