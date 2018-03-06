package thundra

import (
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
	"os"
	"thundra-agent-go/constants"
)

var ShouldSendAsync string

func init() {
	ShouldSendAsync = os.Getenv(constants.Thundra_Lambda_Publish_Cloudwatch_Enable)
}

func sendReport(msg Message) {
	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	if ShouldSendAsync == "true" {
		sendAsync(b)
	} else {
		collect(msg)
	}
}

func sendAsync(msg []byte) {
	fmt.Println("Sending ASYNC request")
	fmt.Println(string(msg))
}

func sendHttpReq(msg []Message) {
	b,_ := json.Marshal(&msg)
	req, _ := http.NewRequest("POST", constants.CollectorUrl, bytes.NewBuffer(b))
	req.Header.Set("Authorization", "ApiKey "+ApiKey)
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("Sending HTTP request")
	fmt.Println(msg)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("ApiKey:", ApiKey)
}
