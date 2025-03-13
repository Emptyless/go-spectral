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

	// CurrentWorkingDirectory needed for the Resolve implementation
	CurrentWorkingDirectory string
}

// Dirname of argument
func (p *Path) Dirname(call goja.FunctionCall) goja.Value {
	dirname := path.Dir(call.Argument(0).String())
	return p.r.ToValue(dirname)
}

// Resolve args to path
func (p *Path) Resolve(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		return goja.Undefined()
	}

	var res string
	for _, arg := range call.Arguments {
		stringarg := arg.String()
		if path.IsAbs(stringarg) {
			res = arg.String()
			continue
		}

		res = path.Join(res, stringarg)
	}

	if !path.IsAbs(res) {
		res = path.Join(p.CurrentWorkingDirectory, res)
	}

	return p.r.ToValue(res)
}

// Relative path. Not implemented.
func (p *Path) Relative(_ goja.FunctionCall) goja.Value {
	panic("not implemented")
}

// Extname is the name of the extension of the argument
func (p *Path) Extname(call goja.FunctionCall) goja.Value {
	return p.r.ToValue(path.Ext(call.Argument(0).String()))
}

// IsURL checks if the input is a URL. Not implemented.
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
func Require(currentWorkingDirectory string) func(runtime *goja.Runtime, module *goja.Object) {
	return func(runtime *goja.Runtime, module *goja.Object) {
		s := &Path{
			r:                       runtime,
			CurrentWorkingDirectory: currentWorkingDirectory,
		}
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
}

// Enable the path package
func Enable(runtime *goja.Runtime, registry *require.Registry, _ *require.RequireModule, currentWorkingDirectory string) {
	registry.RegisterNativeModule("node:"+ModuleName, Require(currentWorkingDirectory))
	registry.RegisterNativeModule(ModuleName, Require(currentWorkingDirectory))
	_ = runtime.Set(ModuleName, require.Require(runtime, ModuleName))
}
