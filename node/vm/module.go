package vm

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// ModuleName of the vm package
const ModuleName = "vm"

// Enable vm package with a no-op
func Enable(_ *goja.Runtime, registry *require.Registry, _ *require.RequireModule) {
	f := func(_ *goja.Runtime, _ *goja.Object) {}
	registry.RegisterNativeModule("node:"+ModuleName, f)
	registry.RegisterNativeModule(ModuleName, f)
}
