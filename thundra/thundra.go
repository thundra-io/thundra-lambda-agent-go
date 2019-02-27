package thundra

import (
	"fmt"

	"github.com/thundra-io/thundra-lambda-agent-go/agent"
	ip "github.com/thundra-io/thundra-lambda-agent-go/invocation"
	mp "github.com/thundra-io/thundra-lambda-agent-go/metric"
	lp "github.com/thundra-io/thundra-lambda-agent-go/log"
	tp "github.com/thundra-io/thundra-lambda-agent-go/trace"
)

var agentInstance *agent.Agent

// Logger is main thundra logger
var Logger = lp.Logger

func addDefaultPlugins(a *agent.Agent) *agent.Agent {
	a.AddPlugin(ip.New()).
		AddPlugin(mp.New()).
		AddPlugin(tp.New()).
		AddPlugin(lp.New())

	return a
}

// Wrap wraps the given handler function so that the 
// thundra agent integrates with given handler
func Wrap(handler interface{}) interface{} {
	if agentInstance == nil {
		fmt.Println("thundra.go: agentInstance is nil")
		return handler
	}

	return agentInstance.Wrap(handler)
}

func init() {
	agentInstance = addDefaultPlugins(agent.New())
}
