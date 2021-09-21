# Lambda Go Agent [![OpenTracing Badge](https://img.shields.io/badge/OpenTracing-enabled-blue.svg)](http://opentracing.io) [![Thundra CI Check](https://github.com/thundra-io/thundra-lambda-agent-go/actions/workflows/build.yml/badge.svg)](https://github.com/thundra-io/thundra-lambda-agent-go/actions/workflows/build.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/thundra-io/thundra-lambda-agent-go)](https://goreportcard.com/report/github.com/thundra-io/thundra-lambda-agent-go)

Trace your AWS lambda functions with async monitoring by [Thundra](https://www.thundra.io/)!

Check out [Thundra docs](https://apm.docs.thundra.io) for more information.

### Usage

In order to trace your lambda usages with Thundra all you need to do is wrap your function.
```go
import "github.com/thundra-io/thundra-lambda-agent-go/v2/thundra"	// with go modules enabled (GO111MODULE=on or outside GOPATH) for version >= v2.3.1
import "github.com/thundra-io/thundra-lambda-agent-go/thundra"         // with go modules disabled
```

```go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	// thundra go agent import here
)

// Your lambda handler
func handler() (string, error) {
	return "Hello, Thundra!", nil
}

func main() {
	// Wrap your lambda handler with Thundra
	lambda.Start(thundra.Wrap(handler))
}
```
Later just build and deploy your executable to AWS as regular. Test your function on lambda console and visit [Thundra](https://console.thundra.io/) to observe your function metrics.

#### Environment variables

| Name                                                  |   Type     |    Default Value                |
|:------------------------------------------------------|:----------:|:-------------------------------:|
| thundra_applicationProfile                            |   string   |    default                      |
| thundra_agent_lambda_disable                          |   bool     |    false                        |
| thundra_agent_lambda_timeout_margin                   |   number   |    200                          |
| thundra_agent_lambda_report_rest_baseUrl              |   string   |    https://api.thundra.io/v1    |
| thundra_agent_lambda_trace_disable                    |   bool     |    false                        |
| thundra_agent_lambda_metric_disable                   |   bool     |    false                        |
| thundra_agent_lambda_log_disable                      |   bool     |    false                        |
| thundra_log_logLevel                                  |   string   |    TRACE                        |
| thundra_agent_lambda_trace_request_skip               |   bool     |    false                        |
| thundra_agent_lambda_trace_response_skip              |   bool     |    false                        |
| thundra_agent_lambda_report_rest_trustAllCertificates |   bool     |    false                        |
| thundra_agent_lambda_debug_enable                     |   bool     |    false                        |
| thundra_agent_lambda_warmup_warmupAware               |   bool     |    false                        |


### Async Monitoring

Check out our [docs](https://docs.thundra.io/docs/how-to-setup-async-monitoring) to see how to configure Thundra and async monitoring to visualize your functions in [Thundra](https://www.thundra.io/).

## Warmup Support
You can cut down cold starts easily by deploying our lambda function [`thundra-lambda-warmup`](https://github.com/thundra-io/thundra-lambda-warmup).

Our agent handles warmup requests automatically so you don't need to make any code changes.

You just need to deploy `thundra-lambda-warmup` once, then you can enable warming up for your lambda by
* setting its environment variable `thundra_agent_lambda_warmup_warmupAware` **true** OR
* adding its name to `thundra-lambda-warmup`'s environment variable `thundra_lambda_warmup_function`.

Check out [this part](https://thundra.readme.io/docs/how-to-warmup) in our docs for more information.

## Thundra Go Agent Integrations

### AWS SDK

Thundra's Go agent provides you with the capability to trace AWS SDK by wrapping the session object that AWS SDK provides.
You can easily start using it by simply wrapping your session objects as shown in the following code:

```go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
	thundraaws "github.com/thundra-io/thundra-lambda-agent-go/wrappers/aws"
)

// Your lambda handler
func handler() (string, error) {
	// Create a new session object
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Wrap it using the thundraaws.Wrap method
	sess = thundraaws.Wrap(sess)

	// Create a new client using the wrapped session
	svc := dynamodb.New(sess)

	// Use the client as normal, Thundra will automatically
	// create spans for the AWS SDK calls
	svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"AlbumTitle": {
				S: aws.String("Somewhat Famous"),
			},
			"Artist": {
				S: aws.String("No One You Know"),
			},
			"SongTitle": {
				S: aws.String("Call Me Today"),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String("Music"),
	})

	return "Hello, Thundra!", nil
}

func main() {
	// Wrap your lambda handler with Thundra
	lambda.Start(thundra.Wrap(handler))
}
```

### Mongo SDK

```go
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
	thundramongo "github.com/thundra-io/thundra-lambda-agent-go/wrappers/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	Author string
	Text   string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(thundramongo.NewCommandMonitor()))

	collection := client.Database("test").Collection("posts")
	post := Post{"Mike", "My first blog!"}

	// Insert post to mongodb
	collection.InsertOne(context.TODO(), post)

	return events.APIGatewayProxyResponse{
		Body:       "res",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(thundra.Wrap(handler))
}
```

### Mysql SDK

```go
package main

import (
	"context"
	"database/sql"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-sql-driver/mysql"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
	thundrardb "github.com/thundra-io/thundra-lambda-agent-go/wrappers/rdb"
)

func handler(ctx context.Context) {

	// Register wrapped driver for tracing
	// Note that driver name registered should be different than "mysql" as it is already registered on init method of
	// this package otherwise it panics.
	sql.Register("thundra-mysql", thundrardb.Wrap(&mysql.MySQLDriver{}))

	// Get the database handle with registered driver name
	db, err := sql.Open("thundra-mysql", "user:userpass@tcp(docker.for.mac.localhost:3306)/db")

	if err != nil {
		// Just for example purpose. You should use proper error handling instead of panic
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		thundra.Logger.Error(err)
	}
	defer rows.Close()
}

func main() {
	lambda.Start(thundra.Wrap(handler))
}
```

### PostgreSQL SDK

```go
package main

import (
	"context"
	"database/sql"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lib/pq"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
	thundrardb "github.com/thundra-io/thundra-lambda-agent-go/wrappers/rdb"
)

func handler(ctx context.Context) {
	// Register wrapped driver for tracing
	// Note that driver name registered should be different than "postgres"
	// as it is already registered on init method of this package
	// otherwise it panics.
	sql.Register("thundra-postgres", thundrardb.Wrap(&pq.Driver{}))

	// Get the database handle with registered driver name
	db, err := sql.Open("thundra-postgres", "postgres://user:userpass@docker.for.mac.localhost:5432/db?sslmode=disable")

	_, err = db.Exec("CREATE table IF NOT EXISTS test(id int, type text)")

	if err != nil {
		// Just for example purpose. You should use proper error handling instead of panic
		panic(err.Error())
	}
	defer db.Close()

	// Perform a query
	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		thundra.Logger.Error(err)
	}
	defer rows.Close()
}

func main() {
	lambda.Start(thundra.Wrap(handler))
}

```

### Redis SDK

```go
package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
	thundraredigo "github.com/thundra-io/thundra-lambda-agent-go/wrappers/redis/redigo"
)

func test() {
	conn, err := thundraredigo.Dial("tcp", "docker.for.mac.localhost:6379")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	ret, _ := conn.Do("SET", "mykey", "hello")
	fmt.Printf("%s\n", ret)

	ret, _ = conn.Do("GET", "mykey")
	fmt.Printf("%s\n", ret)

}

func main() {
	lambda.Start(thundra.Wrap(test))
}
```

### Http SDK

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
	thundrahttp "github.com/thundra-io/thundra-lambda-agent-go/wrappers/http"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	thundraHttpClient := thundrahttp.Wrap(http.Client{})

	resp, err := thundraHttpClient.Get("URL")
	if err == nil {
		fmt.Println(resp)
	}
	fmt.Println(err)

	return events.APIGatewayProxyResponse{
		Body:       "res",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(thundra.Wrap(handler))
}
```

### Elasticsearch SDK

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thundra-io/thundra-lambda-agent-go/thundra"
	thundraelastic "github.com/thundra-io/thundra-lambda-agent-go/wrappers/elastic/olivere"
	"gopkg.in/olivere/elastic.v6"
)

func getClient() *elastic.Client {

	var httpClient = thundraelastic.Wrap(&http.Client{})
	var client, _ = elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetHttpClient(httpClient),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	return client
}

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type Tweet struct {
	User     string   `json:"user"`
	Message  string   `json:"message"`
	Retweets int      `json:"retweets"`
	Image    string   `json:"image,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Location string   `json:"location,omitempty"`
}

func handler() {

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client := getClient()

	// Index a tweet (using JSON serialization)
	tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
	put1, err := client.Index().
		Index("thundra").
		Type("users").
		Id("id").
		BodyJson(tweet1).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// Get tweet with specified ID
	get1, err := client.Get().
		Index("twitter").
		Type("tweet").
		Id("1").
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("twitter").Do(ctx)
	if err != nil {
		panic(err)
	}

	// Search with a term query
	termQuery := elastic.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index("twitter").   // search in index "twitter"
		Query(termQuery).   // specify the query
		Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10).   // take documents 0-9
		Pretty(true).       // pretty print request and response JSON
		Do(ctx)             // execute
	if err != nil {
		// Handle error
		panic(err)
	}

}

func main() {
	lambda.Start(thundra.Wrap(handler))
}
```
