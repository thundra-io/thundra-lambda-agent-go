package redis

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

func BeforeCall(span *tracer.RawSpan, host string, port string, commandName string, command string) {
	span.ClassName = constants.ClassNames["REDIS"]
	span.DomainName = constants.DomainNames["CACHE"]
	if len(commandName) == 0 {
		commandName = strings.Split(command, " ")[0]
	}
	commandName = strings.ToUpper(commandName)
	// Set span tags
	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:      constants.RedisCommandTypes[commandName],
		constants.DBTags["DB_INSTANCE"]:           host,
		constants.DBTags["DB_STATEMENT_TYPE"]:     constants.RedisCommandTypes[commandName],
		constants.DBTags["DB_TYPE"]:               "redis",
		constants.RedisTags["REDIS_HOST"]:         host,
		constants.RedisTags["REDIS_COMMAND_TYPE"]: commandName,
	}

	tags[constants.SpanTags["TRIGGER_DOMAIN_NAME"]] = constants.AwsLambdaApplicationDomain
	tags[constants.SpanTags["TRIGGER_CLASS_NAME"]] = constants.AwsLambdaApplicationClass
	tags[constants.SpanTags["TRIGGER_OPERATION_NAMES"]] = []string{application.FunctionName}
	tags[constants.SpanTags["TOPOLOGY_VERTEX"]] = true

	if !config.MaskRedisCommand {
		tags[constants.DBTags["DB_STATEMENT"]] = command
		tags[constants.RedisTags["REDIS_COMMAND"]] = command
	}

	span.Tags = tags
}

func GetRedisCommand(commandName string, args ...interface{}) string {
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
