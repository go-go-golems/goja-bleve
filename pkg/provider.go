package pkg

import (
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
)

const ProviderPackageID = "goja-bleve"

// RegisterProvider exposes the bleve native module through the xgoja provider
// registry. Phase 7 will add configuration once host path and lifecycle policy
// decisions are implemented.
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
