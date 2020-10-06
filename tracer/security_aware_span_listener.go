package tracer

import (
	"encoding/json"
	"errors"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/utils"
	logger "log"
)

var defaultSecurityMessage = "Operation was blocked due to security configuration"

type SecurityAwareSpanListener struct {
	block     bool
	whitelist *[]Operation
	blacklist *[]Operation
}

func (s *SecurityAwareSpanListener) OnSpanStarted(span *spanImpl) {
	if !s.isExternalOperation(span) {
		return
	}

	if s.blacklist != nil {
		for _, op := range *s.blacklist {
			if op.matches(span) {
				s.handleSecurityIssue(span)
				return
			}
		}
	}

	if s.whitelist != nil {
		for _, op := range *s.whitelist {
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
	ClassName string              `json:"className"`
	Tags      map[string][]string `json:"tags"`
}

func (o *Operation) matches(span *spanImpl) bool {
	var matched = true

	if o.ClassName != "" {
		matched = o.ClassName == "*" || o.ClassName == span.raw.ClassName
	}

	if matched && len(o.Tags) > 0 {
		for key, value := range o.Tags {
			if span.raw.GetTag(key) != nil {
				if utils.StringContains(value, "*") {
					continue
				}
				if !utils.StringContains(value, span.raw.GetTag(key).(string)) {
					matched = false
					break
				}
			}
		}
	}

	return matched
}

func NewSecurityAwareSpanListener(config map[string]interface{}) ThundraSpanListener {
	spanListener := &SecurityAwareSpanListener{}

	if block, ok := config["block"].(bool); ok {
		spanListener.block = block
	}

	if whitelist, ok := config["whitelist"].([]interface{}); ok {
		var wl []Operation
		for _, value := range whitelist {
			op := mapToOperation(value)
			wl = append(wl, op)
		}

		spanListener.whitelist = &wl
	}

	if blacklist, ok := config["blacklist"].([]interface{}); ok {
		var bl []Operation
		for _, value := range blacklist {
			op := mapToOperation(value)
			bl = append(bl, op)
		}
		spanListener.blacklist = &bl
	}

	return spanListener
}

func mapToOperation(opMap interface{}) Operation {
	jsonBody, err := json.Marshal(opMap)
	if err != nil {
		logger.Println("Error on marshal security operation:", err)
		return Operation{}
	}

	op := Operation{}
	if err := json.Unmarshal(jsonBody, &op); err != nil {
		logger.Println("Error on marshal security operation:", err)
		return op
	}

	return op
}
