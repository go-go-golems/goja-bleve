package bleve

import (
	"fmt"

	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/modules"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	blevemodule "github.com/go-go-golems/goja-bleve/pkg"
)

const PackageID = "goja-bleve"

var moduleNames = []string{
	blevemodule.ModuleName,
}

// Register exposes goja-bleve modules as xgoja provider modules.
func Register(registry *providerapi.Registry) error {
	entries := make([]providerapi.Entry, 0, len(moduleNames))
	for _, name := range moduleNames {
		mod := modules.GetModule(name)
		if mod == nil {
			return fmt.Errorf("bleve module %q is not registered", name)
		}
		entries = append(entries, nativeModuleEntry(mod))
	}
	return registry.Package(PackageID, entries...)
}

func nativeModuleEntry(mod modules.NativeModule) providerapi.Module {
	return providerapi.Module{
		Name:        mod.Name(),
		DefaultAs:   mod.Name(),
		Description: mod.Doc(),
		New: func(providerapi.ModuleContext) (require.ModuleLoader, error) {
			return mod.Loader, nil
		},
	}
}
