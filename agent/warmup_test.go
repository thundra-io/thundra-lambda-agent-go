package agent

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckAndHandleWarmupRequest(t *testing.T) {
	type testStruct struct {
		age  int
		name string
	}
	structType := reflect.TypeOf(testStruct{})
	assert.True(t, checkAndHandleWarmupRequest(testStruct{}, structType))
	assert.False(t, checkAndHandleWarmupRequest(testStruct{age: 10, name: "Thundra"}, structType))

	strType := reflect.TypeOf("")
	assert.True(t, checkAndHandleWarmupRequest("", strType))
	assert.True(t, checkAndHandleWarmupRequest("#warmup wait=100", strType))
	assert.False(t, checkAndHandleWarmupRequest("notWarmup", strType))

}

func TestIsZeroEvent(t *testing.T) {

	strType := reflect.TypeOf("")
	assert.False(t, isZeroEvent("NotZero", strType))
	assert.True(t, isZeroEvent("", strType))

	type testStruct struct {
		age  int
		name string
	}
	structType := reflect.TypeOf(testStruct{})

	assert.False(t, isZeroEvent(testStruct{age: 10, name: "Thundra"}, structType))
	assert.True(t, isZeroEvent(testStruct{}, structType))
}
