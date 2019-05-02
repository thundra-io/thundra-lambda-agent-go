package thundrardb

import (
	"context"
	"database/sql"
	"strings"

	"database/sql/driver"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/trace"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type IntegrationSuite struct {
	suite.Suite
	db *sql.DB
}

func newSuite(t *testing.T, driver driver.Driver, dsn string, driverName string) *IntegrationSuite {
	sql.Register("test"+driverName, Wrap(driver))
	db, err := sql.Open("test"+driverName, dsn)
	if err != nil {
		t.Errorf("Could not open connection")
	}

	return &IntegrationSuite{db: db}
}

func (s *IntegrationSuite) TestRdbIntegration(t *testing.T, query string, driverName string, args ...interface{}) {

	t.Run("testDbQuery", func(t *testing.T) {
		tp := trace.New()
		// Actual call
		res, err := s.db.Query(query, args...)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		span := tp.Recorder.GetSpans()[0]
		assert.Equal(t, constants.ClassNames[driverName], span.ClassName)
		assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
		assert.Equal(t, operationToType[strings.ToLower(strings.Split(query, " ")[0])], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
		assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
		assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
		assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		assert.Equal(t, strings.ToLower(driverName), span.Tags[constants.DBTags["DB_TYPE"]])
		assert.Equal(t, query, span.Tags[constants.DBTags["DB_STATEMENT"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		// Clear tracer
		tp.Reset()
	})

	t.Run("testDbQueryContext", func(t *testing.T) {
		tp := trace.New()
		// Create the parent span
		ctx := context.Background()
		parentSpan, ctx := opentracing.StartSpanFromContext(ctx, "parentSpan")
		parentSpanRaw, _ := tracer.GetRaw(parentSpan)
		// Actual call
		res, err := s.db.QueryContext(ctx, query, args...)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		// Get the span created for db query call
		span := tp.Recorder.GetSpans()[1]
		// Check parent span is set
		assert.Equal(t, parentSpanRaw.Context.SpanID, span.ParentSpanID)

		assert.Equal(t, constants.ClassNames[driverName], span.ClassName)
		assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
		assert.Equal(t, operationToType[strings.ToLower(strings.Split(query, " ")[0])], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
		assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
		assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
		assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		assert.Equal(t, strings.ToLower(driverName), span.Tags[constants.DBTags["DB_TYPE"]])
		assert.Equal(t, query, span.Tags[constants.DBTags["DB_STATEMENT"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		tp.Reset()
	})

	t.Run("testDbExec", func(t *testing.T) {
		tp := trace.New()
		// Actual call
		res, err := s.db.Exec(query, args...)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		span := tp.Recorder.GetSpans()[0]
		assert.Equal(t, constants.ClassNames[driverName], span.ClassName)
		assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
		assert.Equal(t, operationToType[strings.ToLower(strings.Split(query, " ")[0])], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
		assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
		assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
		assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		assert.Equal(t, strings.ToLower(driverName), span.Tags[constants.DBTags["DB_TYPE"]])
		assert.Equal(t, query, span.Tags[constants.DBTags["DB_STATEMENT"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		// Clear tracer
		tp.Reset()
	})

	t.Run("testDbExecContext", func(t *testing.T) {
		tp := trace.New()
		// Create the parent span
		ctx := context.Background()
		parentSpan, ctx := opentracing.StartSpanFromContext(ctx, "parentSpan")
		parentSpanRaw, _ := tracer.GetRaw(parentSpan)
		// Actual call
		res, err := s.db.ExecContext(ctx, query, args...)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		// Get the span created for db query call
		span := tp.Recorder.GetSpans()[1]
		// Check parent span is set
		assert.Equal(t, parentSpanRaw.Context.SpanID, span.ParentSpanID)

		assert.Equal(t, constants.ClassNames[driverName], span.ClassName)
		assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
		assert.Equal(t, operationToType[strings.ToLower(strings.Split(query, " ")[0])], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
		assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
		assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
		assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		assert.Equal(t, strings.ToLower(driverName), span.Tags[constants.DBTags["DB_TYPE"]])
		assert.Equal(t, query, span.Tags[constants.DBTags["DB_STATEMENT"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		tp.Reset()
	})

	t.Run("testDbExecContext", func(t *testing.T) {
		tp := trace.New()
		// Create the parent span
		ctx := context.Background()
		parentSpan, ctx := opentracing.StartSpanFromContext(ctx, "parentSpan")
		parentSpanRaw, _ := tracer.GetRaw(parentSpan)
		// Actual call
		res, err := s.db.ExecContext(ctx, query, args...)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		// Get the span created for db query call
		span := tp.Recorder.GetSpans()[1]
		// Check parent span is set
		assert.Equal(t, parentSpanRaw.Context.SpanID, span.ParentSpanID)

		assert.Equal(t, constants.ClassNames[driverName], span.ClassName)
		assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
		assert.Equal(t, operationToType[strings.ToLower(strings.Split(query, " ")[0])], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
		assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
		assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
		assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		assert.Equal(t, strings.ToLower(driverName), span.Tags[constants.DBTags["DB_TYPE"]])
		assert.Equal(t, query, span.Tags[constants.DBTags["DB_STATEMENT"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		tp.Reset()
	})

	t.Run("testStatementQuery", func(t *testing.T) {
		tp := trace.New()
		stmt, err := s.db.Prepare(query)
		assert.Nil(t, err)
		assert.NotNil(t, stmt)

		// Check prepare event is not recorded
		assert.Empty(t, tp.Recorder.GetSpans())

		stmt.Query(args...)

		// Get the span created for statement query call
		span := tp.Recorder.GetSpans()[0]

		assert.Equal(t, constants.ClassNames[driverName], span.ClassName)
		assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
		assert.Equal(t, operationToType[strings.ToLower(strings.Split(query, " ")[0])], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
		assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
		assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
		assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		assert.Equal(t, strings.ToLower(driverName), span.Tags[constants.DBTags["DB_TYPE"]])
		assert.Equal(t, query, span.Tags[constants.DBTags["DB_STATEMENT"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		tp.Reset()
	})

	t.Run("testStatementQueryContext", func(t *testing.T) {
		tp := trace.New()

		// Create the parent span
		ctx := context.Background()
		parentSpan, ctx := opentracing.StartSpanFromContext(ctx, "parentSpan")
		parentSpanRaw, _ := tracer.GetRaw(parentSpan)

		stmt, err := s.db.PrepareContext(ctx, query)
		assert.Nil(t, err)
		assert.NotNil(t, stmt)

		// Check prepare event is not recorded
		assert.Len(t, tp.Recorder.GetSpans(), 1)

		stmt.QueryContext(ctx, args...)

		// Get the span created for statement query call
		span := tp.Recorder.GetSpans()[1]

		// Check parent span is set
		assert.Equal(t, parentSpanRaw.Context.SpanID, span.ParentSpanID)

		assert.Equal(t, constants.ClassNames[driverName], span.ClassName)
		assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
		assert.Equal(t, operationToType[strings.ToLower(strings.Split(query, " ")[0])], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
		assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
		assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
		assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		assert.Equal(t, strings.ToLower(driverName), span.Tags[constants.DBTags["DB_TYPE"]])
		assert.Equal(t, query, span.Tags[constants.DBTags["DB_STATEMENT"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		tp.Reset()
	})

	t.Run("testStatementExec", func(t *testing.T) {
		tp := trace.New()
		stmt, err := s.db.Prepare(query)
		assert.Nil(t, err)
		assert.NotNil(t, stmt)

		// Check prepare event is not recorded
		assert.Empty(t, tp.Recorder.GetSpans())

		stmt.Exec(args...)

		// Get the span created for statement query call
		span := tp.Recorder.GetSpans()[0]

		assert.Equal(t, constants.ClassNames[driverName], span.ClassName)
		assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
		assert.Equal(t, operationToType[strings.ToLower(strings.Split(query, " ")[0])], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
		assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
		assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
		assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		assert.Equal(t, strings.ToLower(driverName), span.Tags[constants.DBTags["DB_TYPE"]])
		assert.Equal(t, query, span.Tags[constants.DBTags["DB_STATEMENT"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		tp.Reset()
	})

	t.Run("testStatementExecContext", func(t *testing.T) {
		tp := trace.New()

		// Create the parent span
		ctx := context.Background()
		parentSpan, ctx := opentracing.StartSpanFromContext(ctx, "parentSpan")
		parentSpanRaw, _ := tracer.GetRaw(parentSpan)

		stmt, err := s.db.PrepareContext(ctx, query)
		assert.Nil(t, err)
		assert.NotNil(t, stmt)

		// Check prepare event is not recorded
		assert.Len(t, tp.Recorder.GetSpans(), 1)

		stmt.ExecContext(ctx, args...)

		// Get the span created for statement query call
		span := tp.Recorder.GetSpans()[1]

		// Check parent span is set
		assert.Equal(t, parentSpanRaw.Context.SpanID, span.ParentSpanID)

		assert.Equal(t, constants.ClassNames[driverName], span.ClassName)
		assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
		assert.Equal(t, operationToType[strings.ToLower(strings.Split(query, " ")[0])], span.Tags[constants.SpanTags["OPERATION_TYPE"]])
		assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
		assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
		assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		assert.Equal(t, strings.ToLower(driverName), span.Tags[constants.DBTags["DB_TYPE"]])
		assert.Equal(t, query, span.Tags[constants.DBTags["DB_STATEMENT"]])
		assert.Equal(t, strings.ToUpper(strings.Split(query, " ")[0]), span.Tags[constants.DBTags["DB_STATEMENT_TYPE"]])
		tp.Reset()
	})

}
