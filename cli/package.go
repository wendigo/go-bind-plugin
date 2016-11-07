package cli

import "unsafe"
import "plugin"

type packageSymbols map[string]interface{}

type _Plugin struct {
	path    string
	_       chan struct{}
	symbols packageSymbols
}

type pluginStructure struct {
	path      string
	functions []*function
	variables []*variable
}

type function struct {
}

type variable struct {
}

type packageImport struct {
}

func loadPlugin(path string) (*pluginStructure, error) {
	p, err := plugin.Open(path)

	if err != nil {
		return nil, err
	}

	structure := (*_Plugin)(unsafe.Pointer(p))

	plug, err := inspectPlugin(structure)

	if err != nil {
		return nil, err
	}

	return plug, nil
}

func inspectPlugin(p *_Plugin) (*pluginStructure, error) {

	return &pluginStructure{
		path:      p.path,
		functions: getFunctions(p.symbols),
		variables: getVariables(p.symbols),
	}, nil
}

func getFunctions(symbols packageSymbols) []*function {
	return make([]*function, 0)
}

func getVariables(symbols packageSymbols) []*variable {
	return make([]*variable, 0)
}

func (p *pluginStructure) Path() string {
	return p.path
}

func (p *pluginStructure) Functions() []*function {
	ret := make([]*function, 0)

	return ret
}

func (p *pluginStructure) Variables() []*variable {
	ret := make([]*variable, 0)

	return ret
}
