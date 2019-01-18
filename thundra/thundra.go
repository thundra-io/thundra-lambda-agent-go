package thundra

import (
	"fmt"

	"github.com/thundra-io/thundra-lambda-agent-go/agent"
	ip "github.com/thundra-io/thundra-lambda-agent-go/invocation"
	mp "github.com/thundra-io/thundra-lambda-agent-go/metric"
	lp "github.com/thundra-io/thundra-lambda-agent-go/thundra_log"
	tp "github.com/thundra-io/thundra-lambda-agent-go/trace"
)

var agentInstance *agent.Agent

func AddDefaultPlugins(a *agent.Agent) *agent.Agent {
	a.AddPlugin(ip.New()).
		AddPlugin(mp.New()).
		AddPlugin(tp.New()).
		AddPlugin(lp.New())

	return a
}

func GetAgent() *agent.Agent {
	return agentInstance
}

func Wrap(handler interface{}) interface{} {
	if agentInstance == nil {
		fmt.Println("thundra.go: agentInstance is nil")
		return nil
	}

	return agentInstance.Wrap(handler)
}

func init() {
	agentInstance = AddDefaultPlugins(agent.New())
}
