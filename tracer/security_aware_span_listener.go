package tracer

import (
	"errors"
	"reflect"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

var defaultSecurityMessage = "Operation was blocked due to security configuration"

type SecurityAwareSpanListener struct {
	block     bool
	whitelist []Operation
	blacklist []Operation
}

func (s *SecurityAwareSpanListener) OnSpanStarted(span *spanImpl) {
	if !s.isExternalOperation(span) {
		return
	}

	for _, op := range s.blacklist {
		if op.matches(span) {
			s.handleSecurityIssue(span)
			return
		}
	}

	if len(s.whitelist) > 0 {
		for _, op := range s.whitelist {
			if op.matches(span) {
				return
			}
		}
		s.handleSecurityIssue(span)
	}

}

func (s *SecurityAwareSpanListener) OnSpanFinished(span *spanImpl) {
	return
}

func (s *SecurityAwareSpanListener) PanicOnError() bool {
	return true
}

func (s *SecurityAwareSpanListener) handleSecurityIssue(span *spanImpl) {
	if s.block {
		err := errors.New(defaultSecurityMessage)
		span.SetTag(constants.SecurityTags["BLOCKED"], true)
		span.SetTag(constants.SecurityTags["VIOLATED"], true)
		utils.SetSpanError(span, err)
		panic(err)
	} else {
		span.SetTag(constants.SecurityTags["VIOLATED"], true)
	}
}

func (s *SecurityAwareSpanListener) isExternalOperation(span *spanImpl) bool {
	return span.raw.GetTag(constants.SpanTags["TOPOLOGY_VERTEX"]) == true
}

type Operation struct {
	className string
	tags      map[string]interface{}
}

func (o *Operation) matches(span *spanImpl) bool {
	var matched = true

	if o.className != "" {
		matched = o.className == span.raw.ClassName
	}

	if matched && len(o.tags) > 0 {
		for key, value := range o.tags {
			if span.raw.GetTag(key) != nil {

				rt := reflect.TypeOf(value)
				if rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
					switch value.(type) {
					case []string:
						matched = utils.StringContains(value.([]string), span.raw.GetTag(key).(string))
						break
					case []int64:
						matched = utils.Int64Contains(value.([]int64), span.raw.GetTag(key).(int64))
						break
					case []float64:
						matched = utils.Float64Contains(value.([]float64), span.raw.GetTag(key).(float64))
						break
					case []interface{}:
						matched = utils.Contains(value.([]interface{}), span.raw.GetTag(key))
						break
					default:
						panic("unexpected value type")

					}
				} else if (rt.Kind() != reflect.Slice || rt.Kind() != reflect.Array) &&
					span.raw.GetTag(key) != value {
					matched = false
					break
				}
			}
		}
	}

	return matched
}
