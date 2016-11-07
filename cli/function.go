package cli

import (
	"fmt"
	"reflect"
	"strings"
)

type function struct {
	Name           string
	Signature      string
	FunctionHeader string
	ArgumentsCount int
	IsVariadic     bool
}

func (f *function) String() string {
	return f.Name
}

func (f *function) TrimmedSignature() string {
	return strings.TrimLeft(f.FunctionHeader, "func ")
}

func (f *function) ArgumentsCall() string {
	var arguments []string

	for i := 0; i < f.ArgumentsCount; i++ {
		if f.IsVariadic && i == f.ArgumentsCount-1 {
			arguments = append(arguments, fmt.Sprintf("in%d...", i))
		} else {
			arguments = append(arguments, fmt.Sprintf("in%d", i))
		}
	}

	return strings.Join(arguments, ", ")
}

func (p *pluginStructure) analyzeFunctions() error {
	for name, pointer := range p.Symbols {
		typ := reflect.TypeOf(pointer)

		if typ.Kind() == reflect.Func {
			p.Functions = append(p.Functions, &function{
				Name:           name,
				Signature:      p.getFunctionSignature(typ, false),
				FunctionHeader: p.getFunctionSignature(typ, true),
				ArgumentsCount: typ.NumIn(),
				IsVariadic:     typ.IsVariadic(),
			})
		}
	}

	return nil
}
