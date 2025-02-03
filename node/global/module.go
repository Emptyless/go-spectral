package global

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

const ModuleName = "global"

// Require "global" by setting global to an empty object
func Require(runtime *goja.Runtime, module *goja.Object) {
	exports := module.Get("exports").(*goja.Object) //nolint:forcetypeassert // based on reference implementation in library
	_ = exports.Set("global", runtime.NewObject())
}

// Enable the global package with an empty object
func Enable(runtime *goja.Runtime, registry *require.Registry, _ *require.RequireModule) {
	registry.RegisterNativeModule("node:"+ModuleName, Require)
	registry.RegisterNativeModule(ModuleName, Require)
	_ = runtime.Set(ModuleName, require.Require(runtime, ModuleName))
}
