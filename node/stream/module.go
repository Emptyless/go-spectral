package stream

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

// ModuleName of the stream package
const ModuleName = "stream"

// Callback to handle a stream
type Callback struct {
	// This object
	This goja.Value
	// Function to call
	Function func(goja.FunctionCall) goja.Value
}

// Stream holds the goja.Runtime and the various methods handled by the Node stream package
type Stream struct {
	r                *goja.Runtime
	data             []goja.Value
	onData           []*Callback
	onDataOffset     map[*Callback]int
	onceError        []*Callback
	onceErrorHandled map[*Callback]bool
	onceEnd          []*Callback
	onceEndHandled   map[*Callback]bool
	ended            bool
	error            goja.Value
}

// PassThrough returns an object implementing Write/Read methods
func (s *Stream) PassThrough(_ goja.ConstructorCall) *goja.Object {
	stream := s.r.NewObject()
	_ = stream.Set("once", s.Once)
	_ = stream.Set("on", s.On)
	_ = stream.Set("write", s.Write)
	_ = stream.Set("end", s.End)
	_ = stream.Set("push", s.Push)

	s.onceErrorHandled = make(map[*Callback]bool)
	s.onceEndHandled = make(map[*Callback]bool)
	s.onDataOffset = make(map[*Callback]int)

	return stream
}

// Write calls _write function on the callee with the _writecall callback
func (s *Stream) Write(call goja.FunctionCall) goja.Value {
	exp := call.This.Export().(map[string]any)
	_write := exp["_write"].(func(goja.FunctionCall) goja.Value)
	_writecall := goja.FunctionCall{
		This:      call.This,
		Arguments: []goja.Value{call.Arguments[0], s.r.ToValue("utf8"), s.r.ToValue(s.OnWrite)},
	}

	return _write(_writecall)
}

// OnWrite handle using Stream.Handle
func (s *Stream) OnWrite(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) > 0 {
		s.error = call.Argument(0)
		s.Handle()
	}

	return goja.Null()
}

// Push data onto Stream
func (s *Stream) Push(call goja.FunctionCall) goja.Value {
	s.data = append(s.data, call.Argument(0))
	return goja.Null()
}

// End Stream
func (s *Stream) End(_ goja.FunctionCall) goja.Value {
	s.ended = true
	s.Handle()
	return goja.Null()
}

// Handle Stream by checking if there is an error and returning if-so. Handle tracks the
// callbacks registered on the Stream to ensure that they are only called Once.
func (s *Stream) Handle() { //nolint:cyclop // accepted complexity
	if s.error != nil {
		for _, cb := range s.onceError {
			_, ok := s.onceErrorHandled[cb]
			if !ok {
				call := goja.FunctionCall{
					This:      cb.This,
					Arguments: []goja.Value{s.error},
				}
				cb.Function(call)
				s.onceErrorHandled[cb] = true
			}
		}

		return
	}

	for _, cb := range s.onData {
		offset, ok := s.onDataOffset[cb]
		if !ok {
			offset = -1
		}
		for i, data := range s.data {
			if i <= offset {
				continue
			}

			call := goja.FunctionCall{
				This:      cb.This,
				Arguments: []goja.Value{data},
			}

			cb.Function(call)
			s.onDataOffset[cb] = i
		}
	}

	if s.ended {
		for _, cb := range s.onceEnd {
			_, ok := s.onceEndHandled[cb]
			if !ok {
				call := goja.FunctionCall{
					This:      cb.This,
					Arguments: []goja.Value{s.error},
				}
				cb.Function(call)
				s.onceEndHandled[cb] = true
			}
		}

		return
	}
}

// On event (e.g. 'data') invoke a callback
func (s *Stream) On(call goja.FunctionCall) goja.Value {
	if call.Argument(0).String() == "data" {
		s.onData = append(s.onData, &Callback{
			This:     call.This,
			Function: call.Argument(1).Export().(func(functionCall goja.FunctionCall) goja.Value),
		})
		s.Handle()

		return goja.Null()
	}

	panic("stream.on called with unknown value")
}

// Once some event occurs (e.g. error, end) invoke a callback
func (s *Stream) Once(call goja.FunctionCall) goja.Value {
	arg := call.Argument(0).String()
	if arg == "error" {
		s.onceError = append(s.onceError, &Callback{
			This:     call.This,
			Function: call.Argument(1).Export().(func(functionCall goja.FunctionCall) goja.Value),
		})
		s.Handle()

		return goja.Null()
	} else if arg == "end" {
		s.onceEnd = append(s.onceEnd, &Callback{
			This:     call.This,
			Function: call.Argument(1).Export().(func(functionCall goja.FunctionCall) goja.Value),
		})
		s.Handle()

		return goja.Null()
	}

	panic("stream.once called with unknown value")
}

// Readable implementation
type Readable struct {
	r *goja.Runtime
}

// Constructor of Readable
func (r *Readable) Constructor(_ goja.ConstructorCall) *goja.Object {
	return r.r.NewObject()
}

// Equals is always true
func (r *Readable) Equals(_ goja.FunctionCall) goja.Value {
	return r.r.ToValue(true)
}

// ToString of Readable
func (r *Readable) ToString(_ goja.FunctionCall) goja.Value {
	return r.r.ToValue("toString of Stream")
}

// readable always returns true
func (r *Readable) readable(_ goja.FunctionCall) goja.Value {
	return r.r.ToValue(true)
}

// Require Stream
// @See https://nodejs.org/api/stream.html#class-streamreadable
func Require(runtime *goja.Runtime, module *goja.Object) {
	s := &Stream{r: runtime}
	runtime.ToValue(s)

	r := &Readable{r: runtime}

	readable := runtime.NewObject()

	proto := runtime.NewObject()
	_ = proto.SetPrototype(runtime.NewObject())
	_ = proto.DefineDataProperty("constructor", runtime.ToValue(r.Constructor), goja.FLAG_TRUE, goja.FLAG_TRUE, goja.FLAG_FALSE)
	_ = proto.Set("equals", r.Equals)
	_ = proto.Set("toString", r.ToString)

	_ = readable.SetPrototype(proto)
	_ = readable.Set("readable", r.readable)
	_ = readable.Set("prototype", proto)

	exports := module.Get("exports").(*goja.Object)
	_ = exports.Set("Readable", readable)
	_ = exports.Set("PassThrough", s.PassThrough)
}

func Enable(runtime *goja.Runtime, registry *require.Registry, _ *require.RequireModule) {
	registry.RegisterNativeModule("node:"+ModuleName, Require)
	registry.RegisterNativeModule(ModuleName, Require)
	_ = runtime.Set("Stream", require.Require(runtime, ModuleName))
}
