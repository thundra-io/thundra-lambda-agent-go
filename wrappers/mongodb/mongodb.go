package thundramongo

import (
	"context"
	"strings"
	"sync"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/application"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/config"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/utils"
	"go.mongodb.org/mongo-driver/event"
)

type commandMonitor struct {
	sync.Mutex
	spans map[spanKey]opentracing.Span
}

type spanKey struct {
	ConnectionID string
	RequestID    int64
}

// NewCommandMonitor returns a new event.CommandMonitor for tracing commands
func NewCommandMonitor() *event.CommandMonitor {
	cm := commandMonitor{spans: make(map[spanKey]opentracing.Span)}
	return &event.CommandMonitor{Started: cm.started, Succeeded: cm.succeeded, Failed: cm.failed}
}

func (c *commandMonitor) started(ctx context.Context, event *event.CommandStartedEvent) {
	span, _ := opentracing.StartSpanFromContext(ctx, event.DatabaseName)
	rawSpan, ok := tracer.GetRaw(span)
	if !ok {
		return
	}

	// Store span to use it on command finished
	c.Lock()
	c.spans[spanKey{event.ConnectionID, event.RequestID}] = span
	c.Unlock()

	beforeCall(rawSpan, event)
	tracer.OnSpanStarted(span)
}

func (c *commandMonitor) succeeded(ctx context.Context, event *event.CommandSucceededEvent) {
	c.finished(ctx, &event.CommandFinishedEvent, "")
}

func (c *commandMonitor) failed(ctx context.Context, event *event.CommandFailedEvent) {
	c.finished(ctx, &event.CommandFinishedEvent, event.Failure)
}

func (c *commandMonitor) finished(ctx context.Context, event *event.CommandFinishedEvent, failure string) {
	key := spanKey{event.ConnectionID, event.RequestID}

	c.Lock()
	// Retrieve span set in command started
	span, ok := c.spans[key]
	if ok {
		delete(c.spans, key)
	}
	c.Unlock()
	if !ok {
		return
	}

	if failure != "" {
		utils.SetSpanError(span, failure)
	}
	span.Finish()
}

func beforeCall(span *tracer.RawSpan, event *event.CommandStartedEvent) {
	span.ClassName = constants.ClassNames["MONGODB"]
	span.DomainName = constants.DomainNames["DB"]

	host, port := "", "27017"
	if len(event.ConnectionID) > 0 {
		host = strings.Split(event.ConnectionID, "[")[0]

		if len(strings.Split(host, ":")) > 1 {
			port = strings.Split(host, ":")[1]
			host = strings.Split(host, ":")[0]
		}
	}

	collectionValue := event.Command.Lookup(event.CommandName)
	collectionName, _ := collectionValue.StringValueOK()

	// Set span tags
	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:          constants.MongoDBCommandTypes[strings.ToUpper(event.CommandName)],
		constants.DBTags["DB_TYPE"]:                   "mongodb",
		constants.DBTags["DB_HOST"]:                   host,
		constants.DBTags["DB_PORT"]:                   port,
		constants.DBTags["DB_INSTANCE"]:               event.DatabaseName,
		constants.MongoDBTags["MONGODB_COMMAND_NAME"]: strings.ToUpper(event.CommandName),
		constants.MongoDBTags["MONGODB_COLLECTION"]:   collectionName,
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
	}

	if !config.MaskMongoDBCommand {
		if event.Command != nil {
			command := event.Command.String()
			size := len(command)
			if size > constants.DefaultMongoDBSizeLimit {
				size = constants.DefaultMongoDBSizeLimit
			}

			tags[constants.MongoDBTags["MONGODB_COMMAND"]] = event.Command.String()[:size]
		} else {
			tags[constants.MongoDBTags["MONGODB_COMMAND"]] = ""
		}
	}

	span.Tags = tags
}
