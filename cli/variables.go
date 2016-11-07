package cli

import (
	"fmt"
	"reflect"
)

type variable struct {
	Name      string
	Typ       string
	PkgPath   string
	Signature string
}

func (v *variable) String() string {
	return fmt.Sprintf("name:%s, signature: %s", v.Name, v.Signature)
}

func (v *variable) Declaration() string {
	return fmt.Sprintf("%s *%s", v.Name, v.Signature)
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

	return nil
}
