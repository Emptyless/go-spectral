package os

import (
	"runtime"
	"testing"

	"github.com/dop251/goja"
	noderequire "github.com/dop251/goja_nodejs/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOS_Platform(t *testing.T) {
	t.Parallel()
	// Arrange
	r := goja.New()
	os := &OS{r: r}

	// Act
	res := os.Platform(goja.FunctionCall{})

	// Assert
	assert.Equal(t, r.ToValue(runtime.GOOS), res)
}

func TestOS_CPUs(t *testing.T) {
	t.Parallel()
	// Arrange
	r := goja.New()
	os := &OS{r: r}

	// Act
	res := os.CPUs(goja.FunctionCall{})

	// Assert
	assert.Equal(t, []map[string]any{
		{
			"model": "vCPU",
			"speed": 2000,
			"times": map[string]any{
				"user": 1,
				"nice": 0,
				"sys":  1,
				"idle": 0,
				"irq":  0,
			},
		},
	}, res.Export())
}

func TestRequire(t *testing.T) {
	t.Parallel()
	// Arrange
	r := goja.New()
	module := r.NewObject()
	exports := r.NewObject()
	_ = module.Set("exports", exports)

	// Act
	Require(r, module)

	// Assert
	assert.NotNil(t, exports.Get("platform"))
	assert.NotNil(t, exports.Get("cpus"))
}

func TestEnable(t *testing.T) {
	t.Parallel()
	// Arrange
	r := goja.New()
	registry := noderequire.NewRegistry()
	requireModule := registry.Enable(r)

	// Act
	Enable(r, registry, requireModule)

	// Assert
	res, err := requireModule.Require(ModuleName)

	// Act
	require.NoError(t, err)
	assert.NotNil(t, res)
}
