package thundra

import (
	"log"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/agent"
	ip "github.com/thundra-io/thundra-lambda-agent-go/v2/invocation"
	lp "github.com/thundra-io/thundra-lambda-agent-go/v2/log"
	mp "github.com/thundra-io/thundra-lambda-agent-go/v2/metric"
	tp "github.com/thundra-io/thundra-lambda-agent-go/v2/trace"
)

var agentInstance *agent.Agent

// Logger is main thundra logger
var Logger = lp.Logger

func addDefaultPlugins(a *agent.Agent) *agent.Agent {
	a.AddPlugin(ip.New()).
		AddPlugin(mp.New()).
		AddPlugin(tp.GetInstance()).
		AddPlugin(lp.New())

	return a
}

// Wrap wraps the given handler function so that the
// thundra agent integrates with given handler
func Wrap(handler interface{}) interface{} {
	if agentInstance == nil {
		log.Println("thundra.go: agentInstance is nil")
		return handler
	}

	return agentInstance.Wrap(handler)
}

func init() {
	agentInstance = addDefaultPlugins(agent.New())
}
