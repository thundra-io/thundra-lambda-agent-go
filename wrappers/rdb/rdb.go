package thundrardb

import (
	"context"
	"database/sql/driver"
	"errors"
	"reflect"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
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
	getOperationName(dsn string) string
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

// Close closes connection in wrapper
func (c *ConnWrapper) Close() error {
	return c.Conn.Close()
}

// Begin starts and returns a new transaction
func (c *ConnWrapper) Begin() (driver.Tx, error) {
	return c.Conn.Begin()
}

// BeginTx starts and returns a new transaction with context and and TxOptions
func (c *ConnWrapper) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	if connBeginTx, ok := c.Conn.(driver.ConnBeginTx); ok {
		return connBeginTx.BeginTx(ctx, opts)
	}

	return c.Conn.Begin()
}

// Exec wraps the driver.execer.Exec and starts a new span.
func (c *ConnWrapper) Exec(query string, args []driver.Value) (res driver.Result, err error) {
	span, _ := opentracing.StartSpanFromContext(
		emptyCtx,
		c.integration.getOperationName(c.dsn),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		c.integration.beforeCall(query, rawSpan, c.dsn)
	}
	tracer.OnSpanStarted(span)

	if execer, ok := c.Conn.(driver.Execer); ok {
		res, err = execer.Exec(query, args)
		if err != nil {
			utils.SetSpanError(span, err)
		}
		return
	}

	return nil, driver.ErrSkip
}

// ExecContext wraps the driver.ExecerContext.ExecContext and starts a new span.
// The newly created span will be a child of the span
// whose context is is passed using the ctx parameter
func (c *ConnWrapper) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (res driver.Result, err error) {
	span, ctxWithSpan := opentracing.StartSpanFromContext(
		ctx,
		c.integration.getOperationName(c.dsn),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		c.integration.beforeCall(query, rawSpan, c.dsn)
	}
	tracer.OnSpanStarted(span)

	if execContext, ok := c.Conn.(driver.ExecerContext); ok {
		res, err = execContext.ExecContext(ctxWithSpan, query, args)
		if err != nil {
			utils.SetSpanError(span, err)
		}
		return
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

	res, err = c.Conn.(driver.Execer).Exec(query, dargs)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return
}

// Ping wraps driver.Pinger.Ping
func (c *ConnWrapper) Ping(ctx context.Context) (err error) {
	if pinger, ok := c.Conn.(driver.Pinger); ok {
		return pinger.Ping(ctx)
	}
	return nil
}

// Query wraps the driver.Queryer.Query and starts a new span.
func (c *ConnWrapper) Query(query string, args []driver.Value) (rows driver.Rows, err error) {
	span, _ := opentracing.StartSpanFromContext(
		emptyCtx,
		c.integration.getOperationName(c.dsn),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		c.integration.beforeCall(query, rawSpan, c.dsn)
	}
	tracer.OnSpanStarted(span)

	if queryer, ok := c.Conn.(driver.Queryer); ok {
		rows, err = queryer.Query(query, args)
		if err != nil {
			utils.SetSpanError(span, err)
		}
		return
	}

	return nil, driver.ErrSkip
}

// QueryContext wraps the driver.Queryer.Query and starts a new span.
// The newly created span will be a child of the span
// whose context is is passed using the ctx parameter
func (c *ConnWrapper) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	span, ctxWithSpan := opentracing.StartSpanFromContext(
		ctx,
		c.integration.getOperationName(c.dsn),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		c.integration.beforeCall(query, rawSpan, c.dsn)
	}
	tracer.OnSpanStarted(span)

	if queryerContext, ok := c.Conn.(driver.QueryerContext); ok {
		res, err := queryerContext.QueryContext(ctxWithSpan, query, args)
		if err != nil {
			utils.SetSpanError(span, err)
		}
		return res, err
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

	res, err := c.Query(query, dargs)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return res, err
}

// Close wraps driver.Stmt.Close
func (s StmtWrapper) Close() (err error) {
	return s.Stmt.Close()
}

// NumInput wraps driver.Stmt.NumInput
func (s StmtWrapper) NumInput() int {
	return s.Stmt.NumInput()
}

// Exec wraps the driver.Stmt.Exec and starts a new span.
func (s StmtWrapper) Exec(args []driver.Value) (res driver.Result, err error) {
	span, _ := opentracing.StartSpanFromContext(
		s.ctx,
		s.integration.getOperationName(s.dsn),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		s.integration.beforeCall(s.query, rawSpan, s.dsn)
	}
	tracer.OnSpanStarted(span)

	res, err = s.Stmt.Exec(args)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return

}

// ExecContext wraps the driver.StmtExecContext.ExecContext and starts a new span.
// The newly created span will be a child of the span
// whose context is is passed using the ctx parameter
func (s StmtWrapper) ExecContext(ctx context.Context, args []driver.NamedValue) (res driver.Result, err error) {
	span, ctxWithSpan := opentracing.StartSpanFromContext(
		ctx,
		s.integration.getOperationName(s.dsn),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		s.integration.beforeCall(s.query, rawSpan, s.dsn)
	}
	tracer.OnSpanStarted(span)

	if stmtExecContext, ok := s.Stmt.(driver.StmtExecContext); ok {
		res, err = stmtExecContext.ExecContext(ctxWithSpan, args)
		if err != nil {
			utils.SetSpanError(span, err)
		}
		return
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

	res, err = s.Stmt.Exec(dargs)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return
}

// Query wraps the driver.Stmt.Query and starts a new span.
func (s StmtWrapper) Query(args []driver.Value) (rows driver.Rows, err error) {
	span, _ := opentracing.StartSpanFromContext(
		s.ctx,
		s.integration.getOperationName(s.dsn),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		s.integration.beforeCall(s.query, rawSpan, s.dsn)
	}
	tracer.OnSpanStarted(span)

	rows, err = s.Stmt.Query(args)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return
}

// QueryContext wraps the driver.StmtQueryContext.QueryContext and starts a new span.
// The newly created span will be a child of the span
// whose context is is passed using the ctx parameter
func (s StmtWrapper) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	span, ctxWithSpan := opentracing.StartSpanFromContext(
		ctx,
		s.integration.getOperationName(s.dsn),
	)

	defer span.Finish()

	rawSpan, ok := tracer.GetRaw(span)
	if ok {
		s.integration.beforeCall(s.query, rawSpan, s.dsn)
	}
	tracer.OnSpanStarted(span)

	if stmtQueryContext, ok := s.Stmt.(driver.StmtQueryContext); ok {
		rows, err = stmtQueryContext.QueryContext(ctxWithSpan, args)
		return
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

	rows, err = s.Query(dargs)
	if err != nil {
		utils.SetSpanError(span, err)
	}
	return
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
