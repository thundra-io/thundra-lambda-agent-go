package thundraelastic

import (
	"context"
	"net/http"
	"testing"

	elasticv6 "github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/config"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/trace"
)

func getClient() *elasticv6.Client {
	var httpClient = Wrap(&http.Client{})
	var client, _ = elasticv6.NewClient(
		elasticv6.SetURL("http://localhost:9200"),
		elasticv6.SetHttpClient(httpClient),
		elasticv6.SetSniff(false),
		elasticv6.SetHealthcheck(false),
	)
	return client
}

func TestCreateIndex(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	client := getClient()

	// Create a new index
	client.CreateIndex("twitter").Do(context.Background())

	// Get the span created
	span := tp.Recorder.GetSpans()[0]

	assert.Equal(t, constants.ClassNames["ELASTICSEARCH"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, "/twitter", span.OperationName)
	assert.ElementsMatch(t, []string{"localhost:9200"}, span.Tags[constants.EsTags["ES_HOSTS"]])
	assert.Equal(t, "PUT", span.Tags[constants.EsTags["ES_METHOD"]])
	assert.Equal(t, "/twitter", span.Tags[constants.EsTags["ES_URI"]])
	assert.Equal(t, "", span.Tags[constants.EsTags["ES_BODY"]])

	assert.Equal(t, "elasticsearch", span.Tags[constants.DBTags["DB_TYPE"]])

	assert.Equal(t, "PUT", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, []string{""}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	tp.Reset()
}

func TestIndex(t *testing.T) {
	config.EsIntegrationUrlPathDepth = 3
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	client := getClient()

	// Add Document
	client.Index().
		Index("twitter").Type("_docs").Id("1").
		BodyString(`{"user": "kimchy", "message": "trying out Elasticsearch"}`).Do(context.Background())

	// Get the span created
	span := tp.Recorder.GetSpans()[0]

	assert.Equal(t, constants.ClassNames["ELASTICSEARCH"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, "/twitter/_docs/1", span.OperationName)
	assert.ElementsMatch(t, []string{"localhost:9200"}, span.Tags[constants.EsTags["ES_HOSTS"]])
	assert.Equal(t, "PUT", span.Tags[constants.EsTags["ES_METHOD"]])
	assert.Equal(t, "/twitter/_docs/1", span.Tags[constants.EsTags["ES_URI"]])
	assert.Equal(t, "/twitter/_docs/1", span.Tags[constants.EsTags["ES_NORMALIZED_URI"]])
	assert.Equal(t, `{"user": "kimchy", "message": "trying out Elasticsearch"}`, span.Tags[constants.EsTags["ES_BODY"]])

	assert.Equal(t, "elasticsearch", span.Tags[constants.DBTags["DB_TYPE"]])

	assert.Equal(t, "PUT", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, []string{""}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	tp.Reset()
}

func TestGetDoc(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	client := getClient()

	// Get Document
	client.Get().Index("twitter").Type("_docs").Id("1").Do(context.Background())

	// Get the span created
	span := tp.Recorder.GetSpans()[0]

	assert.Equal(t, constants.ClassNames["ELASTICSEARCH"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, "/twitter/_docs/1", span.OperationName)
	assert.ElementsMatch(t, []string{"localhost:9200"}, span.Tags[constants.EsTags["ES_HOSTS"]])
	assert.Equal(t, "GET", span.Tags[constants.EsTags["ES_METHOD"]])
	assert.Equal(t, "/twitter/_docs/1", span.Tags[constants.EsTags["ES_URI"]])
	assert.Nil(t, span.Tags[constants.EsTags["ES_BODY"]])

	assert.Equal(t, "elasticsearch", span.Tags[constants.DBTags["DB_TYPE"]])

	assert.Equal(t, "GET", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, []string{""}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	tp.Reset()
}

func TestDeleteDoc(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	client := getClient()

	// Delete Document
	client.Delete().Index("twitter").Type("_docs").Id("1").Do(context.Background())

	// Get the span created
	span := tp.Recorder.GetSpans()[0]

	assert.Equal(t, constants.ClassNames["ELASTICSEARCH"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, "/twitter/_docs/1", span.OperationName)
	assert.ElementsMatch(t, []string{"localhost:9200"}, span.Tags[constants.EsTags["ES_HOSTS"]])
	assert.Equal(t, "DELETE", span.Tags[constants.EsTags["ES_METHOD"]])
	assert.Equal(t, "/twitter/_docs/1", span.Tags[constants.EsTags["ES_URI"]])
	assert.Nil(t, span.Tags[constants.EsTags["ES_BODY"]])

	assert.Equal(t, "elasticsearch", span.Tags[constants.DBTags["DB_TYPE"]])

	assert.Equal(t, "DELETE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, []string{""}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	tp.Reset()
}

func TestRefresh(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	client := getClient()

	// Delete Document
	client.Delete().Index("twitter").Type("_docs").Id("1").Do(context.Background())

	// Get the span created
	span := tp.Recorder.GetSpans()[0]

	assert.Equal(t, constants.ClassNames["ELASTICSEARCH"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, "/twitter/_docs/1", span.OperationName)
	assert.ElementsMatch(t, []string{"localhost:9200"}, span.Tags[constants.EsTags["ES_HOSTS"]])
	assert.Equal(t, "DELETE", span.Tags[constants.EsTags["ES_METHOD"]])
	assert.Equal(t, "/twitter/_docs/1", span.Tags[constants.EsTags["ES_URI"]])
	assert.Nil(t, span.Tags[constants.EsTags["ES_BODY"]])

	assert.Equal(t, "elasticsearch", span.Tags[constants.DBTags["DB_TYPE"]])

	assert.Equal(t, "DELETE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, []string{""}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	tp.Reset()
}

func TestErrorNotExistentURL(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	var httpClient = Wrap(&http.Client{})
	var client, _ = elasticv6.NewClient(
		elasticv6.SetURL("http://localhost:9201"),
		elasticv6.SetHttpClient(httpClient),
		elasticv6.SetSniff(false),
		elasticv6.SetHealthcheck(false),
	)

	// Delete Document
	client.Delete().Index("twitter").Type("_docs").Id("1").Do(context.Background())

	// Get the span created
	span := tp.Recorder.GetSpans()[0]

	assert.Equal(t, constants.ClassNames["ELASTICSEARCH"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, "/twitter/_docs/1", span.OperationName)
	assert.ElementsMatch(t, []string{"localhost:9201"}, span.Tags[constants.EsTags["ES_HOSTS"]])
	assert.Equal(t, "DELETE", span.Tags[constants.EsTags["ES_METHOD"]])
	assert.Equal(t, "/twitter/_docs/1", span.Tags[constants.EsTags["ES_URI"]])
	assert.Nil(t, span.Tags[constants.EsTags["ES_BODY"]])

	assert.Equal(t, "elasticsearch", span.Tags[constants.DBTags["DB_TYPE"]])

	assert.Equal(t, "DELETE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, []string{""}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	assert.True(t, span.Tags[constants.AwsError].(bool))
	assert.Equal(t, "OpError", span.Tags[constants.AwsErrorKind].(string))

	tp.Reset()
}

func TestMaskBody(t *testing.T) {
	config.MaskEsBody = true
	config.EsIntegrationUrlPathDepth = 2
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	var httpClient = Wrap(&http.Client{})
	var client, _ = elasticv6.NewClient(
		elasticv6.SetURL("http://localhost:9200"),
		elasticv6.SetHttpClient(httpClient),
		elasticv6.SetSniff(false),
		elasticv6.SetHealthcheck(false),
	)

	// Delete Document
	client.Delete().Index("twitter").Type("_docs").Id("1").Do(context.Background())

	// Get the span created
	span := tp.Recorder.GetSpans()[0]

	assert.Equal(t, constants.ClassNames["ELASTICSEARCH"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)
	assert.Equal(t, "/twitter/_docs", span.OperationName)
	assert.ElementsMatch(t, []string{"localhost:9200"}, span.Tags[constants.EsTags["ES_HOSTS"]])
	assert.Equal(t, "DELETE", span.Tags[constants.EsTags["ES_METHOD"]])
	assert.Equal(t, "/twitter/_docs/1", span.Tags[constants.EsTags["ES_URI"]])
	assert.Equal(t, "/twitter/_docs", span.Tags[constants.EsTags["ES_NORMALIZED_URI"]])
	assert.Equal(t, "elasticsearch", span.Tags[constants.DBTags["DB_TYPE"]])
	assert.Equal(t, "DELETE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])

	assert.Equal(t, []string{""}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	assert.Nil(t, span.Tags[constants.EsTags["ES_BODY"]])

	tp.Reset()
}
