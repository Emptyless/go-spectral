package perfhooks

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// ModuleName of the perf_hooks package
const ModuleName = "perf_hooks"

// Enable perf_hooks no-op
func Enable(_ *goja.Runtime, registry *require.Registry, _ *require.RequireModule) {
	f := func(_ *goja.Runtime, _ *goja.Object) {}
	registry.RegisterNativeModule("node:"+ModuleName, f)
	registry.RegisterNativeModule(ModuleName, f)
}
