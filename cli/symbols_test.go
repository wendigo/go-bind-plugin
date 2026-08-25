package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"testing"
)

func TestIsExportedGoIdentifier(t *testing.T) {
	cases := []struct {
		name     string
		expected bool
	}{
		{"Foo", true},
		{"F", true},
		{"Foo123", true},
		{"Foo_bar", true},
		{"Über", true},
		{"", false},
		{"foo", false},
		{"_Foo", false},
		{"1Foo", false},
		{"Foo.Bar", false},
		{"Foo-bar", false},
		{"Foo[int]", false},
		{"µfoo", false},
	}

	for _, c := range cases {
		if actual := isExportedGoIdentifier(c.name); actual != c.expected {
			t.Errorf("isExportedGoIdentifier(%q) = %v, expected %v", c.name, actual, c.expected)
		}
	}
}

func lookupIn(known map[string]plugin.Symbol) func(string) (plugin.Symbol, error) {
	return func(name string) (plugin.Symbol, error) {
		if sym, ok := known[name]; ok {
			return sym, nil
		}

		return nil, fmt.Errorf("symbol %s not found", name)
	}
}

func TestResolveSymbolsElectsPluginPackage(t *testing.T) {
	var newFunc, sumFunc plugin.Symbol = func() {}, func() {}

	names := []string{
		"example.com/plug.New",
		"example.com/plug.Sum",
		"example.com/dep.New", // dependency exporting a name that collides with a plugin symbol
		"example.com/dep.Internal",
		"example.com/plug..inittask",
		"go:link.thispluginpath",
		"runtime.text",
		"type:example.com/plug.T",
	}

	pkg, symbols := resolveSymbols(names, lookupIn(map[string]plugin.Symbol{
		"New": newFunc,
		"Sum": sumFunc,
	}))

	if pkg != "example.com/plug" {
		t.Errorf("elected package = %q, expected %q", pkg, "example.com/plug")
	}

	if len(symbols) != 2 {
		t.Errorf("resolved %d symbols, expected 2: %v", len(symbols), symbols)
	}

	for _, name := range []string{"New", "Sum"} {
		if _, ok := symbols[name]; !ok {
			t.Errorf("symbol %q missing from resolved set %v", name, symbols)
		}
	}
}

func TestResolveSymbolsIgnoresUnresolvableCandidates(t *testing.T) {
	names := []string{
		"example.com/plug.Exported",
		"example.com/plug.Removed", // in the symbol table but not resolvable via Lookup
	}

	pkg, symbols := resolveSymbols(names, lookupIn(map[string]plugin.Symbol{
		"Exported": func() {},
	}))

	if pkg != "example.com/plug" {
		t.Errorf("elected package = %q, expected %q", pkg, "example.com/plug")
	}

	if len(symbols) != 1 {
		t.Errorf("resolved %d symbols, expected 1: %v", len(symbols), symbols)
	}
}

func TestResolveSymbolsWithoutExportedSymbols(t *testing.T) {
	names := []string{
		"example.com/plug..inittask",
		"example.com/plug.internal",
		"noDotsAtAll",
	}

	pkg, symbols := resolveSymbols(names, lookupIn(nil))

	if pkg != "" {
		t.Errorf("elected package = %q, expected none", pkg)
	}

	if len(symbols) != 0 {
		t.Errorf("resolved %d symbols, expected none: %v", len(symbols), symbols)
	}
}

func TestResolveSymbolsBreaksTiesDeterministically(t *testing.T) {
	names := []string{
		"example.com/b.Foo",
		"example.com/a.Bar",
	}

	lookup := lookupIn(map[string]plugin.Symbol{
		"Foo": func() {},
		"Bar": func() {},
	})

	for i := 0; i < 10; i++ {
		if pkg, _ := resolveSymbols(names, lookup); pkg != "example.com/a" {
			t.Fatalf("elected package = %q, expected %q", pkg, "example.com/a")
		}
	}
}

func TestSharedObjectSymbolsRejectsNonSharedObject(t *testing.T) {
	path := filepath.Join(t.TempDir(), "not_a_plugin.so")

	if err := os.WriteFile(path, []byte("definitely not a shared object"), 0600); err != nil {
		t.Fatal(err)
	}

	if _, err := sharedObjectSymbols(path); err == nil {
		t.Error("expected an error for a non-shared-object file")
	}
}

func TestSharedObjectSymbolsRejectsMissingFile(t *testing.T) {
	if _, err := sharedObjectSymbols(filepath.Join(t.TempDir(), "missing.so")); err == nil {
		t.Error("expected an error for a missing file")
	}
}
