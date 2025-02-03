package events

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// ModuleName of the events package
const ModuleName = "events"

// EventEmitter holds the goja.Runtime to convert Go>JS
type EventEmitter struct {
	r *goja.Runtime
}

// Constructor of EventEmitter
func (e *EventEmitter) Constructor(_ goja.ConstructorCall) *goja.Object {
	obj := e.r.NewObject()
	_ = obj.Set("setMaxListeners", e.SetMaxListeners)
	_ = obj.Set("once", e.Once)
	_ = obj.Set("emit", e.Emit)
	_ = obj.Set("off", e.Emit)

	return obj
}

// SetMaxListeners no-op
func (e *EventEmitter) SetMaxListeners(_ goja.FunctionCall) goja.Value {
	return goja.Undefined()
}

// Once no-op. Normally this contacts the pluginDriver to get all unfulfilled hooks and emits them
func (e *EventEmitter) Once(_ goja.FunctionCall) goja.Value {
	return goja.Undefined()
}

// Emit no-op
func (e *EventEmitter) Emit(_ goja.FunctionCall) goja.Value {
	return goja.Undefined()
}

// Off no-op
func (e *EventEmitter) Off(_ goja.FunctionCall) goja.Value {
	return goja.Undefined()
}

// Require the events package
func Require(runtime *goja.Runtime, module *goja.Object) {
	eventEmitter := &EventEmitter{r: runtime}
	exports := module.Get("exports").(*goja.Object) //nolint:forcetypeassert // based on reference implementation in library
	_ = exports.Set("EventEmitter", eventEmitter.Constructor)
}

// Enable events package with a stub implementation
func Enable(runtime *goja.Runtime, registry *require.Registry, _ *require.RequireModule) {
	registry.RegisterNativeModule("node:"+ModuleName, Require)
	registry.RegisterNativeModule(ModuleName, Require)
	_ = runtime.Set(ModuleName, require.Require(runtime, ModuleName))
}
