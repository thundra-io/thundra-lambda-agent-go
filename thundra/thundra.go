package thundra

import (
	"context"
	"encoding/json"
	"sync"
	"reflect"
	"fmt"
	"ThundraGo/thundra/plugins"
)

type Thundra struct {
	pluginDictionary map[string]plugins.Plugin
	plugins          []plugins.Plugin
}

//TODO Should be singleton
func New(pluginNames []string) *Thundra {
	th := new(Thundra)
	th.pluginDictionary = make(map[string]plugins.Plugin)

	th.pluginDictionary["trace"] = &plugins.Trace{}

	for _, pN := range pluginNames {
		th.addPlugin(pN)
	}

	return th
}

func (th *Thundra) addPlugin(pluginName string) {
	plugin := th.pluginDictionary[pluginName]
	th.plugins = append(th.plugins, plugin)
}

func (th *Thundra) executePreHooks(ctx context.Context, request json.RawMessage) {
	var wg sync.WaitGroup
	wg.Add(len(th.plugins))
	for _, plugin := range th.plugins {
		go plugin.BeforeExecution(ctx, request, &wg)
	}
	wg.Wait()
}

func (th *Thundra) executePostHooks(ctx context.Context, request json.RawMessage, response interface{}, error interface{}) {
	var wg sync.WaitGroup
	wg.Add(len(th.plugins))
	for _, plugin := range th.plugins {
		go plugin.AfterExecution(ctx, request, response, error, &wg)
	}
	wg.Wait()
}

type ThundraLambdaHandler func(context.Context, json.RawMessage) (interface{}, error)

func Handle(handler interface{}, thundra *Thundra) ThundraLambdaHandler {

	handlerType := reflect.TypeOf(handler)
	handlerValue := reflect.ValueOf(handler)

	takesContext, _ := validateArguments(handlerType)

	return func(ctx context.Context, payload json.RawMessage) (interface{}, error) {

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

		thundra.executePreHooks(ctx,payload)
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

		thundra.executePostHooks(ctx,payload,val, err)

		return val, err
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

