package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

// Agent is thundra agent implementation
type Agent struct {
	Plugins       []plugin.Plugin
	Reporter      reporter
	WarmUp        bool
	TimeoutMargin time.Duration
}

// New is used to collect basic invocation data with thundra. Use NewBuilder and AddPlugin to access full functionality.
func New() *Agent {
	w := determineWarmup()
	g := determineTimeoutMargin()
	determineAPIKey()

	a := &Agent{
		Reporter: &reporterImpl{
			reported: new(uint32),
		},
		WarmUp:        w,
		TimeoutMargin: g,
		Plugins:       []plugin.Plugin{},
	}

	return a
}

// AddPlugin is used to enable plugins on thundra. You can use Trace, Metrics and Log plugins.
// You need to initialize a plugin object and pass it as a parameter in order to enable it.
// e.g. AddPlugin(trace.New())
func (a *Agent) AddPlugin(plugin plugin.Plugin) *Agent {
	if plugin.IsEnabled() {
		a.Plugins = append(a.Plugins, plugin)
	}

	return a
}

// SetReporter sets agent reporter
func (a *Agent) SetReporter(r reporter) *Agent {
	a.Reporter = r
	return a
}

// ExecutePreHooks contains necessary works that should be done before user's handler
func (a *Agent) ExecutePreHooks(ctx context.Context, request json.RawMessage) context.Context {
	a.Reporter.FlushFlag()
	plugin.TraceID = plugin.GenerateNewID()
	plugin.TransactionID = plugin.GenerateNewID()

	updatedCtx := ctx
	for _, p := range a.Plugins {
		updatedCtx = p.BeforeExecution(updatedCtx, request)
	}

	return updatedCtx
}

// ExecutePostHooks contains necessary works that should be done after user's handler
func (a *Agent) ExecutePostHooks(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) {
	// Skip if it is already reported
	if *a.Reporter.Reported() == 1 {
		return
	}
	for _, p := range a.Plugins {
		messages := p.AfterExecution(ctx, request, response, err)
		a.Reporter.Collect(messages)
	}
	a.Reporter.Report()
	a.Reporter.ClearData()
}

// CatchTimeout is checks for a timeout event and sends report if lambda is timedout
func (a *Agent) CatchTimeout(ctx context.Context, payload json.RawMessage) {
	deadline, _ := ctx.Deadline()
	if deadline.IsZero() {
		return
	}

	var timeoutMargin time.Duration

	if a.TimeoutMargin != 0 {
		timeoutMargin = a.TimeoutMargin
	} else {
		timeoutMargin = defaultTimeoutMargin * time.Millisecond
	}

	timeoutDuration := deadline.Add(-timeoutMargin)

	if time.Now().After(timeoutDuration) {
		return
	}

	timer := time.NewTimer(time.Until(timeoutDuration))
	timeoutChannel := timer.C

	select {
	case <-timeoutChannel:
		fmt.Println("Function is timed out")
		a.ExecutePostHooks(ctx, payload, nil, timeoutError{})
		return
	case <-ctx.Done():
		// close timeoutChannel
		timer.Stop()
		return
	}
}
