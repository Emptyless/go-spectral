package constants

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// ModuleName of the "constants" package
const ModuleName = "constants"

// Enable constants using a no-op
func Enable(_ *goja.Runtime, registry *require.Registry, _ *require.RequireModule) {
	f := func(_ *goja.Runtime, _ *goja.Object) {}
	registry.RegisterNativeModule("node:"+ModuleName, f)
	registry.RegisterNativeModule(ModuleName, f)
}
