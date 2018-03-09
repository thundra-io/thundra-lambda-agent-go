package thundra

import (
	"context"
	"encoding/json"
	"sync"
	"reflect"
	"fmt"
	"runtime/debug"
	"os"
	"thundra-agent-go/constants"
)

type thundra struct {
	plugins   []Plugin
	collector Collector
}

var ApiKey string

func init() {
	discoverPlugins()
}

func CreateNew(pluginNames []string) *thundra {
	c := new(collectorImpl)
	return createNewWithCollector(pluginNames, c)
}

func createNewWithCollector(pluginNames []string, collector Collector) *thundra {
	th := new(thundra)
	th.collector = collector
	for _, pN := range pluginNames {
		if pf := pluginDictionary[pN]; pf != nil {
			p := pf.Create()
			var i interface{} = p
			cp, ok := i.(CollecterAwarePlugin)
			if ok {
				cp.SetCollector(collector)
			}
			th.addPlugin(p)
		} else {
			fmt.Println("Invalid Plugin Name: %s ", pN)
		}
	}
	ApiKey = os.Getenv(constants.THUNDRA_API_KEY)
	return th
}

func (th *thundra) addPlugin(plugin Plugin) {
	th.plugins = append(th.plugins, plugin)
}

func (th *thundra) executePreHooks(ctx context.Context, request json.RawMessage) {
	th.collector.clear()
	var wg sync.WaitGroup
	wg.Add(len(th.plugins))
	for _, plugin := range th.plugins {
		go plugin.BeforeExecution(ctx, request, &wg)
	}
	wg.Wait()
}

func (th *thundra) executePostHooks(ctx context.Context, request json.RawMessage, response interface{}, error interface{}) {
	var wg sync.WaitGroup
	wg.Add(len(th.plugins))
	for _, plugin := range th.plugins {
		go plugin.AfterExecution(ctx, request, response, error, &wg)
	}
	wg.Wait()
	th.collector.report()
	th.collector.clear()
}

func (th *thundra) onPanic(ctx context.Context, request json.RawMessage, panic *ThundraPanic) {
	var wg sync.WaitGroup
	wg.Add(len(th.plugins))
	for _, plugin := range th.plugins {
		go plugin.OnPanic(ctx, request, panic, &wg)
	}
	wg.Wait()
	th.collector.report()
	th.collector.clear()
}

type thundraLambdaHandler func(context.Context, json.RawMessage) (interface{}, error)

func thundraErrorHandler(e error) thundraLambdaHandler {
	return func(ctx context.Context, event json.RawMessage) (interface{}, error) {
		return nil, e
	}
}

func WrapLambdaHandler(handler interface{}, thundra *thundra) thundraLambdaHandler {
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
				panicInfo := ThundraPanic{
					//TODO pass error only
					ErrMessage: err.(error).Error(),
					StackTrace: string(debug.Stack()), //fmt.Sprintf("%s: %s", err, debug.Stack()),
					ErrType:    getErrorType(err),
				}
				thundra.onPanic(ctx, payload, &panicInfo)
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
