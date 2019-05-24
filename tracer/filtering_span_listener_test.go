package tracer

import (
	"reflect"
	"testing"

	ot "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func TestCreation(t *testing.T) {

	f1 := ThundraSpanFilter{}
	f2 := ThundraSpanFilter{}
	filterer := ThundraSpanFilterer{}
	filterer.AddFilter(&f1)
	filterer.AddFilter(&f2)
	listener := LatencyInjectorSpanListener{Delay: 370}
	fsl := FilteringSpanListener{Listener: &listener, Filterer: &filterer}

	assert.Equal(t, &listener, fsl.Listener)
}

func TestFilters(t *testing.T) {

	f1 := ThundraSpanFilter{DomainName: "test"}
	filterer := ThundraSpanFilterer{}
	filterer.AddFilter(&f1)
	listener := ErrorInjectorSpanListener{}
	fsl := FilteringSpanListener{Listener: &listener, Filterer: &filterer}

	span := &spanImpl{}
	span.raw.DomainName = "test"
	errorOccured := false
	func() {
		defer func() {
			if recover() != nil {
				errorOccured = true
			}
		}()
		fsl.OnSpanStarted(span)
	}()

	assert.True(t, errorOccured)

	errorOccured = false
	span.raw.DomainName = "test2"
	func() {
		defer func() {
			if recover() != nil {
				errorOccured = true
			}
		}()
		fsl.OnSpanStarted(span)
	}()

	assert.False(t, errorOccured)
}

func TestTagsFromConfig(t *testing.T) {
	config := map[string]string{
		"listener":              "ErrorInjectorSpanListener",
		"config.injectOnFinish": "true",
		"filter1.className":     "AWS-SQS",
		"filter1.domainName":    "Messaging",
		"filter1.tag.test":      "3",
	}

	fsl := NewFilteringSpanListener(config).(*FilteringSpanListener)

	span := &spanImpl{}
	span.raw.Tags = ot.Tags{}
	span.raw.Tags["test"] = 3
	span.raw.ClassName = "AWS-SQS"
	span.raw.DomainName = "Messaging"
	errorOccured := false
	func() {
		defer func() {
			if recover() != nil {
				errorOccured = true
			}
		}()
		fsl.OnSpanFinished(span)
	}()

	assert.True(t, errorOccured)
}

func TestNewFilteringListenerFromConfig(t *testing.T) {
	config := map[string]string{
		"listener":               "ErrorInjectorSpanListener",
		"config.errorMessage":    "You have a very funny name!",
		"config.injectOnFinish":  "true",
		"config.injectCountFreq": "3",
		"filter1.className":      "AWS-SQS",
		"filter1.domainName":     "Messaging",
		"filter2.className":      "HTTP",
		"filter2.tag.http.host":  "foo.com",
	}

	fsl := NewFilteringSpanListener(config).(*FilteringSpanListener)

	assert.Equal(t, reflect.TypeOf(&ErrorInjectorSpanListener{}), reflect.TypeOf(fsl.Listener))
	assert.Equal(t, "You have a very funny name!", fsl.Listener.(*ErrorInjectorSpanListener).ErrorMessage)
	assert.Equal(t, true, fsl.Listener.(*ErrorInjectorSpanListener).InjectOnFinish)
	assert.Equal(t, int64(3), fsl.Listener.(*ErrorInjectorSpanListener).InjectCountFreq)

	assert.Equal(t, 2, len(fsl.Filterer.(*ThundraSpanFilterer).spanFilters))

	f1 := fsl.Filterer.(*ThundraSpanFilterer).spanFilters[0]
	f2 := fsl.Filterer.(*ThundraSpanFilterer).spanFilters[1]

	assert.ElementsMatch(t, []string{"AWS-SQS", "HTTP"}, []string{f1.(*ThundraSpanFilter).ClassName, f2.(*ThundraSpanFilter).ClassName})
	assert.ElementsMatch(t, []interface{}{"foo.com", nil}, []interface{}{f1.(*ThundraSpanFilter).Tags["http.host"], f2.(*ThundraSpanFilter).Tags["http.host"]})

}
