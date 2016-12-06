package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var wrapperImports = []string{
	"plugin",
	"reflect",
	"fmt",
	"strings",
}

var wrapperImportsChecksum = []string{
	"encoding/hex",
	"crypto/sha256",
	"os",
	"io",
}

// Config for Cli
type Config struct {
	// Path to the plugin
	PluginPath string
	// Plugin's package (will be used to compile plugin when PluginPath does not exist)
	PluginPackage string
	// Path where generated wrapper will be written
	OutputPath string

	// Name of the struct/interface that will be generated
	OutputName string

	// Dereference variables exported by plugin
	DereferenceVariables bool

	// Generate SHA256 for existing/compiled plugin and write it to the Bind* function
	CheckSha256 bool

	// Format generated code using gofmt -s -w
	FormatCode bool

	// Rebuild plugin from sources even if *.so exists at PluginPath
	ForcePluginRebuild bool

	// Package name of the generated file
	OutputPackage string

	// Do not export plugin's variables
	HideVariables bool

	// Generate interface and return it in Bind* function instead of struct (enables HideVariables)
	AsInterface bool
}

// String formats config as command line arguments
func (c Config) String() string {
	var commandLine []string

	commandLine = append(commandLine, fmt.Sprintf(
		"-plugin-path %s -plugin-package %s -output-name %s -output-path %s -output-package %s",
		c.PluginPath,
		c.PluginPackage,
		c.OutputName,
		c.OutputPath,
		c.OutputPackage,
	))

	if c.CheckSha256 {
		commandLine = append(commandLine, "-sha256")
	}

	if c.DereferenceVariables {
		commandLine = append(commandLine, "-dereference-vars")
	}

	if c.ForcePluginRebuild {
		commandLine = append(commandLine, "-rebuild")
	}

	if c.AsInterface {
		commandLine = append(commandLine, "-interface")
	}

	if c.HideVariables {
		commandLine = append(commandLine, "-hide-vars")
	}

	return strings.Join(commandLine, " ")
}

// Cli is responsible for generating plugin wrapper, can be initialized with New()
type Cli struct {
	config Config
	logger *log.Logger
}

type buildInfo struct {
	Date    string
	Command string
}

// New creates new plugin wrapper generator
func New(config Config, logger *log.Logger) (*Cli, error) {

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &Cli{
		config: config,
		logger: logger,
	}, nil
}

// GenerateFile generates wrapper file for plugin
func (c *Cli) GenerateFile() error {
	var imports = wrapperImports

	if c.config.CheckSha256 {
		imports = append(imports, wrapperImportsChecksum...)
	}

	if !c.pluginExists(c.config.PluginPath) || c.config.ForcePluginRebuild {
		if err := c.buildPluginFromSources(c.config.PluginPath, c.config.PluginPackage); err != nil {
			return fmt.Errorf("could not build plugin from sources: %s", err)
		}
	}

	c.logger.Printf("Loading and analyzing plugin from: %s", c.config.PluginPath)
	structure, err := loadPlugin(c.config.PluginPath, imports)
	if err != nil {
		return fmt.Errorf("could not load plugin from %s: %s", c.config.PluginPath, err)
	}

	if structure.SymbolsLen() == 0 {
		return fmt.Errorf("plugin %s does not export any symbols", c.config.PluginPath)
	}

	outputPackage := c.config.OutputPackage
	if outputPackage == "" {
		outputPackage = c.getOutputPackage(c.config.OutputPath)
	}

	outputFile, err := c.createOutputFile(c.config.OutputPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %s", err)
	}

	tpl, err := template.New("generate").Parse(generateFileTemplate)

	if err != nil {
		return err
	}

	c.logger.Printf("Generating output wrapper: %s...", c.config.OutputPath)
	if err := tpl.Execute(outputFile, struct {
		Config        Config
		Plugin        pluginStructure
		Build         buildInfo
		OutputPackage string
	}{
		Config: c.config,
		Plugin: *structure,
		Build: buildInfo{
			Date:    time.Now().String(),
			Command: c.buildCommandArgs(),
		},
		OutputPackage: outputPackage,
	}); err != nil {
		return err
	}

	if c.config.FormatCode {
		c.logger.Printf("Formatting generated file with gofmt -s -w %s", c.config.OutputPath)
		if err := c.formatOutputCode(c.config.OutputPath); err != nil {
			return fmt.Errorf("could not format output code: %s", err)
		}
	}

	c.logger.Printf("Generated wrapper %s in file %s", c.config.OutputName, c.config.OutputPath)

	return nil
}

func (c *Cli) getOutputPackage(path string) string {
	directory := c.getOutputDirectory(path)

	if directory == "." {
		return "main"
	}

	parts := strings.Split(directory, "/")
	return parts[len(parts)-1]
}

func (c *Cli) getOutputDirectory(path string) string {
	return filepath.Clean(strings.TrimSuffix(path, filepath.Base(path)))
}

func (c *Cli) createOutputDir(path string) error {
	if info, err := os.Stat(path); err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil && !info.IsDir() {
		return fmt.Errorf("output path %s exists and is not a directory", path)
	}

	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}

	return nil
}

func (c *Cli) createOutputFile(path string) (*os.File, error) {
	if err := c.createOutputDir(c.getOutputDirectory(path)); err != nil {
		return nil, err
	}

	return os.Create(path)
}

func (c *Cli) pluginExists(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func (c *Cli) formatOutputCode(path string) error {
	return exec.Command("gofmt", "-s", "-w", path).Run()
}

func (c *Cli) buildPluginFromSources(pluginPath string, pluginPackage string) error {
	c.logger.Printf("Building plugin %s from package %s", pluginPath, pluginPackage)

	// Check if plugin output path exists
	if err := c.createOutputDir(c.getOutputDirectory(pluginPath)); err != nil {
		return err
	}

	command := []string{"build", "-o", pluginPath, "-buildmode=plugin", pluginPackage}

	c.logger.Printf("Running: go %s", strings.Join(command, " "))

	cmd := exec.Command("go", command...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()

	return cmd.Run()
}

func (c *Cli) buildCommandArgs() string {
	return fmt.Sprintf("go-bind-plugin %s", c.config.String())
}

func validateConfig(config *Config) error {
	if config.PluginPath == "" && config.PluginPackage == "" {
		return fmt.Errorf("either PluginPath or PluginPackage must be provided")
	}

	if config.ForcePluginRebuild && config.PluginPackage == "" {
		return fmt.Errorf("PluginPackage must be provided in order to build a plugin")
	}

	if config.OutputName == "" {
		config.OutputName = "PluginWrapper"
	}

	if config.OutputPath == "" {
		config.OutputPath = "plugin_wrapper.go"
	}

	if config.PluginPackage != "" && config.PluginPath == "" {
		config.PluginPath = "plugin.so"
	}

	if config.AsInterface {
		config.HideVariables = true
	}

	return nil
}
