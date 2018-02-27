package thundra

import (
	"context"
	"encoding/json"
	"sync"
	"reflect"
	"fmt"
	"ThundraGo/thundra/plugins"
	"github.com/aws/aws-lambda-go/lambda/messages"
)

type thundra struct {
	pluginDictionary map[string]plugins.Plugin
	plugins          []plugins.Plugin
}

var instance *thundra

func GetInstance(pluginNames []string) *thundra{
	if instance == nil{
		instance = createNew(pluginNames)
	}
	return instance
}

func createNew(pluginNames []string) *thundra{
	th := new(thundra)
	//TODO remove pluginDictionary to out
	th.pluginDictionary = make(map[string]plugins.Plugin)
	th.pluginDictionary["trace"] = &plugins.Trace{}

	for _, pN := range pluginNames {
		th.addPlugin(pN)
	}

	return th
}

func (th *thundra) addPlugin(pluginName string) {
	plugin := th.pluginDictionary[pluginName]
	th.plugins = append(th.plugins, plugin)
}

func (th *thundra) executePreHooks(ctx context.Context, request json.RawMessage) {
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
}

func (th *thundra) onError(ctx context.Context, request json.RawMessage, panic *messages.InvokeResponse_Error) {
	fmt.Println("Error Occured")
	fmt.Println("ctx: ",ctx)
	fmt.Println("request: ",request)
	fmt.Println("panic" ,*panic)
	fmt.Println("errbitt")
}

type ThundraLambdaHandler func(context.Context, json.RawMessage) (interface{}, error)

func WrapLambdaHandler(handler interface{}, thundra *thundra) ThundraLambdaHandler {

	handlerType := reflect.TypeOf(handler)
	handlerValue := reflect.ValueOf(handler)

	takesContext, _ := validateArguments(handlerType)

	return func(ctx context.Context, payload json.RawMessage) (interface{}, error) {
		defer func() {
			if err := recover(); err != nil {
				panicInfo := GetPanicInfo(err)
				responseError := &messages.InvokeResponse_Error{
					Message:    panicInfo.Message,
					Type:       getErrorType(err),
					StackTrace: panicInfo.StackTrace,
					ShouldExit: true,
				}
				thundra.onError(ctx,payload,responseError)
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

func getErrorType(err interface{}) string {
	errorType := reflect.TypeOf(err)
	if errorType.Kind() == reflect.Ptr {
		return errorType.Elem().Name()
	}
	return errorType.Name()
}
