package thundra

import (
	"context"
	"encoding/json"
	"sync"
	"fmt"
	"reflect"
	"runtime/debug"

	"thundra-agent-go/plugin"
)

type thundra struct {
	plugins  []plugin.Plugin
	reporter Reporter
	apiKey   string
}

type LambdaFunction func(context.Context, json.RawMessage) (interface{}, error)

func Wrap(handler interface{}, thundra *thundra) LambdaFunction {
	if handler == nil {
		return thundraErrorHandler(fmt.Errorf("handler is nil"))
	}
	handlerType := reflect.TypeOf(handler)
	handlerValue := reflect.ValueOf(handler)

	if handlerType.Kind() != reflect.Func {
		return thundraErrorHandler(fmt.Errorf("handler kind %s is not %s", handlerType.Kind(), reflect.Func))
	}

	takesContext, err := validateArguments(handlerType)

	if err != nil {
		return thundraErrorHandler(err)
	}

	if err := validateReturns(handlerType); err != nil {
		return thundraErrorHandler(err)
	}

	return func(ctx context.Context, payload json.RawMessage) (interface{}, error) {
		defer func() {
			if err := recover(); err != nil {
				stackTrace := debug.Stack()
				thundra.onPanic(ctx, payload, err, stackTrace)
				panic(err)
			}
		}()
		var args []reflect.Value
		if takesContext {
			args = append(args, reflect.ValueOf(ctx))
		}

		if (handlerType.NumIn() == 1 && !takesContext) || handlerType.NumIn() == 2 {
			newEventType := handlerType.In(handlerType.NumIn() - 1)
			newEvent := reflect.New(newEventType)

			if err := json.Unmarshal(payload, newEvent.Interface()); err != nil {
				return nil, err
			}
			args = append(args, newEvent.Elem())
		}

		thundra.executePreHooks(ctx, payload)
		response := handlerValue.Call(args)

		var err error
		if len(response) > 0 {
			if errVal, ok := response[len(response)-1].Interface().(error); ok {
				err = errVal
			}
		}
		var val interface{}
		if len(response) > 1 {
			val = response[0].Interface()
		}

		if err != nil {
			val = nil
		}

		thundra.executePostHooks(ctx, payload, val, err)

		return val, err
	}
}

func (th *thundra) executePreHooks(ctx context.Context, request json.RawMessage) {
	th.reporter.Clear()
	var wg sync.WaitGroup
	wg.Add(len(th.plugins))
	for _, p := range th.plugins {
		go p.BeforeExecution(ctx, request, &wg)
	}
	wg.Wait()
}

func (th *thundra) executePostHooks(ctx context.Context, request json.RawMessage, response interface{}, error interface{}) {
	var wg sync.WaitGroup
	wg.Add(len(th.plugins))
	for _, p := range th.plugins {
		go func() {
			data, dType := p.AfterExecution(ctx, request, response, error, &wg)
			msg := prepareMessage(data, dType, th.apiKey)
			th.reporter.Collect(msg)
		}()
	}
	wg.Wait()
	th.reporter.Report(th.apiKey)
	th.reporter.Clear()
}

func (th *thundra) onPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) {
	var wg sync.WaitGroup
	wg.Add(len(th.plugins))
	for _, p := range th.plugins {
		go func() {
			data, dType := p.OnPanic(ctx, request, err, stackTrace, &wg)
			msg := prepareMessage(data, dType, th.apiKey)
			th.reporter.Collect(msg)
		}()
	}
	wg.Wait()
	th.reporter.Report(th.apiKey)
	th.reporter.Clear()
}

func prepareMessage(data interface{}, dataType string, apiKey string) plugin.Message {

	return plugin.Message{
		data,
		dataType,
		apiKey,
		DataFormatVersion,
	}
}

func thundraErrorHandler(e error) LambdaFunction {
	return func(ctx context.Context, event json.RawMessage) (interface{}, error) {
		return nil, e
	}
}

func validateArguments(handler reflect.Type) (bool, error) {
	handlerTakesContext := false
	if handler.NumIn() > 2 {
		return false, fmt.Errorf("handlers may not take more than two arguments, but handler takes %d", handler.NumIn())
	} else if handler.NumIn() > 0 {
		contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
		argumentType := handler.In(0)
		handlerTakesContext = argumentType.Implements(contextType)
		if handler.NumIn() > 1 && !handlerTakesContext {
			return false, fmt.Errorf("handler takes two arguments, but the first is not Context. got %s", argumentType.Kind())
		}
	}

	return handlerTakesContext, nil
}

func validateReturns(handler reflect.Type) error {
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if handler.NumOut() > 2 {
		return fmt.Errorf("handler may not return more than two values")
	} else if handler.NumOut() > 1 {
		if !handler.Out(1).Implements(errorType) {
			return fmt.Errorf("handler returns two values, but the second does not implement error")
		}
	} else if handler.NumOut() == 1 {
		if !handler.Out(0).Implements(errorType) {
			return fmt.Errorf("handler returns a single value, but it does not implement error")
		}
	}
	return nil
}
