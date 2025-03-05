package util

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	gojautil "github.com/dop251/goja_nodejs/util" // init global registry import
)

// ModuleName of the util package
const ModuleName = "util"

// Util  holds the goja.Runtime for value conversion and a forward to the util.Util implementation
type Util struct {
	r *goja.Runtime
}

// Inspect the value of the goja.FunctionCall. Not implemented
func (util *Util) Inspect(_ goja.FunctionCall) goja.Value {
	panic("not implemented")
}

// Require the util package (with a forward reference to the gojautil.Util package)
func Require(runtime *goja.Runtime, module *goja.Object) {
	gojautil.Require(runtime, module) // set the format method

	s := &Util{r: runtime}
	runtime.ToValue(s)

	exports := module.Get("exports").(*goja.Object) //nolint:forcetypeassert // based on library reference implementation
	_ = exports.Set("inspect", s.Inspect)
}

// Enable the util package
func Enable(runtime *goja.Runtime, registry *require.Registry, _ *require.RequireModule) {
	registry.RegisterNativeModule("node:"+ModuleName, Require)
	registry.RegisterNativeModule(ModuleName, Require)
	_ = runtime.Set(ModuleName, require.Require(runtime, ModuleName))
}
