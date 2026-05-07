package tinyflags

import (
	"context"
	"fmt"
	"reflect"
)

var (
	runnableType = reflect.TypeFor[Runnable]()
	contextType  = reflect.TypeFor[context.Context]()
	errorType    = reflect.TypeFor[error]()
)

// wrapCommandBuilder validates one builder function and adapts it to the internal runner shape.
func wrapCommandBuilder(builder any) commandBuilder {
	if builder == nil {
		panic("tinyflags: command builder cannot be nil")
	}

	value := reflect.ValueOf(builder)
	typ := value.Type()
	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("tinyflags: command builder must be a function, got %T", builder))
	}
	if typ.NumIn() != 0 {
		panic("tinyflags: command builder must not accept arguments")
	}
	if typ.NumOut() < 1 || typ.NumOut() > 2 {
		panic("tinyflags: command builder must return Runnable or (Runnable, error)")
	}
	if !typ.Out(0).Implements(runnableType) {
		panic(fmt.Sprintf("tinyflags: command builder must return a Runnable, got %s", typ.Out(0)))
	}
	if typ.NumOut() == 2 && !typ.Out(1).Implements(errorType) {
		panic(fmt.Sprintf("tinyflags: command builder second return must be an error, got %s", typ.Out(1)))
	}

	return func() (Runnable, error) {
		results := value.Call(nil)
		if len(results) == 2 && !results[1].IsZero() {
			return nil, results[1].Interface().(error)
		}
		if isNilBuilderResult(results[0]) {
			return nil, fmt.Errorf("tinyflags: command builder returned a nil Runnable")
		}

		runner, ok := results[0].Interface().(Runnable)
		if !ok {
			return nil, fmt.Errorf("tinyflags: command builder returned %T, which does not implement Runnable", results[0].Interface())
		}
		return runner, nil
	}
}

// wrapCommandRunner validates one handler plus bound arguments and adapts them to a runnable builder.
func wrapCommandRunner(handler any, bindings ...any) commandBuilder {
	value := reflect.ValueOf(handler)
	if !value.IsValid() {
		panic("tinyflags: command handler cannot be nil")
	}

	spec := parseRunHandler(value.Type())
	return func() (Runnable, error) {
		args, err := resolveRunBindings(spec, bindings)
		if err != nil {
			return nil, err
		}
		return commandHandlerRunner{
			handler: value,
			args:    args,
			spec:    spec,
		}, nil
	}
}

type runHandlerSpec struct {
	paramTypes    []reflect.Type
	paramCount    int
	injectContext bool
	returnsError  bool
}

// parseRunHandler validates one registered command handler signature.
func parseRunHandler(handlerType reflect.Type) runHandlerSpec {
	if handlerType == nil || handlerType.Kind() != reflect.Func {
		panic("tinyflags: command handler must be a function")
	}

	spec := runHandlerSpec{}
	firstArg := 0
	if handlerType.NumIn() > 0 && handlerType.In(0).Implements(contextType) {
		spec.injectContext = true
		firstArg = 1
	}
	spec.paramCount = handlerType.NumIn() - firstArg
	for i := firstArg; i < handlerType.NumIn(); i++ {
		spec.paramTypes = append(spec.paramTypes, handlerType.In(i))
	}

	if handlerType.NumOut() > 1 {
		panic("tinyflags: command handler must return no values or one error")
	}
	if handlerType.NumOut() == 1 {
		if !handlerType.Out(0).Implements(errorType) {
			panic(fmt.Sprintf("tinyflags: command handler return type must be error, got %s", handlerType.Out(0)))
		}
		spec.returnsError = true
	}

	return spec
}

type commandHandlerRunner struct {
	handler reflect.Value
	args    []reflect.Value
	spec    runHandlerSpec
}

// Run executes one registered command handler with its parsed argument values.
func (r commandHandlerRunner) Run(ctx context.Context) error {
	callArgs := make([]reflect.Value, 0, len(r.args)+1)
	if r.spec.injectContext {
		callArgs = append(callArgs, reflect.ValueOf(ctx))
	}
	callArgs = append(callArgs, r.args...)

	results := r.handler.Call(callArgs)
	if !r.spec.returnsError || len(results) == 0 || results[0].IsZero() {
		return nil
	}
	return results[0].Interface().(error)
}

// resolveRunBindings freezes one handler invocation's bound arguments after parsing.
func resolveRunBindings(spec runHandlerSpec, bindings []any) ([]reflect.Value, error) {
	if len(bindings) != spec.paramCount {
		return nil, fmt.Errorf("tinyflags: command handler expects %d bound arguments, got %d", spec.paramCount, len(bindings))
	}

	args := make([]reflect.Value, 0, len(bindings))
	for i, binding := range bindings {
		value, err := resolveRunBinding(binding, spec.paramTypes[i])
		if err != nil {
			return nil, fmt.Errorf("tinyflags: binding %d: %w", i+1, err)
		}
		args = append(args, value)
	}
	return args, nil
}

// resolveRunBinding converts one registered binding into the concrete parameter value a handler needs.
func resolveRunBinding(binding any, paramType reflect.Type) (reflect.Value, error) {
	if binding == nil {
		return reflect.Value{}, fmt.Errorf("nil binding is not supported for %s", paramType)
	}

	value := reflect.ValueOf(binding)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}, fmt.Errorf("nil pointer binding cannot satisfy %s", paramType)
		}
		value = value.Elem()
	}

	if value.Type().AssignableTo(paramType) {
		return freezeRunBindingValue(value), nil
	}
	if value.Type().ConvertibleTo(paramType) {
		return freezeRunBindingValue(value.Convert(paramType)), nil
	}

	return reflect.Value{}, fmt.Errorf("cannot use %s as %s", value.Type(), paramType)
}

// freezeRunBindingValue copies one resolved binding so later flag mutations do not leak into one parsed runner.
func freezeRunBindingValue(value reflect.Value) reflect.Value {
	frozen := reflect.New(value.Type()).Elem()
	frozen.Set(value)
	return frozen
}

// isNilBuilderResult reports whether a reflect value contains a nil value for nilable kinds.
func isNilBuilderResult(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
