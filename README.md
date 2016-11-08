# go-bind-plugin [![Build Status](https://travis-ci.org/wendigo/go-bind-plugin.svg?branch=master)](https://travis-ci.org/wendigo/go-bind-plugin)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/wendigo/go-bind-plugin/badge.svg?branch=master)](https://coveralls.io/github/wendigo/go-bind-plugin?branch=master)&nbsp;[![GoDoc](https://godoc.org/github.com/wendigo/go-bind-plugin/cli?status.svg)](https://godoc.org/github.com/wendigo/go-bind-plugin/cli)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/wendigo/go-bind-plugin)](https://goreportcard.com/report/github.com/wendigo/go-bind-plugin)

> TL;DR: See end-to-end example in [go-bind-plugin-example](https://github.com/wendigo/go-bind-plugin-example).


**go-bind-plugin** is `go:generate` tool for building [golang 1.8 plugins](https://tip.golang.org/pkg/plugin) and generating wrappers around exported symbols (functions and variables).

## Usage

```
go get -u github.com/wendigo/go-bind-plugin
go-bind-plugin -help
```

Available flags:

```
Usage of go-bind-plugin:
  -dereference-vars
    	Dereference plugin variables
  -format
    	Format generated output file with gofmt (default true)
  -output-name string
    	Output struct name (default "PluginAPI")
  -output-package string
    	Output package (can be derived from output-path) (default "main")
  -output-path string
    	Output file path (default "plugin_api.go")
  -plugin-package string
    	Plugin package url (as accepted by go get)
  -plugin-path string
    	Path to plugin (.so file)
  -rebuild
    	Rebuild plugin on every run
  -sha256
    	Write plugin's sha256 checksum to wrapper and validate it when loading it
```

### Example

`
//go:generate go-bind-plugin -format -plugin-package github.com/plugin_test/plug -rebuild -sha256 -dereference-vars -output-name TestPlugin -output-path tmp/plugin.go -plugin-path tmp/plugin.so -output-package wrapper
`

**go-bind-plugin** will do following things on invocation:

- build plugin to `tmp/plugin.so` (even if plugin exists it will be rebuilded) from package `github.com/plugin_test/plug` (must exist in $GOPATH or vendor/)
- create wrapper struct `wrapper.TestPlugin` in `tmp/plugin.go`
- dereference variables exposed by the plugin in the generated wrapper
- format generated code with `gofmt -s -w`
- write sha256 checksum to `tmp/plugin.go` that will be validated when plugin is loaded via `wrapper.BindTestPlugin(path string) (*TestPlugin, error)`

### Wrapper API example (for -output-name "PluginAPI")

`BindPluginAPI(path string) (*PluginAPI, error)` - loads plugin from `path` and wraps it with `type PluginAPI struct {}`:
  - all functions exposed in the plugin are exposed as methods on struct `PluginAPI`
  - all variables references exposed in the plugin are exposed as fields on struct `PluginAPI` (if `-dereference-vars` is used fields are not references to plugin's variables)

`func (*PluginAPI) String() string` - provides nice textual representation of the wrapper

### Example generated wrapper information

```
Wrapper info:
	- Generated on: 2016-11-08 16:15:07.513150982 +0100 CET
	- Command: go-bind-plugin -plugin-path ./internal/test_fixtures/generated/basic_plugin/plugin.so -plugin-package ./internal/test_fixtures/basic_plugin -output-name TestWrapper -output-path ./internal/test_fixtures/generated/basic_plugin/plugin.go -output-package main -sha256 true -format true -rebuild true

Plugin info:
	- package: github.com/wendigo/go-bind-plugin/internal/test_fixtures/basic_plugin
	- sha256 sum: 55aa13402686f3200f5067604c04ce8d365e7cf2095d8278b2ff52ae26df7e6d
	- size: 1232572 bytes

Exported functions (3):
	- ReturningInt32 func() (int32)
	- ReturningStringSlice func() ([]string)
	- ReturningIntArray func() ([3]int32)

Exported variables (0):

Plugin imports:
```
