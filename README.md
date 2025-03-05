# Go-Spectral

_This tool is under active development and must be considered alpha. It's API may be changed in a breaking way until a
1.0 version is released. Submit issues to the Github issue tracker if found._

`go-spectral` is a Go wrapper for [stoplightio/spectral](https://github.com/stoplightio/spectral) using the pure Go
ECMAScript 5.1 implementation [Goja](https://github.com/dop251/goja) with Node support
using [goja_nodejs](https://github.com/dop251/goja_nodejs). The additional Node system calls are translated (roughly) to
either no-op's (if the function during preliminary testing was not used) or the respective Go SDK methods.

### Import

```
$ go get github.com/Emptyless/go-spectral@latest
```

### Quickstart

```go
package main

import (
	"fmt"

	"github.com/Emptyless/go-spectral"
)

func main() {
	output, err := gospectral.Lint([]string{"./openapi.yaml"}, "./.spectral.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}
```

### Options

To customise the behavior of `Lint`, additional options can be supplied:

- `WithWorkingDirectory`: sets the working directory used to load system files (e.g. .spectral.yaml)
- `WithFS`: sets the `Config.FS` to load documents and rulesets from. This can be useful when using e.g. `embed.FS` as a
  means to bundle specs/rulesets.
- `WithDist`: sets the `Config.Dist` to a custom supplied value. This can be useful for using a specific version of the
  source and/or bundling it on your own.
- `WithScript`: sets the `Config.Script` to a custom value
- `WithBeforeModule`: is evaluated before any module is enabled and can be used to change e.g. the Loader. If a non-nil
  `Enable.Fn` is returned, it is used instead of the provided Enable
- `WithAfterModule`: is evaluated after any module is enabled and can be used to change the current runtime state

### TODO's

- [x] get basic structure of the wrapper working
- [ ] add tests

### Mentions

- [stoplightio/spectral](https://github.com/stoplightio/spectral): JSON/YAML linter with support for OpenAPI
- [dop251/goja](https://github.com/dop251/goja): Go ECMAScript 5.1 Implementation
- [dop251/goja_nodejs](https://github.com/dop251/goja_nodejs): Goja implementations of NodeJS packages

