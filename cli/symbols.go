package cli

import (
	"debug/elf"
	"debug/macho"
	"fmt"
	"plugin"
	"slices"
	"strings"
	"unicode"
)

// pluginSymbols enumerates symbols exported by the plugin using its shared
// object symbol table and the public plugin.Lookup API. It returns the
// plugin's package path along with the resolved symbols.
//
// Symbols are named "<pkgpath>.<Name>" in the symbol table. Candidates whose
// name part is an exported Go identifier are resolved with Lookup (which only
// succeeds for the plugin's own package) and grouped by package prefix; the
// prefix with the most resolved symbols is the plugin's package.
func pluginSymbols(p *plugin.Plugin, path string) (string, packageSymbols, error) {
	names, err := sharedObjectSymbols(path)
	if err != nil {
		return "", nil, fmt.Errorf("could not read plugin symbol table: %w", err)
	}

	pkg, symbols := resolveSymbols(names, p.Lookup)

	return pkg, symbols, nil
}

// resolveSymbols implements the package election described on pluginSymbols;
// lookup resolves a bare symbol name within the plugin's own package.
func resolveSymbols(names []string, lookup func(string) (plugin.Symbol, error)) (string, packageSymbols) {
	candidates := make(map[string][]string)

	for _, full := range names {
		idx := strings.LastIndex(full, ".")
		if idx <= 0 {
			continue
		}

		pkg, name := full[:idx], full[idx+1:]
		if !isExportedGoIdentifier(name) {
			continue
		}

		candidates[pkg] = append(candidates[pkg], name)
	}

	resolved := make(map[string]packageSymbols)

	for pkg, symbolNames := range candidates {
		symbols := packageSymbols{}

		for _, name := range symbolNames {
			if sym, err := lookup(name); err == nil {
				symbols[name] = sym
			}
		}

		if len(symbols) > 0 {
			resolved[pkg] = symbols
		}
	}

	packages := make([]string, 0, len(resolved))
	for pkg := range resolved {
		packages = append(packages, pkg)
	}
	slices.Sort(packages)

	var pluginPackage string
	for _, pkg := range packages {
		if pluginPackage == "" || len(resolved[pkg]) > len(resolved[pluginPackage]) {
			pluginPackage = pkg
		}
	}

	if pluginPackage == "" {
		return "", packageSymbols{}
	}

	return pluginPackage, resolved[pluginPackage]
}

// sharedObjectSymbols lists raw symbol names from an ELF or Mach-O shared object.
func sharedObjectSymbols(path string) ([]string, error) {
	if f, err := elf.Open(path); err == nil {
		defer f.Close()

		symbols, err := f.DynamicSymbols()
		if err != nil {
			return nil, err
		}

		names := make([]string, 0, len(symbols))
		for _, s := range symbols {
			names = append(names, s.Name)
		}

		return names, nil
	}

	if f, err := macho.Open(path); err == nil {
		defer f.Close()

		if f.Symtab == nil {
			return nil, fmt.Errorf("no symbol table in %s", path)
		}

		names := make([]string, 0, len(f.Symtab.Syms))
		for _, s := range f.Symtab.Syms {
			// Mach-O symbol names carry a leading underscore
			names = append(names, strings.TrimPrefix(s.Name, "_"))
		}

		return names, nil
	}

	return nil, fmt.Errorf("%s is not an ELF or Mach-O shared object", path)
}

func isExportedGoIdentifier(name string) bool {
	for i, r := range name {
		if i == 0 && !unicode.IsUpper(r) {
			return false
		}

		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}

	return name != ""
}
