package thundra

import (
	"context"
	"encoding/json"
	"sync"
	"fmt"
	"reflect"
	"os"

	"thundra-agent-go/plugin"
)

var apiKey string

type thundra struct {
	plugins  []plugin.Plugin
	reporter reporter
}

func init() {
	apiKey = os.Getenv(plugin.ThundraApiKey)
}

type LambdaFunction func(context.Context, json.RawMessage) (interface{}, error)

func WrapLambdaHandler(handler interface{}, thundra *thundra) LambdaFunction {
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
				//TODO pass error only
				/*panicInfo := ThundraPanic{
					ErrMessage: err.(error).Error(),
					StackTrace: string(debug.Stack()), //fmt.Sprintf("%s: %s", err, debug.Stack()),
					ErrType:    getErrorType(err),
				}
				thundra.onPanic(ctx, payload, &panicInfo)*/
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

		thundra.executePostHooks(ctx, payload, val, err)

		return val, err
	}
}

func (th *thundra) executePreHooks(ctx context.Context, request json.RawMessage) {
	th.reporter.clear()
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
			msg := p.AfterExecution(ctx, request, response, error, &wg)
			th.reporter.collect(msg)
		}()
	}
	wg.Wait()
	th.reporter.report()
	th.reporter.clear()
}

func (th *thundra) onPanic(ctx context.Context, request json.RawMessage, panic interface{}) {
	var wg sync.WaitGroup
	wg.Add(len(th.plugins))
	for _, p := range th.plugins {
		go func() {
			msg := p.OnPanic(ctx, request, panic, &wg)
			th.reporter.collect(msg)
		}()
	}
	wg.Wait()
	th.reporter.report()
	th.reporter.clear()
}

func thundraErrorHandler(e error) LambdaFunction {
	return func(ctx context.Context, event json.RawMessage) (interface{}, error) {
		return nil, e
	}
}

//Taken from Amazon Inc
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

func getErrorType(err interface{}) string {
	errorType := reflect.TypeOf(err)
	if errorType.Kind() == reflect.Ptr {
		return errorType.Elem().Name()
	}
	return errorType.Name()
}