# go-bind-plugin [![Build Status](https://travis-ci.org/wendigo/go-bind-plugin.svg?branch=master)](https://travis-ci.org/wendigo/go-bind-plugin)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/wendigo/go-bind-plugin/badge.svg?branch=master)](https://coveralls.io/github/wendigo/go-bind-plugin?branch=master)&nbsp;[![GoDoc](https://godoc.org/github.com/wendigo/go-bind-plugin/cli?status.svg)](https://godoc.org/github.com/wendigo/go-bind-plugin/cli)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/wendigo/go-bind-plugin)](https://goreportcard.com/report/github.com/wendigo/go-bind-plugin)

> TL;DR: See end-to-end example in [go-bind-plugin-example](https://github.com/wendigo/go-bind-plugin-example).


**go-bind-plugin** is `go:generate` tool for building [golang 1.8 plugins](https://tip.golang.org/pkg/plugin) and generating wrappers around exported symbols (functions and variables).

## What is it?

**go-bind-plugin** generates neat API around symbols exported by a `*.so` plugin built with `go build -buildmode=plugin` in upcoming go 1.8. [plugin.Plugin](https://tip.golang.org/pkg/plugin/#Plugin) holds information about exported symbols as `map[string]interface{}`.

**go-bind-plugins** uses reflection to find out actual types of symbols and generates typed API wrapping plugin with additional functionalities (like dereferencing exported variables and checking SHA256 sum).

*Note: Basic usage does not require plugin sources as wrapper can be generated using only* `*.so` *file.*

## Why should I use it?

In example if plugin exports `func AddTwoInts(a, b int) int` and `var BuildVersion string` instead of using [Plugin.Lookup](https://tip.golang.org/pkg/plugin/#Plugin.Lookup) directly:

```go
plug, err := plugin.Open("plugin.so")

if err != nil {
  panic(err)
}

symbol, err := plug.Lookup("AddTwoInts")
if err != nil {
  panic("AddTwoInts was not found in a plugin")
}

if typed, ok := symbol.(func(int, int) int); ok {
  result := typed(10, 20)
} else {
  panic("AddTwoInts has different type than exported by plugin")
}

symbol, err := plug.Lookup("BuildVersion")
if err != nil {
  panic("BuildVersion was not found in a plugin")
}

if typed, ok := symbol.(*string); ok {
  fmt.Println(*typed)
} else {
  panic("BuildVersion is not a string reference")
}
```

you can just simply do:

```go
plug, err := pluginapi.BindPluginAPI("plugin.so") // plug is *plugin_api.PluginAPI

if err != nil {
  panic(err)
}

result := plug.AddTwoInts(10, 20)
fmt.Println(plug.BuildVersion) // or fmt.Println(*plug.BuildVersion) if -dereference-vars is not used
```

`pluginapi.BindPluginAPI()` ensures that plugin exports required symbols and their types are correct.

## Usage

```
go get -u github.com/wendigo/go-bind-plugin
go-bind-plugin -help

Usage of go-bind-plugin:
  -dereference-vars
    	Dereference plugin variables
  -format
    	Format generated output file with gofmt (default true)
  -hide-vars
    	Do not export plugin variables
  -interface
    	Generate and return interface instead of struct (turns on -hide-vars)
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
//go:generate go-bind-plugin -format -plugin-package github.com/plugin_test/plug -rebuild -sha256 -dereference-vars -output-name TestPlugin -output-path tmp/plugin.go -plugin-path tmp/plugin.so -output-package pluginapi
`

**go-bind-plugin** will do following things on invocation:

- build plugin to `tmp/plugin.so` (even if plugin exists it will be rebuilded) from `github.com/plugin_test/plug` source (must exist in $GOPATH or vendor/)
- generate wrapper struct `wrapper.TestPlugin` in `tmp/plugin.go`
- dereference variables in the generated wrapper
- format generated wrapper with `gofmt -s -w`
- write sha256 checksum to `tmp/plugin.go` that will be checked when loading plugin with `pluginapi.BindTestPlugin(path string) (*TestPlugin, error)`

### Wrapper API example (for -output-name "PluginAPI")

`BindPluginAPI(path string) (*PluginAPI, error)` - loads plugin from `path` and wraps it with `type PluginAPI struct`:
  - all functions exported by the plugin are exposed as methods on struct `PluginAPI`
  - all variables exported by the plugin are exposed as fields on struct `PluginAPI` (if `-dereference-vars` is used fields are not references to plugin's variables)

`func (*PluginAPI) String() string` - provides nice textual representation of the wrapper

### Wrapper as interface

When `-interface` is used instead of generating and returning `struct` interface containing all exported symbols is generated. This eases mocking and working with multiple plugins exporting the same API. 

**Note** that `-interface` effectively enables `-hide-vars` so variables won't be exported from the plugin.

### Generated code quality

Generated code passes both `go vet` and `golint` and can be formatted using `gofmt -s -w`. Exported symbols names are not changed in any way so names not following [go naming convention](https://golang.org/doc/effective_go.html) will still be reported by `golint` as invalid.

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
```

## Plugin call overhead

Using `-buildmode=plugin` with generated wrapper seems not to add overhead when calling methods on a wrapper (creating plugin instance and loading `*.so` file is constant cost).

```go
BenchmarkCallOverhead/plugin-8         	30000000	        58.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkCallOverhead/plugin-8         	30000000	        59.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkCallOverhead/plugin-8         	30000000	        54.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkCallOverhead/native-8         	20000000	        59.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkCallOverhead/native-8         	20000000	        59.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkCallOverhead/native-8         	20000000	        59.7 ns/op	       0 B/op	       0 allocs/op
```
