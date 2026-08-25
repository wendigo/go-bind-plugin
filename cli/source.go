package cli

import (
	"fmt"
	"go/types"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

// analyzeSource builds the plugin structure by type-checking the plugin
// package sources instead of loading the compiled plugin into the process.
// pluginPath must point at the already built shared object; it is only
// read to record the plugin's size and checksum.
func analyzeSource(pluginPackage string, pluginPath string, imports []string) (*pluginStructure, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes,
	}, pluginPackage)

	if err != nil {
		return nil, fmt.Errorf("could not load package %s: %w", pluginPackage, err)
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("expected pattern %s to match exactly one package, got %d", pluginPackage, len(pkgs))
	}

	pkg := pkgs[0]

	if len(pkg.Errors) > 0 {
		return nil, fmt.Errorf("could not load package %s: %s", pluginPackage, pkg.Errors[0])
	}

	if pkg.Types == nil {
		return nil, fmt.Errorf("no type information for package %s", pluginPackage)
	}

	stat, err := os.Stat(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("could not stat file: %w", err)
	}

	shaSum, err := fileChecksum(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("could not calculate plugin checksum: %w", err)
	}

	ps := &pluginStructure{
		Symbols:      packageSymbols{},
		Package:      pkg.PkgPath,
		Size:         stat.Size(),
		Sha256:       shaSum,
		importsNames: map[string]string{},
	}

	for _, imp := range imports {
		ps.getNamedPkgImport(imp)
	}

	scope := pkg.Types.Scope()

	// scope.Names() is sorted, so functions and variables come out ordered by name
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)

		if !obj.Exported() {
			continue
		}

		switch obj := obj.(type) {
		case *types.Func:
			sig := obj.Type().(*types.Signature)

			// generic functions are not plugin symbols
			if sig.TypeParams().Len() > 0 {
				continue
			}

			ps.Symbols[name] = nil
			ps.Functions = append(ps.Functions, &function{
				Name:           name,
				Signature:      ps.getTypesFunctionSignature(sig, false),
				FunctionHeader: ps.getTypesFunctionSignature(sig, true),
				ArgumentsCount: sig.Params().Len(),
				ReturnsVoid:    sig.Results().Len() == 0,
				IsVariadic:     sig.Variadic(),
			})
		case *types.Var:
			ps.Symbols[name] = nil
			ps.Variables = append(ps.Variables, &variable{
				Name:      name,
				Typ:       obj.Type().String(),
				Signature: ps.getTypesSignature(obj.Type(), false),
			})
		}
	}

	ps.Imports = ps.imports()
	ps.NamedImports = ps.namedImports()

	return ps, nil
}

// getTypesSignature is the go/types counterpart of getVariableSignature
func (p *pluginStructure) getTypesSignature(t types.Type, isVariadic bool) string {
	if obj := typeName(t); obj != nil {
		var pkgImport string

		if obj.Pkg() != nil {
			pkgImport = p.getNamedPkgImport(obj.Pkg().Path())
		}

		if isVariadic {
			return fmt.Sprintf("...%s%s", pkgImport, obj.Name())
		}

		return fmt.Sprintf("%s%s", pkgImport, obj.Name())
	}

	switch t := t.(type) {
	case *types.Basic:
		if isVariadic {
			return fmt.Sprintf("...%s", t.Name())
		}

		return t.Name()
	case *types.Signature:
		return p.getTypesFunctionSignature(t, false)
	case *types.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), p.getTypesSignature(t.Elem(), false))
	case *types.Slice:
		if isVariadic {
			return fmt.Sprintf("...%s", p.getTypesSignature(t.Elem(), false))
		}

		return fmt.Sprintf("[]%s", p.getTypesSignature(t.Elem(), false))
	case *types.Chan:
		return fmt.Sprintf("%s %s", chanDirString(t.Dir()), p.getTypesSignature(t.Elem(), false))
	case *types.Map:
		return fmt.Sprintf("map[%s]%s", p.getTypesSignature(t.Key(), false), p.getTypesSignature(t.Elem(), false))
	case *types.Pointer:
		return fmt.Sprintf("*%s", p.getTypesSignature(t.Elem(), false))
	}

	return types.TypeString(t, p.typesQualifier)
}

// getTypesFunctionSignature is the go/types counterpart of getFunctionSignature
func (p *pluginStructure) getTypesFunctionSignature(sig *types.Signature, namedParams bool) string {
	var in []string
	var out []string

	params := sig.Params()

	for i := range params.Len() {
		signature := p.getTypesSignature(params.At(i).Type(), sig.Variadic() && i == params.Len()-1)

		if namedParams {
			signature = fmt.Sprintf("in%d %s", i, signature)
		}

		in = append(in, signature)
	}

	results := sig.Results()

	for i := range results.Len() {
		out = append(out, p.getTypesSignature(results.At(i).Type(), false))
	}

	var outParams string

	if results.Len() > 0 {
		outParams = fmt.Sprintf(" (%s)", strings.Join(out, ", "))
	}

	return fmt.Sprintf("func(%s)%s", strings.Join(in, ", "), outParams)
}

func (p *pluginStructure) typesQualifier(pkg *types.Package) string {
	return strings.TrimSuffix(p.getNamedPkgImport(pkg.Path()), ".")
}

// typeName returns the defining object for named and alias types, nil otherwise
func typeName(t types.Type) *types.TypeName {
	if named, ok := t.(interface{ Obj() *types.TypeName }); ok {
		return named.Obj()
	}

	return nil
}

func chanDirString(dir types.ChanDir) string {
	switch dir {
	case types.RecvOnly:
		return "<-chan"
	case types.SendOnly:
		return "chan<-"
	default:
		return "chan"
	}
}
