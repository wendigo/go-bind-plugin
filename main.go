package main

import "github.com/wendigo/go-bind-plugin/cli"

func main() {
	cli := cli.New(cli.Config{
		PluginPath:    "plugin.so",
		PluginPackage: "github.com/wendigo/plugin_test",
		OutputName:    "PluginWrapper",
		OutputPath:    ".",
	})

	cli.Generate()
}
