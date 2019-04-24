package tgoredis

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"

	"github.com/go-redis/redis"
	opentracing "github.com/opentracing/opentracing-go"
	tredis "github.com/thundra-io/thundra-lambda-agent-go/wrappers/redis"
)

type ClientWrapper struct {
	*redis.Client
	host string
	port string
}

var emptyCtx = context.Background()

// NewClient returns a client to the Redis Server specified
// by Options and wrapped by thundra
func NewClient(opt *redis.Options) *ClientWrapper {
	client := redis.NewClient(opt)
	host, port, _ := net.SplitHostPort(opt.Addr)
	return &ClientWrapper{client, host, port}
}

func (c *ClientWrapper) startSpan() opentracing.Span {
	ctx := c.Client.Context()
	span, _ := opentracing.StartSpanFromContext(ctx, c.host)
	return span
}

func setCommand(span *tracer.RawSpan, args []interface{}) {
	commandName := args[0]
	if commandName != "" {
		args[0] = strings.ToUpper(args[0].(string))
	}
	command := tredis.GetRedisCommand("", args...)
	tredis.AfterCall(span, command)
}

func (c *ClientWrapper) Ping() *redis.StatusCmd {
	span := c.startSpan()
	defer span.Finish()
	if rawSpan, ok := tracer.GetRaw(span); ok {
		commandName := "PING"
		tredis.BeforeCall(rawSpan, c.host, c.port, commandName, commandName)
	}
	res := c.Client.Ping()
	if err := res.Err(); err != nil {
		utils.SetSpanError(span, err)
	}
	return res
}

func (c *ClientWrapper) Append(key, value string) *redis.IntCmd {
	span := c.startSpan()
	defer span.Finish()
	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		commandName := "APPEND"
		tredis.BeforeCall(rawSpan, c.host, c.port, commandName, commandName)
	}
	res := c.Client.Append(key, value)
	setCommand(rawSpan, res.Args())
	if err := res.Err(); err != nil {
		utils.SetSpanError(span, err)
	}
	return res
}

func (c *ClientWrapper) BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	span := c.startSpan()
	defer span.Finish()
	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		commandName := "BLPOP"
		tredis.BeforeCall(rawSpan, c.host, c.port, commandName, commandName)
	}
	res := c.Client.BLPop(timeout, keys...)
	if err := res.Err(); err != nil {
		utils.SetSpanError(span, err)
	}
	setCommand(rawSpan, res.Args())
	return res
}
