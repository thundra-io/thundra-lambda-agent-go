package tracer

import (
	"errors"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
	"reflect"
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
	className      string
	operationTypes []string
	tags           map[string]interface{}
}

func (o *Operation) matches(span *spanImpl) bool {
	var matched = true

	if o.className != "" {
		matched = o.className == span.raw.ClassName
	}

	if matched && len(o.operationTypes) > 0 {
		matched = utils.Contains(o.operationTypes, o.getOperationType(span))
	}

	if matched && len(o.tags) > 0 {
		for key, value := range o.tags {
			rt := reflect.TypeOf(value)

			if (rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array) &&
				utils.Contains(value.([]string), span.raw.GetTag(key).(string)) {
				matched = false
				break
			} else if (rt.Kind() != reflect.Slice || rt.Kind() != reflect.Array) &&
				span.raw.GetTag(key) != value {
				matched = false
				break
			}
		}
	}

	return matched

}

func (o *Operation) getOperationType(span *spanImpl) string {
	return span.raw.GetTag(constants.SpanTags["OPERATION_TYPE"]).(string)
}
