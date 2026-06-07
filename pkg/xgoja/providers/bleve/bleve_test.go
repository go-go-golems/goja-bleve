package bleve

import (
	"testing"

	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	blevemodule "github.com/go-go-golems/goja-bleve/pkg"
)

func TestRegisterAddsBleveModule(t *testing.T) {
	registry := providerapi.NewProviderRegistry()
	if err := Register(registry); err != nil {
		t.Fatalf("Register: %v", err)
	}
	entry, ok := registry.ResolveModule(PackageID, blevemodule.ModuleName)
	if !ok {
		t.Fatalf("provider package %q did not register module %q", PackageID, blevemodule.ModuleName)
	}
	if entry.Name != blevemodule.ModuleName || entry.DefaultAs != blevemodule.ModuleName {
		t.Fatalf("unexpected provider module entry: %#v", entry)
	}
	if entry.Description == "" {
		t.Fatalf("expected provider module description")
	}
	if _, err := entry.NewModuleFactory(providerapi.ModuleSetupContext{Name: blevemodule.ModuleName, As: blevemodule.ModuleName}); err != nil {
		t.Fatalf("module factory: %v", err)
	}
}
