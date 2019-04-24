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

func (c *ClientWrapper) BgRewriteAOF() *redis.StatusCmd {
	rc := c.newRedisCall("bgrewriteaof")
	rc.beforeCall()
	res := c.Client.BgRewriteAOF()
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) BgSave() *redis.StatusCmd {
	rc := c.newRedisCall("bgsave")
	rc.beforeCall()
	res := c.Client.BgSave()
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) BitCount(key string, bitCount *redis.BitCount) *redis.IntCmd {
	rc := c.newRedisCall("bitcount")
	rc.beforeCall()
	res := c.Client.BitCount(key, bitCount)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) BitOpAnd(destKey string, keys ...string) *redis.IntCmd {
	rc := c.newRedisCall("bitopand")
	rc.beforeCall()
	res := c.Client.BitOpAnd(destKey, keys...)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) BitOpOr(destKey string, keys ...string) *redis.IntCmd {
	rc := c.newRedisCall("bitopor")
	rc.beforeCall()
	res := c.Client.BitOpOr(destKey, keys...)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) BitOpXor(destKey string, keys ...string) *redis.IntCmd {
	rc := c.newRedisCall("bitopxor")
	rc.beforeCall()
	res := c.Client.BitOpXor(destKey, keys...)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) BitPos(key string, bit int64, pos ...int64) *redis.IntCmd {
	rc := c.newRedisCall("bitpos")
	rc.beforeCall()
	res := c.Client.BitPos(key, bit, pos...)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) DBSize() *redis.IntCmd {
	rc := c.newRedisCall("dbsize")
	rc.beforeCall()
	res := c.Client.DBSize()
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) Decr(key string) *redis.IntCmd {
	rc := c.newRedisCall("decr")
	rc.beforeCall()
	res := c.Client.Decr(key)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) DecrBy(key string, decrement int64) *redis.IntCmd {
	rc := c.newRedisCall("decrby")
	rc.beforeCall()
	res := c.Client.DecrBy(key, decrement)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) Del(keys ...string) *redis.IntCmd {
	rc := c.newRedisCall("delete")
	rc.beforeCall()
	res := c.Client.Del(key)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) Eval(script string, keys []string, args ...interface{}) *redis.Cmd {
	rc := c.newRedisCall("eval")
	rc.beforeCall()
	res := c.Client.Eval(script, keys, args...)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	rc := c.newRedisCall("evalsha")
	rc.beforeCall()
	res := c.Client.EvalSha(sha1, keys, args...)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) Exists(keys ...string) *redis.IntCmd {
	rc := c.newRedisCall("exists")
	rc.beforeCall()
	res := c.Client.Exists(keys...)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	rc := c.newRedisCall("exists")
	rc.beforeCall()
	res := c.Client.Expire(key, expiration)
	rc.afterCall(res)
	return res
}

func (c *ClientWrapper) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	rc := c.newRedisCall("exists")
	rc.beforeCall()
	res := c.Client.ExpireAt(key, tm)
	rc.afterCall(res)
	return res
}
