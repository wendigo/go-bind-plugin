package main

import (
	"flag"
	"os"

	"log"

	"github.com/wendigo/go-bind-plugin/cli"
)

func main() {
	config := cli.Config{}

	flag.StringVar(&config.PluginPath, "plugin-path", "", "Path to plugin (.so file)")
	flag.StringVar(&config.PluginPackage, "plugin-package", "", "Plugin package url (as accepted by go get)")
	flag.StringVar(&config.OutputName, "output-name", "PluginAPI", "Output struct name")
	flag.StringVar(&config.OutputPath, "output-path", "plugin_api.go", "Output file path")
	flag.BoolVar(&config.DereferenceVariables, "dereference-vars", false, "Dereference plugin variables")
	flag.BoolVar(&config.CheckSha256, "sha256", false, "Write plugin's sha256 checksum to wrapper and validate it when loading it")
	flag.BoolVar(&config.FormatCode, "format", true, "Format generated output file with gofmt")
	flag.BoolVar(&config.ForcePluginRebuild, "rebuild", false, "Rebuild plugin on every run")
	flag.StringVar(&config.OutputPackage, "output-package", "", "Output package (can be derived from output-path)")
	flag.Parse()

	logger := log.New(os.Stderr, "go-bind-plugin ", log.Ltime)
	cli, err := cli.New(config, logger)

	if err != nil {
		logger.Fatal(err)
	}

	if err := cli.GenerateFile(); err != nil {
		logger.Fatal(err)
	}
}
