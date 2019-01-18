package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
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

// determineApiKey determines which apiKey to use. if apiKey is set from environment variable, returns that value.
// Otherwise returns the value from builder's setApiKey method. Panic if it's not set by neither.
func determineAPIKey() {
	k := os.Getenv(thundraApiKey)
	if k == "" {
		// TODO remove panics just log
		fmt.Println("Error no APIKey in env variables")
	}

	// Set it globally
	plugin.ApiKey = k
}

// determineWarmup determines which warmup value to use. if warmup is set from environment variable, returns that value.
// Otherwise returns true if it's enabled by builder's enableWarmup method. Default value is false.
func determineWarmup() bool {
	w := os.Getenv(thundraLambdaWarmupWarmupAware)
	b, err := strconv.ParseBool(w)
	if err != nil {
		if w != "" {
			fmt.Println(err, " thundra_lambda_warmup_warmupAware should be set with a boolean.")
		}
		return false
	}
	return b
}

// determineTimeoutMargin fetches thundraLambdaTimeoutMargin if it exist, if not returns default timegap value
func determineTimeoutMargin() time.Duration {
	t := os.Getenv(thundraLambdaTimeoutMargin)
	// environment variable is not set
	if t == "" {
		return time.Duration(defaultTimeoutMargin)
	}

	i, err := strconv.ParseInt(t, 10, 32)

	// environment variable is not set in the correct format
	if err != nil {
		fmt.Println(err, " "+thundraLambdaTimeoutMargin+" should be set with an integer.")
		return time.Duration(defaultTimeoutMargin)
	}

	return time.Duration(i) * time.Millisecond
}

func (a *Agent) ExecutePreHooks(ctx context.Context, request json.RawMessage) {
	a.Reporter.FlushFlag()
	plugin.TraceID = plugin.GenerateNewID()
	plugin.TransactionID = plugin.GenerateNewID()
	var wg sync.WaitGroup
	wg.Add(len(a.Plugins))
	for _, p := range a.Plugins {
		go p.BeforeExecution(ctx, request, &wg)
	}
	wg.Wait()
}

func (a *Agent) ExecutePostHooks(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) {
	// Skip if it is already reported
	if *a.Reporter.Reported() == 1 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(a.Plugins))
	for _, p := range a.Plugins {
		go func(plugin plugin.Plugin) {
			messages := plugin.AfterExecution(ctx, request, response, err)
			a.Reporter.Collect(messages)
			wg.Done()
		}(p)
	}
	wg.Wait()
	a.Reporter.Report()
	a.Reporter.ClearData()
}

func (a *Agent) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) {
	// Skip if it is already reported
	if *a.Reporter.Reported() == 1 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(a.Plugins))
	for _, p := range a.Plugins {
		go func(plugin plugin.Plugin) {
			messages := plugin.OnPanic(ctx, request, err, stackTrace)
			a.Reporter.Collect(messages)
			wg.Done()
		}(p)
	}
	wg.Wait()
	a.Reporter.Report()
	a.Reporter.ClearData()
}

type timeoutError struct{}

func (e timeoutError) Error() string {
	return fmt.Sprintf("Lambda is timed out")
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
