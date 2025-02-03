package events

import (
	"testing"

	"github.com/dop251/goja"
	noderequire "github.com/dop251/goja_nodejs/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventEmitter_Constructor(t *testing.T) {
	t.Parallel()
	// Arrange
	emitter := &EventEmitter{r: goja.New()}

	// Act
	obj := emitter.Constructor(goja.ConstructorCall{})

	// Assert
	actual := obj.Export()
	assert.Len(t, actual, 4)
	assert.Contains(t, actual, "setMaxListeners")
	assert.Contains(t, actual, "once")
	assert.Contains(t, actual, "emit")
	assert.Contains(t, actual, "off")
}

func TestEventEmitter_SetMaxListeners(t *testing.T) {
	t.Parallel()
	// Arrange
	emitter := &EventEmitter{r: goja.New()}

	// Act
	obj := emitter.SetMaxListeners(goja.FunctionCall{})

	// Assert
	assert.Equal(t, goja.Undefined(), obj)
}

func TestEventEmitter_SetOnce(t *testing.T) {
	t.Parallel()
	// Arrange
	emitter := &EventEmitter{r: goja.New()}

	// Act
	obj := emitter.Once(goja.FunctionCall{})

	// Assert
	assert.Equal(t, goja.Undefined(), obj)
}

func TestEventEmitter_SetOff(t *testing.T) {
	t.Parallel()
	// Arrange
	emitter := &EventEmitter{r: goja.New()}

	// Act
	obj := emitter.Off(goja.FunctionCall{})

	// Assert
	assert.Equal(t, goja.Undefined(), obj)
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
	assert.NotNil(t, exports.Get("EventEmitter"))
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
