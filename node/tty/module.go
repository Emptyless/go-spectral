package tty

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// ModuleName of the tty package
const ModuleName = "tty"

// TTY holds the goja.Runtime
type TTY struct {
	r *goja.Runtime
}

// IsAtty always returns false
func (t *TTY) IsAtty(_ goja.FunctionCall) goja.Value {
	return t.r.ToValue(false)
}

// Require the tty package
func Require(runtime *goja.Runtime, module *goja.Object) {
	s := &TTY{r: runtime}
	runtime.ToValue(s)

	exports := module.Get("exports").(*goja.Object) //nolint:forcetypeassert // based on library reference implementation
	_ = exports.Set("isatty", s.IsAtty)
}

func Enable(runtime *goja.Runtime, registry *require.Registry, _ *require.RequireModule) {
	registry.RegisterNativeModule("node:"+ModuleName, Require)
	registry.RegisterNativeModule(ModuleName, Require)
	_ = runtime.Set(ModuleName, require.Require(runtime, ModuleName))
}
