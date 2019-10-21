package tracer

type TagInjectorSpanListener struct {
	tags           map[string]interface{}
	injectOnFinish bool
}

func (t *TagInjectorSpanListener) OnSpanStarted(span *spanImpl) {
	if !t.injectOnFinish {
		t.injectTags(span)
	}
}

func (t *TagInjectorSpanListener) OnSpanFinished(span *spanImpl) {
	if t.injectOnFinish {
		t.injectTags(span)
	}
}

func (t *TagInjectorSpanListener) injectTags(span *spanImpl) {
	if t.tags == nil {
		return
	}

	for k, v := range t.tags {
		span.SetTag(k, v)
	}
}

func (t *TagInjectorSpanListener) PanicOnError() bool {
	return false
}

func NewTagInjectorSpanListener(config map[string]interface{}) ThundraSpanListener {
	listener := &TagInjectorSpanListener{tags: map[string]interface{}{}, injectOnFinish:  false}
	if injectOnFinish, ok := config["injectOnFinish"].(bool); ok {
		listener.injectOnFinish = injectOnFinish
	}
	if tags, ok := config["tags"].(map[string]interface{}); ok {
		listener.tags = tags
	}

	return listener
}
