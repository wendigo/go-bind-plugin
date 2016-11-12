package cli

import (
	"fmt"
	"os"
	"plugin"
	"sort"
	"strings"
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
	Imports      []string
	NamedImports map[string]string

	importsNames map[string]string
}

func loadPlugin(path string, imports []string) (*pluginStructure, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open plugin: %s", err)
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("could not stat file: %s", err)
	}

	shaSum, err := fileChecksum(path)
	if err != nil {
		return nil, fmt.Errorf("could not calculate plugin checksum: %s", err)
	}

	plug := (*_Plugin)(unsafe.Pointer(p))

	ps := &pluginStructure{
		Symbols:      plug.symbols,
		Package:      plug.path,
		Size:         stat.Size(),
		Sha256:       shaSum,
		importsNames: make(map[string]string),
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

	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})

	return ret
}

func (p *pluginStructure) isNamedImport(pkg string, importName string) bool {
	return pkg[strings.LastIndex(pkg, "/")+1:] != importName
}
