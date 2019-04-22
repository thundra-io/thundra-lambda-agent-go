package thundraredis

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
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
		beforeCall(rawSpan, c.host, c.port, commandName, getRedisCommand(commandName, args...))
	}

	reply, err := c.Conn.Do(commandName, args...)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return reply, err
}

func (c connWrapper) Send(commandName string, args ...interface{}) error {
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
		beforeCall(rawSpan, c.host, c.port, commandName, getRedisCommand(commandName, args...))
	}

	err := c.Conn.Send(commandName, args...)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return err
}

func getRedisCommand(commandName string, args ...interface{}) string {
	var b bytes.Buffer
	b.WriteString(commandName)
	for _, arg := range args {
		b.WriteString(" ")
		switch arg := arg.(type) {
		case string:
			b.WriteString(arg)
		case int:
			b.WriteString(strconv.Itoa(arg))
		case int32:
			b.WriteString(strconv.FormatInt(int64(arg), 10))
		case int64:
			b.WriteString(strconv.FormatInt(arg, 10))
		case fmt.Stringer:
			b.WriteString(arg.String())
		}
	}
	return b.String()
}
