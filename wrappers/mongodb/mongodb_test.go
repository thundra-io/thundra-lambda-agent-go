package thundramongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	Author string
	Text   string
}

func TestCommandInsert(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Get client with traced monitor
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(NewCommandMonitor()))

	collection := client.Database("test").Collection("posts")
	post := Post{"Mike", "My first blog!"}

	// Insert post to mongodb
	collection.InsertOne(context.TODO(), post)

	// Get the span created
	span := tp.Recorder.GetSpans()[0]

	assert.Equal(t, "INSERT", span.OperationName)
	assert.Equal(t, constants.ClassNames["MONGODB"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)

	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.DBTags["DB_HOST"]])
	assert.Equal(t, "27017", span.Tags[constants.DBTags["DB_PORT"]])
	assert.Equal(t, "test", span.Tags[constants.DBTags["DB_INSTANCE"]])

	assert.Equal(t, []string{""}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	assert.Equal(t, "INSERT", span.Tags[constants.MongoDBTags["MONGODB_COMMAND_NAME"]])
	assert.Equal(t, "posts", span.Tags[constants.MongoDBTags["MONGODB_COLLECTION"]])

	// Unmarshal command to read document fields
	command := bson.M{}
	bson.UnmarshalExtJSON([]byte(span.Tags[constants.MongoDBTags["MONGODB_COMMAND"]].(string)), false, &command)

	assert.Equal(t, "posts", command["insert"])
	assert.Equal(t, "Mike", command["documents"].(primitive.A)[0].(primitive.M)["author"])
	assert.Equal(t, "My first blog!", command["documents"].(primitive.A)[0].(primitive.M)["text"])

	tp.Reset()
}

func TestCommandUpdate(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Get client with traced monitor
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(NewCommandMonitor()))

	collection := client.Database("test").Collection("posts")

	// Create a filter for update
	filter := bson.D{{"name", "Mike"}}

	// Create update document
	update := bson.D{
		{"$set", bson.D{
			{"text", "My edited blog post!"},
		}},
	}

	// Update document with filter
	collection.UpdateOne(context.TODO(), filter, update)

	// Get the span created
	span := tp.Recorder.GetSpans()[0]

	assert.Equal(t, "UPDATE", span.OperationName)
	assert.Equal(t, constants.ClassNames["MONGODB"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)

	assert.Equal(t, "WRITE", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.DBTags["DB_HOST"]])
	assert.Equal(t, "27017", span.Tags[constants.DBTags["DB_PORT"]])
	assert.Equal(t, "test", span.Tags[constants.DBTags["DB_INSTANCE"]])

	assert.Equal(t, "UPDATE", span.Tags[constants.MongoDBTags["MONGODB_COMMAND_NAME"]])
	assert.Equal(t, "posts", span.Tags[constants.MongoDBTags["MONGODB_COLLECTION"]])

	assert.Equal(t, []string{""}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	// Unmarshal command to read document fields
	command := bson.M{}
	bson.UnmarshalExtJSON([]byte(span.Tags[constants.MongoDBTags["MONGODB_COMMAND"]].(string)), false, &command)

	assert.Equal(t, "posts", command["update"])
	assert.Equal(t, primitive.M(primitive.M{"name": "Mike"}), command["updates"].(primitive.A)[0].(primitive.M)["q"])

	tp.Reset()

}

func TestCommandFailed(t *testing.T) {
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Get client with traced monitor
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(NewCommandMonitor()))

	// Try to perform an unknown command
	client.Database("test").RunCommand(context.TODO(), bson.D{{"unknown_command", 1}})

	// Get the span created
	span := tp.Recorder.GetSpans()[0]

	assert.Equal(t, "UNKNOWN_COMMAND", span.OperationName)
	assert.Equal(t, constants.ClassNames["MONGODB"], span.ClassName)
	assert.Equal(t, constants.DomainNames["DB"], span.DomainName)

	assert.Equal(t, "", span.Tags[constants.SpanTags["OPERATION_TYPE"]])
	assert.Equal(t, "localhost", span.Tags[constants.DBTags["DB_HOST"]])
	assert.Equal(t, "27017", span.Tags[constants.DBTags["DB_PORT"]])
	assert.Equal(t, "test", span.Tags[constants.DBTags["DB_INSTANCE"]])

	assert.Equal(t, "UNKNOWN_COMMAND", span.Tags[constants.MongoDBTags["MONGODB_COMMAND_NAME"]])
	assert.Equal(t, "", span.Tags[constants.MongoDBTags["MONGODB_COLLECTION"]])

	assert.Equal(t, []string{""}, span.Tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]])
	assert.Equal(t, constants.AwsLambdaApplicationDomain, span.Tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]])
	assert.Equal(t, constants.AwsLambdaApplicationClass, span.Tags[constants.SpanTags["TRIGGER_CLASS_NAME"]])
	assert.Equal(t, true, span.Tags[constants.SpanTags["TOPOLOGY_VERTEX"]])

	// Error should be set as command does not exist
	assert.NotNil(t, span.Tags["error"])

	tp.Reset()
}

func TestCommandMasked(t *testing.T) {
	config.MaskMongoDBCommand = true
	// Initilize trace plugin to set GlobalTracer of opentracing
	tp := trace.New()

	// Get client with traced monitor
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(NewCommandMonitor()))

	client.Database("test").RunCommand(context.TODO(), bson.D{{"ping", 1}})

	// Get the span created
	span := tp.Recorder.GetSpans()[0]

	assert.Nil(t, span.Tags["error"])
	assert.Nil(t, span.Tags[constants.MongoDBTags["MONGODB_COMMAND"]])
	tp.Reset()
}
