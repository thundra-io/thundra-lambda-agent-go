package tgoredis

import (
	"bytes"
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

// Pipeliner is used to trace pipelines executed on a Redis server
type Pipeliner struct {
	redis.Pipeliner

	host string
	port string
	ctx  context.Context
	cw   *ClientWrapper
}

type redisCall struct {
	cw       *ClientWrapper
	command  string
	span     opentracing.Span
	pipeline bool
	err      error
}

func (rc *redisCall) beforeCall() {
	if rc.cw == nil {
		return
	}
	ctx := rc.cw.Context()
	span, _ := opentracing.StartSpanFromContext(ctx, rc.cw.host)
	rc.span = span
	if rawSpan, ok := tracer.GetRaw(rc.span); ok {
		commandName := ""
		if rc.pipeline {
			commandName = "pipeline"
		}
		tredis.BeforeCall(rawSpan, rc.cw.host, rc.cw.port, commandName, rc.command)
	}
}

func (rc *redisCall) afterCall() {
	if rc.span == nil {
		return
	}
	defer rc.span.Finish()
	if rc.err != nil {
		utils.SetSpanError(rc.span, rc.err)
	}
	if rawSpan, ok := tracer.GetRaw(rc.span); ok {
		tredis.AfterCall(rawSpan, rc.command)
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

// Pipeline creates a Pipeline from a ClientWrapper
func (cw *ClientWrapper) Pipeline() redis.Pipeliner {
	cw.mu.RLock()
	ctx := cw.ctx
	cw.mu.RUnlock()
	return &Pipeliner{
		Pipeliner: cw.Client.Pipeline(),
		host:      cw.host,
		port:      cw.port,
		cw:        cw,
		ctx:       ctx,
	}
}

// ExecWithContext calls Pipeline.Exec(). It ensures that the resulting Redis calls
// are traced, and that emitted spans are children of the given Context.
func (p *Pipeliner) ExecWithContext(ctx context.Context) ([]redis.Cmder, error) {
	return p.execWithContext(ctx)
}

// Exec calls Pipeline.Exec() ensuring that the resulting Redis calls are traced.
func (p *Pipeliner) Exec() ([]redis.Cmder, error) {
	return p.execWithContext(p.ctx)
}

func (p *Pipeliner) execWithContext(ctx context.Context) ([]redis.Cmder, error) {
	rc := &redisCall{
		cw: p.cw,
	}
	rc.beforeCall()
	cmds, err := p.Pipeliner.Exec()
	rc.command = multipleCommandString(cmds)
	rc.afterCall()
	return cmds, err
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
			}
			rc.beforeCall()
			err := oldProcess(cmd)
			rc.err = err
			rc.afterCall()
			return err
		}
	}
}

func multipleCommandString(cmds []redis.Cmder) string {
	var b bytes.Buffer
	for _, cmd := range cmds {
		b.WriteString(getCommandString(cmd))
		b.WriteString("\n")
	}
	return b.String()
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
