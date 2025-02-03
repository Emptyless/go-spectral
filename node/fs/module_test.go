package fs

import (
	"embed"
	"os"
	"testing"

	"github.com/dop251/goja"
	noderequire "github.com/dop251/goja_nodejs/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata
var testdata embed.FS

func TestFS_Native_NotImplemented(t *testing.T) {
	t.Parallel()
	// Arrange
	fs := &FS{
		r: goja.New(),
	}

	// Act
	f := func() {
		fs.Native(goja.FunctionCall{})
	}

	// Assert
	assert.Panics(t, f)
}

func TestFS_LStat_CanReadFile(t *testing.T) {
	t.Parallel()
	// Arrange
	vm := goja.New()
	fs := &FS{r: vm}

	expected, lstatErr := os.Lstat("./module.go")
	require.NoError(t, lstatErr)

	var value goja.Value
	var err goja.Value
	callback := func(call goja.FunctionCall) goja.Value {
		err = call.Argument(0)
		value = call.Argument(1)
		return goja.Undefined()
	}

	// Act
	res := fs.LStat(goja.FunctionCall{
		This:      vm.NewObject(),
		Arguments: []goja.Value{vm.ToValue("./module.go"), vm.ToValue(callback)},
	})

	// Assert
	assert.Equal(t, goja.Undefined(), res)
	assert.Equal(t, goja.Null(), err)
	actual := value.Export().(map[string]any)
	assert.Equal(t, expected.Size(), actual["dev"])
	assert.Equal(t, vm.ToValue(false), actual["isDirectory"].(func(call goja.FunctionCall) goja.Value)(goja.FunctionCall{}))
	assert.Equal(t, vm.ToValue(false), actual["isSymbolicLink"].(func(call goja.FunctionCall) goja.Value)(goja.FunctionCall{}))
	assert.Equal(t, vm.ToValue(false), actual["isBlockDevice"].(func(call goja.FunctionCall) goja.Value)(goja.FunctionCall{}))
	assert.Equal(t, vm.ToValue(false), actual["isCharacterDevice"].(func(call goja.FunctionCall) goja.Value)(goja.FunctionCall{}))
	assert.Equal(t, vm.ToValue(false), actual["isFIFO"].(func(call goja.FunctionCall) goja.Value)(goja.FunctionCall{}))
	assert.Equal(t, vm.ToValue(true), actual["isFile"].(func(call goja.FunctionCall) goja.Value)(goja.FunctionCall{}))
	assert.Equal(t, vm.ToValue(false), actual["isSocket"].(func(call goja.FunctionCall) goja.Value)(goja.FunctionCall{}))
}

func TestFS_LStat_FileNotExists(t *testing.T) {
	t.Parallel()
	// Arrange
	vm := goja.New()
	fs := &FS{r: vm}

	var value goja.Value
	var err goja.Value
	callback := func(call goja.FunctionCall) goja.Value {
		err = call.Argument(0)
		value = call.Argument(1)
		return goja.Undefined()
	}

	// Act
	res := fs.LStat(goja.FunctionCall{
		This:      vm.NewObject(),
		Arguments: []goja.Value{vm.ToValue("./doesnotexist"), vm.ToValue(callback)},
	})

	// Assert
	assert.Equal(t, goja.Undefined(), res)
	assert.Contains(t, err.Export(), "no such file or directory")
	assert.Equal(t, goja.Null(), value)
}

func TestFS_ReadFile(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	fs := &FS{
		r:          runtime,
		FileSystem: testdata,
	}

	var err goja.Value
	var value goja.Value
	callback := func(call goja.FunctionCall) goja.Value {
		err = call.Argument(0)
		value = call.Argument(1)
		return goja.Undefined()
	}

	// Act
	res := fs.ReadFile(goja.FunctionCall{
		Arguments: []goja.Value{runtime.ToValue("testdata/file.yaml"), runtime.ToValue(callback)},
	})

	// Assert
	assert.Equal(t, goja.Undefined(), res)
	assert.Nil(t, err.Export())
	assert.Equal(t, "key: value", value.Export())
}

func TestFS_PromiseReadFile(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	fs := &FS{
		r:                       runtime,
		CurrentWorkingDirectory: "testdata",
		FileSystem:              nil,
	}

	// Act
	res := fs.PromiseReadFile(goja.FunctionCall{Arguments: []goja.Value{runtime.ToValue("testdata/file.yaml")}})
	actual := res.Export().(*goja.Promise).Result()

	// Assert
	assert.Equal(t, "key: value", actual.Export())
}

func TestRequire(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	module := runtime.NewObject()
	exports := runtime.NewObject()
	_ = module.Set("exports", exports)

	// Act
	Require(&FS{r: runtime})(runtime, module)

	// Assert
	assert.NotNil(t, exports.Get("realpath"))
	assert.NotNil(t, exports.Get("promises"))
	assert.NotNil(t, exports.Get("lstat"))
	assert.NotNil(t, exports.Get("readFile"))
}

func TestEnable(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	registry := noderequire.NewRegistry()
	requireModule := registry.Enable(runtime)

	// Act
	Enable(runtime, registry, requireModule, ".", nil)

	// Assert
	res, err := requireModule.Require(ModuleName)

	// Act
	require.NoError(t, err)
	assert.NotNil(t, res)
}
