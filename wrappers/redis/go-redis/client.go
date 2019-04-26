package tgoredis

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"

	"github.com/go-redis/redis"
	opentracing "github.com/opentracing/opentracing-go"
	tredis "github.com/thundra-io/thundra-lambda-agent-go/wrappers/redis"
)

// ClientWrapper wraps the *redis.Client to trace redis calls
type ClientWrapper struct {
	*redis.Client

	host string
	port string
	mu   sync.RWMutex
	ctx  context.Context
}

type redisCall struct {
	cw      *ClientWrapper
	command string
	span    opentracing.Span
}

func (rc *redisCall) beforeCall() {
	if rc.cw == nil {
		return
	}
	ctx := rc.cw.Context()
	span, _ := opentracing.StartSpanFromContext(ctx, rc.cw.host)
	rc.span = span
	if rawSpan, ok := tracer.GetRaw(rc.span); ok {
		tredis.BeforeCall(rawSpan, rc.cw.host, rc.cw.port, "", rc.command)
	}
}

func (rc *redisCall) afterCall(err error) {
	if rc.span == nil {
		return
	}
	defer rc.span.Finish()
	if err != nil {
		utils.SetSpanError(rc.span, err)
	}
}

var emptyCtx = context.Background()

// NewClient returns a client to the Redis Server specified
// by Options and wrapped by thundra
func NewClient(opt *redis.Options) *ClientWrapper {
	return WrapClient(redis.NewClient(opt))
}

// WrapClient wraps the given *redis.Clint and returns a new
// *ClientWrapper that can be usedas the redis client and traced by
func WrapClient(c *redis.Client) *ClientWrapper {
	opt := c.Options()
	host, port, err := net.SplitHostPort(opt.Addr)
	if err != nil {
		host = opt.Addr
		port = "6379"
	}
	cw := &ClientWrapper{
		Client: c,
		host:   host,
		port:   port,
		ctx:    emptyCtx,
	}
	c.WrapProcess(processWrapper(cw))
	return cw
}

// WithContext sets the given context to the ClientWrapper,
// so that new spans created using this context to have correct parent
func (cw *ClientWrapper) WithContext(ctx context.Context) *ClientWrapper {
	cw.mu.Lock()
	cw.ctx = ctx
	cw.mu.Unlock()
	return cw
}

// Context returns the current active context of the ClientWrapper
func (cw *ClientWrapper) Context() context.Context {
	cw.mu.RLock()
	ctx := cw.ctx
	cw.mu.RUnlock()
	return ctx
}

func processWrapper(cw *ClientWrapper) func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
	return func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			raw := getCommandString(cmd)
			rc := &redisCall{
				cw:      cw,
				command: raw,
				span:    nil,
			}
			rc.beforeCall()
			err := oldProcess(cmd)
			rc.afterCall(err)
			return err
		}
	}
}

func getCommandString(cmd redis.Cmder) string {
	switch v := cmd.(type) {
	case fmt.Stringer:
		return strings.Split(v.String(), ":")[0]
	case interface{ String() (string, error) }:
		str, err := v.String()
		if err == nil {
			return strings.Split(str, ":")[0]
		}
	}
	args := cmd.Args()
	if len(args) == 0 {
		return ""
	}
	if str, ok := args[0].(string); ok {
		return strings.Split(str, ":")[0]
	}
	return ""
}
