package global

import (
	"testing"

	"github.com/dop251/goja"
	noderequire "github.com/dop251/goja_nodejs/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	assert.NotNil(t, exports.Get("global"))
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
