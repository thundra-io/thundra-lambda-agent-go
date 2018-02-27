package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"time"
	"ThundraGo/thundra"
)

type MyEvent struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

func HandleLambdaEvent(event MyEvent) (MyResponse, error) {
	time.Sleep(15 * time.Millisecond)
	testFunc()
	return MyResponse{Message: fmt.Sprintf("%s is %d years old!", event.Name, event.Age)}, nil
}

func testFunc() {
	time.Sleep(12 * time.Millisecond)
}

func main() {
	th := thundra.GetInstance([]string{"trace"})

	lambda.Start(thundra.WrapLambdaHandler(HandleLambdaEvent, th))
	lambda.Start(HandleLambdaEvent)
}
