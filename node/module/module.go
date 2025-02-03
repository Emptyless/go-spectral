package module

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// ModuleName of the module package
const ModuleName = "module"

// Module holds the goja.Runtime
type Module struct {
	r             *goja.Runtime
	requireModule *require.RequireModule
}

// CreateRequire using the require.RequireModule
func (m *Module) CreateRequire(call goja.FunctionCall) goja.Value {
	return m.r.ToValue(func(_ goja.FunctionCall) goja.Value {
		v, _ := m.requireModule.Require(call.Argument(0).String())
		return v
	})
}

// Require the module package
func Require(runtime *goja.Runtime, module *goja.Object) {
	s := &Module{r: runtime}
	runtime.ToValue(s)

	exports := module.Get("exports").(*goja.Object)
	_ = exports.Set("createRequire", s.CreateRequire)
}

// Enable the module package
func Enable(runtime *goja.Runtime, registry *require.Registry, _ *require.RequireModule) {
	registry.RegisterNativeModule("node:"+ModuleName, Require)
	registry.RegisterNativeModule(ModuleName, Require)
	_ = runtime.Set(ModuleName, require.Require(runtime, ModuleName))
}
