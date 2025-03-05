package gospectral

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"strings"
	"sync"
	"unsafe"

	"github.com/dop251/goja"
	noderequire "github.com/dop251/goja_nodejs/require"
)

// Config to instantiate the runner
type Config struct {
	// Dist is the source JS file exporting
	// 1) function: formatOutput
	// 2) promise: lint output using global variables lintDocuments and lintRuleset
	//
	// See the index.js for more details
	Dist []byte

	// Script that evaluates the Dist, i.e. by the default
	// by using the output of the lint promise and wrapping
	// it in a formatOutput promise
	Script []byte

	// FS to use when loading document(s) and ruleset. If set, when a file is loaded (e.g. an OpenAPI document or
	// a ruleset) it is first searched in the FS *without* the WorkingDirectory reference. If the file is not found
	// in the FS, continue the search relative to the WorkingDirectory. This mechanism allows to bundle static files
	// (e.g. rulesets) and also have a runtime component.
	FS fs.FS

	// WorkingDirectory, defaults to os.Getcwd() if ""
	WorkingDirectory string

	// BeforeModule hook to customize behavior before (or instead of) enabling a module
	BeforeModule BeforeModule

	// AfterModule hook to customize runtime or registry state
	AfterModule AfterModule
}

// Output is a slice of Rule
type Output []Rule

// Rule is an instance of a failure during Lint
type Rule struct {
	Source   string   `json:"source"`
	Code     string   `json:"code"`
	Path     []string `json:"path"`
	Message  string   `json:"message"`
	Severity int      `json:"severity"`
	Range    struct {
		Start struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"start"`
		End struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"end"`
	}
}

// lintDocuments global variable name when providing a custom dist
const lintDocuments = "lintDocuments"

const lintRuleset = "lintRuleset"

// Option that can be supplied to modify the Config
type Option func(config *Config) error

// lock the Lint method to a single invocation at once
var lock sync.Mutex

// Lint OpenAPI documents (e.g. openapi.yaml) with a Spectral ruleset, e.g. `extends: ["spectral:oas"]`
func Lint(documents []string, ruleset string, options ...Option) (Output, error) { //nolint:cyclop // accepted complexity
	// make the Lint method somewhat thread safe (depends on global variables in node packages)
	lock.Lock()
	defer lock.Unlock()

	// instantiate default Config
	cfg := &Config{
		Dist:         DefaultDist(),
		Script:       DefaultScript(),
		BeforeModule: nil,
		AfterModule:  nil,
	}

	// apply Option's
	for _, option := range options {
		if err := option(cfg); err != nil {
			return nil, err
		}
	}

	// Set working directory if ""
	if cfg.WorkingDirectory == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		cfg.WorkingDirectory = wd
	}

	// Set default BeforeModule if nil such that node:fs and node:process can use the working directory and/or virtual file system
	if cfg.BeforeModule == nil {
		cfg.BeforeModule = DefaultBeforeModule(cfg.WorkingDirectory, cfg.FS)
	}

	// initiate runtime with NodeJS modules
	runtime := goja.New()
	registry := noderequire.NewRegistry(noderequire.WithLoader(func(path string) ([]byte, error) {
		if path == DistName {
			return cfg.Dist, nil
		}

		return noderequire.DefaultSourceLoader(path)
	}))
	require, loadModulesErr := LoadModules(runtime, registry, cfg.BeforeModule, cfg.AfterModule)
	if loadModulesErr != nil {
		return nil, loadModulesErr
	}

	// set the __dirname global to the working directory
	if err := runtime.GlobalObject().Set("__dirname", runtime.ToValue(cfg.WorkingDirectory)); err != nil {
		return nil, err
	}

	// set the lintDocuments global variable
	if err := runtime.GlobalObject().Set(lintDocuments, runtime.ToValue(documents)); err != nil {
		return nil, err
	}

	// set the lintRuleset global variable
	if err := runtime.GlobalObject().Set(lintRuleset, runtime.ToValue(ruleset)); err != nil {
		return nil, err
	}

	if err := EnableDist(require); err != nil {
		return nil, &EvaluateError{Err: err}
	}

	// run the script
	v, err := runtime.RunString(string(cfg.Script))
	if err != nil {
		return nil, &EvaluateError{Err: err}
	}

	// if the result is a goja.Promise, wait for completion
	value := v.Export()
	promise, ok := value.(*goja.Promise)
	if ok {
		for promise.State() == goja.PromiseStatePending {
			continue
		}
		if promise.State() == goja.PromiseStateRejected {
			return nil, fmt.Errorf("%s: %w", promise.Result().String(), ErrPromiseRejected)
		}
		value = promise.Result().String()
	}

	if _, ok := value.(string); !ok {
		return nil, fmt.Errorf("invalid value type of Lint result '%T': %w", value, ErrUnknownReturn)
	}

	var output Output
	if err := json.Unmarshal([]byte(value.(string)), &output); err != nil {
		return nil, err
	}

	return output, nil
}

// WithWorkingDirectory sets the working directory used to load system files (e.g. .spectral.yaml)
func WithWorkingDirectory(workingDirectory string) Option {
	return func(config *Config) error {
		config.WorkingDirectory = workingDirectory

		return nil
	}
}

// WithFS sets the Config.FS to load documents and rulesets from. This can be useful when using e.g. embed.FS as a
// means to bundle specs/rulesets.
func WithFS(fs fs.FS) Option {
	return func(config *Config) error {
		config.FS = fs

		return nil
	}
}

// ErrUnknownReturn when the result of the operation is not a string type
var ErrUnknownReturn = errors.New("unknown return")

// ErrPromiseRejected when the JS promise is rejected
var ErrPromiseRejected = errors.New("promise rejected")

// EvaluateError translates various failure cases in an easier to understand format
type EvaluateError struct {
	Err error
}

// Error implementation of EvaluateError
func (e EvaluateError) Error() string {
	var exception *goja.Exception
	if !errors.As(e.Err, &exception) {
		return e.Err.Error()
	}

	value := exception.Value()
	export := value.Export()
	if stacks := exception.Stack(); len(stacks) == 0 {
		return fmt.Sprintf("%v\n%s", export, exception.String())
	} else if stack := stacks[0]; reflect.ValueOf(stack).IsZero() {
		return fmt.Sprintf("%v\n%s", export, exception.String())
	} else if prg := reflect.ValueOf(stack).FieldByName("prg"); !prg.IsZero() {
		f := prg.Elem().FieldByName("src")
		field := f.Elem().FieldByName("src")
		src := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface().(string)
		position := stack.Position()
		lines := strings.Split(src, "\n")
		line := lines[position.Line-1]

		return fmt.Sprintf("%v\n...\n%s\n...\n\n%s", export, line, exception.Error())
	} else if stack.FuncName() == "github.com/dop251/goja_nodejs/require.(*RequireModule).require-fm" {
		return fmt.Sprintf("%v\nfailed to import '%s'\n\n%s", export, stacks[1].FuncName(), exception.Error())
	}

	return fmt.Sprintf("%v\n%s", export, exception.String())
}
