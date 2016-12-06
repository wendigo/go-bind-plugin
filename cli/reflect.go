package cli

import (
	"fmt"
	"reflect"
	"strings"
)

func (p *pluginStructure) getVariableSignature(variable reflect.Type, isVariadic bool) string {
	if variable.Name() != "" {
		if isVariadic {
			return fmt.Sprintf("...%s%s", p.getNamedPkgImport(variable.PkgPath()), variable.Name())
		}

		return fmt.Sprintf("%s%s", p.getNamedPkgImport(variable.PkgPath()), variable.Name())
	}

	switch variable.Kind() {
	case reflect.Func:
		return p.getFunctionSignature(variable, false)
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", variable.Len(), p.getVariableSignature(variable.Elem(), false))
	case reflect.Slice:
		if isVariadic {
			return fmt.Sprintf("...%s", p.getVariableSignature(variable.Elem(), false))
		}

		return fmt.Sprintf("[]%s", p.getVariableSignature(variable.Elem(), false))

	case reflect.Chan:
		return fmt.Sprintf("%s %s", variable.ChanDir().String(), p.getVariableSignature(variable.Elem(), false))
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", p.getVariableSignature(variable.Key(), false), p.getVariableSignature(variable.Elem(), false))
	case reflect.Ptr:
		return fmt.Sprintf("*%s", p.getVariableSignature(variable.Elem(), false))
	}

	return variable.String()

}

func (p *pluginStructure) getNamedPkgImport(pkg string) string {

	if pkg == "" {
		return ""
	}

	importName := pkg[strings.LastIndex(pkg, "/")+1:]
	nextIndex := 0

	for {
		var nextImportName string

		if nextIndex == 0 {
			nextImportName = importName
		} else {
			nextImportName = fmt.Sprintf("%s_%d", importName, nextIndex)
		}

		if current, ok := p.importsNames[nextImportName]; ok {
			// We already have this package
			if current == pkg {
				return nextImportName + "."
			}

			nextIndex++
			continue
		} else {
			p.importsNames[nextImportName] = pkg
			return nextImportName + "."
		}
	}
}

func (p *pluginStructure) getFunctionSignature(fun reflect.Type, namedParams bool) string {
	var in []string
	var out []string

	for i := 0; i < fun.NumIn(); i++ {
		if namedParams {
			in = append(in, fmt.Sprintf("in%d %s", i, p.getVariableSignature(fun.In(i), fun.IsVariadic() && i == fun.NumIn()-1)))
		} else {
			in = append(in, p.getVariableSignature(fun.In(i), fun.IsVariadic() && i == fun.NumIn()-1))
		}
	}

	for i := 0; i < fun.NumOut(); i++ {
		out = append(out, p.getVariableSignature(fun.Out(i), false))
	}

	var outParams string

	if fun.NumOut() > 0 {
		outParams = fmt.Sprintf(" (%s)", strings.Join(out, ", "))
	}

	return fmt.Sprintf("func(%s)%s", strings.Join(in, ", "), outParams)
}
