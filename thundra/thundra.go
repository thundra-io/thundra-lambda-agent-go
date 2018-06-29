package thundra

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type thundra struct {
	plugins       []plugin.Plugin
	reporter      reporter
	apiKey        string
	warmup        bool
	timeoutMargin time.Duration
}

func (t *thundra) executePreHooks(ctx context.Context, request json.RawMessage) {
	t.reporter.Clear()
	plugin.GenerateNewTransactionId()
	var wg sync.WaitGroup
	wg.Add(len(t.plugins))
	for _, p := range t.plugins {
		go p.BeforeExecution(ctx, request, &wg)
	}
	wg.Wait()
}

func (t *thundra) executePostHooks(ctx context.Context, request json.RawMessage, response interface{}, error interface{}) {
	var wg sync.WaitGroup
	wg.Add(len(t.plugins))
	for _, p := range t.plugins {
		go func(plugin plugin.Plugin) {
			data, dType := plugin.AfterExecution(ctx, request, response, error)
			messages := prepareMessages(data, dType, t.apiKey)
			t.reporter.Collect(messages)
			wg.Done()
		}(p)
	}
	wg.Wait()
	t.reporter.Report(t.apiKey)
	t.reporter.Clear()
}

func (t *thundra) onPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) {
	var wg sync.WaitGroup
	wg.Add(len(t.plugins))
	for _, p := range t.plugins {
		go func(plugin plugin.Plugin) {
			data, dType := plugin.OnPanic(ctx, request, err, stackTrace)
			messages := prepareMessages(data, dType, t.apiKey)
			t.reporter.Collect(messages)
			wg.Done()
		}(p)
	}
	wg.Wait()
	t.reporter.Report(t.apiKey)
	t.reporter.Clear()
}

func prepareMessages(data []interface{}, dataType string, apiKey string) []interface{} {
	var messages []interface{}
	for _, d := range data {
		m := plugin.Message{
			Data:              d,
			Type:              dataType,
			ApiKey:            apiKey,
			DataFormatVersion: dataFormatVersion,
		}
		messages = append(messages, m)
	}
	return messages
}

type timeoutError struct{}

func (e timeoutError) Error() string {
	return fmt.Sprintf("Lambda is timed out")
}

// catchTimeout is checks for a timeout event and sends report if lambda is timedout
func (t *thundra) catchTimeout(ctx context.Context, payload json.RawMessage) {
	deadline, _ := ctx.Deadline()
	if deadline.IsZero() {
		return
	}

	var timeoutMargin time.Duration

	if t.timeoutMargin != 0 {
		timeoutMargin = t.timeoutMargin
	} else {
		timeoutMargin = defaultTimeoutMargin * time.Millisecond
	}

	timeoutDuration := deadline.Add(-timeoutMargin)

	if time.Now().After(timeoutDuration) {
		return
	}

	timeoutChannel := time.After(time.Until(timeoutDuration))

	select {
	case <-timeoutChannel:
		fmt.Println("Function is timed out")
		t.executePostHooks(ctx, payload, nil, timeoutError{})
		return
	case <-ctx.Done():
		return
	}
}
