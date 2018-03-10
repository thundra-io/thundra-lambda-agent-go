package thundra

import (
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
	"os"
)

var ShouldSendAsync string

func init() {
	ShouldSendAsync = os.Getenv(ThundraLambdaPublishCloudwatchEnable)
}

func sendReport(collector collector, msg Message) {
	if ShouldSendAsync == "true" {
		sendAsync(msg)
	} else {
		collector.collect(msg)
	}
}

func sendAsync(msg Message) {
	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Sending ASYNC request")
	fmt.Println(string(b))
}

func sendHttpReq(msg []Message) {
	b, _ := json.Marshal(&msg)
	req, _ := http.NewRequest("POST", collectorUrl, bytes.NewBuffer(b))
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
	fmt.Println("response Status:", resp.Status)
	//TODO if resp.status == 401 unauthorized : ApiKey is missing
	fmt.Println("ApiKey:", ApiKey)
	resp.Body.Close()
}
