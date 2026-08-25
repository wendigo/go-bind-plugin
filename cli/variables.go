package cli

import (
	"reflect"
	"slices"
	"strings"
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
		if typ.Kind() == reflect.Pointer {
			p.Variables = append(p.Variables, &variable{
				Name:      name,
				Typ:       typ.Elem().String(),
				Signature: p.getVariableSignature(typ.Elem(), false),
			})
		}
	}

	slices.SortFunc(p.Variables, func(a, b *variable) int {
		return strings.Compare(a.Name, b.Name)
	})

	return nil
}
