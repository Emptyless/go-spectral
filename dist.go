package gospectral

import (
	"embed"

	"github.com/dop251/goja_nodejs/require"
)

// DistName used to import built.js
// e.g. with `var spectral = require('./dist/built.js');`
const DistName = "./dist/built.js"

//go:embed dist/built.js
var dist []byte

// Licenses includes the built.js.LICENSE.txt and oss-licenses.json file produced by Webpack for the dist
//
//go:embed dist/built.js.LICENSE.txt dist/oss-licenses.json
var Licenses embed.FS

// DefaultDist returning the transpiled source with the library
func DefaultDist() []byte {
	return dist
}

// WithDist sets the Config.Dist to a custom supplied value. This can be useful for
// using a specific version of the source and/or bundling it on your own.
func WithDist(dist []byte) Option {
	return func(config *Config) error {
		config.Dist = dist

		return nil
	}
}

func EnableDist(require *require.RequireModule) error {
	_, err := require.Require(DistName)
	return err
}
