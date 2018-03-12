package thundra

import (
	"sync"
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
	"os"
)

type reporter interface {
	collect(msg interface{})
	report()
	clear()
}

type reporterImpl struct {
	messageQueue []interface{}
}

var shouldSendAsync = os.Getenv(ThundraLambdaPublishCloudwatchEnable)
var mutex = &sync.Mutex{}

func (c *reporterImpl) collect(msg interface{}) {
	defer mutex.Unlock()
	mutex.Lock()
	if shouldSendAsync == "true" {
		sendAsync(msg)
		return
	}
	c.messageQueue = append(c.messageQueue, msg)
}

func (c *reporterImpl) report() {
	if shouldSendAsync == "false" {
		sendHttpReq(c.messageQueue)
	}
}

func (c *reporterImpl) clear() {
	c.messageQueue = c.messageQueue[:0]
}

func sendAsync(msg interface{}) {
	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Sending ASYNC request")
	fmt.Println(string(b))
}

func sendHttpReq(msg []interface{}) {
	b, _ := json.Marshal(&msg)
	req, _ := http.NewRequest("POST", collectorUrl, bytes.NewBuffer(b))
	req.Header.Set("Authorization", "ApiKey "+apiKey)
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
	resp.Body.Close()
}
