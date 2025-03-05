package process

import (
	"os"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// ModuleName of process package
const ModuleName = "process"

// Process holds the goja.Runtime for converting values and a map of environment values (by os.Environ) when
// the program started
type Process struct {
	r   *goja.Runtime
	env map[string]string

	// CurrentWorkingDirectory used by process
	CurrentWorkingDirectory string
}

// On returns null
func (p *Process) On(_ goja.FunctionCall) goja.Value {
	return goja.Undefined()
}

// Cwd returns the CurrentWorkingDirectory if set or else the os.Getwd
func (p *Process) Cwd(_ goja.FunctionCall) goja.Value {
	if p.CurrentWorkingDirectory != "" {
		return p.r.ToValue(p.CurrentWorkingDirectory)
	}

	cwd, _ := os.Getwd()

	return p.r.ToValue(cwd)
}

// Versions implemented by Node
func (p *Process) Versions() goja.Value {
	return p.r.ToValue(Versions)
}

// Require the process package
func Require(p *Process) func(runtime *goja.Runtime, module *goja.Object) {
	return func(runtime *goja.Runtime, module *goja.Object) {
		for _, e := range os.Environ() {
			if p.env == nil {
				p.env = map[string]string{}
			}

			envKeyValue := strings.SplitN(e, "=", 2) //nolint:mnd // split in key=value, two parts
			p.env[envKeyValue[0]] = envKeyValue[1]
		}

		o := module.Get("exports").(*goja.Object) //nolint:forcetypeassert // based on library reference implementation
		_ = o.Set("env", runtime.ToValue(p.env))
		_ = o.Set("on", p.On)
		_ = o.Set("versions", p.Versions())
		_ = o.Set("version", runtime.ToValue(Version))
		_ = o.Set("cwd", p.Cwd)
		_ = o.Set("stdout", runtime.NewObject())
		_ = o.Set("stderr", runtime.NewObject())
		_ = o.Set("argv", runtime.ToValue([]string{"spectral"}))
	}
}

// Enable the process package
func Enable(runtime *goja.Runtime, registry *require.Registry, _ *require.RequireModule, currentWorkingDirectory string) {
	p := &Process{
		r:                       runtime,
		env:                     make(map[string]string),
		CurrentWorkingDirectory: currentWorkingDirectory,
	}

	registry.RegisterNativeModule("node:"+ModuleName, Require(p))
	registry.RegisterNativeModule(ModuleName, Require(p))
	_ = runtime.Set("process", require.Require(runtime, ModuleName))
}
