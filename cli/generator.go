package cli

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	PluginPath    string
	PluginPackage string
	OutputPath    string
	OutputPackage string
	OutputName    string
}

type Cli struct {
	config Config
	logger *log.Logger
}

func New(config Config) *Cli {
	return &Cli{
		config: config,
		logger: log.New(os.Stderr, "", 0),
	}
}

func (c *Cli) Generate() error {
	structure, _ := loadPlugin(c.config.PluginPath)
	fmt.Println(structure)

	return nil
}
