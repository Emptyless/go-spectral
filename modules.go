package gospectral

import (
	"errors"
	"io/fs"

	"github.com/Emptyless/go-spectral/node/assert"
	"github.com/Emptyless/go-spectral/node/constants"
	"github.com/Emptyless/go-spectral/node/crypto"
	"github.com/Emptyless/go-spectral/node/events"
	nodefs "github.com/Emptyless/go-spectral/node/fs"
	"github.com/Emptyless/go-spectral/node/global"
	"github.com/Emptyless/go-spectral/node/http"
	"github.com/Emptyless/go-spectral/node/https"
	"github.com/Emptyless/go-spectral/node/module"
	osmodule "github.com/Emptyless/go-spectral/node/os"
	"github.com/Emptyless/go-spectral/node/path"
	perfhooks "github.com/Emptyless/go-spectral/node/perf_hooks"
	"github.com/Emptyless/go-spectral/node/process"
	"github.com/Emptyless/go-spectral/node/punycode"
	"github.com/Emptyless/go-spectral/node/stream"
	"github.com/Emptyless/go-spectral/node/tty"
	"github.com/Emptyless/go-spectral/node/util"
	"github.com/Emptyless/go-spectral/node/vm"
	"github.com/Emptyless/go-spectral/node/zlib"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/buffer"
	"github.com/dop251/goja_nodejs/console"
	noderequire "github.com/dop251/goja_nodejs/require"
	"github.com/dop251/goja_nodejs/url"
	log "github.com/sirupsen/logrus"
)

// Enable module with Name
type Enable struct {
	// Name of the Module to Enable
	Name string

	// Fn to Enable Module
	Fn func(runtime *goja.Runtime, registry *noderequire.Registry, requireModule *noderequire.RequireModule)
}

// BeforeModule is evaluated before any module is enabled and can be used to change e.g. the Loader. If a non-nil Enable.Fn
// is returned, it is used instead of the provided Enable
type BeforeModule = func(enable Enable, runtime *goja.Runtime, registry *noderequire.Registry, requireModule *noderequire.RequireModule) (func(runtime *goja.Runtime, registry *noderequire.Registry, requireModule *noderequire.RequireModule), error)

// DefaultBeforeModule implementation of BeforeModule is constructed using a curried function containing the working directory and a possibly virtual filesystem
func DefaultBeforeModule(currentWorkingDirectory string, fileSystem fs.FS) BeforeModule {
	return func(enable Enable, _ *goja.Runtime, _ *noderequire.Registry, _ *noderequire.RequireModule) (func(runtime *goja.Runtime, registry *noderequire.Registry, requireModule *noderequire.RequireModule), error) {
		switch enable.Name {
		case process.ModuleName:
			return func(runtime *goja.Runtime, registry *noderequire.Registry, requireModule *noderequire.RequireModule) {
				process.Enable(runtime, registry, requireModule, currentWorkingDirectory)
			}, nil
		case nodefs.ModuleName:
			return func(runtime *goja.Runtime, registry *noderequire.Registry, requireModule *noderequire.RequireModule) {
				nodefs.Enable(runtime, registry, requireModule, currentWorkingDirectory, fileSystem)
			}, nil
		default:
			return enable.Fn, nil
		}
	}
}

// AfterModule is evaluated after any module is enabled and can be used to change the current runtime state
type AfterModule = func(name string, runtime *goja.Runtime, registry *noderequire.Registry, requireModule *noderequire.RequireModule) error

// WithBeforeModule sets the Config.BeforeModule
func WithBeforeModule(before BeforeModule) Option {
	return func(config *Config) error {
		config.BeforeModule = before

		return nil
	}
}

// WithAfterModule sets the Config.AfterModule
func WithAfterModule(after AfterModule) Option {
	return func(config *Config) error {
		config.AfterModule = after

		return nil
	}
}

// Enables slice to LoadModules
func Enables() []Enable {
	return []Enable{
		{Name: util.ModuleName, Fn: util.Enable},
		{Name: stream.ModuleName, Fn: stream.Enable},
		{Name: http.ModuleName, Fn: http.Enable},
		{Name: https.ModuleName, Fn: https.Enable},
		{Name: zlib.ModuleName, Fn: zlib.Enable},
		{Name: url.ModuleName, Fn: func(runtime *goja.Runtime, _ *noderequire.Registry, _ *noderequire.RequireModule) {
			url.Enable(runtime)
		}},
		{Name: global.ModuleName, Fn: global.Enable},
		{Name: nodefs.ModuleName, Fn: func(_ *goja.Runtime, _ *noderequire.Registry, _ *noderequire.RequireModule) {
			panic(process.ModuleName + " relies on working directory and FS and must be provided with a BeforeLoader")
		}},
		{Name: vm.ModuleName, Fn: vm.Enable},
		{Name: console.ModuleName, Fn: func(runtime *goja.Runtime, _ *noderequire.Registry, _ *noderequire.RequireModule) {
			console.Enable(runtime)
		}},
		{Name: process.ModuleName, Fn: func(_ *goja.Runtime, _ *noderequire.Registry, _ *noderequire.RequireModule) {
			panic(process.ModuleName + " relies on working directory and must be provided with a BeforeLoader")
		}},
		{Name: module.ModuleName, Fn: module.Enable},
		{Name: perfhooks.ModuleName, Fn: perfhooks.Enable},
		{Name: crypto.ModuleName, Fn: crypto.Enable},
		{Name: assert.ModuleName, Fn: assert.Enable},
		{Name: path.ModuleName, Fn: path.Enable},
		{Name: osmodule.ModuleName, Fn: osmodule.Enable},
		{Name: buffer.ModuleName, Fn: func(runtime *goja.Runtime, _ *noderequire.Registry, _ *noderequire.RequireModule) {
			buffer.Enable(runtime)
		}},
		{Name: tty.ModuleName, Fn: tty.Enable},
		{Name: constants.ModuleName, Fn: constants.Enable},
		{Name: punycode.ModuleName, Fn: punycode.Enable},
		{Name: events.ModuleName, Fn: events.Enable},
	}
}

// ErrBeforeModule if the BeforeModule call failed
var ErrBeforeModule = errors.New("failed to run BeforeModule hook")

// ErrAfterModule if the AfterModule call failed
var ErrAfterModule = errors.New("failed to run AfterModule hook")

// LoadModules replacing NodeJS functionality required by Dist
func LoadModules(runtime *goja.Runtime, registry *noderequire.Registry, beforeModule BeforeModule, afterModule AfterModule) (*noderequire.RequireModule, error) {
	log.Info("loading modules")
	require := registry.Enable(runtime)

	for _, enable := range Enables() {
		if beforeModule != nil {
			loader, err := beforeModule(enable, runtime, registry, require)
			if err != nil {
				log.Infof("failed to run BeforeModule hook for %s", enable.Name)
				return nil, errors.Join(ErrBeforeModule, err)
			}

			loader(runtime, registry, require)
		} else {
			enable.Fn(runtime, registry, require)
		}

		if afterModule != nil {
			if err := afterModule(enable.Name, runtime, registry, require); err != nil {
				log.Infof("failed to run AfterModule hook for %s", enable.Name)
				return nil, errors.Join(ErrAfterModule, err)
			}
		}
	}

	return require, nil
}
