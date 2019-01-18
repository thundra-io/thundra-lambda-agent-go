package thundra

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	ip "github.com/thundra-io/thundra-lambda-agent-go/invocation"
	mp "github.com/thundra-io/thundra-lambda-agent-go/metric"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
	lp "github.com/thundra-io/thundra-lambda-agent-go/thundra_log"
	tp "github.com/thundra-io/thundra-lambda-agent-go/trace"
)

var agent *thundra

type thundra struct {
	plugins       []plugin.Plugin
	reporter      reporter
	warmup        bool
	timeoutMargin time.Duration
}

// New is used to collect basic invocation data with thundra. Use NewBuilder and AddPlugin to access full functionality.
func New() *thundra {
	w := determineWarmup()
	g := determineTimeoutMargin()
	determineAPIKey()

	t := &thundra{
		reporter: &reporterImpl{
			reported: new(uint32),
		},
		warmup:        w,
		timeoutMargin: g,
		plugins:       []plugin.Plugin{},
	}

	return t
}

func (t *thundra) AddDefaultPlugins() *thundra {
	t.AddPlugin(ip.New()).
		AddPlugin(mp.New()).
		AddPlugin(tp.New()).
		AddPlugin(lp.New())

	return t
}

// AddPlugin is used to enable plugins on thundra. You can use Trace, Metrics and Log plugins.
// You need to initialize a plugin object and pass it as a parameter in order to enable it.
// e.g. AddPlugin(trace.New())
func (t *thundra) AddPlugin(plugin plugin.Plugin) *thundra {
	if plugin.IsEnabled() {
		t.plugins = append(t.plugins, plugin)
	}

	return t
}

func (t *thundra) SetReporter(r reporter) *thundra {
	t.reporter = r
	return t
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

func (t *thundra) executePreHooks(ctx context.Context, request json.RawMessage) {
	t.reporter.FlushFlag()
	plugin.TraceID = plugin.GenerateNewID()
	plugin.TransactionID = plugin.GenerateNewID()
	var wg sync.WaitGroup
	wg.Add(len(t.plugins))
	for _, p := range t.plugins {
		go p.BeforeExecution(ctx, request, &wg)
	}
	wg.Wait()
}

func (t *thundra) executePostHooks(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) {
	// Skip if it is already reported
	if *t.reporter.Reported() == 1 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(t.plugins))
	for _, p := range t.plugins {
		go func(plugin plugin.Plugin) {
			messages := plugin.AfterExecution(ctx, request, response, err)
			t.reporter.Collect(messages)
			wg.Done()
		}(p)
	}
	wg.Wait()
	t.reporter.Report()
	t.reporter.ClearData()
}

func (t *thundra) onPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) {
	// Skip if it is already reported
	if *t.reporter.Reported() == 1 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(t.plugins))
	for _, p := range t.plugins {
		go func(plugin plugin.Plugin) {
			messages := plugin.OnPanic(ctx, request, err, stackTrace)
			t.reporter.Collect(messages)
			wg.Done()
		}(p)
	}
	wg.Wait()
	t.reporter.Report()
	t.reporter.ClearData()
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

	timer := time.NewTimer(time.Until(timeoutDuration))
	timeoutChannel := timer.C

	select {
	case <-timeoutChannel:
		fmt.Println("Function is timed out")
		t.executePostHooks(ctx, payload, nil, timeoutError{})
		return
	case <-ctx.Done():
		// close timeoutChannel
		timer.Stop()
		return
	}
}

// GetAgent returns agent instance
func GetAgent() *thundra {
	return agent
}

func init() {
	agent = New().AddDefaultPlugins()
}
