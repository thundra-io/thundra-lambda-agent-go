package thundra

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type Builder interface {
	AddPlugin(plugin.Plugin) Builder
	SetReporter(Reporter) Builder
	SetAPIKey(string) Builder
	EnableWarmup() Builder
	Build() *thundra
}

type builder struct {
	plugins  []plugin.Plugin
	reporter Reporter
	apiKey   string
	warmup   bool
}

// AddPlugin is used to enable plugins on thundra. Trace and Metrics are two plugins that are supported currently.
// You need to initialize a plugin object and pass it as a parameter in order to enable it.
func (b *builder) AddPlugin(plugin plugin.Plugin) Builder {
	b.plugins = append(b.plugins, plugin)
	return b
}

// SetReporter is used
func (b *builder) SetReporter(reporter Reporter) Builder {
	b.reporter = reporter
	return b
}

// SetAPIKey is used to set ApiKey to use Thundra. See https://docs.thundra.io/docs/api-keys to learn how you can
// generate your own api key.
func (b *builder) SetAPIKey(apiKey string) Builder {
	b.apiKey = apiKey
	return b
}

// EnableWarmup enables warming up to reduce cold starts of your lambda. See https://docs.thundra.io/docs/how-to-warmup
// to learn how you can configure thundra-lambda-warmup.
func (b *builder) EnableWarmup() Builder {
	b.warmup = true
	return b
}

// Build returns the thundra object that you will pass to thundra.Wrap() function.
func (b *builder) Build() *thundra {
	if b.reporter == nil {
		b.reporter = &reporterImpl{}
	}

	k := determineApiKey(b.apiKey)
	w := determineWarmup(b.warmup)

	return &thundra{
		plugins:  b.plugins,
		reporter: b.reporter,
		apiKey:   k,
		warmup:   w,
	}
}

// NewBuilder returns a new thundra builder.
func NewBuilder() Builder {
	return &builder{}
}

// determineApiKey determines which apiKey to use. if apiKey is set from environment variable, returns that value.
// Otherwise returns the value from builder's setApiKey method. Panic if it's not set by neither.
func determineApiKey(builderApiKey string) string {
	k, err := checkApiKey()
	if err != nil {
		if builderApiKey == "" {
			panic(err)
		}
		k = builderApiKey
	}
	return k
}

// checkApiKey is used to fetch the apiKey value from environment variable
func checkApiKey() (string, error) {
	k := os.Getenv(thundraApiKey)
	if k == "" {
		return "", errors.New("thundraApiKey is not set")
	}
	return k, nil
}

// determineWarmup determines which warmup value to use. if warmup is set from environment variable, returns that value.
// Otherwise returns true if it's enabled by builder's enableWarmup method. Default value is false.
func determineWarmup(builderWarmup bool) bool {
	w, err := checkWarmup()
	if err != nil {
		w = builderWarmup
	}
	return w
}

// checkWarmup fetches the warmup value from environment variable
func checkWarmup() (bool, error) {
	w := os.Getenv(thundraLambdaWarmupWarmupAware);
	b, err := strconv.ParseBool(w)
	if err != nil {
		fmt.Println(err, " thundra_lambda_warmup_warmupAware should be set with a boolean.")
		return false, err
	}
	return b, nil
}
