package os

import (
	"runtime"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// ModuleName of the os package
const ModuleName = "os"

// OS holds the goja.Runtime for value conversion
type OS struct {
	r *goja.Runtime
}

// Platform returns the runtime.GOOS
func (o *OS) Platform(_ goja.FunctionCall) goja.Value {
	return o.r.ToValue(runtime.GOOS)
}

// CPUs returns a vCPU 2Ghz processor
func (o *OS) CPUs(_ goja.FunctionCall) goja.Value {
	//nolint:mnd // random values without influence
	return o.r.ToValue([]map[string]any{
		{
			"model": "vCPU",
			"speed": 2000,
			"times": map[string]any{
				"user": 1,
				"nice": 0,
				"sys":  1,
				"idle": 0,
				"irq":  0,
			},
		},
	})
}

// Require the os package
func Require(runtime *goja.Runtime, module *goja.Object) {
	s := &OS{r: runtime}
	runtime.ToValue(s)

	exports := module.Get("exports").(*goja.Object)
	_ = exports.Set("platform", s.Platform)
	_ = exports.Set("cpus", s.CPUs)
}

// Enable the os package
func Enable(runtime *goja.Runtime, registry *require.Registry, _ *require.RequireModule) {
	registry.RegisterNativeModule("node:"+ModuleName, Require)
	registry.RegisterNativeModule(ModuleName, Require)
	_ = runtime.Set(ModuleName, require.Require(runtime, ModuleName))
}
