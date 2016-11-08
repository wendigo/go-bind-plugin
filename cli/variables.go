package cli

import (
	"reflect"
	"sort"
)

type variable struct {
	Name      string
	Typ       string
	PkgPath   string
	Signature string
}

func (p *pluginStructure) analyzeVariables() error {
	for name, pointer := range p.Symbols {
		typ := reflect.TypeOf(pointer)

		// variables are always pointers
		if typ.Kind() == reflect.Ptr {
			p.Variables = append(p.Variables, &variable{
				Name:      name,
				Typ:       typ.Elem().String(),
				Signature: p.getVariableSignature(typ.Elem(), false),
			})
		}
	}

	sort.Slice(p.Variables, func(i, j int) bool {
		return p.Variables[i].Name < p.Variables[j].Name
	})

	return nil
}
