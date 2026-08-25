package cli

import (
	"fmt"
	"os"
	"plugin"
	"slices"
	"strings"
)

type packageSymbols map[string]any

type pluginStructure struct {
	Package string
	Size    int64
	Sha256  string
	Symbols packageSymbols

	Functions    []*function
	Variables    []*variable
	Imports      []string
	NamedImports map[string]string

	importsNames map[string]string
}

func loadPlugin(path string, imports []string) (*pluginStructure, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open plugin: %w (loading a plugin requires a cgo-enabled go-bind-plugin built with the same toolchain as the plugin; pass -plugin-package to analyze plugin sources instead)", err)
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("could not stat file: %w", err)
	}

	shaSum, err := fileChecksum(path)
	if err != nil {
		return nil, fmt.Errorf("could not calculate plugin checksum: %w", err)
	}

	pluginPackage, symbols, err := pluginSymbols(p, path)
	if err != nil {
		return nil, err
	}

	ps := &pluginStructure{
		Symbols:      symbols,
		Package:      pluginPackage,
		Size:         stat.Size(),
		Sha256:       shaSum,
		importsNames: map[string]string{},
	}

	for _, pkg := range imports {
		ps.getNamedPkgImport(pkg)
	}

	err2 := ps.analyze()

	if err2 != nil {
		return nil, err2
	}

	ps.Imports = ps.imports()
	ps.NamedImports = ps.namedImports()

	return ps, nil
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

func (p *pluginStructure) namedImports() map[string]string {
	var ret = make(map[string]string)

	for name, imp := range p.importsNames {
		if p.isNamedImport(imp, name) {
			ret[name] = imp
		}
	}

	return ret
}

func (p *pluginStructure) imports() []string {
	var ret []string

	for name, imp := range p.importsNames {
		if p.isNamedImport(imp, name) {
			ret = append(ret, fmt.Sprintf("%s %q", name, imp))
		} else {
			ret = append(ret, fmt.Sprintf("%q", imp))
		}
	}

	slices.Sort(ret)

	return ret
}

func (p *pluginStructure) isNamedImport(pkg string, importName string) bool {
	return pkg[strings.LastIndex(pkg, "/")+1:] != importName
}
