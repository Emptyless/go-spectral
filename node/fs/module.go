package fs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/sirupsen/logrus"
)

// ModuleName of the "fs" package
const ModuleName = "fs"

// FS struct initialization
type FS struct {
	r *goja.Runtime

	// CurrentWorkingDirectory set in the process module. Not that this must NOT be used to prepend paths
	// but is rather to resolve files relative to CurrentWorkingDirectory when also a FileSystem is provided.
	// This achieves the effect that when FileSystem is used, the 'proper' structure of that fileystem is preserved.
	//
	// An example of this would be having an embed.FS which is a directory containing the ruleset:
	// go-embed testdata
	// var fs embed.FS
	//
	// On the process module, the CurrentWorkingDirectory is set and provides the Cwd method. The JavaScript code
	// already uses this on paths that are then provided to fs so the CurrentWorkingDirectory is already part of that
	// This directory however most likely doesn't exist in embedded so that prefix is then stripped of such that
	// the ruleset can still be referenced now using testdata/ruleset.yaml
	CurrentWorkingDirectory string

	// FileSystem if not nil is used to resolve e.g. embedded filesystems (embed.FS). These can contain bundled
	// OpenAPI specs and/or rulesets. The files are resolved by trimming the CurrentWorkingDirectory prefix
	// from paths and if there is a match that file is used. In case of no match, the search continues on the
	// system file system using os.ReadFile.
	FileSystem fs.FS
}

// Native asynchronous realpath. Not implemented
func (f *FS) Native(_ goja.FunctionCall) goja.Value {
	panic("not implemented")
}

// LStat translated to os.Lstat but returning false for
// isDirectory
// isSymbolicLink
// isBlockDevice
// isCharacterDevice
// isFIFO
// isSocket
// and always true for isFile. Note that this most likely breaks glob patterns with directories and
// this method should be properly implemented
// TODO fix LStat such that the returned values are correct for also directories / symlinks
func (f *FS) LStat(call goja.FunctionCall) goja.Value {
	filePath := call.Argument(0).String()
	cb := call.Argument(1).Export().(func(goja.FunctionCall) goja.Value)
	stats := f.r.NewObject()

	file, openFileErr := f.openFile(filePath)
	if openFileErr != nil {
		cb(goja.FunctionCall{
			This:      call.This,
			Arguments: []goja.Value{f.r.ToValue(openFileErr.Error()), goja.Null()},
		})

		return goja.Undefined()
	}

	lstat, lstatErr := file.Stat()
	if lstatErr != nil {
		cb(goja.FunctionCall{
			This:      call.This,
			Arguments: []goja.Value{f.r.ToValue(lstatErr.Error()), goja.Null()},
		})

		return goja.Undefined()
	}

	_ = stats.Set("dev", lstat.Size())
	_ = stats.Set("isDirectory", func(_ goja.FunctionCall) goja.Value {
		return f.r.ToValue(lstat.IsDir())
	})
	_ = stats.Set("isSymbolicLink", func(_ goja.FunctionCall) goja.Value {
		return f.r.ToValue(false)
	})
	_ = stats.Set("isBlockDevice", func(_ goja.FunctionCall) goja.Value {
		return f.r.ToValue(false)
	})
	_ = stats.Set("isCharacterDevice", func(_ goja.FunctionCall) goja.Value {
		return f.r.ToValue(false)
	})
	_ = stats.Set("isFIFO", func(_ goja.FunctionCall) goja.Value {
		return f.r.ToValue(false)
	})
	_ = stats.Set("isFile", func(_ goja.FunctionCall) goja.Value {
		return f.r.ToValue(true)
	})
	_ = stats.Set("isSocket", func(_ goja.FunctionCall) goja.Value {
		return f.r.ToValue(false)
	})

	cb(goja.FunctionCall{
		This:      call.This,
		Arguments: []goja.Value{goja.Null(), stats},
	})

	return goja.Undefined()
}

// ReadFile with callback
func (f *FS) ReadFile(call goja.FunctionCall) goja.Value {
	filePath := call.Argument(0).String()
	var cb func(functionCall goja.FunctionCall) goja.Value
	if len(call.Arguments) == 2 { //nolint:mnd // function optionally has two arguments
		cb = call.Argument(1).Export().(func(goja.FunctionCall) goja.Value)
	} else {
		cb = call.Argument(2).Export().(func(goja.FunctionCall) goja.Value) //nolint:mnd // select second argument
	}

	file, openFileErr := f.openFile(filePath)
	if openFileErr != nil {
		cb(goja.FunctionCall{
			This:      call.This,
			Arguments: []goja.Value{f.r.ToValue(openFileErr.Error()), goja.Null()},
		})

		return goja.Undefined()
	}

	b, readAllErr := io.ReadAll(file)
	if readAllErr != nil {
		cb(goja.FunctionCall{
			This:      call.This,
			Arguments: []goja.Value{f.r.ToValue(readAllErr.Error()), goja.Null()},
		})

		return goja.Undefined()
	}

	cb(goja.FunctionCall{
		This:      call.This,
		Arguments: []goja.Value{goja.Null(), f.r.ToValue(string(b))},
	})

	return goja.Undefined()
}

// PromiseReadFile using os.ReadFile
func (f *FS) PromiseReadFile(call goja.FunctionCall) goja.Value {
	promise, resolve, reject := f.r.NewPromise()
	filePath := call.Argument(0).String()
	file, openFileErr := f.openFile(filePath)
	if openFileErr != nil {
		_ = reject(openFileErr.Error())
		return f.r.ToValue(promise)
	}

	b, readAllErr := io.ReadAll(file)
	if readAllErr != nil {
		_ = reject(readAllErr.Error())
		return f.r.ToValue(promise)
	}

	_ = resolve(string(b))

	return f.r.ToValue(promise)
}

// openFile either through embedded FileSystem or system fs
func (f *FS) openFile(filePath string) (fs.File, error) {
	if f.FileSystem != nil {
		// prepend that the f.FileSystem is stored at the root of the f.CurrentWorkingDirectory
		rel, relErr := filepath.Rel(f.CurrentWorkingDirectory, filePath)
		if relErr != nil {
			return nil, relErr
		}

		file, openErr := f.FileSystem.Open(rel)
		if openErr != nil {
			logrus.Warnf("fs.Open: failed to open file from embedded FileSystem: %v\n", openErr)
		} else {
			return file, nil
		}
	}

	return os.Open(filePath)
}

// Require fs package
func Require(s *FS) func(runtime *goja.Runtime, module *goja.Object) {
	return func(runtime *goja.Runtime, module *goja.Object) {
		runtime.ToValue(s)
		realpath := runtime.NewObject()
		_ = realpath.Set("native", s.Native)

		promises := runtime.NewObject()
		_ = promises.Set("readFile", s.PromiseReadFile)

		exports := module.Get("exports").(*goja.Object)
		_ = exports.Set("realpath", realpath)
		_ = exports.Set("promises", promises)
		_ = exports.Set("lstat", s.LStat)
		_ = exports.Set("readFile", s.ReadFile)
	}
}

// Enable fs package
func Enable(runtime *goja.Runtime, registry *require.Registry, _ *require.RequireModule, currentWorkingDirectory string, fileSystem fs.FS) {
	s := &FS{
		r:                       runtime,
		CurrentWorkingDirectory: currentWorkingDirectory,
		FileSystem:              fileSystem,
	}

	registry.RegisterNativeModule("node:"+ModuleName, Require(s))
	registry.RegisterNativeModule(ModuleName, Require(s))
	_ = runtime.Set("fs", require.Require(runtime, ModuleName))
}
