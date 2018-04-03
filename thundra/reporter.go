package thundra

import (
	"sync"
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
	"os"
	"io/ioutil"
)

type Reporter interface {
	Collect(msg interface{})
	Report(apiKey string)
	Clear()
}

type reporterImpl struct {
	messageQueue []interface{}
}

var shouldSendAsync string
var mutex = &sync.Mutex{}

func init() {
	shouldSendAsync = os.Getenv(thundraLambdaPublishCloudwatchEnable)
}

func (c *reporterImpl) Collect(msg interface{}) {
	defer mutex.Unlock()
	mutex.Lock()
	if shouldSendAsync == "true" {
		sendAsync(msg)
		return
	}
	c.messageQueue = append(c.messageQueue, msg)
}

func (c *reporterImpl) Report(apiKey string) {
	if shouldSendAsync == "false" || shouldSendAsync == "" {
		sendHttpReq(c.messageQueue, apiKey)
	}
}

func (c *reporterImpl) Clear() {
	c.messageQueue = c.messageQueue[:0]
}

func sendAsync(msg interface{}) {
	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Sending ASYNC request to Thundra collector")
	fmt.Println(string(b))
}

func sendHttpReq(mesageQueue interface{}, apiKey string) {
	b, _ := json.Marshal(&mesageQueue)
	req, _ := http.NewRequest("POST", collectorUrl, bytes.NewBuffer(b))
	req.Header.Set("Authorization", "ApiKey "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("Sending HTTP request to Thundra collector")
	fmt.Println(mesageQueue)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	resp.Body.Close()
}
