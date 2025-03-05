package path

import (
	"testing"

	"github.com/dop251/goja"
	noderequire "github.com/dop251/goja_nodejs/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPath_Dirname(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	p := &Path{r: runtime}

	// Act
	res := p.Dirname(goja.FunctionCall{Arguments: []goja.Value{runtime.ToValue("testdata/file.yaml")}})

	// Assert
	assert.Equal(t, "testdata", res.Export())
}

func TestPath_Resolve(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	p := &Path{r: runtime}

	// Act
	res := p.Resolve(goja.FunctionCall{Arguments: []goja.Value{runtime.ToValue("testdata"), runtime.ToValue("file.yaml")}})

	// Assert
	assert.Equal(t, "testdata/file.yaml", res.Export())
}

func TestPath_Relative_NotImplementedShouldPanic(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	p := &Path{r: runtime}

	// Act
	f := func() {
		p.Relative(goja.FunctionCall{})
	}

	// Assert
	assert.Panics(t, f)
}

func TestPath_Extname(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	p := &Path{r: runtime}

	// Act
	res := p.Extname(goja.FunctionCall{Arguments: []goja.Value{runtime.ToValue("testdata/file.yaml")}})

	// Assert
	assert.Equal(t, ".yaml", res.Export())
}

func TestPath_IsURL(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	p := &Path{r: runtime}

	// Act
	f := func() {
		p.IsURL(goja.FunctionCall{})
	}

	// Assert
	assert.Panics(t, f)
}

func TestPath_Basename(t *testing.T) {
	t.Parallel()
	// Arrange
	runtime := goja.New()
	p := &Path{r: runtime}

	// Act
	res := p.Basename(goja.FunctionCall{Arguments: []goja.Value{runtime.ToValue("testdata/file.yaml")}})

	// Assert
	assert.Equal(t, "file.yaml", res.Export())
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
	assert.NotNil(t, exports.Get("posix"))
	assert.NotNil(t, exports.Get("dirname"))
	assert.NotNil(t, exports.Get("resolve"))
	assert.NotNil(t, exports.Get("relative"))
	assert.NotNil(t, exports.Get("extname"))
	assert.NotNil(t, exports.Get("basename"))
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
