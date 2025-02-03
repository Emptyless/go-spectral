package module

import (
	"testing"

	"github.com/dop251/goja"
	noderequire "github.com/dop251/goja_nodejs/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule_CreateRequire(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	registry := noderequire.NewRegistry()
	requireModule := registry.Enable(runtime)

	registry.RegisterNativeModule("package", func(runtime *goja.Runtime, module *goja.Object) {
		// do something with object
		exports := module.Get("exports").(*goja.Object)
		_ = exports.Set("key", runtime.ToValue("value"))
	})

	module := &Module{
		r:             runtime,
		requireModule: requireModule,
	}

	// Act
	createRequire := module.CreateRequire(goja.FunctionCall{Arguments: []goja.Value{runtime.ToValue("package")}})
	pkg := createRequire.Export().(func(call goja.FunctionCall) goja.Value)(goja.FunctionCall{})

	// Assert
	assert.NotNil(t, createRequire)
	assert.NotNil(t, pkg)
	actual := pkg.Export().(map[string]any)
	assert.Equal(t, "value", actual["key"])
}

func TestRequire(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	module := runtime.NewObject()
	exports := runtime.NewObject()
	_ = module.Set("exports", exports)

	// Act
	Require(runtime, module)

	// Assert
	assert.NotNil(t, exports.Get("createRequire"))
}

func TestEnable(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	registry := noderequire.NewRegistry()
	requireModule := registry.Enable(runtime)

	// Act
	Enable(runtime, registry, requireModule)

	// Assert
	res, err := requireModule.Require(ModuleName)

	// Act
	require.NoError(t, err)
	assert.NotNil(t, res)
}
