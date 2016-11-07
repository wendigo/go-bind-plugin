package cli

import (
	"fmt"
	"os"
	"plugin"
	"unsafe"
)

type packageSymbols map[string]interface{}

type _Plugin struct {
	path    string
	_       chan struct{}
	symbols packageSymbols
}

type pluginStructure struct {
	Package string
	Size    int64
	Sha256  string
	Symbols packageSymbols

	Functions    []*function
	Variables    []*variable
	ImportsNames map[string]string
}

func loadPlugin(path string, imports []string) (*pluginStructure, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Could not open plugin: %s", err)
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("Could not check file size: %s", err)
	}

	shaSum, err := fileChecksum(path)
	if err != nil {
		return nil, fmt.Errorf("Could not calculate plugin checksum: %s", err)
	}

	plug := (*_Plugin)(unsafe.Pointer(p))

	ps := &pluginStructure{
		Symbols:      plug.symbols,
		Package:      plug.path,
		ImportsNames: make(map[string]string),
		Size:         stat.Size(),
		Sha256:       shaSum,
	}

	for _, pkg := range imports {
		ps.getNamedPkgImport(pkg)
	}

	err2 := ps.analyze()

	if err2 != nil {
		return nil, err2
	}

	return ps, nil
}

func (p *pluginStructure) String() string {
	return fmt.Sprintf("plugin %s\n{\n\tfunctions=%s,\n\tvariables=%s\n}\n", p.Package, p.Functions, p.Variables)
}

func (p *pluginStructure) SymbolsLen() int {
	return len(p.Symbols)
}

func (p *pluginStructure) analyze() error {
	if err := p.analyzeVariables(); err != nil {
		return err
	}

	if err := p.analyzeFunctions(); err != nil {
		return err
	}

	return nil
}
