package path

import (
	"path"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// ModuleName of the path package
const ModuleName = "path"

// Path holds the goja.Runtime for value conversion
type Path struct {
	r *goja.Runtime
}

// Dirname of argument
func (p *Path) Dirname(call goja.FunctionCall) goja.Value {
	dirname := path.Dir(call.Argument(0).String())
	return p.r.ToValue(dirname)
}

// Resolve args to path
func (p *Path) Resolve(call goja.FunctionCall) goja.Value {
	args := make([]string, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.String()
	}

	return p.r.ToValue(path.Join(args...))
}

// Relative path. Not implemented.
func (p *Path) Relative(_ goja.FunctionCall) goja.Value {
	panic("not implemented")
}

// Extname is the name of the extension of the argument
func (p *Path) Extname(call goja.FunctionCall) goja.Value {
	return p.r.ToValue(path.Ext(call.Argument(0).String()))
}

// IsURL checks if the input is an URL. Not implemented.
func (p *Path) IsURL(_ goja.FunctionCall) goja.Value {
	panic("not implemented")
}

// Basename of the argument.
func (p *Path) Basename(call goja.FunctionCall) goja.Value {
	base := path.Base(call.Argument(0).String())
	if len(call.Arguments) == 2 { //nolint:mnd // function optionally has two arguments
		base = strings.TrimSuffix(base, call.Argument(1).String())
	}

	return p.r.ToValue(base)
}

// Require the path package
func Require(runtime *goja.Runtime, module *goja.Object) {
	s := &Path{r: runtime}
	runtime.ToValue(s)

	p := runtime.NewObject()
	exports := module.Get("exports").(*goja.Object) //nolint:forcetypeassert // based on library reference implementation
	_ = exports.Set("posix", p)
	for _, o := range []*goja.Object{p, exports} { // set all methods to be equivalent for posix/non-posix
		_ = o.Set("dirname", s.Dirname)
		_ = o.Set("resolve", s.Resolve)
		_ = o.Set("relative", s.Relative)
		_ = o.Set("extname", s.Extname)
		_ = o.Set("basename", s.Basename)
	}
}

// Enable the path package
func Enable(runtime *goja.Runtime, registry *require.Registry, _ *require.RequireModule) {
	registry.RegisterNativeModule("node:"+ModuleName, Require)
	registry.RegisterNativeModule(ModuleName, Require)
	_ = runtime.Set(ModuleName, require.Require(runtime, ModuleName))
}
