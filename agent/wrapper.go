package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

type lambdaFunction func(context.Context, json.RawMessage) (interface{}, error)

// Wrap is used for wrapping your lambda functions and start monitoring it by following the thundra objects settings
// It wraps your lambda function and return a new lambda function. By that, AWS will be able to run this function
// and Thundra will be able to collect monitoring data from your function.
func (a *Agent) Wrap(handler interface{}) interface{} {
	if config.ThundraDisabled {
		return handler
	}

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
				a.ExecutePostHooks(ctx, payload, nil, err)
				panic(err)
			}
		}()

		// Timeout handler
		go a.CatchTimeout(ctx, payload)

		var args []reflect.Value
		var elem reflect.Value

		if (handlerType.NumIn() == 1 && !takesContext) || handlerType.NumIn() == 2 {
			newEventType := handlerType.In(handlerType.NumIn() - 1)
			newEvent := reflect.New(newEventType)

			if err := json.Unmarshal(payload, newEvent.Interface()); err != nil {
				return nil, err
			}

			elem = newEvent.Elem()
			ctx = utils.SetEventTypeToContext(ctx, newEventType)

			if a.WarmUp && checkAndHandleWarmupRequest(elem, newEventType) {
				return nil, nil
			}
		}

		ctxAfterPreHooks := a.ExecutePreHooks(ctx, payload)

		if takesContext {
			args = append(args, reflect.ValueOf(ctxAfterPreHooks))
		}

		if elem.IsValid() {
			args = append(args, elem)
		}

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

		a.ExecutePostHooks(ctxAfterPreHooks, payload, &val, err)

		return val, err
	}
}

func thundraErrorHandler(e error) lambdaFunction {
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
