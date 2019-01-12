package ext

import (
	ot "github.com/opentracing/opentracing-go"
)

func ClassName(className string) ot.StartSpanOption {
	return ot.Tag {
		Key:   ClassNameKey,
		Value: className,
	}
}

func DomainName(domainName string) ot.StartSpanOption {
	return ot.Tag {
		Key: DomainNameKey,
		Value: domainName,
	}
}
