package gospectral

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/.spectral.yaml testdata/openapi.yaml
var bundle embed.FS

func TestLint_WorksUsingWorkingDirectory(t *testing.T) {
	t.Parallel()
	// Act
	output, err := Lint([]string{"./openapi.yaml"}, "./.spectral.yaml", WithWorkingDirectory("./testdata"))

	// Assert
	require.NoError(t, err)
	assert.Empty(t, output)
}

func TestLint_WorksWithEmbedFS(t *testing.T) {
	t.Parallel()
	// Act
	output, err := Lint([]string{"./testdata/openapi.yaml"}, "./testdata/.spectral.yaml", WithFS(bundle))

	// Assert
	require.NoError(t, err)
	assert.Empty(t, output)
}

func TestLint_WorksWithEmbedFSAndWorkingDirectory(t *testing.T) {
	t.Parallel()
	// Act
	output, err := Lint([]string{"./openapi.yaml"}, "./testdata/.spectral.yaml", WithFS(bundle), WithWorkingDirectory("./testdata"))

	// Assert
	require.NoError(t, err)
	assert.Empty(t, output)
}

func TestLint_ReportsErrorUsingWorkingDirectory(t *testing.T) {
	t.Parallel()
	// Act
	output, err := Lint([]string{"./openapi-without-contact.yaml"}, "./.spectral.yaml", WithWorkingDirectory("./testdata"))

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, output)
}

func TestLint_ReportsErrorWithEmbedFS(t *testing.T) {
	t.Parallel()
	// Act
	output, err := Lint([]string{"./testdata/openapi-without-contact.yaml"}, "./testdata/.spectral.yaml", WithFS(bundle))

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, output)
}

func TestLint_ReturnsErrorOnNotExistingPath(t *testing.T) {
	t.Parallel()
	// Act
	output, err := Lint([]string{"./doesnotexist.yaml"}, "./.spectral.yaml", WithFS(bundle))

	// Assert
	require.Error(t, err)
	assert.Nil(t, output)
}
