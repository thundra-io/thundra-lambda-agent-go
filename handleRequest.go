package main

import (
	"ThundraGo/thundra"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

func HandleLambdaEvent(event MyEvent) (MyResponse,error) {
	return MyResponse{Message: fmt.Sprintf("%s is %d years old!", event.Name, event.Age)},nil
}

func MyTestFunc(event MyEvent){

	//th := thundra.New([]string{"trace"},[]string{"trace"})

	//th.executePreHooks(args)
	/*executePreHooksValue := reflect.ValueOf(th.executePreHooks)
	executePreHooksValue.Call(args)*/
}

func main() {
	plugins := []string{"trace"}
	th := thundra.New(plugins)
	lambda.Start(thundra.Handle(HandleLambdaEvent, th))
}
