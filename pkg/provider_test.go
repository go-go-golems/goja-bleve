package pkg

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/modules"
	"github.com/go-go-golems/go-go-goja/pkg/tsgen/spec"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
)

func TestNativeModuleRegistersWithDefaultRegistryAndTypescript(t *testing.T) {
	mod := modules.GetModule(ModuleName)
	if mod == nil {
		t.Fatalf("module %q was not registered in modules.DefaultRegistry", ModuleName)
	}
	if mod.Name() != ModuleName {
		t.Fatalf("module name = %q", mod.Name())
	}
	typed, ok := mod.(modules.TypeScriptDeclarer)
	if !ok {
		t.Fatalf("module %q does not implement TypeScriptDeclarer", ModuleName)
	}
	ts := typed.TypeScriptModule()
	if ts == nil || ts.Name != ModuleName {
		t.Fatalf("TypeScriptModule = %#v", ts)
	}
	if len(ts.RawDTS) == 0 || len(ts.Functions) == 0 {
		t.Fatalf("TypeScriptModule should expose raw interfaces and functions: %#v", ts)
	}
	if !typescriptFunctionNamesContain(ts.Functions, []string{"mapping", "field", "search", "create", "match", "matchAll"}) {
		t.Fatalf("TypeScript functions missing expected exports: %#v", ts.Functions)
	}
}

func TestRegisterProviderResolvesBleveModuleAndCanRequire(t *testing.T) {
	registry := providerapi.NewProviderRegistry()
	if err := RegisterProvider(registry); err != nil {
		t.Fatalf("RegisterProvider: %v", err)
	}
	entry, ok := registry.ResolveModule(ProviderPackageID, ModuleName)
	if !ok {
		t.Fatalf("provider did not register module %q", ModuleName)
	}
	if entry.DefaultAs != ModuleName {
		t.Fatalf("DefaultAs = %q", entry.DefaultAs)
	}
	loader, err := entry.NewModuleFactory(providerapi.ModuleSetupContext{Name: ModuleName, As: ModuleName})
	if err != nil {
		t.Fatalf("provider module factory: %v", err)
	}

	vm := goja.New()
	requireRegistry := require.NewRegistry()
	requireRegistry.RegisterNativeModule(ModuleName, loader)
	requireRegistry.Enable(vm)
	value, err := vm.RunString(`
		const bleve = require("bleve");
		({ version: bleve.version, search: typeof bleve.search, vectorSupport: typeof bleve.vectorSupport });
	`)
	if err != nil {
		t.Fatalf("require provider module: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["version"] != "0.1.0" || got["search"] != "function" || got["vectorSupport"] != "boolean" {
		t.Fatalf("unexpected provider require result: %#v", got)
	}
}

func typescriptFunctionNamesContain(functions []spec.Function, names []string) bool {
	seen := map[string]bool{}
	for _, fn := range functions {
		seen[fn.Name] = true
	}
	for _, name := range names {
		if !seen[name] {
			return false
		}
	}
	return true
}
