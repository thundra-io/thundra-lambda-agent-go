package thundrardb

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"

	"github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

// DriverWrapper wraps sql driver.Driver
type DriverWrapper struct {
	Driver driver.Driver
}

// ConnWrapper wraps sql driver.Conn
type ConnWrapper struct {
	ctx         context.Context
	Conn        driver.Conn
	dsn         string
	integration rdbIntegration
}

// StmtWrapper wraps sql driver.Stmt with context
type StmtWrapper struct {
	ctx         context.Context
	Stmt        driver.Stmt
	dsn         string
	integration rdbIntegration
	query       string
}

type rdbIntegration interface {
	beforeCall(query string, span *tracer.RawSpan, dsn string)
	afterCall(query string, span *tracer.RawSpan, dsn string)
	getOperationName(query string) string
}

var emptyCtx = context.Background()

// Wrap wraps the SQL driver
// Returned driver should be registered to be able to use
func Wrap(d driver.Driver) driver.Driver {
	return &DriverWrapper{d}
}

var operationToType = map[string]string{
	"select": "READ",
	"insert": "WRITE",
	"update": "WRITE",
	"delete": "DELETE",
}

// Open opens a new connection and wraps connection for mysql and postgresql
// Returns connection as it is for other drivers
func (d *DriverWrapper) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	switch connType := reflect.TypeOf(conn).String(); connType {
	case "*mysql.mysqlConn":
		return &ConnWrapper{Conn: conn, dsn: name, integration: &mysqlIntegration{}}, nil

	case "*pq.conn":
		return &ConnWrapper{Conn: conn, dsn: name, integration: &postgresqlIntegration{}}, nil
	}
	return conn, nil
}

// Prepare creates a prepared statement and wraps it
func (c *ConnWrapper) Prepare(query string) (driver.Stmt, error) {
	stmt, err := c.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}

	return StmtWrapper{Stmt: stmt, ctx: emptyCtx, query: query, dsn: c.dsn, integration: c.integration}, nil
}

// PrepareContext creates a prepared statement and wraps it
func (c *ConnWrapper) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	if connPrepareCtx, ok := c.Conn.(driver.ConnPrepareContext); ok {
		stmt, err := connPrepareCtx.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}

		return StmtWrapper{ctx: ctx, Stmt: stmt, query: query, dsn: c.dsn, integration: c.integration}, nil
	}

	return c.Prepare(query)
}

func (c *ConnWrapper) Close() error {
	return c.Conn.Close()
}

func (c *ConnWrapper) Begin() (driver.Tx, error) {
	return c.Conn.Begin()
}

func (c *ConnWrapper) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	if connBeginTx, ok := c.Conn.(driver.ConnBeginTx); ok {
		return connBeginTx.BeginTx(ctx, opts)
	}

	return c.Conn.Begin()
}

func (c *ConnWrapper) Exec(query string, args []driver.Value) (driver.Result, error) {
	span, _ := opentracing.StartSpanFromContext(
		emptyCtx,
		c.integration.getOperationName(query),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		c.integration.beforeCall(query, rawSpan, c.dsn)
	}

	if execer, ok := c.Conn.(driver.Execer); ok {
		return execer.Exec(query, args)
	}

	return nil, driver.ErrSkip
}

func (c *ConnWrapper) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (r driver.Result, err error) {
	span, ctxWithSpan := opentracing.StartSpanFromContext(
		ctx,
		c.integration.getOperationName(query),
	)

	if err != nil {
		fmt.Println(err)
	}
	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		c.integration.beforeCall(query, rawSpan, c.dsn)
	}

	if execContext, ok := c.Conn.(driver.ExecerContext); ok {
		return execContext.ExecContext(ctxWithSpan, query, args)
	}

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return c.Conn.(driver.Execer).Exec(query, dargs)
}

func (c *ConnWrapper) Ping(ctx context.Context) (err error) {
	if pinger, ok := c.Conn.(driver.Pinger); ok {
		return pinger.Ping(ctx)
	}
	return nil
}

func (c *ConnWrapper) Query(query string, args []driver.Value) (driver.Rows, error) {
	span, _ := opentracing.StartSpanFromContext(
		emptyCtx,
		c.integration.getOperationName(query),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		c.integration.beforeCall(query, rawSpan, c.dsn)
	}
	if queryer, ok := c.Conn.(driver.Queryer); ok {
		return queryer.Query(query, args)
	}

	return nil, driver.ErrSkip
}

func (c *ConnWrapper) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	span, ctxWithSpan := opentracing.StartSpanFromContext(
		ctx,
		c.integration.getOperationName(query),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		c.integration.beforeCall(query, rawSpan, c.dsn)
	}

	if queryerContext, ok := c.Conn.(driver.QueryerContext); ok {
		return queryerContext.QueryContext(ctxWithSpan, query, args)
	}

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return c.Query(query, dargs)
}

func (s StmtWrapper) Close() (err error) {
	return s.Stmt.Close()
}

func (s StmtWrapper) NumInput() int {
	return s.Stmt.NumInput()
}

func (s StmtWrapper) Exec(args []driver.Value) (res driver.Result, err error) {
	span, _ := opentracing.StartSpanFromContext(
		s.ctx,
		s.integration.getOperationName(s.query),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		s.integration.beforeCall(s.query, rawSpan, s.dsn)
	}

	return s.Stmt.Exec(args)

}

func (s StmtWrapper) ExecContext(ctx context.Context, args []driver.NamedValue) (res driver.Result, err error) {
	span, ctxWithSpan := opentracing.StartSpanFromContext(
		ctx,
		s.integration.getOperationName(s.query),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		s.integration.beforeCall(s.query, rawSpan, s.dsn)
	}

	if stmtExecContext, ok := s.Stmt.(driver.StmtExecContext); ok {
		return stmtExecContext.ExecContext(ctxWithSpan, args)
	}

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return s.Stmt.Exec(dargs)
}

func (s StmtWrapper) Query(args []driver.Value) (rows driver.Rows, err error) {
	span, _ := opentracing.StartSpanFromContext(
		s.ctx,
		s.integration.getOperationName(s.query),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		s.integration.beforeCall(s.query, rawSpan, s.dsn)
	}

	return s.Stmt.Query(args)
}

func (s StmtWrapper) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	span, ctxWithSpan := opentracing.StartSpanFromContext(
		ctx,
		s.integration.getOperationName(s.query),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		s.integration.beforeCall(s.query, rawSpan, s.dsn)
	}

	if stmtQueryContext, ok := s.Stmt.(driver.StmtQueryContext); ok {
		return stmtQueryContext.QueryContext(ctxWithSpan, args)
	}

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return s.Query(dargs)
}

func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}
