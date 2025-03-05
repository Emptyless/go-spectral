package process

import (
	"os"
	"testing"

	"github.com/dop251/goja"
	noderequire "github.com/dop251/goja_nodejs/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcess_On(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	process := &Process{r: runtime}

	// Act
	res := process.On(goja.FunctionCall{})

	// Assert
	assert.Equal(t, goja.Undefined(), res)
}

func TestProcess_Cwd(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	process := &Process{r: runtime}
	expected, err := os.Getwd()
	require.NoError(t, err)

	// Act
	res := process.Cwd(goja.FunctionCall{})

	// Assert
	assert.Equal(t, expected, res.Export())
}

func TestProcess_Versions(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	process := &Process{r: runtime}

	// Act
	res := process.Versions()

	// Assert
	assert.Equal(t, Versions, res.Export())
}

func TestRequire(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	p := &Process{r: runtime}

	module := runtime.NewObject()
	exports := runtime.NewObject()
	_ = module.Set("exports", exports)

	// Act
	Require(p)(runtime, module)

	// Assert
	assert.NotNil(t, exports.Get("env"))
	assert.NotNil(t, exports.Get("on"))
	assert.NotNil(t, exports.Get("versions"))
	assert.NotNil(t, exports.Get("version"))
	assert.NotNil(t, exports.Get("cwd"))
	assert.NotNil(t, exports.Get("stdout"))
	assert.NotNil(t, exports.Get("stderr"))
	assert.NotNil(t, exports.Get("argv"))
}

func TestEnable(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	registry := noderequire.NewRegistry()
	requireModule := registry.Enable(runtime)

	// Act
	Enable(runtime, registry, requireModule, "")

	// Assert
	res, err := requireModule.Require(ModuleName)

	// Act
	require.NoError(t, err)
	assert.NotNil(t, res)
}
