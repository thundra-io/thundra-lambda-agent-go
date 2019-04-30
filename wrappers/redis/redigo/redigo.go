package tredigo

import (
	"context"
	"net"
	"net/url"

	"github.com/gomodule/redigo/redis"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
	tredis "github.com/thundra-io/thundra-lambda-agent-go/wrappers/redis"
)

type connWrapper struct {
	redis.Conn
	host string
	port string
}

var emptyCtx = context.Background()

// Dial wraps redis.Dial and returns a wrapped connection
func Dial(network, address string, options ...redis.DialOption) (redis.Conn, error) {
	conn, err := redis.Dial(network, address, options...)
	if err != nil {
		return nil, err
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	return connWrapper{conn, host, port}, nil
}

// DialURL wraps redis.DialURL and returns a wrapped connection
func DialURL(rawurl string, options ...redis.DialOption) (redis.Conn, error) {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	host, port, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		host = parsedURL.Host
		port = "6379"
	}
	if host == "" {
		host = "localhost"
	}
	conn, err := redis.DialURL(rawurl, options...)
	return connWrapper{conn, host, port}, err
}

// Do wraps the redis.Conn.Do and starts a new span. If context.Context is provided as last argument,
// the newly created span will be a child of span with the passed context. Otherwise, new span will be
// created with an empty context.
func (c connWrapper) Do(commandName string, args ...interface{}) (interface{}, error) {
	ctx := emptyCtx
	if n := len(args); n > 0 {
		var ok bool
		if ctx, ok = args[n-1].(context.Context); ok {
			args = args[:n-1]
		} else {
			ctx = emptyCtx
		}
	}

	span, _ := opentracing.StartSpanFromContext(
		ctx,
		c.host,
	)
	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		tredis.BeforeCall(rawSpan, c.host, c.port, commandName, tredis.GetRedisCommand(commandName, args...))
	}

	reply, err := c.Conn.Do(commandName, args...)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return reply, err
}