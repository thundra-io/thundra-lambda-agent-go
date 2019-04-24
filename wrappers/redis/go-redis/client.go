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

type redisResult interface {
	Args() []interface{}
	Err() error
}
type redisCall struct {
	cw   *ClientWrapper
	cn   string
	span opentracing.Span
}

func (rc *redisCall) beforeCall() {
	if rc.cw == nil {
		return
	}
	ctx := rc.cw.Context()
	span, _ := opentracing.StartSpanFromContext(ctx, rc.cw.host)
	rc.span = span
	if rawSpan, ok := tracer.GetRaw(rc.span); ok {
		tredis.BeforeCall(rawSpan, rc.cw.host, rc.cw.port, rc.cn, rc.cn)
	}
}

func (rc *redisCall) afterCall(res redisResult) {
	defer rc.span.Finish()
	if err := res.Err(); err != nil {
		utils.SetSpanError(rc.span, err)
	}
	if rawSpan, ok := tracer.GetRaw(rc.span); ok {
		args := res.Args()
		commandName := args[0]
		if commandName.(string) != "" {
			args[0] = strings.ToUpper(args[0].(string))
		}
		command := tredis.GetRedisCommand("", args...)
		tredis.AfterCall(rawSpan, command)
	}
}

var emptyCtx = context.Background()

// NewClient returns a client to the Redis Server specified
// by Options and wrapped by thundra
func NewClient(opt *redis.Options) *ClientWrapper {
	client := redis.NewClient(opt)
	host, port, _ := net.SplitHostPort(opt.Addr)
	return &ClientWrapper{client, host, port}
}

func (c *ClientWrapper) newRedisCall(cn string) *redisCall {
	return &redisCall{
		cw: c,
		cn: cn,
	}
}

func (c *ClientWrapper) Ping() *redis.StatusCmd {
	rc := c.newRedisCall("ping")
	rc.beforeCall()
	res := c.Client.Ping()
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) Append(key, value string) *redis.IntCmd {
	rc := c.newRedisCall("append")
	rc.beforeCall()
	res := c.Client.Append(key, value)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	rc := c.newRedisCall("blpop")
	rc.beforeCall()
	res := c.Client.BLPop(timeout, keys...)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	rc := c.newRedisCall("brpop")
	rc.beforeCall()
	res := c.Client.BRPop(timeout, keys...)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) BRPopLPush(source, destination string, timeout time.Duration) *redis.StringCmd {
	rc := c.newRedisCall("brpoplpush")
	rc.beforeCall()
	res := c.Client.BRPopLPush(source, destination, timeout)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) BZPopMax(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	rc := c.newRedisCall("bzpopmax")
	rc.beforeCall()
	res := c.Client.BZPopMax(timeout, keys...)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) BZPopMin(timeout time.Duration, keys ...string) *redis.ZWithKeyCmd {
	rc := c.newRedisCall("bzpopmin")
	rc.beforeCall()
	res := c.Client.BZPopMin(timeout, keys...)
	rc.afterCall(res)
	return res
}
