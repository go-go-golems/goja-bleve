package pkg

import (
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
)

const ProviderPackageID = "goja-bleve"

// RegisterProvider exposes the bleve native module through the xgoja provider
// registry. The module currently has no provider-level configuration schema:
// index path policy and runtime cleanup are host-application responsibilities.
// The JavaScript API still makes lifecycle explicit through index.close().
func RegisterProvider(registry *providerapi.Registry) error {
	return registry.Package(ProviderPackageID, providerapi.Module{
		Name:        ModuleName,
		DefaultAs:   ModuleName,
		Description: module{}.Doc(),
		New: func(providerapi.ModuleContext) (require.ModuleLoader, error) {
			return NewLoader(), nil
		},
	})
}
