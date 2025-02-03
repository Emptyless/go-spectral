package gospectral

import _ "embed"

//go:embed script.js
var script []byte

// DefaultScript returns the default script that evaluates the DefaultDist
func DefaultScript() []byte {
	return script
}

// WithScript sets the Config.Script to a custom value
func WithScript(script []byte) Option {
	return func(config *Config) error {
		config.Script = script

		return nil
	}
}
